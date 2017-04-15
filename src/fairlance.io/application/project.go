package application

import (
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

func whenProjectBelongsToUserByID(next http.Handler) http.Handler {
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

func withProjectByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)

		project, err := appContext.ProjectRepository.getByID(id)
		if err != nil {
			log.Printf("could not find project: %v", err)
			respond.With(w, r, http.StatusNotFound, fmt.Errorf("could not find project"))
			return
		}

		context.Set(r, "project", project)

		next.ServeHTTP(w, r)
	})
}

func createProjectFromJobApplication() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)

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
			Deadline:           deadline,
			Hours:              jobApplication.Hours,
			PerHour:            jobApplication.HourPrice,
			ClientAgreed:       false,
			FreelancersToAgree: []uint{jobApplication.FreelancerID},
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

func agreeToContractTerms() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		user := context.Get(r, "user").(*User)
		userType := context.Get(r, "userType").(string)
		project := context.Get(r, "project").(*Project)

		var freelancersToAgree = project.Contract.FreelancersToAgree
		var clientAgreed = project.Contract.ClientAgreed
		if userType == "client" {
			clientAgreed = true
		} else if userType == "freelancer" {
			freelancersToAgree = removeFromUINTSlice(freelancersToAgree, user.ID)
		}

		err := appContext.ProjectRepository.updateContract(project.Contract, map[string]interface{}{
			"clientAgreed":       clientAgreed,
			"freelancersToAgree": freelancersToAgree,
		})
		if err != nil {
			log.Printf("could not update contract: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update contract"))
			return
		}

		respond.With(w, r, http.StatusOK, project)
	})
}
