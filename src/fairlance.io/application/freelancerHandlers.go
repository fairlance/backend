package main

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func IndexFreelancer(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancers, err := appContext.FreelancerRepository.GetAllFreelancers()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancers)
}

func AddFreelancer(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*User)
	freelancer := &Freelancer{User: *user}
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.AddFreelancer(freelancer); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancer)
}

func GetFreelancer(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	var id = context.Get(r, "id").(uint)
	freelancer, err := appContext.FreelancerRepository.GetFreelancer(id)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancer)
}

func DeleteFreelancer(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	var id = context.Get(r, "id").(uint)
	if err := appContext.FreelancerRepository.DeleteFreelancer(id); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}

func AddFreelancerReference(w http.ResponseWriter, r *http.Request) {
	var reference = context.Get(r, "reference").(*Reference)
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.ReferenceRepository.AddReference(reference); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}

func AddFreelancerReview(w http.ResponseWriter, r *http.Request) {
	var review = context.Get(r, "review").(*Review)
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.AddReview(review); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}

func FreelancerReviewHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var review Review
		if err := decoder.Decode(&review); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if id != review.FreelancerId {
			respond.With(w, r, http.StatusBadRequest, "Freelancer id must match the id in the body!")
			return
		}

		if ok, err := govalidator.ValidateStruct(review); ok == false || err != nil {
			errs := govalidator.ErrorsByField(err)
			respond.With(w, r, http.StatusBadRequest, errs)
			return
		}

		context.Set(r, "review", &review)
		next.ServeHTTP(w, r)
	})
}

func FreelancerReferenceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var reference Reference
		if err := decoder.Decode(&reference); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		if id != reference.FreelancerId {
			respond.With(w, r, http.StatusBadRequest, "Freelancer id must match the id in the body!")
			return
		}

		if ok, err := govalidator.ValidateStruct(reference); ok == false || err != nil {
			errs := govalidator.ErrorsByField(err)
			respond.With(w, r, http.StatusBadRequest, errs)
			return
		}

		context.Set(r, "reference", &reference)
		next.ServeHTTP(w, r)
	})
}
