package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"net/http"
)

type appHandler struct {
	context *appContext
	handle  func(*appContext, http.ResponseWriter, *http.Request) error
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: make this configurable
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := ah.handle(ah.context, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func indexHandler(context *appContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"Method not allowed! Use GET"})
		return nil
	}

	users, err := getAllRegisteredUsers(context)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	return nil
}

func registerHandler(context *appContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"Method not allowed! Use POST"})
		return nil
	}

	email := r.FormValue("email")

	if email != "" {

		err := addRegisteredUser(context, email)
		if err != nil {
			if mgo.IsDup(err) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(struct {
					Error string `json:"error"`
				}{"Email exists!"})
				return nil
			}

			return err
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Email string `json:"email"`
		}{email})
		return nil
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{"Email missing!"})
	return nil
}
