package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fairlance/backend/models"
	"github.com/gorilla/context"
	respond "gopkg.in/matryer/respond.v1"
)

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

func withExtensionWhenBelongsToProject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		project := context.Get(r, "project").(*Project)
		extensionID := context.Get(r, "extension_id").(uint)
		extension, err := appContext.ProjectRepository.getExtension(extensionID)
		if err != nil {
			log.Printf("extension not found (id %d): %v", extensionID, err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("extension not found: %d", extensionID))
			return
		}
		if extension.ContractID != project.ContractID {
			log.Printf("extension does not belong to the project: extension %d, contract %d", extensionID, project.ContractID)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("extension does not belong to the project"))
			return
		}
		context.Set(r, "extension", extension)
		next.ServeHTTP(w, r)
	})
}

func agreeToExtensionTerms() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		user := context.Get(r, "user").(*models.User)
		userType := context.Get(r, "userType").(string)
		extension := context.Get(r, "extension").(*Extension)

		if userType == "client" {
			extension.ClientAgreed = true
		} else if userType == "freelancer" {
			extension.FreelancersAgreed = append(extension.FreelancersAgreed, user.ID)
		}
		err := appContext.ProjectRepository.updateExtension(extension)
		if err != nil {
			log.Printf("could not update extension: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update extension"))
			return
		}
		respond.With(w, r, http.StatusOK, extension)
	})
}

func setProposalToProjectContractExtension() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appContext := context.Get(r, "context").(*ApplicationContext)
		project := context.Get(r, "project").(*Project)
		proposal := context.Get(r, "proposal").(*Proposal)
		extension := context.Get(r, "extension").(*Extension)
		err := appContext.ProjectRepository.setContractExtensionProposal(extension, proposal)
		if err != nil {
			log.Printf("could not set proposal: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not set proposal"))
			return
		}
		if err := appContext.MessagingDispatcher.sendProjectContractExtensionProposalAdded(project.ID, extension, proposal); err != nil {
			log.Printf("could not sendProjectContractProposalAdded: %v", err)
		}
		respond.With(w, r, http.StatusOK, project)
	})
}
