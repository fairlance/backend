package main

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
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

func IndexHandler(context *appContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"Method not allowed! Use GET"})
		return nil
	}

	users, err := context.userRepository.GetAllRegisteredUsers()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	return nil
}

func RegisterHandler(context *appContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"Method not allowed! Use POST"})
		return nil
	}

	email := r.FormValue("email")

	if email != "" {

		// validate email first
		if !govalidator.IsEmail(email) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				Error string `json:"error"`
			}{"Email not valid!"})
			return nil
		}

		err := context.userRepository.AddRegisteredUser(email)
		if err != nil {
			if mgo.IsDup(err) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(struct {
					Error string `json:"error"`
				}{"Email exists!"})
				return nil
			}

			return err
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(struct {
			Email string `json:"email"`
		}{email})
		context.mailer.SendWelcomeMessage(email)
		return nil
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{"Email missing!"})
	return nil
}
