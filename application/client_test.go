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

func TestIndexClient(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.GetAllClientsCall.Returns.Clients = []Client{
		Client{
			User: User{
				Model:     Model{ID: 1},
				FirstName: "firstname",
				LastName:  "lastname",
				Email:     "email@mail.com",
				Password:  "password",
			},
		},
		Client{
			User: User{
				Model:     Model{ID: 2},
				FirstName: "firstname2",
				LastName:  "lastname2",
				Email:     "email2@mail.com",
				Password:  "password2",
			},
		},
	}
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	getAllClients().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Client
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestIndexClientWithError(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.GetAllClientsCall.Returns.Error = errors.New("Clients kabooom")
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	getAllClients().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestGetClientByIDWithError(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.GetClientCall.Returns.Error = errors.New("Clients kabooom")
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()
	context.Set(r, "id", uint(1))

	getClientByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusNotFound)
}

func TestGetClientByID(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.GetClientCall.Returns.Client = &Client{
		User: User{
			Model: Model{
				ID: 1,
			},
		},
	}
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()
	context.Set(r, "id", uint(1))

	getClientByID().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body["id"], 1)
	is.Equal(clientRepositoryMock.GetClientCall.Receives.ID, uint(1))
}

func TestAddClientWithError(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.AddClientCall.Returns.Error = errors.New("bummer")
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	user := &User{
		Model: Model{
			ID: 1,
		},
	}
	context.Set(r, "userToAdd", user)

	addClient().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}

func TestAddClient(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	user := &User{
		Model: Model{
			ID: 1,
		},
	}
	context.Set(r, "userToAdd", user)

	addClient().ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(clientRepositoryMock.AddClientCall.Receives.Client.User.ID, 1)
	var body map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	userMap := body["user"].(map[string]interface{})
	is.Equal(userMap["id"], 1)
	is.Equal(body["type"], "client")
}

var badBodyUpdateClientByID = []struct {
	in string
}{
	{""},
	{"{bad json}"},
}

func TestUpdateClientByIDWithBadBody(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}
	for _, data := range badBodyUpdateClientByID {
		w := httptest.NewRecorder()
		r := getRequest(userContext, data.in)
		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { nextCalled = true })
		withClientUpdateFromRequest(next).ServeHTTP(w, r)
		is.Equal(nextCalled, false)
		is.Equal(w.Code, http.StatusBadRequest)
	}
}

func TestUpdateClientByIDWithNonExistingClient(t *testing.T) {
	is := isHelper.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.GetClientCall.Returns.Error = errors.New("nope")
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}
	w := httptest.NewRecorder()
	r := getRequest(userContext, "")
	context.Set(r, "id", uint(1))
	context.Set(r, "clientUpdate", &ClientUpdate{})
	updateClientByID().ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusNotFound)
}

func TestUpdateClientByIDWithFailedUpdate(t *testing.T) {
	is := isHelper.New(t)
	userContext := &ApplicationContext{
		ClientRepository: &ClientRepositoryMock{},
	}
	w := httptest.NewRecorder()
	r := getRequest(userContext, "{}")
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { nextCalled = true })
	withClientUpdateFromRequest(next).ServeHTTP(w, r)
	is.Equal(nextCalled, false)
	is.Equal(w.Code, http.StatusBadRequest)
}

func TestUpdateClientByID(t *testing.T) {
	is := isHelper.New(t)
	userContext := &ApplicationContext{
		ClientRepository: &ClientRepositoryMock{},
	}
	r := getRequest(userContext, `
    {
        "timezone": "UTC",
        "about": "about",
		"birthdate": "27.07.2007",
		"image": "url"
    }`)
	w := httptest.NewRecorder()
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { nextCalled = true })
	withClientUpdateFromRequest(next).ServeHTTP(w, r)
	is.Equal(nextCalled, true)
	clientUpdate := context.Get(r, "clientUpdate").(*ClientUpdate)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(clientUpdate.Timezone, "UTC")
	is.Equal(clientUpdate.Birthdate, "27.07.2007")
	is.Equal(clientUpdate.Image, "url")
	is.Equal(clientUpdate.About, "about")
}
