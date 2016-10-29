package application

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func IndexClient(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	clients, err := appContext.ClientRepository.GetAllClients()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, clients)
}

// GetClientByID handler
func GetClientByID(id uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		client, err := appContext.ClientRepository.GetClient(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, client)
	})
}

func AddClient(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*User)
	client := &Client{User: *user}
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.ClientRepository.AddClient(client); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, struct {
		User User   `json:"user"`
		Type string `json:"type"`
	}{
		User: client.User,
		Type: "client",
	})
}

// UpdateClientByID handler
func UpdateClientByID(id uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var body struct {
			Timezone string `json:"timezone"`
			Payment  string `json:"payment"`
			Industry string `json:"industry"`
		}

		if err := decoder.Decode(&body); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		var appContext = context.Get(r, "context").(*ApplicationContext)
		client, err := appContext.ClientRepository.GetClient(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		client.Timezone = body.Timezone
		client.Payment = body.Payment
		client.Industry = body.Industry

		if err := appContext.ClientRepository.UpdateClient(&client); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}
