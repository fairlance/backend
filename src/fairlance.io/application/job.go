package application

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func IndexJob(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	jobs, err := appContext.JobRepository.GetAllJobs()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, jobs)
}

func AddJob(job *Job) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.JobRepository.AddJob(job); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, job)
	})
}

// GetJobByID handler
func GetJobByID(id uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		job, err := appContext.JobRepository.GetJob(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, job)
	})
}

type WithJob struct {
	next func(job *Job) http.Handler
}

func (withJob WithJob) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var body struct {
		Name     string  `json:"name" valid:"required"`
		Summary  string  `json:"summary" valid:"required"`
		Details  string  `json:"details" valid:"required"`
		ClientID uint    `json:"clientId" valid:"required"`
		IsActive bool    `json:"isActive"`
		Tags     stringList `json:"tags"`
		Links    stringList `json:"links"`
	}

	if err := decoder.Decode(&body); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	// https://github.com/asaskevich/govalidator/issues/133
	// https://github.com/asaskevich/govalidator/issues/112
	if len(body.Tags) > 10 {
		respond.With(w, r, http.StatusBadRequest, errors.New("Max of 10 tags are allowed."))
		return
	}

	if ok, err := govalidator.ValidateStruct(body); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		respond.With(w, r, http.StatusBadRequest, errs)
		return
	}

	job := &Job{
		Name:     body.Name,
		Summary:  body.Summary,
		Details:  body.Details,
		ClientID: body.ClientID,
		IsActive: body.IsActive,
		Tags:     body.Tags,
		Links:    body.Links,
	}

	withJob.next(job).ServeHTTP(w, r)
}

type ApplyForJobHandler struct {
	jobID          uint
	jobApplication *JobApplication
}

func (afjh ApplyForJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	afjh.jobApplication.JobID = afjh.jobID
	if err := appContext.JobRepository.AddJobApplication(afjh.jobApplication); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, afjh.jobApplication)
}

type WithJobApplication struct {
	next func(jobApplication *JobApplication) http.Handler
}

func (withJobApplication WithJobApplication) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var jobApplication JobApplication
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	// use JobApplication
	// JobApplication should not contain JobID
	var body struct {
		Message          string  `json:"message" valid:"required"`
		Samples          uintList   `json:"samples" valid:"required"`
		DeliveryEstimate int     `json:"deliveryEstimate" valid:"required"`
		Milestones       stringList `json:"milestones" valid:"required"`
		HourPrice        float64 `json:"hourPrice" valid:"required"`
		Hours            int     `json:"hours" valid:"required"`
		FreelancerID     uint    `json:"freelancerId" valid:"required"`
	}

	if err := decoder.Decode(&body); err != nil {
		respond.With(w, r, http.StatusBadRequest, errors.New("Invalid JSON"))
		return
	}

	if ok, err := govalidator.ValidateStruct(body); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		respond.With(w, r, http.StatusBadRequest, errs)
		return
	}

	jobApplication = JobApplication{
		Message:          body.Message,
		Milestones:       body.Milestones,
		Samples:          body.Samples,
		DeliveryEstimate: body.DeliveryEstimate,
		Hours:            body.Hours,
		HourPrice:        body.HourPrice,
		FreelancerID:     body.FreelancerID,
	}

	withJobApplication.next(&jobApplication).ServeHTTP(w, r)
}
