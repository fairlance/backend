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

func CreateToken(claims map[string]interface{}, secret string, duration time.Duration) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	jwtClaims := token.Claims.(jwt.MapClaims)
	// Set some claims
	for k, v := range claims {
		jwtClaims[k] = v
	}
	jwtClaims["exp"] = time.Now().Add(duration).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

func login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var body map[string]string
		if err := decoder.Decode(&body); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		email, emailOk := body["email"]
		password, passwordOk := body["password"]

		if !emailOk || !passwordOk {
			respond.With(w, r, http.StatusBadRequest, errors.New("Provide email and password."))
			return
		}

		var appContext = context.Get(r, "context").(*ApplicationContext)
		user, userType, err := appContext.UserRepository.CheckCredentials(email, password)
		if err != nil {
			respond.With(w, r, http.StatusUnauthorized, err)
			return
		}

		claims := make(map[string]interface{})
		claims["user"] = user
		tokenString, err := CreateToken(claims, appContext.JwtSecret, time.Hour*8)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, struct {
			UserID uint   `json:"id"`
			Token  string `json:"token"`
			Type   string `json:"type"`
		}{
			UserID: user.Model.ID,
			Token:  tokenString,
			Type:   userType,
		})
	})
}
