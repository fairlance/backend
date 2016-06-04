package main

import (
	"net/http"

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
