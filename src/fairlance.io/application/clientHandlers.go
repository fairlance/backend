package main

import (
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

func GetClient(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	var id = context.Get(r, "id").(uint)
	client, err := appContext.ClientRepository.GetClient(id)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, err)
		return
	}

	respond.With(w, r, http.StatusOK, client)
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
