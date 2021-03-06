package application

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/fairlance/backend/models"
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
		var user = context.Get(r, "userToAdd").(*User)
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

func withClientUpdateFromRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientUpdate ClientUpdate
		if err := json.NewDecoder(r.Body).Decode(&clientUpdate); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()
		if ok, err := govalidator.ValidateStruct(clientUpdate); ok == false || err != nil {
			respond.With(w, r, http.StatusBadRequest, models.GovalidatorErrors{Err: err})
			return
		}
		context.Set(r, "clientUpdate", &clientUpdate)
		next.ServeHTTP(w, r)
	})
}

func updateClientByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		var clientUpdate = context.Get(r, "clientUpdate").(*ClientUpdate)
		var id = context.Get(r, "id").(uint)
		client, err := appContext.ClientRepository.GetClient(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}
		client.Image = clientUpdate.Image
		client.Birthdate = clientUpdate.Birthdate
		client.About = clientUpdate.About
		client.Timezone = clientUpdate.Timezone
		client.ProfileCompleted = true
		if err := appContext.ClientRepository.UpdateClient(client); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}

func withClientFromJobByID(handler http.Handler) http.Handler {
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
