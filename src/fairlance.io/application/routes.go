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
		IdHandler(http.HandlerFunc(GetFreelancer)),
	},
	Route{
		"UpdateFreelancer",
		"POST",
		"/freelancer/{id}",
		IdHandler(FreelancerUpdateHandler(http.HandlerFunc(AddFreelancerUpdates))),
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		AuthHandler(IdHandler(http.HandlerFunc(DeleteFreelancer))),
	},
	Route{
		"AddFreelancerReference",
		"PUT",
		"/freelancer/{id}/reference",
		IdHandler(FreelancerReferenceHandler(http.HandlerFunc(AddFreelancerReference))),
	},
	Route{
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		IdHandler(FreelancerReviewHandler(http.HandlerFunc(AddFreelancerReview))),
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
		IdHandler(http.HandlerFunc(GetClient)),
	},
	Route{
		"UpdateClient",
		"POST",
		"/client/{id}",
		IdHandler(http.HandlerFunc(UpdateClient)),
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
		IdHandler(http.HandlerFunc(GetJob)),
	},

	Route{
		"Info",
		"GET",
		"/info/",
		http.HandlerFunc(Info),
	},
}
