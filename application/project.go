package application

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fairlance/backend/models"
	"github.com/gorilla/context"
	respond "gopkg.in/matryer/respond.v1"

	"fmt"
)

const (
	projectStatusFinalizingTerms = "finalizing_terms"
	projectStatusPendingFunds    = "pending_funds"
	projectStatusInProgress      = "in_progress"
	projectStatusPendingFinished = "pending_finished"
	projectStatusDone            = "done"
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

func getAllProjectsForFreelancer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var user = context.Get(r, "user").(*models.User)
		projects, err := appContext.ProjectRepository.getAllProjectsForFreelancer(user.ID)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, projects)
	})
}

func getAllProjectsForClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var user = context.Get(r, "user").(*models.User)
		projects, err := appContext.ProjectRepository.getAllProjectsForClient(user.ID)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, projects)
	})
}

func whenProjectBelongsToClientByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*models.User)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		project, err := appContext.ProjectRepository.getByID(id)
		if err != nil {
			log.Printf("whenProjectBelongsToClientByID, client(%d), project(%d): %v", user.ID, id, err)
			respond.With(w, r, http.StatusNotFound, err)
			return
		}
		if project.ClientID == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		respond.With(w, r, http.StatusForbidden, fmt.Errorf("client %d does not work on the project: %d", user.ID, id))
	})
}

func whenProjectBelongsToFreelancerByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*models.User)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		project, err := appContext.ProjectRepository.getByID(id)
		if err != nil {
			log.Printf("whenProjectBelongsToFreelancerByID: %v", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		for _, frelancer := range project.Freelancers {
			if frelancer.ID == user.ID {
				next.ServeHTTP(w, r)
				return
			}
		}
		respond.With(w, r, http.StatusForbidden, fmt.Errorf("freelancer %d does not work on the project: %d", user.ID, id))
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
		project := NewProject(job, jobApplication)
		if err = appContext.ProjectRepository.add(project); err != nil {
			log.Printf("create project: %v", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		if err = appContext.JobRepository.update(job); err != nil {
			log.Printf("deactivate job: %v", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		if err := appContext.Indexer.Delete("jobs", fmt.Sprint(job.ID)); err != nil {
			log.Printf("job could not be deleted from searcher: %v", err)
		}
		if err := appContext.NotificationDispatcher.notifyJobApplicationAccepted(jobApplication, project); err != nil {
			log.Printf("could not notifyJobApplicationAccepted: %v", err)
		}
		respond.With(w, r, http.StatusOK, project)
	})
}

// todo: !!!
func agreeToContractTerms() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		user := context.Get(r, "user").(*models.User)
		project := context.Get(r, "project").(*Project)
		if user.Type == "client" {
			project.ClientAgreed = true
		} else if user.Type == "freelancer" {
			project.FreelancersAgreed = append(project.FreelancersAgreed, user.ID)
		}
		if err := appContext.ProjectRepository.update(project); err != nil {
			log.Printf("could not update project contract: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project contract"))
			return
		}
		if err := appContext.MessagingDispatcher.sendContractAccepted(project, user); err != nil {
			log.Printf("could not sendContractAccepted: %v", err)
		}
		if project.canBeStarted() {
			project.mergeProposalToContract()
			project.Status = projectStatusPendingFunds
			if err := appContext.ProjectRepository.update(project); err != nil {
				log.Printf("could not update project: %v", err)
				respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project"))
				return
			}
			if err := appContext.MessagingDispatcher.sendProjectStateChanged(project); err != nil {
				log.Printf("could not sendProjectStateChanged: %v", err)
			}
		}

		respond.With(w, r, http.StatusOK, project)
	})
}

func projectFunded() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		project := context.Get(r, "project").(*Project)
		project.Status = projectStatusInProgress
		if err := appContext.ProjectRepository.update(project); err != nil {
			log.Printf("could not update project: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project"))
			return
		}
		if err := appContext.MessagingDispatcher.sendProjectStateChanged(project); err != nil {
			log.Printf("could not sendProjectStateChanged: %v", err)
		}
		respond.With(w, r, http.StatusOK, project)
	})
}

func withProposal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*models.User)
		proposal := &Proposal{}
		if err := json.NewDecoder(r.Body).Decode(proposal); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()
		proposal.Time = time.Now()
		proposal.UserID = user.ID
		proposal.UserType = user.Type
		context.Set(r, "proposal", proposal)
		next.ServeHTTP(w, r)
	})
}

func setProposalToProjectContract() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		project := context.Get(r, "project").(*Project)
		proposal := context.Get(r, "proposal").(*Proposal)
		user := context.Get(r, "user").(*models.User)
		if err := appContext.ProjectRepository.setContractProposal(project.Contract, proposal); err != nil {
			log.Printf("could not set proposal: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not set proposal"))
			return
		}
		if err := appContext.MessagingDispatcher.sendProjectContractProposalAdded(project.ID, proposal, user); err != nil {
			log.Printf("could not sendProjectContractProposalAdded: %v", err)
		}
		respond.With(w, r, http.StatusOK, project)
	})
}

func freelancerFinishProject() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		user := context.Get(r, "user").(*models.User)
		project := context.Get(r, "project").(*Project)
		if contains(project.FreelancersConcluded, user.ID) {
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("freelancer already concluded"))
			return
		}
		project.FreelancersConcluded = append(project.FreelancersConcluded, user.ID)
		if err := appContext.ProjectRepository.update(project); err != nil {
			log.Printf("could not update project: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project"))
			return
		}
		if err := appContext.MessagingDispatcher.sendProjectFinishedByFreelancer(project, user); err != nil {
			log.Printf("could not sendProjectFinishedByFreelancer: %v", err)
		}
		if project.allFreelancersConcluded() {
			project.Status = projectStatusPendingFinished
			if err := appContext.ProjectRepository.update(project); err != nil {
				log.Printf("could not update project status to pending_finished: %v", err)
				respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project status to pending_finished"))
				return
			}
			if err := appContext.MessagingDispatcher.sendProjectStateChanged(project); err != nil {
				log.Printf("could not sendProjectStateChanged to pending_finished: %v", err)
			}
		}
		respond.With(w, r, http.StatusOK, project)
	})
}

func projectDone() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		user := context.Get(r, "user").(*models.User)
		project := context.Get(r, "project").(*Project)
		project.ClientConcluded = true
		if err := appContext.ProjectRepository.update(project); err != nil {
			log.Printf("could not update project: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project"))
			return
		}
		if err := appContext.MessagingDispatcher.sendProjectDone(project, user); err != nil {
			log.Printf("could not sendProjectDone: %v", err)
		}
		if project.allUsersConcluded() {
			project.Status = projectStatusDone
			if err := appContext.ProjectRepository.update(project); err != nil {
				log.Printf("could not update project status to done: %v", err)
				respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update project status to done"))
				return
			}
			if err := appContext.MessagingDispatcher.sendProjectStateChanged(project); err != nil {
				log.Printf("could not sendProjectStateChanged to done: %v", err)
			}
			if err := appContext.PaymentDispatcher.execute(project.ID); err != nil {
				log.Printf("could not execute payment for project %d: %v", project.ID, err)
				respond.With(w, r, http.StatusFailedDependency, fmt.Errorf("payment could not be execued"))
				return
			}
		}
		respond.With(w, r, http.StatusOK, project)
	})
}

func fundedProject() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		user := context.Get(r, "user").(*models.User)
		project := context.Get(r, "project").(*Project)
		if err := appContext.Mailer.SendProjectFunded(
			project.ID,
			project.Name,
			user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		); err != nil {
			respond.With(w, r, http.StatusFailedDependency, err)
			return
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}
