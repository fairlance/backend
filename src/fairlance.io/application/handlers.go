package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var body map[string]string
	if err := decoder.Decode(&body); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}
	email := body["email"]
	password := body["password"]

	if email == "" || password == "" {
		respond.With(w, r, http.StatusUnauthorized, errors.New("Provide username and password."))
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	err := appContext.FreelancerRepository.CheckCredentials(email, password)
	if err != nil {
		respond.With(w, r, http.StatusUnauthorized, err)
		return
	}

	freelancer, err := appContext.FreelancerRepository.GetFreelancerByEmail(email)
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["user"] = freelancer.getRepresentationMap()
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(appContext.JwtSecret))
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, struct {
		UserId uint   `json:"id"`
		Token  string `json:"token"`
	}{
		UserId: freelancer.ID,
		Token:  tokenString,
	})
}

func Info(w http.ResponseWriter, r *http.Request) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		respond.With(w, r, http.StatusNotFound, err)
	}
	info, err := ioutil.ReadFile(dir + "/application_info.txt")
	if err != nil {
		respond.With(w, r, http.StatusNotFound, "No info file found!")
		return
	}
	respond.With(w, r, http.StatusOK, string(info))
}
