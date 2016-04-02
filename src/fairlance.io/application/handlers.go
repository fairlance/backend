package application

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var body map[string]string
	if err := decoder.Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	email := body["email"]
	password := body["password"]

	var appContext = context.Get(r, "context").(*ApplicationContext)
	err := appContext.FreelancerRepository.CheckCredentials(email, password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err)
		return
	}

	freelancer, err := appContext.FreelancerRepository.GetFreelancerByEmail(email)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err)
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
		WriteError(w, http.StatusUnauthorized, err)
		return
	}

	json.NewEncoder(w).Encode(struct {
		UserId int    `json:"id"`
		Token  string `json:"token"`
	}{
		UserId: freelancer.Id,
		Token:  tokenString,
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi"))
}

func IndexFreelancer(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancers, err := appContext.FreelancerRepository.GetAllFreelancers()
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(freelancers)
}

func IndexProject(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	projects, err := appContext.ProjectRepository.GetAllProjects()
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func IndexClient(w http.ResponseWriter, r *http.Request) {

	var appContext = context.Get(r, "context").(*ApplicationContext)
	clients, err := appContext.ClientRepository.GetAllClients()
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(clients)
}

func NewFreelancer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var body map[string]string
	if err := decoder.Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	freelancer := Freelancer{
		FirstName: body["firstName"],
		LastName:  body["lastName"],
		Password:  body["password"],
		Email:     body["email"],
		Created:   time.Now(),
	}

	if ok, err := govalidator.ValidateStruct(freelancer); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.AddFreelancer(freelancer); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(freelancer)
}

func NewFreelancerReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		WriteError(w, http.StatusBadRequest, errors.New("Id not provided."))
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var reference Reference
	if err := decoder.Decode(&reference); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if ok, err := govalidator.ValidateStruct(reference); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.addReference(id, reference); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFreelancer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	freelancer := Freelancer{}

	if vars["id"] == "" {
		WriteError(w, http.StatusBadRequest, errors.New("Id not provided."))
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancer, err = appContext.FreelancerRepository.GetFreelancer(id)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(freelancer)
}

func DeleteFreelancer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		WriteError(w, http.StatusBadRequest, errors.New("Id not provided."))
		return
	}

	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.DeleteFreelancer(vars["id"]); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	body, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		err.Error(),
	})

	w.WriteHeader(status)
	w.Write(body)
}
