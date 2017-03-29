package application

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	respond "gopkg.in/matryer/respond.v1"

	"fmt"
)

const (
	projectStatusWorking         = "working"
	projectStatusFinilazingTerms = "finalizing_terms"
	projectStatusPending         = "pending"
	projectStatusArchived        = "archived"
	projectStatusCanceled        = "canceled"
)

func getAllProjects() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		projects, err := appContext.ProjectRepository.getAllProjects()
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, projects)
	})
}

func getAllProjectsForUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var user = context.Get(r, "user").(*User)
		var userType = context.Get(r, "userType").(string)
		var projects []Project
		var err error
		switch userType {
		case "freelancer":
			projects, err = appContext.ProjectRepository.getAllProjectsForFreelancer(user.ID)
		case "client":
			projects, err = appContext.ProjectRepository.getAllProjectsForClient(user.ID)
		default:
			err = fmt.Errorf("found type '%s' unrecognized", userType)
			if err != nil {
				respond.With(w, r, http.StatusBadRequest, err)
				return
			}
		}
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, projects)
	})
}

func getProjectByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		project, err := appContext.ProjectRepository.getByID(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, project)
	})
}

func createProjectFromJobApplication() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var jobApplicationID = context.Get(r, "job_application_id").(uint)
		jobApplication, err := appContext.JobRepository.GetJobApplication(jobApplicationID)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		job, err := appContext.JobRepository.GetJob(jobApplication.JobID)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		deadlineWithTime := time.Now().Add(time.Hour * 24 * time.Duration(jobApplication.DeliveryEstimate))
		deadline := time.Date(deadlineWithTime.Year(), deadlineWithTime.Month(), deadlineWithTime.Day()+1, 0, 0, 0, 0, deadlineWithTime.Location())

		project := Project{
			Name:            job.Name,
			Description:     job.Details,
			ClientID:        job.ClientID,
			Status:          projectStatusPending,
			Deadline:        deadline,
			WorkhoursPerDay: jobApplication.Hours,
			PerHour:         jobApplication.HourPrice,
			Freelancers: []Freelancer{
				*jobApplication.Freelancer,
			},
		}

		err = appContext.ProjectRepository.add(&project)
		if err != nil {
			log.Printf("create project: %v\n", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		err = appContext.JobRepository.DeleteJobApplication(jobApplication)
		if err != nil {
			log.Printf("delete job application: %v\n", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, project)
	})
}
