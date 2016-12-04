package application

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func IndexFreelancer(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancers, err := appContext.FreelancerRepository.GetAllFreelancers()
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, freelancers)
}

func AddFreelancer(user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		freelancer := &Freelancer{User: *user}
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.FreelancerRepository.AddFreelancer(freelancer); err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, struct {
			User User   `json:"user"`
			Type string `json:"type"`
		}{
			User: freelancer.User,
			Type: "freelancer",
		})
	})
}

// GetFreelancerByID handler
func GetFreelancerByID(id uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		freelancer, err := appContext.FreelancerRepository.GetFreelancer(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, freelancer)
	})
}

// DeleteFreelancerByID handler
func DeleteFreelancerByID(id uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.FreelancerRepository.DeleteFreelancer(id); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}

type withReference struct {
	next func(reference *Reference) http.Handler
}

func (wr withReference) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var reference Reference
	if err := decoder.Decode(&reference); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	if ok, err := govalidator.ValidateStruct(reference); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		respond.With(w, r, http.StatusBadRequest, errs)
		return
	}

	wr.next(&reference).ServeHTTP(w, r)
}

type addFreelancerReferenceByID struct {
	freelancerID uint
	reference    *Reference
}

func (afrbi addFreelancerReferenceByID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.ReferenceRepository.AddReference(afrbi.freelancerID, afrbi.reference); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}

type withReview struct {
	next func(review *Review) http.Handler
}

func (wr withReview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var review Review
	if err := decoder.Decode(&review); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	if ok, err := govalidator.ValidateStruct(review); ok == false || err != nil {
		errs := govalidator.ErrorsByField(err)
		respond.With(w, r, http.StatusBadRequest, errs)
		return
	}

	wr.next(&review).ServeHTTP(w, r)
}

type addFreelancerReviewByID struct {
	freelancerID uint
	review       *Review
}

func (afrbi addFreelancerReviewByID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	if err := appContext.FreelancerRepository.AddReview(afrbi.freelancerID, afrbi.review); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}

type withFreelancerUpdate struct {
	next func(freelancerUpdate *FreelancerUpdate) http.Handler
}

func (wfu withFreelancerUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var freelancerUpdate FreelancerUpdate

	if err := decoder.Decode(&freelancerUpdate); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	// https://github.com/asaskevich/govalidator/issues/133
	// https://github.com/asaskevich/govalidator/issues/112
	if len(freelancerUpdate.Skills) > 20 {
		respond.With(w, r, http.StatusBadRequest, errors.New("max of 20 skills are allowed"))
		return
	}

	wfu.next(&freelancerUpdate).ServeHTTP(w, r)
}

type updateFreelancerHandler struct {
	freelancerID     uint
	freelancerUpdate *FreelancerUpdate
}

func (ufh updateFreelancerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var appContext = context.Get(r, "context").(*ApplicationContext)
	freelancer, err := appContext.FreelancerRepository.GetFreelancer(ufh.freelancerID)
	if err != nil {
		respond.With(w, r, http.StatusNotFound, err)
		return
	}

	freelancer.Skills = ufh.freelancerUpdate.Skills
	freelancer.Timezone = ufh.freelancerUpdate.Timezone
	freelancer.IsAvailable = ufh.freelancerUpdate.IsAvailable
	freelancer.HourlyRateFrom = ufh.freelancerUpdate.HourlyRateFrom
	freelancer.HourlyRateTo = ufh.freelancerUpdate.HourlyRateTo

	if err := appContext.FreelancerRepository.UpdateFreelancer(&freelancer); err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	respond.With(w, r, http.StatusOK, nil)
}
