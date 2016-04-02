package application

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
		"Index",
		"GET",
		"/",
		http.HandlerFunc(Index),
		[]string{"GET"},
	},

	Route{
		"IndexFreelancer",
		"GET",
		"/freelancer/",
		http.HandlerFunc(IndexFreelancer),
		[]string{"GET"},
	},
	Route{
		"NewFreelancer",
		"POST",
		"/freelancer/new",
		http.HandlerFunc(NewFreelancer),
		[]string{"OPTIONS", "POST"},
	},
	Route{
		"GetFreelancer",
		"GET",
		"/freelancer/{id}",
		AuthHandler(http.HandlerFunc(GetFreelancer)),
		[]string{"GET"},
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		AuthHandler(http.HandlerFunc(DeleteFreelancer)),
		[]string{"OPTIONS", "DELETE"},
	},

	Route{
		"NewFreelancerReference",
		"POST",
		"/freelancer/{id}/reference/new",
		http.HandlerFunc(NewFreelancerReference),
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
}
