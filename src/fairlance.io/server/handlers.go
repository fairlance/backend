package main

import (
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type appHandler struct {
	context *appContext
	handle  func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: set content type, and other headers
	status, err := ah.handle(ah.context, w, r)
	// handle errors
	if err != nil {
		// todo: implement better error handling
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			log.Printf("HTTP %d: %q", status, err)
		default:
			log.Printf("HTTP %d: %q", status, err)
		}
	}
}

func IndexHandler(context *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		// todo: err is never used
		return http.StatusNotFound, errors.New("Bad method.")
	}

	users, err := getAllRegisteredUsers(context)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return http.StatusInternalServerError, err
	}

	json.NewEncoder(w).Encode(users)
	return http.StatusOK, nil
}

func RegisterHandler(context *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "POST" {
		// todo: err is never used
		return http.StatusNotFound, errors.New("Bad method.")
	}

	email := r.FormValue("email")

	if email != "" {

		err := addRegisteredUser(context, email)
		if err != nil {
			if mgo.IsDup(err) {
				json.NewEncoder(w).Encode(struct {
					Error string `json:"error"`
				}{"Email exists!"})
				return http.StatusBadRequest, err
			}

			json.NewEncoder(w).Encode(err)
			return http.StatusInternalServerError, err
		}

		json.NewEncoder(w).Encode(struct {
			Email string `json:"email"`
		}{email})
		return http.StatusOK, nil
	}

	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{"Email missing!"})
	return http.StatusBadRequest, nil
}
