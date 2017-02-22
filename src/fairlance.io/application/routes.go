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
		login(),
	},

	Route{
		"IndexFreelancer",
		"GET",
		"/freelancer",
		getAllFreelancers(),
	},
	Route{
		"RegisterFreelancer",
		"PUT",
		"/freelancer/new",
		withUser(addFreelancer()),
	},
	Route{
		"GetFreelancer",
		"GET",
		"/freelancer/{id}",
		withID(getFreelancerByID()),
	},
	Route{
		"UpdateFreelancer",
		"POST",
		"/freelancer/{id}",
		withID(withFreelancerUpdate(updateFreelancerByID())),
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		authHandler(withID(deleteFreelancerByID())),
	},
	Route{
		"AddFreelancerReference",
		"PUT",
		"/freelancer/{id}/reference",
		withID(withReference(addFreelancerReferenceByID())),
	},
	Route{
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		withID(withReview(addFreelancerReviewByID())),
	},

	Route{
		"IndexProject",
		"GET",
		"/project",
		authHandler(getAllProjectsForUser()),
	},
	Route{
		"GetProject",
		"GET",
		"/project/{id}",
		authHandler(withID(getProjectByID())),
	},

	Route{
		"IndexClient",
		"GET",
		"/client",
		getAllClients(),
	},
	Route{
		"RegisterClient",
		"PUT",
		"/client/new",
		withUser(addClient()),
	},
	Route{
		"GetClient",
		"GET",
		"/client/{id}",
		withID(getClientByID()),
	},
	Route{
		"UpdateClient",
		"POST",
		"/client/{id}",
		withID(updateClientByID()),
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
		withID(withJobApplication(addJobApplicationByID())),
	},
	Route{
		"DeleteJobApplication",
		"DELETE",
		"/job_application/{id}",
		authHandler(withID(deleteJobApplicationByID())),
	},
}
