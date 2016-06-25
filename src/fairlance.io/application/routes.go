package main

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
		"POST",
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
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		AuthHandler(IdHandler(http.HandlerFunc(DeleteFreelancer))),
	},
	Route{
		"AddFreelancerReference",
		"POST",
		"/freelancer/{id}/reference",
		IdHandler(FreelancerReferenceHandler(http.HandlerFunc(AddFreelancerReference))),
	},
	Route{
		"AddFreelancerReview",
		"POST",
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
		"POST",
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
		"IndexJob",
		"GET",
		"/job/",
		http.HandlerFunc(IndexJob),
	},
	Route{
		"Info",
		"GET",
		"/info/",
		http.HandlerFunc(Info),
	},
}
