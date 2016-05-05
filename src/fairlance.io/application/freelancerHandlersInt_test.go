package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func xTestIndexFreelancerWhenEmpty(t *testing.T) {
	setUp()
	is := is.New(t)

	w := httptest.NewRecorder()
	r := getRequest("GET", "")
	IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data []interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, []interface{}{})
}

func xTestIndexFreelancerWithFreelancers(t *testing.T) {
	setUp()
	is := is.New(t)
	AddFreelancerToDB()
	AddFreelancerToDB()

	w := httptest.NewRecorder()
	r := getRequest("GET", "")
	IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data []interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(len(data), 2)
}

func xTestAddFreelancer(t *testing.T) {
	setUp()
	is := is.New(t)

	w := httptest.NewRecorder()
	r := getRequest("POST", "")
	context.Set(r, "freelancer", GetMockFreelancer())

	AddFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.NotEqual(data["id"], 0)
	is.Equal(data["firstName"], "Pera")
	is.Equal(data["lastName"], "Peric")
	is.Equal(data["email"], "pera@gmail.com")
	is.Equal(data["title"], "Dev")
	is.Equal(data["hourlyRateFrom"], 12)
	is.Equal(data["hourlyRateTo"], 22)
	is.Equal(data["timeZone"], "CET")

	freelancers := GetFreelancersFromDB()
	is.Equal(len(freelancers), 1)
	is.NotEqual(freelancers[0].ID, 0)
	is.Equal(freelancers[0].FirstName, "Pera")
	is.Equal(freelancers[0].LastName, "Peric")
	is.Equal(freelancers[0].Email, "pera@gmail.com")
	is.Equal(freelancers[0].Title, "Dev")
	is.Equal(freelancers[0].HourlyRateFrom, 12)
	is.Equal(freelancers[0].HourlyRateTo, 22)
	is.Equal(freelancers[0].TimeZone, "CET")
}

func xTestDeleteFreelancer(t *testing.T) {
	setUp()
	is := is.New(t)
	id := AddFreelancerToDB()

	w := httptest.NewRecorder()
	r := getRequest("POST", "")
	context.Set(r, "id", id)

	DeleteFreelancer(w, r)

	var data map[string]interface{}
	is.Equal(w.Code, http.StatusOK)
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))

	is.Equal(len(GetFreelancersFromDB()), 0)
}

func GetMockFreelancer() *Freelancer {
	return NewFreelancer(
		"Pera",
		"Peric",
		"Dev",
		"$2a$10$VJ8H9EYOIj9mnyW5mUm/nOWUrz/Rkak4/Ov3Lnw1GsAm4gmYU6sQu",
		"pera@gmail.com",
		12,
		22,
		"CET",
	)
}

func AddFreelancerToDB() uint {
	f := GetMockFreelancer()
	appContext.FreelancerRepository.AddFreelancer(f)
	return f.ID
}

func GetFreelancersFromDB() []Freelancer {
	f, _ := appContext.FreelancerRepository.GetAllFreelancers()
	return f
}
