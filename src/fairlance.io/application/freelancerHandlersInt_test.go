package application_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	app "fairlance.io/application"
	"github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestIndexFreelancerWhenEmpty(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIndexFreelancerWhenEmpty in short mode")
	}
	setUp()
	is := is.New(t)

	w := httptest.NewRecorder()
	r := getRequest("GET", "")
	app.IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data []interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, []interface{}{})
}

func TestIndexFreelancerWithFreelancers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestIndexFreelancerWithFreelancers in short mode")
	}
	setUp()
	is := is.New(t)
	AddFreelancerToDB()
	AddFreelancerToDB()

	w := httptest.NewRecorder()
	r := getRequest("GET", "")
	app.IndexFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data []interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(len(data), 2)
}

func TestAddFreelancer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAddFreelancer in short mode")
	}
	setUp()
	is := is.New(t)

	w := httptest.NewRecorder()
	r := getRequest("PUT", "")
	context.Set(r, "user", GetMockUser())

	app.AddFreelancer(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))

	user := data["user"].(map[string]interface{})
	is.Equal(data["type"], "freelancer")
	is.NotEqual(user["id"], 0)
	is.Equal(user["firstName"], "Pera")
	is.Equal(user["lastName"], "Peric")
	is.True(strings.HasSuffix(user["email"].(string), "pera@gmail.com"))

	freelancers := GetFreelancersFromDB()
	is.Equal(len(freelancers), 1)
	is.NotEqual(freelancers[0].ID, 0)
	is.Equal(freelancers[0].FirstName, "Pera")
	is.Equal(freelancers[0].LastName, "Peric")
	is.True(strings.HasSuffix(freelancers[0].Email, "pera@gmail.com"))
}

func TestDeleteFreelancer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestDeleteFreelancer in short mode")
	}
	setUp()
	is := is.New(t)
	id := AddFreelancerToDB()

	w := httptest.NewRecorder()
	r := getRequest("POST", "")
	context.Set(r, "id", id)

	app.DeleteFreelancer(w, r)

	var data map[string]interface{}
	is.Equal(w.Code, http.StatusOK)
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))

	is.Equal(len(GetFreelancersFromDB()), 0)
}

func GetMockUser() *app.User {
	var email string
	rand.Seed(time.Now().UTC().UnixNano())
	email = strconv.Itoa(rand.Intn(100)) + "pera@gmail.com"
	return &app.User{
		FirstName: "Pera",
		LastName:  "Peric",
		Password:  "$2a$10$VJ8H9EYOIj9mnyW5mUm/nOWUrz/Rkak4/Ov3Lnw1GsAm4gmYU6sQu",
		Email:     email,
	}
}

func AddFreelancerToDB() uint {
	u := GetMockUser()
	f := &app.Freelancer{User: *u}
	appContext.FreelancerRepository.AddFreelancer(f)
	return f.ID
}

func GetFreelancersFromDB() []app.Freelancer {
	f, _ := appContext.FreelancerRepository.GetAllFreelancers()
	return f
}
