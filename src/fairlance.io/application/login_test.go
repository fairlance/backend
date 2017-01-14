package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	emptyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

var badBodyLoginTestData = []struct {
	in string
}{
	{"{}"},
	{""},
	{"{bad json}"},
	{`{"email":"no@password.com"}`},
	{`{"password":"noemail"}`},
}

func TestLoginWithBadEmailAndPassword(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &UserRepositoryMock{}
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	for _, data := range badBodyLoginTestData {
		r := getRequest(userContext, data.in)
		w := httptest.NewRecorder()

		login().ServeHTTP(w, r)

		is.Equal(w.Code, http.StatusBadRequest)
	}
}

func TestLoginWhenUnauthorized(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &UserRepositoryMock{}
	userRepositoryMock.CheckCredentialsCall.Returns.Error = errors.New("unauthorized")
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	r := getRequest(userContext, `
		{
			"email": "not@important.com",
			"password": "notimportant"
		}
	`)
	w := httptest.NewRecorder()

	login().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusUnauthorized)
}

func TestLoginCallsCheckCredentialsWithCorrectEmailAndPassword(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &UserRepositoryMock{}
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	r := getRequest(userContext, `
		{
			"email": "email@mail.com",
			"password": "password"
		}
	`)
	w := httptest.NewRecorder()

	login().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(userRepositoryMock.CheckCredentialsCall.Receives.Email, "email@mail.com")
	is.Equal(userRepositoryMock.CheckCredentialsCall.Receives.Password, "password")
}

func TestLogin(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &UserRepositoryMock{}
	userRepositoryMock.CheckCredentialsCall.Returns.UserType = "freelancer"
	userRepositoryMock.CheckCredentialsCall.Returns.User = User{
		Model: Model{
			ID: 1,
		},
		FirstName: "firstname",
		LastName:  "lastname",
		Password:  "password",
		Email:     "email@mail.com",
	}
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	r := getRequest(userContext, `
		{
			"email": "not@important.com",
			"password": "notimportant"
		}
	`)
	w := httptest.NewRecorder()

	login().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body["id"], 1)
	is.OK(body["token"])
	is.Equal(body["type"], "freelancer")
}
