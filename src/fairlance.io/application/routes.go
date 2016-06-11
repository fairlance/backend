package main

import (
	"net/http"
)

type Route struct {
	Name           string
	Method         string
	Pattern        string
	Handler        http.Handler
	AllowedMethods []string
}

type Routes []Route

var routes = Routes{
	Route{
		"Login",
		"POST",
		"/login",
		http.HandlerFunc(Login),
		[]string{"OPTIONS", "POST"},
	},

	Route{
		"IndexFreelancer",
		"GET",
		"/freelancer/",
		http.HandlerFunc(IndexFreelancer),
		[]string{"GET"},
	},
	Route{
		"RegisterFreelancer",
		"POST",
		"/freelancer/new",
		RegisterFreelancerHandler(http.HandlerFunc(AddFreelancer)),
		[]string{"OPTIONS", "POST"},
	},
	Route{
		"GetFreelancer",
		"GET",
		"/freelancer/{id}",
		IdHandler(http.HandlerFunc(GetFreelancer)),
		[]string{"GET"},
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		AuthHandler(IdHandler(http.HandlerFunc(DeleteFreelancer))),
		[]string{"OPTIONS", "DELETE"},
	},
	Route{
		"AddFreelancerReference",
		"POST",
		"/freelancer/{id}/reference",
		IdHandler(FreelancerReferenceHandler(http.HandlerFunc(AddFreelancerReference))),
		[]string{"OPTIONS", "POST"},
	},
	Route{
		"AddFreelancerReview",
		"POST",
		"/freelancer/{id}/review",
		IdHandler(FreelancerReviewHandler(http.HandlerFunc(AddFreelancerReview))),
		[]string{"OPTIONS", "POST"},
	},

	Route{
		"IndexProject",
		"GET",
		"/project/",
		http.HandlerFunc(IndexProject),
		[]string{"GET"},
	},

	Route{
		"IndexClient",
		"GET",
		"/client/",
		http.HandlerFunc(IndexClient),
		[]string{"GET"},
	},

	Route{
		"IndexJob",
		"GET",
		"/job/",
		http.HandlerFunc(IndexJob),
		[]string{"GET"},
	},
	Route{
		"Info",
		"GET",
		"/info/",
		http.HandlerFunc(Info),
		[]string{"GET"},
	},
}
