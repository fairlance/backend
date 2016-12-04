package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

func TestIndexFreelancer(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetAllFreelancersCall.Returns.Freelancers = []Freelancer{
		Freelancer{
			User: User{
				Model: Model{
					ID: 1,
				},
			},
		},
		Freelancer{
			User: User{
				Model: Model{
					ID: 2,
				},
			},
		},
	}
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Freelancer
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestIndexFreelancerWithError(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetAllFreelancersCall.Returns.Error = errors.New("bb")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestAddFreelancer(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	user := &User{
		Model: Model{
			ID: 1,
		},
		FirstName: "first",
		LastName:  "last",
		Email:     "email@mail.com",
	}

	AddFreelancer(user).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.Model.ID, 1)
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.FirstName, "first")
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.LastName, "last")
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.Email, "email@mail.com")
}

func TestAddFreelancerWithError(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.AddFreelancerCall.Returns.Error = errors.New("dum dum dum duuummmm")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	AddFreelancer(&User{}).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestGetFreelancerByID(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Freelancer = Freelancer{
		User: User{
			Model: Model{
				ID: 1,
			},
		},
	}
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	GetFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(freelancerRepositoryMock.GetFreelancerCall.Receives.ID, 1)
	var body Freelancer
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body.Model.ID, 1)
}

func TestGetFreelancerByIDWithError(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Error = errors.New("oopsy daisy")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	GetFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
}

func TestDeleteFreelancerByID(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	DeleteFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
}

func TestDeleteFreelancerByIDWithError(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.DeleteFreelancerCall.Returns.Error = errors.New("oopsy daisy")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	DeleteFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestWithFreelancerUpdate(t *testing.T) {
	var freelancerContext = &ApplicationContext{}
	is := is.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"hourlyRateFrom": 11,
		"hourlyRateTo": 22,
        "isAvailable": true,
        "timezone": "timez",
		"skills": ["one", "two"]
	}`
	r := getRequest(freelancerContext, requestBody)

	next := func(freelancerUpdate *FreelancerUpdate) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(freelancerUpdate.HourlyRateFrom, 11)
			is.Equal(freelancerUpdate.HourlyRateTo, 22)
			is.Equal(freelancerUpdate.IsAvailable, true)
			is.Equal(freelancerUpdate.Timezone, "timez")
			is.Equal(freelancerUpdate.Skills[0], "one")
			is.Equal(freelancerUpdate.Skills[1], "two")
		})
	}

	withFreelancerUpdate{next}.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
}

func TestWithFreelancerUpdateWithErrorMaxSkills(t *testing.T) {
	is := is.New(t)
	var freelancerContext = &ApplicationContext{}
	w := httptest.NewRecorder()

	skills := []string{}
	for i := 0; i < 21; i++ {
		skills = append(skills, fmt.Sprintf("tag%d", i))
	}

	requestBody := `{
		"hourlyRateFrom": 11,
		"hourlyRateTo": 22,
        "isAvailable": true,
        "timezone": "timez",
		"skills": ["` + strings.Join(skills, `","`) + `"]
	}`
	r := getRequest(freelancerContext, requestBody)

	next := func(freelancerUpdate *FreelancerUpdate) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}

	withFreelancerUpdate{next}.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

var badBodyWithFreelancerUpdate = []struct {
	in string
}{
	{""},
	{"{bad json}"},
}

func TestWithFreelancerUpdateWithBadBody(t *testing.T) {
	is := is.New(t)
	var freelancerContext = &ApplicationContext{}

	for _, data := range badBodyWithFreelancerUpdate {
		w := httptest.NewRecorder()
		r := getRequest(freelancerContext, data.in)

		next := func(freelancerUpdate *FreelancerUpdate) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		}

		withFreelancerUpdate{next}.ServeHTTP(w, r)

		is.Equal(w.Code, http.StatusBadRequest)
	}
}

func TestUpdateFreelancerHandler(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Freelancer = Freelancer{
		User: User{
			Model: Model{
				ID: 1,
			},
		},
	}
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	updateFreelancerHandler{1, &FreelancerUpdate{
		HourlyRateFrom: 11,
		HourlyRateTo:   22,
		IsAvailable:    true,
		Timezone:       "timez",
		Skills:         stringList{"one", "two"},
	}}.ServeHTTP(w, r)

	freelancer := freelancerRepositoryMock.UpdateFreelancerCall.Receives.Freelancer
	is.Equal(w.Code, http.StatusOK)
	is.Equal(freelancer.HourlyRateFrom, 11)
	is.Equal(freelancer.HourlyRateTo, 22)
	is.Equal(freelancer.IsAvailable, true)
	is.Equal(freelancer.Timezone, "timez")
	is.Equal(freelancer.Skills[0], "one")
	is.Equal(freelancer.Skills[1], "two")
}

func TestUpdateFreelancerHandlerFailedUpdate(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Freelancer = Freelancer{
		User: User{
			Model: Model{
				ID: 1,
			},
		},
	}
	freelancerRepositoryMock.UpdateFreelancerCall.Returns.Error = errors.New("bad updataa")
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	updateFreelancerHandler{1, &FreelancerUpdate{}}.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestUpdateFreelancerHandlerNotExistingFreelancer(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Error = errors.New("freelancer mia")
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	updateFreelancerHandler{1, &FreelancerUpdate{}}.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
	is.Equal(freelancerRepositoryMock.GetFreelancerCall.Receives.ID, 1)
}

func TestWithReview(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, `
	{
		"title":        "title",
		"clientID":     2,
		"content":      "content",
		"freelancerID": 3,
		"jobID":        4,
		"rating":       5.6
	}`)

	next := func(review *Review) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(review.Title, "title")
			is.Equal(review.ClientID, 2)
			is.Equal(review.Content, "content")
			is.Equal(review.FreelancerID, 3)
			is.Equal(review.JobID, 4)
			is.Equal(review.Rating, 5.6)
		})
	}

	withReview{next}.ServeHTTP(w, r)
}

func TestAddFreelancerReviewByID(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	review := Review{
		Title:    "title",
		ClientID: 2,
		Content:  "content",
		JobID:    4,
		Rating:   5.6,
	}

	addFreelancerReviewByID{3, &review}.ServeHTTP(w, r)
	receivedReview := freelancerRepositoryMock.AddReviewCall.Receives.Review

	is.Equal(receivedReview.Title, "title")
	is.Equal(receivedReview.Content, "content")
	is.Equal(receivedReview.Rating, 5.6)
	is.Equal(receivedReview.ClientID, 2)
	is.Equal(receivedReview.JobID, 4)
	is.Equal(freelancerRepositoryMock.AddReviewCall.Receives.ID, 3)
}

func TestWithReference(t *testing.T) {
	is := is.New(t)
	referenceRepositoryMock := &ReferenceRepositoryMock{}
	var freelancerContext = &ApplicationContext{
		ReferenceRepository: referenceRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, `
	{
		"content":      "content",
		"title":		"title",
		"media":		{
			"image":	"image",
			"video":	"video"
		}
	}`)

	next := func(reference *Reference) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(reference.Title, "title")
			is.Equal(reference.Content, "content")
			is.Equal(reference.Media.Image, "image")
			is.Equal(reference.Media.Video, "video")
		})
	}

	withReference{next}.ServeHTTP(w, r)
}

func TestAddFreelancerReferenceByID(t *testing.T) {
	is := is.New(t)
	referenceRepositoryMock := &ReferenceRepositoryMock{}
	var freelancerContext = &ApplicationContext{
		ReferenceRepository: referenceRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	reference := Reference{
		Title:   "title",
		Content: "content",
		Media: Media{
			Image: "image",
			Video: "video",
		},
	}

	addFreelancerReferenceByID{3, &reference}.ServeHTTP(w, r)
	receivedReference := referenceRepositoryMock.AddReferenceCall.Receives.Reference

	is.Equal(receivedReference.Title, "title")
	is.Equal(receivedReference.Content, "content")
	is.Equal(receivedReference.Media.Image, "image")
	is.Equal(receivedReference.Media.Video, "video")
	is.Equal(referenceRepositoryMock.AddReferenceCall.Receives.ID, 3)
}
