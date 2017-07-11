package application

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	"gopkg.in/matryer/respond.v1"
)

func getAllFreelancers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appContext = context.Get(r, "context").(*ApplicationContext)
		freelancers, err := appContext.FreelancerRepository.GetAllFreelancers()
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, freelancers)
	})
}

func addFreelancer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user = context.Get(r, "userToAdd").(*User)
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

func getFreelancerByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		freelancer, err := appContext.FreelancerRepository.GetFreelancer(id)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}

		respond.With(w, r, http.StatusOK, freelancer)
	})
}

func deleteFreelancerByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.FreelancerRepository.DeleteFreelancer(id); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}

func withReference(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		context.Set(r, "reference", &reference)

		handler.ServeHTTP(w, r)
	})
}

func addFreelancerReferenceByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var freelancerID = context.Get(r, "id").(uint)
		var reference = context.Get(r, "reference").(*Reference)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.ReferenceRepository.AddReference(freelancerID, reference); err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}

func withReview(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		context.Set(r, "review", &review)

		handler.ServeHTTP(w, r)
	})
}

func addFreelancerReviewByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var freelancerID = context.Get(r, "id").(uint)
		var review = context.Get(r, "review").(*Review)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		if err := appContext.FreelancerRepository.AddReview(freelancerID, review); err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, nil)
	})
}

func withFreelancerUpdate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var freelancerUpdate FreelancerUpdate
		if err := decoder.Decode(&freelancerUpdate); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		if ok, err := govalidator.ValidateStruct(freelancerUpdate); ok == false || err != nil {
			errs := govalidator.ErrorsByField(err)
			respond.With(w, r, http.StatusBadRequest, errs)
			return
		}
		// https://github.com/asaskevich/govalidator/issues/133
		// https://github.com/asaskevich/govalidator/issues/112
		if len(freelancerUpdate.Skills) > 20 {
			respond.With(w, r, http.StatusBadRequest, errors.New("max of 20 skills are allowed"))
			return
		}
		context.Set(r, "freelancerUpdate", &freelancerUpdate)
		handler.ServeHTTP(w, r)
	})
}

func updateFreelancerByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var freelancerID = context.Get(r, "id").(uint)
		var freelancerUpdate = context.Get(r, "freelancerUpdate").(*FreelancerUpdate)
		var appContext = context.Get(r, "context").(*ApplicationContext)
		freelancer, err := appContext.FreelancerRepository.GetFreelancer(freelancerID)
		if err != nil {
			respond.With(w, r, http.StatusNotFound, err)
			return
		}
		freelancer.Image = freelancerUpdate.Image
		freelancer.About = freelancerUpdate.About
		freelancer.Timezone = freelancerUpdate.Timezone
		freelancer.PayPalEmail = freelancerUpdate.PayPalEmail
		freelancer.Phone = freelancerUpdate.Phone
		freelancer.AdditionalFiles = freelancerUpdate.AdditionalFiles
		freelancer.Skills = freelancerUpdate.Skills
		freelancer.PortfolioItems = freelancerUpdate.PortfolioItems
		freelancer.PortfolioLinks = freelancerUpdate.PortfolioLinks
		freelancer.Birthdate = freelancerUpdate.Birthdate
		freelancer.ProfileCompleted = true
		if err := appContext.FreelancerRepository.UpdateFreelancer(&freelancer); err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}
