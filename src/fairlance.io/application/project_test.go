package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

func TestIndexProject(t *testing.T) {
	is := is.New(t)
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

	IndexProject(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Project
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestIndexProjectWithError(t *testing.T) {
	is := is.New(t)
	projectRepositoryMock := &ProjectRepositoryMock{}
	projectRepositoryMock.GetAllProjectsCall.Returns.Error = errors.New("nein")
	userContext := &ApplicationContext{
		ProjectRepository: projectRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	IndexProject(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}
