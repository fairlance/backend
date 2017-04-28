package application

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestWithUser(t *testing.T) {
	is := is.New(t)
	userContext := &ApplicationContext{}
	r := getRequest(userContext, `
		{
            "firstName": "firstname",
            "lastName": "lastname",
			"email": "email@mail.com",
			"password": "password",
			"image": "http://url.for/image"
		}
	`)
	w := httptest.NewRecorder()

	nextCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	withUser(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(nextCalled, true)
	u := context.Get(r, "user").(*User)
	is.Equal(u.FirstName, "firstname")
	is.Equal(u.LastName, "lastname")
	is.Equal(u.Email, "email@mail.com")
	is.Equal(u.Password, "password")
	is.Equal(u.Image, "http://url.for/image")
}

var badBodyWithUserTestData = []struct {
	in string
}{
	{""},
	{"{bad json}"},
}

func TestWithUserWithBadBody(t *testing.T) {
	is := is.New(t)
	userContext := &ApplicationContext{}
	for _, data := range badBodyWithUserTestData {
		r := getRequest(userContext, data.in)
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Should not be called")
		})

		withUser(handler).ServeHTTP(w, r)

		is.Equal(w.Code, http.StatusBadRequest)
		var body map[string]interface{}
		is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	}
}

func TestWithUserWithNotAllDataInBody(t *testing.T) {
	is := is.New(t)
	userContext := &ApplicationContext{}
	r := getRequest(userContext, `{}`)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not be called")
	})

	withUser(handler).ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
	var body map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body["firstName"], "non zero value required")
	is.Equal(body["lastName"], "non zero value required")
	is.Equal(body["password"], "non zero value required")
	is.Equal(body["email"], "non zero value required")
	is.Equal(body["image"], "non zero value required")
}
