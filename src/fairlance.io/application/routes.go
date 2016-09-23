package application

import (
	"net/http"
)

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
		"/freelancer/",
		http.HandlerFunc(IndexFreelancer),
	},
	Route{
		"RegisterFreelancer",
		"PUT",
		"/freelancer/new",
		RegisterUserHandler(http.HandlerFunc(AddFreelancer)),
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
		WithID{AddFreelancerUpdatesByID},
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
		WithID{AddFreelancerReferenceByID},
	},
	Route{
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		WithID{AddFreelancerReviewByID},
	},

	Route{
		"IndexProject",
		"GET",
		"/project/",
		http.HandlerFunc(IndexProject),
	},

	Route{
		"IndexClient",
		"GET",
		"/client/",
		http.HandlerFunc(IndexClient),
	},
	Route{
		"RegisterClient",
		"PUT",
		"/client/new",
		RegisterUserHandler(http.HandlerFunc(AddClient)),
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
		"/job/",
		http.HandlerFunc(IndexJob),
	},
	Route{
		"NewJob",
		"PUT",
		"/job/new",
		NewJobHandler(http.HandlerFunc(AddJob)),
	},
	Route{
		"GetJob",
		"GET",
		"/job/{id}",
		WithID{GetJobByID},
	},

	Route{
		"Info",
		"GET",
		"/info/",
		http.HandlerFunc(Info),
	},
}
