package application

import (
    "net/http"
    "encoding/json"
    "time"
    "github.com/asaskevich/govalidator"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hi"))
}

func IndexFreelancer(w http.ResponseWriter, r *http.Request) {

    var appContext = context.Get(r, "context").(*ApplicationContext)
    freelancers, err := appContext.FreelancerRepository.GetAllFreelancers()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.NewEncoder(w).Encode(freelancers)
}

func NewFreelancer(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()

    var body map[string]string
    if err := decoder.Decode(&body); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    freelancer := Freelancer{
        FirstName: body["firstName"],
        LastName: body["lastName"],
        Password: body["password"],
        Email: body["email"],
        Created: time.Now(),
    }

    if ok, err := govalidator.ValidateStruct(freelancer); ok == false || err != nil {
        errs := govalidator.ErrorsByField(err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errs)
        return
    }

    var appContext = context.Get(r, "context").(*ApplicationContext)
    if err := appContext.FreelancerRepository.AddFreelancer(freelancer); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(freelancer)
}

func GetFreelancer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    freelancer := Freelancer{}

    if vars["id"] == "" {
        w.Write([]byte("Id not provided."))
        return
    }

    var appContext = context.Get(r, "context").(*ApplicationContext)
    freelancer, err := appContext.FreelancerRepository.GetFreelancer(vars["id"])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(freelancer)
}

func DeleteFreelancer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    if vars["id"] == "" {
        w.Write([]byte("Id not provided."))
        return
    }

    var appContext = context.Get(r, "context").(*ApplicationContext)
    if err := appContext.FreelancerRepository.DeleteFreelancer(vars["id"]); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.WriteHeader(http.StatusOK)
}