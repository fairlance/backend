package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/models"

	respond "gopkg.in/matryer/respond.v1"

	"strings"

	isHelper "github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestJobIndexJob(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GettAllJobsCall.Returns.Jobs =
		[]Job{
			Job{
				Model: Model{
					ID: 1,
				},
				Name:      "Name1",
				Summary:   "Summary1",
				Details:   "Details1",
				ClientID:  1,
				PriceFrom: 100,
				PriceTo:   150,
				Examples: []File{
					{
						Name: "example",
						URL:  "www.example.com",
					},
				},
				Attachments: []File{
					{
						Name: "attachment",
						URL:  "www.attachment.com",
					},
				},
			},
		}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")

	getAllJobs().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Job
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(len(body), 1)
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[0].Name, "Name1")
	is.Equal(body[0].Summary, "Summary1")
	is.Equal(body[0].Details, "Details1")
	is.Equal(body[0].PriceFrom, 100)
	is.Equal(body[0].PriceTo, 150)
	is.Equal(len(body[0].Examples), 1)
	is.Equal(body[0].Examples[0].Name, "example")
	is.Equal(body[0].Examples[0].URL, "www.example.com")
	is.Equal(len(body[0].Attachments), 1)
	is.Equal(body[0].Attachments[0].Name, "attachment")
	is.Equal(body[0].Attachments[0].URL, "www.attachment.com")
}

func TestGetJobsForClient(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")

	context.Set(r, "user", &models.User{
		ID:   1,
		Type: "client",
	})

	getAllJobsForClient().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.GetAllJobsForClientCall.Receives.ID, 1)
}

func TestJobIndexJobError(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GettAllJobsCall.Returns.Error = errors.New("error")
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")

	getAllJobs().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobAddJob(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.AddJobCall.Updates.ID = 4
	var indexName string
	var documentID string
	var document *Job
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
		Indexer: &testIndexer{
			indexCallback: func(index, docID string, doc interface{}) error {
				indexName = index
				documentID = docID
				document = doc.(*Job)
				return nil
			},
			deleteCallback: func(index, docID string) error { return nil },
		},
	}
	job := &Job{
		Name:      "Name1",
		Summary:   "Summary1",
		Details:   "Details1",
		ClientID:  1,
		PriceFrom: 100,
		PriceTo:   150,
		Examples: []File{
			{
				Name: "example",
				URL:  "www.example.com",
			},
		},
		Attachments: []File{
			{
				Name: "attachment",
				URL:  "www.attachment.com",
			},
		},
	}
	jobRepositoryMock.GetJobCall.Returns.Job = job

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "job", job)
	context.Set(r, "user", &models.User{
		ID: 1,
	})

	addJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Name, "Name1")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Summary, "Summary1")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Details, "Details1")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.ClientID, 1)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.PriceFrom, 100)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.PriceTo, 150)
	is.Equal(len(jobRepositoryMock.AddJobCall.Receives.Job.Examples), 1)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Examples[0].Name, "example")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Examples[0].URL, "www.example.com")
	is.Equal(len(jobRepositoryMock.AddJobCall.Receives.Job.Attachments), 1)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Attachments[0].Name, "attachment")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Attachments[0].URL, "www.attachment.com")
	is.Equal(indexName, "jobs")
	is.Equal(documentID, "4")
	is.Equal(document.Name, "Name1")
}

func TestJobGetJobForClientReceivesTheRightParameters(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(1))
	context.Set(r, "user", &models.User{
		ID:   2,
		Type: "client",
	})

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.GetJobForClientCall.Receives.ID, uint(1))
	is.Equal(jobRepositoryMock.GetJobForClientCall.Receives.ClientID, uint(2))
}

func TestJobGetJobForFreelancerReceivesTheRightParameters(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(1))
	context.Set(r, "user", &models.User{
		ID:   2,
		Type: "freelancer",
	})

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.GetJobForFreelancerCall.Receives.ID, uint(1))
	is.Equal(jobRepositoryMock.GetJobForFreelancerCall.Receives.FreelancerID, uint(2))
}

func TestJobGetJobForClient(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GetJobForClientCall.Returns.Job = &Job{
		Model: Model{
			ID: 123456789,
		},
		Name:      "Name1",
		Summary:   "Summary1",
		Details:   "Details1",
		ClientID:  1,
		PriceFrom: 100,
		PriceTo:   150,
	}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(123456789))
	context.Set(r, "user", &models.User{
		ID:   1,
		Type: "client",
	})

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body Job
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body.Model.ID, uint(123456789))
	is.Equal(body.Name, "Name1")
	is.Equal(body.Summary, "Summary1")
	is.Equal(body.Details, "Details1")
	is.Equal(body.PriceFrom, 100)
	is.Equal(body.PriceTo, 150)
}

func TestJobGetJobForFreelancer(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GetJobForFreelancerCall.Returns.Job = &Job{
		Model: Model{
			ID: 123456789,
		},
		Name:      "Name1",
		Summary:   "Summary1",
		Details:   "Details1",
		ClientID:  1,
		PriceFrom: 100,
		PriceTo:   150,
	}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(123456789))
	context.Set(r, "user", &models.User{
		ID:   1,
		Type: "freelancer",
	})

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body Job
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body.Model.ID, uint(123456789))
	is.Equal(body.Name, "Name1")
	is.Equal(body.Summary, "Summary1")
	is.Equal(body.Details, "Details1")
	is.Equal(body.PriceFrom, 100)
	is.Equal(body.PriceTo, 150)
}

func TestJobGetJobByIDError(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GetJobForClientCall.Returns.Error = errors.New("Blah")
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(0))
	context.Set(r, "user", &models.User{
		ID:   1,
		Type: "client",
	})

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
}

func TestJobAddJobError(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.AddJobCall.Returns.Error = errors.New("Oooops")
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "job", &Job{})
	context.Set(r, "user", &models.User{})

	addJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobHandleApplyForJob(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobApplication := &JobApplication{
		JobID:       1,
		Title:       "title",
		HourPrice:   1.1,
		Hours:       1,
		Examples:    []File{{Name: "example", URL: "www.example.com"}},
		Attachments: []File{{Name: "attachment", URL: "www.attachment.com"}},
	}
	jobRepositoryMock.GetJobApplicationCall.Returns.JobApplication = jobApplication
	var notifiedClientID uint
	var notificationType string
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
		NotificationDispatcher: NewNotificationDispatcher(&testNotifier{
			callback: func(notification *dispatcher.Notification) error {
				notifiedClientID = notification.To[0].ID
				notificationType = notification.Type
				return nil
			},
		}),
	}

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(1))
	jobApplication.JobID = 0 // reset jobID because we needed it before
	context.Set(r, "jobApplication", jobApplication)
	context.Set(r, "user", &models.User{ID: 1})
	context.Set(r, "client", &Client{User: User{Model: Model{ID: 22}}})

	addJobApplication().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.JobID, 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.FreelancerID, 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Title, "title")
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.HourPrice, 1.1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Hours, 1)
	is.Equal(len(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Examples), 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Examples[0].Name, "example")
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Examples[0].URL, "www.example.com")
	is.Equal(len(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Attachments), 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Attachments[0].Name, "attachment")
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Attachments[0].URL, "www.attachment.com")
	is.Equal(notifiedClientID, uint(22))
	is.Equal(notificationType, "job_application_added")
}

func TestJobHandleApplyForJobHandlerError(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.AddJobApplicationCall.Returns.Error = errors.New("fuckup")
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(1))
	context.Set(r, "jobApplication", &JobApplication{})
	context.Set(r, "user", &models.User{ID: 1})

	addJobApplication().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobWithJob(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"name": "Name1",
		"details": "Details1",
		"summary": "Summary1",
		"priceFrom": 12,
		"priceTo": 15,
		"deadline": "2017-07-27T00:00:00.000Z",
		"flexibility": 2,
		"examples": [
			{"name": "example", "url": "www.example.com"}
		],
		"attachments": [
			{"name": "attachment", "url": "www.attachment.com"}
		],
		"tags": ["one", "two"]
	}`
	r := getRequest(jobContext, requestBody)

	nextCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	context.Set(r, "user", &models.User{ID: 1})
	withJobFromRequest(handler).ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(nextCalled, true)
	job := context.Get(r, "job").(*Job)
	is.Equal(job.ClientID, 1)
	is.Equal(job.Name, "Name1")
	is.Equal(job.Details, "Details1")
	is.Equal(job.Summary, "Summary1")
	is.Equal(job.PriceFrom, 12)
	is.Equal(job.PriceTo, 15)
	deadline, err := time.Parse(time.RFC3339, "2017-07-27T00:00:00.000Z")
	if err != nil {
		t.Fatal(err)
	}
	is.Equal(job.Deadline, deadline)
	is.Equal(job.DeadlineFlexibility, 2)
	is.Equal(len(job.Examples), 1)
	is.Equal(job.Examples[0].Name, "example")
	is.Equal(job.Examples[0].URL, "www.example.com")
	is.Equal(len(job.Attachments), 1)
	is.Equal(job.Attachments[0].Name, "attachment")
	is.Equal(job.Attachments[0].URL, "www.attachment.com")
	is.Equal(job.Tags[0], "one")
	is.Equal(job.Tags[1], "two")
}

// func TestJobWithJobError(t *testing.T) {
// 	var jobContext = &ApplicationContext{}
// 	is := isHelper.New(t)
// 	w := httptest.NewRecorder()
// 	requestBody := `{}`
// 	r := getRequest(jobContext, requestBody)

// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		t.Error("Should not be called")
// 	})
// 	withJob(handler).ServeHTTP(w, r)

// 	is.Equal(w.Code, http.StatusBadRequest)
// 	var body map[string]string
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
// 	is.Equal(len(body), 4)
// 	is.Equal(body["details"], "non zero value required")
// 	is.Equal(body["name"], "non zero value required")
// 	is.Equal(body["summary"], "non zero value required")
// }

func TestJobWithJobErrorBadTooManyTags(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"name": "Name1",
		"details": "Details1",
		"summary": "Summary1",
		"clientId": 1,
		"tags": ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
	}`
	r := getRequest(jobContext, requestBody)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})
	context.Set(r, "user", &models.User{ID: 1})
	withJobFromRequest(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobWithJobErrorBadJSON(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{bad:json}`
	r := getRequest(jobContext, requestBody)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})
	context.Set(r, "user", &models.User{ID: 1})
	withJobFromRequest(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobWithJobApplication(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"summary": "summary",
		"solution": "solution",
		"deadline": "2017-07-27T00:00:00.000Z",
		"title":"title",
		"hours": 12,
		"hourPrice": 1.1,
		"attachments": [
			{"name": "attachment", "url": "www.attachment.com"}
		],
		"examples": [
			{"name": "example", "url": "www.example.com"}
		]
	}`
	r := getRequest(jobContext, requestBody)

	nextCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	var jobID uint = 1
	var userID uint = 11
	context.Set(r, "id", jobID)
	context.Set(r, "user", &models.User{ID: userID})
	withJobApplicationFromRequest(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(nextCalled, true)
	jobApplication := context.Get(r, "jobApplication").(*JobApplication)
	is.Equal(jobApplication.Summary, "summary")
	is.Equal(jobApplication.Solution, "solution")
	is.Equal(jobApplication.Title, "title")
	deadline, err := time.Parse(time.RFC3339, "2017-07-27T00:00:00.000Z")
	if err != nil {
		t.Fatal(err)
	}
	is.Equal(jobApplication.Deadline, deadline)
	is.Equal(jobApplication.FreelancerID, userID)
	is.Equal(jobApplication.JobID, jobID)
	is.Equal(jobApplication.Hours, 12)
	is.Equal(jobApplication.HourPrice, 1.1)
	is.Equal(len(jobApplication.Examples), 1)
	is.Equal(jobApplication.Examples[0].Name, "example")
	is.Equal(jobApplication.Examples[0].URL, "www.example.com")
	is.Equal(len(jobApplication.Attachments), 1)
	is.Equal(jobApplication.Attachments[0].Name, "attachment")
	is.Equal(jobApplication.Attachments[0].URL, "www.attachment.com")
}

//func TestJobWithJobApplicationError(t *testing.T) {
//	var jobContext = &ApplicationContext{}
//	is := isHelper.New(t)
//	w := httptest.NewRecorder()
//	requestBody := `{}`
//	r := getRequest(jobContext, requestBody)
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		t.Error("Should not be called")
//	})
//	withJobApplication(handler).ServeHTTP(w, r)
//
//	is.Equal(w.Code, http.StatusBadRequest)
//	var body map[string]string
//	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
//	is.Equal(len(body), 7)
//	is.Equal(body["Message"], "non zero value required")
//	is.Equal(body["Milestones"], "non zero value required")
//	is.Equal(body["FreelancerID"], "non zero value required")
//	is.Equal(body["Hours"], "non zero value required")
//	is.Equal(body["HourPrice"], "non zero value required")
//}

func TestJobWithJobApplicationErrorBadJSON(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{bad:json}`
	r := getRequest(jobContext, requestBody)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})
	context.Set(r, "id", uint(1))
	context.Set(r, "user", &models.User{ID: uint(1)})
	withJobApplicationFromRequest(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestDeleteJobApplicationByID(t *testing.T) {
	jobRepositoryMock := &JobRepositoryMock{}
	jobRepositoryMock.GetJobApplicationCall.Returns.JobApplication = &JobApplication{
		Model: Model{
			ID: 1,
		},
		FreelancerID: 2,
	}
	var jobContext = &ApplicationContext{
		JobRepository: jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")

	context.Set(r, "id", uint(1))
	context.Set(r, "user", &models.User{
		ID:   2,
		Type: "freelancer",
	})
	deleteJobApplicationByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.GetJobApplicationCall.Receives.ID, 1)
	is.Equal(jobRepositoryMock.DeleteJobApplicationCall.Receives.ID, 1)
}

var deleteJobApplicationByIDData = []struct {
	inID             uint
	inJobApplication *JobApplication
	inGetError       error
	inDeleteError    error
	inUser           *models.User
	outStatus        int
	out              string
}{
	{
		uint(1),
		&JobApplication{
			Model: Model{
				ID: 1,
			},
			FreelancerID: 2,
		},
		nil,
		nil,
		&models.User{
			ID:   3,
			Type: "freelancer",
		},
		http.StatusBadRequest,
		"freelancer not the owner",
	},
	{
		uint(5),
		&JobApplication{
			Model: Model{
				ID: 1,
			},
			FreelancerID: 2,
		},
		errors.New("bad id"),
		nil,
		&models.User{
			ID:   2,
			Type: "freelancer",
		},
		http.StatusBadRequest,
		"bad id",
	},
	{
		uint(1),
		&JobApplication{
			Model: Model{
				ID: 1,
			},
			FreelancerID: 2,
		},
		nil,
		errors.New("cannot delete"),
		&models.User{
			ID:   2,
			Type: "freelancer",
		},
		http.StatusInternalServerError,
		"cannot delete",
	},
}

func TestDeleteJobApplicationByIDWithError(t *testing.T) {
	jobRepositoryMock := &JobRepositoryMock{}
	for _, testCase := range deleteJobApplicationByIDData {
		jobRepositoryMock.GetJobApplicationCall.Returns.JobApplication = testCase.inJobApplication
		jobRepositoryMock.GetJobApplicationCall.Returns.Error = testCase.inGetError
		jobRepositoryMock.DeleteJobApplicationCall.Returns.Error = testCase.inDeleteError
		var jobContext = &ApplicationContext{
			JobRepository: jobRepositoryMock,
		}
		is := isHelper.New(t)
		w := httptest.NewRecorder()
		r := getRequest(jobContext, "")

		context.Set(r, "id", testCase.inID)
		context.Set(r, "user", testCase.inUser)

		opts := &respond.Options{
			Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
				return status, data.(error).Error()
			},
		}
		opts.Handler(deleteJobApplicationByID()).ServeHTTP(w, r)

		is.Equal(w.Code, testCase.outStatus)
		is.Equal(strings.TrimSpace(w.Body.String()), fmt.Sprintf(`"%s"`, testCase.out))
	}
}

var whenJobApplicationBelongsToClientData = []struct {
	inOK            bool
	inError         error
	inUserType      string
	outIsNextCalled bool
	outStatus       int
	out             string
}{
	{
		inOK:            true,
		inUserType:      "freelancer",
		inError:         nil,
		outIsNextCalled: true,
		outStatus:       http.StatusOK,
		out:             "",
	},
	{
		inOK:            true,
		inUserType:      "client",
		inError:         nil,
		outIsNextCalled: true,
		outStatus:       http.StatusOK,
		out:             "",
	},
	{
		inOK:            false,
		inUserType:      "freelancer",
		inError:         nil,
		outIsNextCalled: false,
		outStatus:       http.StatusForbidden,
		out:             `"job application 1 does not belong to freelancer 2"`,
	},
	{
		inOK:            false,
		inUserType:      "client",
		inError:         nil,
		outIsNextCalled: false,
		outStatus:       http.StatusForbidden,
		out:             `"job application 1 does not belong to client 2"`,
	},
	{
		inOK:            false,
		inUserType:      "freelancer",
		inError:         errors.New("error message"),
		outIsNextCalled: false,
		outStatus:       http.StatusInternalServerError,
		out:             `"error message"`,
	},
}

func TestWhenJobApplicationBelongsToClient(t *testing.T) {
	jobRepositoryMock := &JobRepositoryMock{}
	for _, testCase := range whenJobApplicationBelongsToClientData {
		jobRepositoryMock.JobApplicationBelongsToClientCall.Returns.OK = testCase.inOK
		jobRepositoryMock.JobApplicationBelongsToClientCall.Returns.Error = testCase.inError
		jobRepositoryMock.JobApplicationBelongsToFreelancerCall.Returns.OK = testCase.inOK
		jobRepositoryMock.JobApplicationBelongsToFreelancerCall.Returns.Error = testCase.inError
		var jobContext = &ApplicationContext{
			JobRepository: jobRepositoryMock,
		}
		is := isHelper.New(t)
		w := httptest.NewRecorder()
		r := getRequest(jobContext, "")

		context.Set(r, "id", uint(1))
		context.Set(r, "user", &models.User{
			ID:   2,
			Type: testCase.inUserType,
		})

		opts := &respond.Options{
			Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
				switch data.(type) {
				case string:
					return status, data.(string)
				case error:
					return status, data.(error).Error()
				}

				return status, data
			},
		}

		var isNextCalled bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isNextCalled = true
		})

		opts.Handler(
			whenBasedOnUserType(
				whenJobApplicationBelongsToClientByID,
				whenJobApplicationBelongsToFreelancerByID,
			)(next),
		).ServeHTTP(w, r)

		switch testCase.inUserType {
		case "client":
			is.Equal(jobRepositoryMock.JobApplicationBelongsToClientCall.Receives.ID, 1)
			is.Equal(jobRepositoryMock.JobApplicationBelongsToClientCall.Receives.ClientID, 2)
		case "freelancer":
			is.Equal(jobRepositoryMock.JobApplicationBelongsToFreelancerCall.Receives.ID, 1)
			is.Equal(jobRepositoryMock.JobApplicationBelongsToFreelancerCall.Receives.FreelancerID, 2)
		}
		is.Equal(isNextCalled, testCase.outIsNextCalled)
		is.Equal(w.Code, testCase.outStatus)
		is.Equal(strings.TrimSpace(w.Body.String()), testCase.out)
	}
}
