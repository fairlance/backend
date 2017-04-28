package application

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func getAllClients() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		clients, err := appContext.ClientRepository.GetAllClients()
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, clients)
	})
}

func getClientByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var id = context.Get(r, "id").(uint)
		client, err := appContext.ClientRepository.GetClient(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, client)
	})
}

func addClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user = context.Get(r, "user").(*User)
		client := &Client{User: *user}
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.ClientRepository.AddClient(client); err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, struct {
			User User   `json:"user"`
			Type string `json:"type"`
		}{
			User: client.User,
			Type: "client",
		})
	})
}

func updateClientByID() http.Handler {
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
		var id = context.Get(r, "id").(uint)
		client, err := appContext.ClientRepository.GetClient(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		if body.Timezone != "" {
			client.Timezone = body.Timezone
		}
		if body.Industry != "" {
			client.Payment = body.Payment
		}
		if body.Industry != "" {
			client.Industry = body.Industry
		}

		if err := appContext.ClientRepository.UpdateClient(client); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}

func withClientFromJobID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var jobID = context.Get(r, "id").(uint)

		job, err := appContext.JobRepository.GetJob(jobID)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		context.Set(r, "client", job.Client)

		handler.ServeHTTP(w, r)
	})
}
