package application

import "net/http"

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

var routes = Routes{
	Route{
		"Login",
		"POST",
		"/login",
		http.HandlerFunc(Login),
	},

	Route{
		"IndexFreelancer",
		"GET",
		"/freelancer",
		http.HandlerFunc(IndexFreelancer),
	},
	Route{
		"RegisterFreelancer",
		"PUT",
		"/freelancer/new",
		WithUser{AddFreelancer},
	},
	Route{
		"GetFreelancer",
		"GET",
		"/freelancer/{id}",
		WithID{GetFreelancerByID},
	},
	Route{
		"UpdateFreelancer",
		"POST",
		"/freelancer/{id}",
		WithID{func(freelancerID uint) http.Handler {
			return withFreelancerUpdate{func(freelancerUpdate *FreelancerUpdate) http.Handler {
				return updateFreelancerHandler{freelancerID, freelancerUpdate}
			}}
		}},
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		AuthHandler(WithID{DeleteFreelancerByID}),
	},
	Route{
		"AddFreelancerReference",
		"PUT",
		"/freelancer/{id}/reference",
		WithID{func(freelancerID uint) http.Handler {
			return withReference{func(reference *Reference) http.Handler {
				return addFreelancerReferenceByID{
					freelancerID: freelancerID,
					reference:    reference,
				}
			}}
		}},
	},
	Route{
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		WithID{func(freelancerID uint) http.Handler {
			return withReview{func(review *Review) http.Handler {
				return addFreelancerReviewByID{
					freelancerID: freelancerID,
					review:       review,
				}
			}}
		}},
	},

	Route{
		"IndexProject",
		"GET",
		"/project",
		http.HandlerFunc(IndexProject),
	},
	Route{
		"GetProject",
		"GET",
		"/project/{id}",
		WithID{GetProjectByID},
	},

	Route{
		"IndexClient",
		"GET",
		"/client",
		http.HandlerFunc(IndexClient),
	},
	Route{
		"RegisterClient",
		"PUT",
		"/client/new",
		WithUser{AddClient},
	},
	Route{
		"GetClient",
		"GET",
		"/client/{id}",
		WithID{GetClientByID},
	},
	Route{
		"UpdateClient",
		"POST",
		"/client/{id}",
		WithID{UpdateClientByID},
	},

	Route{
		"IndexJob",
		"GET",
		"/job",
		getAllJobs(),
	},
	Route{
		"NewJob",
		"PUT",
		"/job/new",
		withID(addJob()),
	},
	Route{
		"GetJob",
		"GET",
		"/job/{id}",
		withID(getJob()),
	},
	Route{
		"ApplyForJob",
		"PUT",
		"/job/{id}/apply",
		withID(withJobApplication(applyForJob())),
	},
}
