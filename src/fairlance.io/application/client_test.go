package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

func TestIndexClient(t *testing.T) {
	is := is.New(t)
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

	IndexClient(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body []Client
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body[0].Model.ID, 1)
	is.Equal(body[1].Model.ID, 2)
}

func TestIndexClientWithError(t *testing.T) {
	is := is.New(t)
	clientRepositoryMock := &ClientRepositoryMock{}
	clientRepositoryMock.GetAllClientsCall.Returns.Error = errors.New("Clients kabooom")
	userContext := &ApplicationContext{
		ClientRepository: clientRepositoryMock,
	}

	r := getRequest(userContext, ``)
	w := httptest.NewRecorder()

	IndexClient(w, r)

	is.Equal(w.Code, http.StatusInternalServerError)
}
