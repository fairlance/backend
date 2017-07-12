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
		if job.StartDate.IsZero() {
			job.StartDate = time.Now()
		}
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

func withJob(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var job Job
		if err := decoder.Decode(&job); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if len(job.Tags) > 10 {
			respond.With(w, r, http.StatusBadRequest, errors.New("max of 10 tags are allowed"))
			return
		}

		if ok, err := govalidator.ValidateStruct(job); ok == false || err != nil {
			respond.With(w, r, http.StatusBadRequest, govalidator.ErrorsByField(err))
			return
		}

		context.Set(r, "job", &job)

		handler.ServeHTTP(w, r)
	})
}

func withJobApplication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var jobApplication JobApplication

		if err := decoder.Decode(&jobApplication); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if ok, err := govalidator.ValidateStruct(jobApplication); ok == false || err != nil {
			errs := govalidator.ErrorsByField(err)
			respond.With(w, r, http.StatusBadRequest, errs)
			return
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

func whenJobApplicationBelongsToFreelancer(next http.Handler) http.Handler {
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

func whenJobApplicationBelongsToClient(next http.Handler) http.Handler {
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
