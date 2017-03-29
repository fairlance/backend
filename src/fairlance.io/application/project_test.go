package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	isHelper "github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestGetAllProjects(t *testing.T) {
	is := isHelper.New(t)
	projectRepoMock := &projectRepositoryMock{}
	projectRepoMock.GetAllProjectsCall.Returns.Projects = []Project{
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
		ProjectRepository: projectRepoMock,
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
	projectRepoMock := &projectRepositoryMock{}
	projectRepoMock.GetAllProjectsCall.Returns.Error = errors.New("nein")
	userContext := &ApplicationContext{
		ProjectRepository: projectRepoMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	getAllProjects().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestGetAllProjectsForFreelancer(t *testing.T) {
	is := isHelper.New(t)
	projectRepoMock := &projectRepositoryMock{}
	projectRepoMock.GetAllProjectsForFreelancerCall.Returns.Projects = []Project{
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
		ProjectRepository: projectRepoMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	context.Set(r, "userType", "freelancer")
	context.Set(r, "user", &User{
		Model: Model{
			ID: 1,
		},
	})
	getAllProjectsForUser().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Project
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(projectRepoMock.GetAllProjectsForFreelancerCall.Receives.ID, 1)
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestGetAllProjectsForClient(t *testing.T) {
	is := isHelper.New(t)
	projectRepoMock := &projectRepositoryMock{}
	projectRepoMock.GetAllProjectsForClientCall.Returns.Projects = []Project{
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
		ProjectRepository: projectRepoMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	context.Set(r, "userType", "client")
	context.Set(r, "user", &User{
		Model: Model{
			ID: 1,
		},
	})
	getAllProjectsForUser().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Project
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(projectRepoMock.GetAllProjectsForClientCall.Receives.ID, 1)
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestGetAllProjectsForUserWithError(t *testing.T) {
	is := isHelper.New(t)
	projectRepoMock := &projectRepositoryMock{}
	projectRepoMock.GetAllProjectsCall.Returns.Error = errors.New("nein")
	userContext := &ApplicationContext{
		ProjectRepository: projectRepoMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	getAllProjects().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestProjectGetByID(t *testing.T) {
	projectRepoMock := &projectRepositoryMock{}
	timeNow := time.Now()
	projectRepoMock.GetByIDCall.Returns.Project = Project{
		Model: Model{
			ID: 123456789,
		},
		Name:                "Name1",
		Description:         "Description1",
		ClientID:            1,
		Status:              projectStatusArchived,
		Deadline:            timeNow,
		PerHour:             2,
		WorkhoursPerDay:     3,
		DeadlineFlexibility: 4,
	}
	var appContext = &ApplicationContext{
		ProjectRepository: projectRepoMock,
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
	is.Equal(body.Deadline.Format(time.RFC3339), timeNow.Format(time.RFC3339))
	is.Equal(body.Status, projectStatusArchived)
	is.Equal(body.PerHour, 2)
	is.Equal(body.WorkhoursPerDay, 3)
	is.Equal(body.DeadlineFlexibility, 4)
	is.Equal(projectRepoMock.GetByIDCall.Receives.ID, uint(1))
}

func TestProjectGetByIDError(t *testing.T) {
	projectRepoMock := &projectRepositoryMock{}
	projectRepoMock.GetByIDCall.Returns.Error = errors.New("Blah")
	var appContext = &ApplicationContext{
		ProjectRepository: projectRepoMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(appContext, "")
	context.Set(r, "id", uint(1))

	getProjectByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
}

func TestCreateProjectFromJobApplication(t *testing.T) {
	projectRepoMock := &projectRepositoryMock{}
	jobRepoMock := &JobRepositoryMock{}
	jobRepoMock.GetJobApplicationCall.Returns.JobApplication = &JobApplication{
		Freelancer: &Freelancer{
			User: User{
				Model: Model{
					ID: 22,
				},
			},
		},
		Hours:            62,
		HourPrice:        8,
		DeliveryEstimate: 2,
	}
	deadline := time.Now().Add(time.Hour * 24 * 2)
	expectedDeadline := time.Date(deadline.Year(), deadline.Month(), deadline.Day()+1, 0, 0, 0, 0, deadline.Location())
	jobRepoMock.GetJobCall.Returns.Job = Job{
		Name:     "jobName",
		Details:  "jobDetails",
		ClientID: uint(33),
	}
	var appContext = &ApplicationContext{
		ProjectRepository: projectRepoMock,
		JobRepository:     jobRepoMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(appContext, "")
	context.Set(r, "job_application_id", uint(2))

	createProjectFromJobApplication().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Name, "jobName")
	is.Equal(projectRepoMock.AddCall.Receives.Project.Description, "jobDetails")
	is.Equal(projectRepoMock.AddCall.Receives.Project.ClientID, uint(33))
	is.Equal(projectRepoMock.AddCall.Receives.Project.Deadline, expectedDeadline)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Freelancers[0].ID, uint(22))
	is.Equal(projectRepoMock.AddCall.Receives.Project.WorkhoursPerDay, 62)
	is.Equal(projectRepoMock.AddCall.Receives.Project.PerHour, 8)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Status, projectStatusPending)
}
