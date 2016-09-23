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

func AddJob(w http.ResponseWriter, r *http.Request) {
	job := context.Get(r, "job").(*Job)
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.JobRepository.AddJob(job); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, job)
}

// GetJobByID handler
func GetJobByID(id uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		client, err := appContext.JobRepository.GetJob(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, client)
	})
}

func NewJobHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var body struct {
			Name        string `json:"name" valid:"required"`
			Description string `json:"description" valid:"required"`
			ClientId    uint   `json:"clientId" valid:"required"`
			IsActive    bool   `json:"isActive"`
			Tags        []Tag  `json:"tags"`
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
			Name:        body.Name,
			Description: body.Description,
			ClientId:    body.ClientId,
			IsActive:    body.IsActive,
			Tags:        body.Tags,
		}

		context.Set(r, "job", job)
		next.ServeHTTP(w, r)
	})
}
