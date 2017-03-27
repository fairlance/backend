package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
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

func getAllJobsForUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var user = context.Get(r, "user").(*User)
		var userType = context.Get(r, "userType").(string)
		var jobs []Job

		if userType != "client" {
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("user of type '%s' is not a client", userType))
			return
		}

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
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var job = context.Get(r, "job").(*Job)
		if err := appContext.JobRepository.AddJob(job); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, job)
	})
}

func getJob() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		job, err := appContext.JobRepository.GetJob(id)
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

		var body struct {
			Name        string       `json:"name" valid:"required"`
			Summary     string       `json:"summary" valid:"required"`
			Details     string       `json:"details" valid:"required"`
			ClientID    uint         `json:"clientId" valid:"required"`
			IsActive    bool         `json:"isActive"`
			Tags        stringList   `json:"tags"`
			Attachments []Attachment `json:"attachments"`
			Examples    []Example    `json:"examples"`
		}

		if err := decoder.Decode(&body); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if len(body.Tags) > 10 {
			respond.With(w, r, http.StatusBadRequest, errors.New("max of 10 tags are allowed"))
			return
		}

		if ok, err := govalidator.ValidateStruct(body); ok == false || err != nil {
			respond.With(w, r, http.StatusBadRequest, govalidator.ErrorsByField(err))
			return
		}

		job := &Job{
			Name:        body.Name,
			Summary:     body.Summary,
			Details:     body.Details,
			ClientID:    body.ClientID,
			IsActive:    body.IsActive,
			Tags:        body.Tags,
			Attachments: body.Attachments,
			Examples:    body.Examples,
		}

		context.Set(r, "job", job)

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

func addJobApplicationByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var jobID = context.Get(r, "id").(uint)
		var jobApplication = context.Get(r, "jobApplication").(*JobApplication)

		jobApplication.JobID = jobID
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
			notifyJobApplicationAdded(appContext.Notifier, jobApplication, client.ID)
		}

		respond.With(w, r, http.StatusOK, jobApplication)
	})
}

func deleteJobApplicationByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		user := context.Get(r, "user").(*User)
		userType := context.Get(r, "userType")

		if userType != "freelancer" {
			err := errors.New("user not a freelancer")
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		jobApplication, err := appContext.JobRepository.GetJobApplication(id)
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if jobApplication.FreelancerID != user.ID {
			err := errors.New("user not the owner")
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if err := appContext.JobRepository.DeleteJobApplication(jobApplication); err != nil {
			log.Println("delete job application:", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}
