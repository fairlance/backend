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
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
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
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestAddFreelancer(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
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
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
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
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
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
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	GetFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
}

func TestDeleteFreelancerByID(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	DeleteFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
}

func TestDeleteFreelancerByIDWithError(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.DeleteFreelancerCall.Returns.Error = errors.New("oopsy daisy")
	userContext := &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	DeleteFreelancerByID(1).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestWithFreelancerUpdate(t *testing.T) {
	var jobContext = &ApplicationContext{}
	is := is.New(t)
	w := httptest.NewRecorder()
	requestBody := `{
		"hourlyRateFrom": 11,
		"hourlyRateTo": 22,
        "isAvailable": true,
        "timezone": "timez",
		"skills": ["one", "two"]
	}`
	r := getRequest(jobContext, requestBody)

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
	var jobContext = &ApplicationContext{}
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
	r := getRequest(jobContext, requestBody)

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
	var jobContext = &ApplicationContext{}

	for _, data := range badBodyWithFreelancerUpdate {
		w := httptest.NewRecorder()
		r := getRequest(jobContext, data.in)

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
	var userContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(userContext, ``)

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
	var userContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(userContext, ``)

	updateFreelancerHandler{1, &FreelancerUpdate{}}.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
}

func TestUpdateFreelancerHandlerNotExistingFreelancer(t *testing.T) {
	is := is.New(t)
	freelancerRepositoryMock := &FreelancerRepositoryMock{}
	freelancerRepositoryMock.GetFreelancerCall.Returns.Error = errors.New("freelancer mia")
	var userContext = &ApplicationContext{
		FreelancerRepository: freelancerRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(userContext, ``)

	updateFreelancerHandler{1, &FreelancerUpdate{}}.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
	is.Equal(freelancerRepositoryMock.GetFreelancerCall.Receives.ID, 1)
}
