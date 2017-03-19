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
		authHandler(getAllFreelancers()),
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
		authHandler(withID(getFreelancerByID())),
	},
	Route{
		"UpdateFreelancer",
		"POST",
		"/freelancer/{id}",
		authHandler(withID(withFreelancerUpdate(updateFreelancerByID()))),
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
		authHandler(withID(withReference(addFreelancerReferenceByID()))),
	},
	Route{
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		authHandler(withID(withReview(addFreelancerReviewByID()))),
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
		authHandler(getAllClients()),
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
		authHandler(withID(getClientByID())),
	},
	Route{
		"UpdateClient",
		"POST",
		"/client/{id}",
		authHandler(withID(updateClientByID())),
	},

	Route{
		"IndexJob",
		"GET",
		"/job",
		authHandler(getAllJobsForUser()),
	},
	Route{
		"NewJob",
		"PUT",
		"/job/new",
		authHandler(withID(addJob())),
	},
	Route{
		"GetJob",
		"GET",
		"/job/{id}",
		authHandler(withID(getJob())),
	},
	Route{
		"ApplyForJob",
		"PUT",
		"/job/{id}/apply",
		authHandler(withID(withClientFromJobID(withJobApplication(addJobApplicationByID())))),
	},
	Route{
		"DeleteJobApplication",
		"DELETE",
		"/job_application/{id}",
		authHandler(withID(deleteJobApplicationByID())),
	},
}
