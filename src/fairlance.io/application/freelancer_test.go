package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	isHelper "github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestIndexFreelancer(t *testing.T) {
	is := isHelper.New(t)
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

	getAllFreelancers().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Freelancer
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestIndexFreelancerWithError(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetAllFreelancersCall.Returns.Error = errors.New("bb")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()

	getAllFreelancers().ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusInternalServerError)
	}
}

func TestAddFreelancer(t *testing.T) {
	is := isHelper.New(t)
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
	context.Set(r, "user", user)

	addFreelancer().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.Model.ID, 1)
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.FirstName, "first")
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.LastName, "last")
	is.Equal(freelancerRepositoryMock.AddFreelancerCall.Receives.Freelancer.User.Email, "email@mail.com")
}

func TestAddFreelancerWithError(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.AddFreelancerCall.Returns.Error = errors.New("dum dum dum duuummmm")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()
	context.Set(r, "user", &User{})

	addFreelancer().ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetFreelancerByID(t *testing.T) {
	is := isHelper.New(t)
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
	context.Set(r, "id", uint(1))

	getFreelancerByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(freelancerRepositoryMock.GetFreelancerCall.Receives.ID, uint(1))
	var body Freelancer
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body.Model.ID, 1)
}

func TestGetFreelancerByIDWithError(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Error = errors.New("oopsy daisy")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()
	context.Set(r, "id", uint(1))

	getFreelancerByID().ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusNotFound)
	}
}

func TestDeleteFreelancerByID(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()
	context.Set(r, "id", uint(1))

	deleteFreelancerByID().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusOK)
	}
}

func TestDeleteFreelancerByIDWithError(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.DeleteFreelancerCall.Returns.Error = errors.New("oopsy daisy")
	freelancerContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(freelancerContext, ``)
	w := httptest.NewRecorder()
	context.Set(r, "id", uint(1))

	deleteFreelancerByID().ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusBadRequest)
	}
}

func TestWithFreelancerUpdate(t *testing.T) {
	var freelancerContext = &ApplicationContext{}
	is := isHelper.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"hourlyRateFrom": 11,
		"hourlyRateTo": 22,
		"isAvailable": true,
		"timezone": "timez",
		"skills": ["one", "two"]
	}`
	r := getRequest(freelancerContext, requestBody)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	withFreelancerUpdate(next).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(nextCalled, true)
	freelancerUpdate := context.Get(r, "freelancerUpdate").(*FreelancerUpdate)
	is.Equal(freelancerUpdate.HourlyRateFrom, 11)
	is.Equal(freelancerUpdate.HourlyRateTo, 22)
	is.Equal(freelancerUpdate.IsAvailable, true)
	is.Equal(freelancerUpdate.Timezone, "timez")
	is.Equal(freelancerUpdate.Skills[0], "one")
	is.Equal(freelancerUpdate.Skills[1], "two")
}

func TestWithFreelancerUpdateWithErrorMaxSkills(t *testing.T) {
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
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})

	withFreelancerUpdate(next).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusBadRequest)
	}
}

var badBodyWithFreelancerUpdate = []struct {
	in  string
	out int
}{
	{"", http.StatusBadRequest},
	{"{bad json}", http.StatusBadRequest},
}

func TestWithFreelancerUpdateWithBadBody(t *testing.T) {
	var freelancerContext = &ApplicationContext{}

	for _, data := range badBodyWithFreelancerUpdate {
		w := httptest.NewRecorder()
		r := getRequest(freelancerContext, data.in)
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Should not be called")
		})

		withFreelancerUpdate(next).ServeHTTP(w, r)

		if w.Code != data.out {
			t.Errorf("Bad status code %d, expected %d", w.Code, data.out)
		}
	}
}

func TestUpdateFreelancerHandler(t *testing.T) {
	is := isHelper.New(t)
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
	context.Set(r, "id", uint(1))
	context.Set(r, "freelancerUpdate", &FreelancerUpdate{
		HourlyRateFrom: 11,
		HourlyRateTo:   22,
		IsAvailable:    true,
		Timezone:       "timez",
		Skills:         stringList{"one", "two"},
	})

	updateFreelancerHandler().ServeHTTP(w, r)

	freelancer := freelancerRepositoryMock.UpdateFreelancerCall.Receives.Freelancer
	is.Equal(w.Code, http.StatusOK)
	is.Equal(freelancerRepositoryMock.GetFreelancerCall.Receives.ID, uint(1))
	is.Equal(freelancer.HourlyRateFrom, 11)
	is.Equal(freelancer.HourlyRateTo, 22)
	is.Equal(freelancer.IsAvailable, true)
	is.Equal(freelancer.Timezone, "timez")
	is.Equal(freelancer.Skills[0], "one")
	is.Equal(freelancer.Skills[1], "two")
}

func TestUpdateFreelancerHandlerFailedUpdate(t *testing.T) {
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
	context.Set(r, "id", uint(1))
	context.Set(r, "freelancerUpdate", &FreelancerUpdate{})

	updateFreelancerHandler().ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdateFreelancerHandlerNotExistingFreelancer(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Error = errors.New("freelancer mia")
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)
	context.Set(r, "id", uint(1))
	context.Set(r, "freelancerUpdate", &FreelancerUpdate{})

	updateFreelancerHandler().ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusNotFound)
	}

	if freelancerRepositoryMock.GetFreelancerCall.Receives.ID != 1 {
		t.Errorf("Wrong freelancerID received %d, expected %d", freelancerRepositoryMock.GetFreelancerCall.Receives.ID, 1)
	}
}

func TestWithReview(t *testing.T) {
	is := isHelper.New(t)
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

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	withReview(next).ServeHTTP(w, r)

	is.Equal(nextCalled, true)
	review := context.Get(r, "review").(*Review)
	is.Equal(review.Title, "title")
	is.Equal(review.ClientID, 2)
	is.Equal(review.Content, "content")
	is.Equal(review.FreelancerID, 3)
	is.Equal(review.JobID, 4)
	is.Equal(review.Rating, 5.6)
}

var bodyWithReview = []struct {
	in  string
	out int
}{
	{"", http.StatusBadRequest},
	{"{bad json}", http.StatusBadRequest},
	{`{
		"clientID":     2,
		"content":      "content",
		"jobID":        4,
		"rating":       5.6
	}`, http.StatusBadRequest},
	{`{
		"title":        "title",
		"content":      "content",
		"jobID":        4,
		"rating":       5.6
	}`, http.StatusBadRequest},
	{`{
		"title":        "title",
		"clientID":     2,
		"content":      "content",
		"rating":       5.6
	}`, http.StatusBadRequest},
	{`{
		"title":        "title",
		"clientID":     2,
		"content":      "content",
		"jobID":        4
	}`, http.StatusBadRequest},
	{`{
		"title":        "title",
		"clientID":     "2",
		"content":      "content",
		"jobID":        4,
		"rating":       5.6
	}`, http.StatusBadRequest},
	{`{
		"title":        "title",
		"clientID":     2,
		"content":      "content",
		"jobID":        "4",
		"rating":       5.6
	}`, http.StatusBadRequest},
	{`{
		"title":        "title",
		"clientID":     2,
		"content":      "content",
		"jobID":        4,
		"rating":       "5.6"
	}`, http.StatusBadRequest},
	{`{
		"title":      "no content",
		"clientID":     2,
		"jobID":        4,
		"rating":       5.6
	}`, http.StatusOK},
}

func TestWithReviewWithBadBody(t *testing.T) {
	var freelancerContext = &ApplicationContext{}

	for _, data := range bodyWithReview {
		w := httptest.NewRecorder()
		r := getRequest(freelancerContext, data.in)
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if data.out != http.StatusOK {
				t.Error("Should not be called")
			}
		})

		withReview(next).ServeHTTP(w, r)

		if w.Code != data.out {
			t.Errorf("Bad status code %d, expected %d\nFor request body: %s\nResponse body: %s", w.Code, data.out, data.in, w.Body.String())
		}
	}
}

func TestAddFreelancerReviewByID(t *testing.T) {
	is := isHelper.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	review := &Review{
		Title:    "title",
		ClientID: 2,
		Content:  "content",
		JobID:    4,
		Rating:   5.6,
	}
	context.Set(r, "id", uint(3))
	context.Set(r, "review", review)

	addFreelancerReviewByID().ServeHTTP(w, r)
	receivedReview := freelancerRepositoryMock.AddReviewCall.Receives.Review

	is.Equal(freelancerRepositoryMock.AddReviewCall.Receives.ID, uint(3))
	is.Equal(receivedReview.Title, "title")
	is.Equal(receivedReview.Content, "content")
	is.Equal(receivedReview.Rating, 5.6)
	is.Equal(receivedReview.ClientID, 2)
	is.Equal(receivedReview.JobID, 4)
	is.Equal(freelancerRepositoryMock.AddReviewCall.Receives.ID, 3)
}

func TestAddFreelancerReviewByIDWithError(t *testing.T) {
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.AddReviewCall.Returns.Error = errors.New("review krak")
	var freelancerContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)
	context.Set(r, "id", uint(3))
	context.Set(r, "review", &Review{})

	addFreelancerReviewByID().ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusInternalServerError)
	}
}

func TestWithReference(t *testing.T) {
	is := isHelper.New(t)
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

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	withReference(next).ServeHTTP(w, r)

	reference := context.Get(r, "reference").(*Reference)

	is.Equal(nextCalled, true)
	is.Equal(reference.Title, "title")
	is.Equal(reference.Content, "content")
	is.Equal(reference.Media.Image, "image")
	is.Equal(reference.Media.Video, "video")
}

var bodyWithReference = []struct {
	in  string
	out int
}{
	{"", http.StatusBadRequest},
	{"{bad json}", http.StatusBadRequest},
	{`{
		"content":      "no title",
		"media":		{
			"image":	"image",
			"video":	"video"
		}
	}`, http.StatusBadRequest},
	{`{
		"content":      "content",
		"title":		"title"
	}`, http.StatusOK},
	{`{
		"title":		"title"
	}`, http.StatusOK},
}

func TestWithReferenceWithBadBody(t *testing.T) {
	var freelancerContext = &ApplicationContext{}

	for _, data := range bodyWithReference {
		w := httptest.NewRecorder()
		r := getRequest(freelancerContext, data.in)

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if data.out != http.StatusOK {
				t.Error("Should not be called")
			}
		})

		withReference(next).ServeHTTP(w, r)

		if w.Code != data.out {
			t.Errorf("Bad status code %d, expected %d\nFor request body: %s\nResponse body: %s", w.Code, data.out, data.in, w.Body.String())
		}
	}
}

func TestAddFreelancerReferenceByID(t *testing.T) {
	is := isHelper.New(t)
	referenceRepositoryMock := &ReferenceRepositoryMock{}
	var freelancerContext = &ApplicationContext{
		ReferenceRepository: referenceRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)

	reference := &Reference{
		Title:   "title",
		Content: "content",
		Media: Media{
			Image: "image",
			Video: "video",
		},
	}
	context.Set(r, "id", uint(3))
	context.Set(r, "reference", reference)

	addFreelancerReferenceByID().ServeHTTP(w, r)

	receivedReference := referenceRepositoryMock.AddReferenceCall.Receives.Reference

	is.Equal(referenceRepositoryMock.AddReferenceCall.Receives.ID, uint(3))
	is.Equal(receivedReference.Title, "title")
	is.Equal(receivedReference.Content, "content")
	is.Equal(receivedReference.Media.Image, "image")
	is.Equal(receivedReference.Media.Video, "video")
	is.Equal(referenceRepositoryMock.AddReferenceCall.Receives.ID, 3)
}

func TestAddFreelancerReferenceByIDWithError(t *testing.T) {
	referenceRepositoryMock := &ReferenceRepositoryMock{}
	referenceRepositoryMock.AddReferenceCall.Returns.Error = errors.New("darn it")
	var freelancerContext = &ApplicationContext{
		ReferenceRepository: referenceRepositoryMock,
	}

	w := httptest.NewRecorder()
	r := getRequest(freelancerContext, ``)
	context.Set(r, "id", uint(3))
	context.Set(r, "reference", &Reference{})

	addFreelancerReferenceByID().ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Bad status code %d, expected %d", w.Code, http.StatusInternalServerError)
	}
}
