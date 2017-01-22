package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	isHelper "github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestGetAllProjects(t *testing.T) {
	is := isHelper.New(t)
	projectRepositoryMock := &ProjectRepositoryMock{}
	projectRepositoryMock.GetAllProjectsCall.Returns.Projects = []Project{
		Project{
			Model: Model{
				ID: 1,
			},
		},
		Project{
			Model: Model{
				ID: 2,
			},
		},
	}
	userContext := &ApplicationContext{
		ProjectRepository: projectRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	getAllProjects().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Project
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestGetAllProjectsWithError(t *testing.T) {
	is := isHelper.New(t)
	projectRepositoryMock := &ProjectRepositoryMock{}
	projectRepositoryMock.GetAllProjectsCall.Returns.Error = errors.New("nein")
	userContext := &ApplicationContext{
		ProjectRepository: projectRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	getAllProjects().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

// TODO: getAllProjectsForUser
// func TestGetAllProjectsForUser(t *testing.T) {
// 	is := isHelper.New(t)
// 	projectRepositoryMock := &ProjectRepositoryMock{}
// 	projectRepositoryMock.GetAllProjectsCall.Returns.Projects = []Project{
// 		Project{
// 			Model: Model{
// 				ID: 1,
// 			},
// 		},
// 		Project{
// 			Model: Model{
// 				ID: 2,
// 			},
// 		},
// 	}
// 	userContext := &ApplicationContext{
// 		ProjectRepository: projectRepositoryMock,
// 	}

// 	r := getRequest(userContext, ``)
// 	w := httptest.NewRecorder()

// 	getAllProjects().ServeHTTP(w, r)

// 	is.Equal(w.Code, http.StatusOK)
// 	var body []Project
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
// 	is.Equal(body[0].Model.ID, 1)
// 	is.Equal(body[1].Model.ID, 2)
// }

// func TestGetAllProjectsForUserWithError(t *testing.T) {
// 	is := isHelper.New(t)
// 	projectRepositoryMock := &ProjectRepositoryMock{}
// 	projectRepositoryMock.GetAllProjectsCall.Returns.Error = errors.New("nein")
// 	userContext := &ApplicationContext{
// 		ProjectRepository: projectRepositoryMock,
// 	}

// 	r := getRequest(userContext, ``)
// 	w := httptest.NewRecorder()

// 	getAllProjects().ServeHTTP(w, r)

// 	is.Equal(w.Code, http.StatusInternalServerError)
// }

func TestProjectGetByID(t *testing.T) {
	projectRepositoryMock := ProjectRepositoryMock{}
	projectRepositoryMock.GetByIDCall.Returns.Project = Project{
		Model: Model{
			ID: 123456789,
		},
		Name:        "Name1",
		Description: "Description1",
		ClientID:    1,
		IsActive:    true,
	}
	var appContext = &ApplicationContext{
		ProjectRepository: &projectRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(appContext, "")
	context.Set(r, "id", uint(1))

	getProjectByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body Project
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body.Model.ID, uint(123456789))
	is.Equal(body.Name, "Name1")
	is.Equal(body.Description, "Description1")
	is.Equal(body.IsActive, true)
	is.Equal(projectRepositoryMock.GetByIDCall.Receives.ID, uint(1))
}

func TestProjectGetByIDError(t *testing.T) {
	projectRepositoryMock := ProjectRepositoryMock{}
	projectRepositoryMock.GetByIDCall.Returns.Error = errors.New("Blah")
	var appContext = &ApplicationContext{
		ProjectRepository: &projectRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(appContext, "")
	context.Set(r, "id", uint(1))

	getProjectByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
}
