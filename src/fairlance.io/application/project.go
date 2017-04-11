package application

import (
	"encoding/json"
	"errors"
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

func whenProjectBelongsToUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*User)
		var userType = context.Get(r, "userType").(string)
		var appContext = context.Get(r, "context").(*ApplicationContext)

		project, err := appContext.ProjectRepository.getByID(id)
		if err != nil {
			log.Printf("whenProjectBelongsToUser: %v", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		if userType == "client" && project.ClientID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		if userType == "freelancer" {
			for _, frelancer := range project.Freelancers {
				if frelancer.ID == user.ID {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		respond.With(w, r, http.StatusForbidden, errors.New("user does not work on the project"))
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

func withExtension(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		extension := &Extension{}
		if err := decoder.Decode(extension); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		context.Set(r, "extension", extension)

		next.ServeHTTP(w, r)
	})
}

func addExtensionToProjectContract() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		var extension, ok = context.Get(r, "extension").(*Extension)
		if ok != true {
			log.Println("add extention to project contract: extension not provided")
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("extension could not be created"))
			return
		}

		project, err := appContext.ProjectRepository.getByID(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		extension.ContractID = project.ContractID
		err = appContext.ProjectRepository.addExtension(extension)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, extension)
	})
}

func createProjectFromJobApplication() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id, idOK = context.Get(r, "id").(uint)
		if !idOK {
			log.Printf("createProjectFromJobApplication: job application id not provided")
			respond.With(w, r, http.StatusInternalServerError, nil)
			return
		}
		jobApplication, err := appContext.JobRepository.GetJobApplication(id)
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

		contract := &Contract{
			Deadline: deadline,
			Hours:    jobApplication.Hours,
			PerHour:  jobApplication.HourPrice,
		}

		err = appContext.ProjectRepository.addContract(contract)
		if err != nil {
			log.Printf("create contract: %v\n", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		project := &Project{
			Name:        job.Name,
			Description: job.Details,
			ClientID:    job.ClientID,
			Status:      projectStatusPending,
			Deadline:    deadline,
			Hours:       jobApplication.Hours,
			PerHour:     jobApplication.HourPrice,
			Freelancers: []Freelancer{
				*jobApplication.Freelancer,
			},
			ContractID: contract.ID,
			Contract:   contract,
		}

		err = appContext.ProjectRepository.add(project)
		if err != nil {
			log.Printf("create project: %v\n", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		job.IsActive = false
		err = appContext.JobRepository.DeactivateJob(job)
		if err != nil {
			log.Printf("deactivate job: %v\n", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		notifyJobApplicationAccepted(appContext.Notifier, jobApplication, project)

		respond.With(w, r, http.StatusOK, project)
	})
}
