package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
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
		respond.With(w, r, http.StatusUnauthorized, errors.New("Provide email and password."))
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	user, userType, err := appContext.UserRepository.CheckCredentials(email, password)
	if err != nil {
		respond.With(w, r, http.StatusUnauthorized, err)
		return
	}

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["user"] = user
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
		Type   string `json:"type"`
	}{
		UserId: user.Model.ID,
		Token:  tokenString,
		Type:   userType,
	})
}
