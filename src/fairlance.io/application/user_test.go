package application

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

func TestWithUser(t *testing.T) {
	is := is.New(t)
	userContext := &ApplicationContext{}
	r := getRequest(userContext, `
		{
            "firstName": "firstname",
            "lastName": "lastname",
			"email": "email@mail.com",
			"password": "password"
		}
	`)
	w := httptest.NewRecorder()
	withUser := WithUser{
		next: func(u *User) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				is.Equal(u.FirstName, "firstname")
				is.Equal(u.LastName, "lastname")
				is.Equal(u.Email, "email@mail.com")
				is.Equal(u.Password, "password")
			})
		},
	}

	withUser.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
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
		withUser := WithUser{
			next: func(u *User) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			},
		}

		withUser.ServeHTTP(w, r)

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
	withUser := WithUser{
		next: func(u *User) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		},
	}

	withUser.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusBadRequest)
	var body map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Nil(body["firstName"])
	is.Nil(body["lastName"])
	is.Nil(body["password"])
	is.Nil(body["email"])
}
