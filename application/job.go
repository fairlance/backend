package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fairlance/backend/models"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func getAllJobs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		jobs, err := appContext.JobRepository.GetAllJobs()
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, jobs)
	})
}

func getAllJobsForClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var user = context.Get(r, "user").(*models.User)
		jobs, err := appContext.JobRepository.GetAllJobsForClient(user.ID)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		respond.With(w, r, http.StatusOK, jobs)
	})
}

func addJob() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		job := context.Get(r, "job").(*Job)
		user := context.Get(r, "user").(*models.User)
		job.ClientID = user.ID
		if job.Deadline.IsZero() {
			job.Deadline = time.Now()
		}
		if err := appContext.JobRepository.AddJob(job); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		// get full job to index
		job, err := appContext.JobRepository.GetJob(job.ID)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		if err := appContext.Indexer.Index("jobs", fmt.Sprint(job.ID), job); err != nil {
			log.Printf("job could not be indexed: %v", err)
		}
		respond.With(w, r, http.StatusOK, job)
	})
}

// todo: !!!
func getJob() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		user := context.Get(r, "user").(*models.User)
		var job *Job
		var err error
		switch user.Type {
		case "client":
			job, err = appContext.JobRepository.GetJobForClient(id, user.ID)
		case "freelancer":
			job, err = appContext.JobRepository.GetJobForFreelancer(id, user.ID)
		default:
			log.Printf("getJob: userType not recognized: %s", user.Type)
			respond.With(w, r, http.StatusInternalServerError, nil)
			return
		}
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, job)
	})
}

func withJobFromRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*models.User)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var newJob struct {
			Name                string     `json:"name" valid:"required"`
			Details             string     `json:"details" valid:"required"`
			Summary             string     `json:"summary" valid:"required"`
			PriceFrom           int        `json:"priceFrom" valid:"required"`
			PriceTo             int        `json:"priceTo" valid:"required"`
			Tags                stringList `json:"tags" valid:"required"`
			Deadline            string     `json:"deadline" valid:"required"`
			DeadlineFlexibility int        `json:"flexibility" valid:"required"`
			Attachments         []File     `json:"attachments"`
			Examples            []File     `json:"examples"`
		}
		if err := decoder.Decode(&newJob); err != nil {
			log.Printf("could not decode job: %v", err)
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		if len(newJob.Tags) > 10 {
			log.Printf("could not create job: too many tags")
			respond.With(w, r, http.StatusBadRequest, errors.New("max of 10 tags are allowed"))
			return
		}
		if ok, err := govalidator.ValidateStruct(newJob); ok == false || err != nil {
			log.Printf("could not create job: job not valid")
			respond.With(w, r, http.StatusBadRequest, models.GovalidatorErrors{Err: err})
			return
		}
		if newJob.PriceFrom > newJob.PriceTo {
			log.Printf("could not create job: priceFrom must be larger that priceTo")
			respond.With(w, r, http.StatusBadRequest, errors.New("priceFrom must be larger that priceTo"))
			return
		}
		deadlilne, err := time.Parse(time.RFC3339, newJob.Deadline)
		if err != nil {
			err := fmt.Errorf("could not parse deadline")
			log.Printf("could not create job: %v", err)
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		job := &Job{
			Name:                newJob.Name,
			Summary:             newJob.Summary,
			PriceFrom:           newJob.PriceFrom,
			PriceTo:             newJob.PriceTo,
			Tags:                newJob.Tags,
			Details:             newJob.Details,
			ClientID:            user.ID,
			Deadline:            deadlilne,
			DeadlineFlexibility: newJob.DeadlineFlexibility,
			Attachments:         newJob.Attachments,
			Examples:            newJob.Examples,
		}
		context.Set(r, "job", job)
		handler.ServeHTTP(w, r)
	})
}

func whenFreelancerHasNotAppliedBeforeByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var jobID = context.Get(r, "id").(uint)
		job, err := appContext.JobRepository.GetJob(jobID)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		var user = context.Get(r, "user").(*models.User)
		for _, application := range job.JobApplications {
			if application.FreelancerID == user.ID {
				respond.With(w, r, http.StatusBadRequest, fmt.Errorf("freelancer already applied for the job"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func withJobApplicationFromRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jobID = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*models.User)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var newJobApplication struct {
			Summary             string  `json:"summary" valid:"required"`
			Solution            string  `json:"solution" valid:"required"`
			Deadline            string  `json:"deadline" valid:"required"`
			DeadlineFlexibility int     `json:"flexibility"`
			Title               string  `json:"title" valid:"required"`
			Hours               int     `json:"hours" valid:"required"`
			HourPrice           float64 `json:"hourPrice" valid:"required"`
			Attachments         []File  `json:"attachments"`
			Examples            []File  `json:"examples"`
		}
		if err := decoder.Decode(&newJobApplication); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		if ok, err := govalidator.ValidateStruct(newJobApplication); ok == false || err != nil {
			respond.With(w, r, http.StatusBadRequest, models.GovalidatorErrors{Err: err})
			return
		}
		deadlilne, err := time.Parse(time.RFC3339, newJobApplication.Deadline)
		if err != nil {
			err := fmt.Errorf("could not parse deadline")
			log.Printf("could not create job: %v", err)
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		jobApplication := JobApplication{
			Summary:             newJobApplication.Summary,
			Solution:            newJobApplication.Solution,
			Title:               newJobApplication.Title,
			Deadline:            deadlilne,
			DeadlineFlexibility: newJobApplication.DeadlineFlexibility,
			FreelancerID:        user.ID,
			Hours:               newJobApplication.Hours,
			HourPrice:           newJobApplication.HourPrice,
			Attachments:         newJobApplication.Attachments,
			Examples:            newJobApplication.Examples,
			JobID:               jobID,
		}
		context.Set(r, "jobApplication", &jobApplication)
		handler.ServeHTTP(w, r)
	})
}

func addJobApplication() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var jobID = context.Get(r, "id").(uint)
		var jobApplication = context.Get(r, "jobApplication").(*JobApplication)
		var user = context.Get(r, "user").(*models.User)
		jobApplication.JobID = jobID
		jobApplication.FreelancerID = user.ID
		if err := appContext.JobRepository.AddJobApplication(jobApplication); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		client, ok := context.Get(r, "client").(*Client)
		if ok && client != nil {
			// get full job application with freelancer and e'rythang
			jobApplication, err := appContext.JobRepository.GetJobApplication(jobApplication.ID)
			if err != nil {
				respond.With(w, r, http.StatusInternalServerError, err)
				return
			}
			if err := appContext.NotificationDispatcher.notifyJobApplicationAdded(jobApplication, client.ID); err != nil {
				log.Printf("could not notifyJobApplicationAdded: %v", err)
			}
		}
		respond.With(w, r, http.StatusOK, jobApplication)
	})
}

func deleteJobApplicationByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		user := context.Get(r, "user").(*models.User)
		jobApplication, err := appContext.JobRepository.GetJobApplication(id)
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		if jobApplication.FreelancerID != user.ID {
			err := errors.New("freelancer not the owner")
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		if err := appContext.JobRepository.DeleteJobApplication(jobApplication.ID); err != nil {
			log.Println("delete job application:", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}

func whenJobApplicationBelongsToFreelancerByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*models.User)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		ok, err := appContext.JobRepository.jobApplicationBelongsToFreelancer(id, user.ID)
		if err != nil {
			log.Printf("could not check if job application %d belongs to freelancer %d: %v", id, user.ID, err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		if !ok {
			err := fmt.Errorf("job application %d does not belong to freelancer %d", id, user.ID)
			respond.With(w, r, http.StatusForbidden, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func whenJobApplicationBelongsToClientByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*models.User)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		ok, err := appContext.JobRepository.jobApplicationBelongsToClient(id, user.ID)
		if err != nil {
			log.Printf("could not check if job application %d belongs to client %d: %v", id, user.ID, err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		if !ok {
			err := fmt.Errorf("job application %d does not belong to client %d", id, user.ID)
			respond.With(w, r, http.StatusForbidden, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
