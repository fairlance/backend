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

func TestJobIndexJob(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GettAllJobsCall.Returns.Jobs =
		[]Job{
			Job{
				Model: Model{
					ID: 1,
				},
				Name:     "Name1",
				Summary:  "Summary1",
				Details:  "Details1",
				ClientID: 1,
				IsActive: true,
				Price:    100,
				Examples: []Example{
					{
						Description: "example",
						URL:         "www.example.com",
					},
				},
				Attachments: []Attachment{
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
	is.Equal(body[0].IsActive, true)
	is.Equal(body[0].Price, 100)
	is.Equal(len(body[0].Examples), 1)
	is.Equal(body[0].Examples[0].Description, "example")
	is.Equal(body[0].Examples[0].URL, "www.example.com")
	is.Equal(len(body[0].Attachments), 1)
	is.Equal(body[0].Attachments[0].Name, "attachment")
	is.Equal(body[0].Attachments[0].URL, "www.attachment.com")
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
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "job", &Job{
		Name:     "Name1",
		Summary:  "Summary1",
		Details:  "Details1",
		ClientID: 1,
		IsActive: true,
		Price:    100,
		Examples: []Example{
			{
				Description: "example",
				URL:         "www.example.com",
			},
		},
		Attachments: []Attachment{
			{
				Name: "attachment",
				URL:  "www.attachment.com",
			},
		},
	})

	addJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Name, "Name1")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Summary, "Summary1")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Details, "Details1")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.ClientID, 1)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.IsActive, true)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Price, 100)
	is.Equal(len(jobRepositoryMock.AddJobCall.Receives.Job.Examples), 1)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Examples[0].Description, "example")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Examples[0].URL, "www.example.com")
	is.Equal(len(jobRepositoryMock.AddJobCall.Receives.Job.Attachments), 1)
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Attachments[0].Name, "attachment")
	is.Equal(jobRepositoryMock.AddJobCall.Receives.Job.Attachments[0].URL, "www.attachment.com")
}

func TestJobGetJobByIDReceivesTheRightID(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(1))

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.GetJobCall.Receives.ID, uint(1))
}

func TestJobGetJobByID(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GetJobCall.Returns.Job = Job{
		Model: Model{
			ID: 123456789,
		},
		Name:     "Name1",
		Summary:  "Summary1",
		Details:  "Details1",
		ClientID: 1,
		IsActive: true,
		Price:    100,
	}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(0))

	getJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body Job
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body.Model.ID, uint(123456789))
	is.Equal(body.Name, "Name1")
	is.Equal(body.Summary, "Summary1")
	is.Equal(body.Details, "Details1")
	is.Equal(body.IsActive, true)
	is.Equal(body.Price, 100)
}

func TestJobGetJobByIDError(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	jobRepositoryMock.GetJobCall.Returns.Error = errors.New("Blah")
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	context.Set(r, "id", uint(0))

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

	addJob().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobHandleApplyForJob(t *testing.T) {
	jobRepositoryMock := JobRepositoryMock{}
	var jobContext = &ApplicationContext{
		JobRepository: &jobRepositoryMock,
	}

	is := isHelper.New(t)
	w := httptest.NewRecorder()
	r := getRequest(jobContext, "")
	jobApplication := &JobApplication{
		FreelancerID: 1,
		Message:      "message",
		Samples:      []uint{1, 2},
		Milestones:   []string{"one", "two"},
		HourPrice:    1.1,
		Hours:        1,
		Examples:     []Example{{Description: "example", URL: "www.example.com"}},
		Attachments:  []Attachment{{Name: "attachment", URL: "www.attachment.com"}},
	}
	context.Set(r, "id", uint(1))
	context.Set(r, "jobApplication", jobApplication)

	addJobApplicationByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.JobID, 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.FreelancerID, 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Message, "message")
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Samples, []uint{1, 2})
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Milestones, []string{"one", "two"})
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.HourPrice, 1.1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Hours, 1)
	is.Equal(len(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Examples), 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Examples[0].Description, "example")
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Examples[0].URL, "www.example.com")
	is.Equal(len(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Attachments), 1)
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Attachments[0].Name, "attachment")
	is.Equal(jobRepositoryMock.AddJobApplicationCall.Receives.JobApplication.Attachments[0].URL, "www.attachment.com")
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

	addJobApplicationByID().ServeHTTP(w, r)

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
		"clientId": 1,
		"examples": [
			{"description": "example", "url": "www.example.com"}
		],
		"attachments": [
			{"name": "attachment", "url": "www.attachment.com"}
		]
	}`
	r := getRequest(jobContext, requestBody)

	nextCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	withJob(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(nextCalled, true)
	job := context.Get(r, "job").(*Job)
	is.Equal(job.Name, "Name1")
	is.Equal(job.Details, "Details1")
	is.Equal(job.Summary, "Summary1")
	is.Equal(job.ClientID, 1)
	is.Equal(len(job.Examples), 1)
	is.Equal(job.Examples[0].Description, "example")
	is.Equal(job.Examples[0].URL, "www.example.com")
	is.Equal(len(job.Attachments), 1)
	is.Equal(job.Attachments[0].Name, "attachment")
	is.Equal(job.Attachments[0].URL, "www.attachment.com")
}

func TestJobWithJobError(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{}`
	r := getRequest(jobContext, requestBody)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})
	withJob(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
	var body map[string]string
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(len(body), 4)
	is.Equal(body["ClientID"], "non zero value required")
	is.Equal(body["Details"], "non zero value required")
	is.Equal(body["Name"], "non zero value required")
	is.Equal(body["Summary"], "non zero value required")
}

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
	withJob(handler).ServeHTTP(w, r)

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
	withJob(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestJobWithJobApplication(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"message":"message",
		"milestones": ["one", "two"],
		"samples": [1,2],
		"deliveryEstimate": 3,
		"freelancerId": 1,
		"hours": 1,
		"hourPrice": 1.1,
		"examples": [
			{"description": "example", "url": "www.example.com"}
		],
		"attachments": [
			{"name": "attachment", "url": "www.attachment.com"}
		]
	}`
	r := getRequest(jobContext, requestBody)

	nextCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	withJobApplication(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(nextCalled, true)
	jobApplication := context.Get(r, "jobApplication").(*JobApplication)
	is.Equal(jobApplication.Message, "message")
	is.Equal(jobApplication.Milestones, []string{"one", "two"})
	is.Equal(jobApplication.Samples, []uint{1, 2})
	is.Equal(jobApplication.DeliveryEstimate, 3)
	is.Equal(jobApplication.FreelancerID, 1)
	is.Equal(jobApplication.Hours, 1)
	is.Equal(jobApplication.HourPrice, 1.1)
	is.Equal(len(jobApplication.Examples), 1)
	is.Equal(jobApplication.Examples[0].Description, "example")
	is.Equal(jobApplication.Examples[0].URL, "www.example.com")
	is.Equal(len(jobApplication.Attachments), 1)
	is.Equal(jobApplication.Attachments[0].Name, "attachment")
	is.Equal(jobApplication.Attachments[0].URL, "www.attachment.com")
}

func TestJobWithJobApplicationError(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{}`
	r := getRequest(jobContext, requestBody)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})
	withJobApplication(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
	var body map[string]string
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(len(body), 7)
	is.Equal(body["Message"], "non zero value required")
	is.Equal(body["Milestones"], "non zero value required")
	is.Equal(body["Samples"], "non zero value required")
	is.Equal(body["DeliveryEstimate"], "non zero value required")
	is.Equal(body["FreelancerID"], "non zero value required")
	is.Equal(body["Hours"], "non zero value required")
	is.Equal(body["HourPrice"], "non zero value required")
}

func TestJobWithJobApplicationErrorBadJSON(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{bad:json}`
	r := getRequest(jobContext, requestBody)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})
	withJobApplication(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}
