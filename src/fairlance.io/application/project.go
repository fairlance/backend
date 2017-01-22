package application

import (
	"net/http"

	"fmt"

	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func getAllProjects() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		projects, err := appContext.ProjectRepository.GetAllProjects()
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
		var user = context.Get(r, "user").(map[string]interface{})
		var userID = uint(user["id"].(float64))
		var userType = context.Get(r, "userType").(string)
		var projects []Project
		var err error
		switch userType {
		case "client":
			projects, err = appContext.ProjectRepository.GetAllProjectsForClient(userID)
		case "freelancer":
			projects, err = appContext.ProjectRepository.GetAllProjectsForFreelancer(userID)
		default:
			err = fmt.Errorf("found type '%s' unrecognized", userType)
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
		project, err := appContext.ProjectRepository.GetByID(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, project)
	})
}
