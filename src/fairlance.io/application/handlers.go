package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
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

func RegisterUserHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var body struct {
			FirstName string `json:"firstName" valid:"required"`
			LastName  string `json:"lastName" valid:"required"`
			Password  string `json:"password" valid:"required"`
			Email     string `json:"email" valid:"required,email"`
		}

		if err := decoder.Decode(&body); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if ok, err := govalidator.ValidateStruct(body); ok == false || err != nil {
			errs := govalidator.ErrorsByField(err)
			respond.With(w, r, http.StatusBadRequest, errs)
			return
		}

		user := &User{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Password:  body.Password,
			Email:     body.Email,
		}

		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	})
}
