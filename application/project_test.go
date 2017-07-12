package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/models"

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

	context.Set(r, "user", &models.User{
		ID:   1,
		Type: "freelancer",
	})
	getAllProjectsForFreelancer().ServeHTTP(w, r)

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

	context.Set(r, "user", &models.User{
		ID:   1,
		Type: "client",
	})
	getAllProjectsForClient().ServeHTTP(w, r)

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
	projectRepoMock.GetByIDCall.Returns.Project = &Project{
		Model: Model{
			ID: 123456789,
		},
		Name:        "Name1",
		Description: "Description1",
		ClientID:    1,
		Status:      projectStatusDone,
		Contract: &Contract{
			Deadline:            timeNow,
			PerHour:             2,
			Hours:               3,
			DeadlineFlexibility: 4,
		},
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
	is.Equal(body.Contract.Deadline.Format(time.RFC3339), timeNow.Format(time.RFC3339))
	is.Equal(body.Status, projectStatusDone)
	is.Equal(body.Contract.PerHour, 2)
	is.Equal(body.Contract.Hours, 3)
	is.Equal(body.Contract.DeadlineFlexibility, 4)
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
		Hours:     62,
		HourPrice: 8,
	}
	deadline := time.Now()
	jobRepoMock.GetJobCall.Returns.Job = &Job{
		Model:    Model{ID: 4},
		Name:     "jobName",
		Details:  "jobDetails",
		ClientID: uint(33),
		Deadline: deadline,
	}

	var indexName string
	var documentID string
	var notifiedFreelancerID uint
	var notificationType string

	var appContext = &ApplicationContext{
		ProjectRepository: projectRepoMock,
		JobRepository:     jobRepoMock,
		NotificationDispatcher: NewNotificationDispatcher(&testNotifier{
			callback: func(notification *dispatcher.Notification) error {
				notifiedFreelancerID = notification.To[0].ID
				notificationType = notification.Type
				return nil
			},
		}),
		Indexer: &testIndexer{
			indexCallback: func(index, docID string, doc interface{}) error { return nil },
			deleteCallback: func(index, docID string) error {
				indexName = index
				documentID = docID
				return nil
			},
		},
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(appContext, "")
	context.Set(r, "id", uint(2))

	createProjectFromJobApplication().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Name, "jobName")
	is.Equal(projectRepoMock.AddCall.Receives.Project.Description, "jobDetails")
	is.Equal(projectRepoMock.AddCall.Receives.Project.ClientID, uint(33))
	is.Equal(projectRepoMock.AddCall.Receives.Project.Freelancers[0].ID, uint(22))
	is.Equal(projectRepoMock.AddCall.Receives.Project.Status, projectStatusFinalizingTerms)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Contract.Deadline, deadline)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Contract.Hours, 62)
	is.Equal(projectRepoMock.AddCall.Receives.Project.Contract.PerHour, 8)
	is.Equal(jobRepoMock.DeleteJobCall.Receives.ID, 4)
	is.Equal(notifiedFreelancerID, uint(22))
	is.Equal(notificationType, "job_application_accepted")
	is.Equal(indexName, "jobs")
	is.Equal(documentID, "4")
}

var whenProjectBelongsToUserData = []struct {
	clientID      uint
	freelancerIDs []uint
	userType      string
	userID        uint
	isNextCalled  bool
	status        int
}{
	{
		clientID:      11,
		freelancerIDs: []uint{21, 22},
		userType:      "client",
		userID:        11,
		isNextCalled:  true,
		status:        http.StatusOK,
	},
	{
		clientID:      11,
		freelancerIDs: []uint{21, 22},
		userType:      "freelancer",
		userID:        22,
		isNextCalled:  true,
		status:        http.StatusOK,
	},
	{
		clientID:      11,
		freelancerIDs: []uint{21, 22},
		userType:      "client",
		userID:        12,
		isNextCalled:  false,
		status:        http.StatusForbidden,
	},
	{
		clientID:      11,
		freelancerIDs: []uint{21, 22},
		userType:      "freelancer",
		userID:        23,
		isNextCalled:  false,
		status:        http.StatusForbidden,
	},
	{
		clientID:      11,
		freelancerIDs: []uint{},
		userType:      "freelancer",
		userID:        23,
		isNextCalled:  false,
		status:        http.StatusForbidden,
	},
}

func TestWhenProjectBelongsToUser(t *testing.T) {
	projectRepoMock := &projectRepositoryMock{}

	for _, testCase := range whenProjectBelongsToUserData {
		var freelancers []Freelancer
		for _, fid := range testCase.freelancerIDs {
			freelancers = append(freelancers, Freelancer{User: User{Model: Model{ID: fid}}})
		}

		projectRepoMock.GetByIDCall.Returns.Project = &Project{
			Model:       Model{ID: 1},
			ClientID:    testCase.clientID,
			Freelancers: freelancers,
		}
		var appContext = &ApplicationContext{
			ProjectRepository: projectRepoMock,
		}

		is := isHelper.New(t)
		w := httptest.NewRecorder()
		r := getRequest(appContext, "")
		context.Set(r, "id", uint(2))
		context.Set(r, "user", &models.User{
			ID:   testCase.userID,
			Type: testCase.userType,
		})

		isNextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isNextCalled = true
		})

		whenBasedOnUserType(
			whenProjectBelongsToClientByID,
			whenProjectBelongsToFreelancerByID,
		)(next).ServeHTTP(w, r)

		is.Equal(projectRepoMock.GetByIDCall.Receives.ID, uint(2))
		is.Equal(isNextCalled, testCase.isNextCalled)
		is.Equal(w.Code, testCase.status)
	}
}
