package application

import (
	"encoding/json"

	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"

	"github.com/gorilla/mux"
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
	// Set some claims
	token.Claims["user"] = freelancer.getRepresentationMap()
	token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
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

func Index(w http.ResponseWriter, r *http.Request) {
	respond.With(w, r, http.StatusOK, "Hi")
}

func IndexFreelancer(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancers, err := appContext.FreelancerRepository.GetAllFreelancers()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancers)
}

func NewFreelancer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var body map[string]string
	if err := decoder.Decode(&body); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	freelancer := Freelancer{
		FirstName: body["firstName"],
		LastName:  body["lastName"],
		Password:  body["password"],
		Email:     body["email"],
		Title:     body["title"],
	}

	if ok, err := govalidator.ValidateStruct(freelancer); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		respond.With(w, r, http.StatusBadRequest, errs)
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.AddFreelancer(&freelancer); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancer)
}
func GetFreelancer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	freelancer := Freelancer{}

	if vars["id"] == "" {
		respond.With(w, r, http.StatusBadRequest, errors.New("Id not provided."))
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancer, err = appContext.FreelancerRepository.GetFreelancer(id)
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancer)
}

func DeleteFreelancer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		respond.With(w, r, http.StatusBadRequest, errors.New("Id not provided."))
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.DeleteFreelancer(uint(id)); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}

func IndexProject(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	projects, err := appContext.ProjectRepository.GetAllProjects()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, projects)
}

func IndexClient(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	clients, err := appContext.ClientRepository.GetAllClients()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, clients)
}

//
//func NewFreelancerReference(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//
//	if vars["id"] == "" {
//		WriteError(w, http.StatusBadRequest, errors.New("Id not provided."))
//		return
//	}
//
//	id, err := strconv.Atoi(vars["id"])
//	if err != nil {
//		WriteError(w, http.StatusBadRequest, err)
//		return
//	}
//
//	decoder := json.NewDecoder(r.Body)
//	defer r.Body.Close()
//
//	var reference Reference
//	if err := decoder.Decode(&reference); err != nil {
//		WriteError(w, http.StatusBadRequest, err)
//		return
//	}
//
//	if ok, err := govalidator.ValidateStruct(reference); ok == false || err != nil {
//		errs := govalidator.ErrorsByField(err)
//		w.WriteHeader(http.StatusBadRequest)
//		json.NewEncoder(w).Encode(errs)
//		return
//	}
//
//	var appContext = context.Get(r, "context").(*ApplicationContext)
//	if err := appContext.FreelancerRepository.addReference(id, reference); err != nil {
//		WriteError(w, http.StatusBadRequest, err)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}
