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

	// Route{
	// 	"IndexFreelancer",
	// 	"GET",
	// 	"/freelancer",
	// 	whenLoggedIn(getAllFreelancers()),
	// },
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
		whenLoggedIn(withID(getFreelancerByID())),
	},
	Route{
		"UpdateFreelancer",
		"POST",
		"/freelancer/{id}",
		whenLoggedIn(whenFreelancer(withID(whenIDBelongsToUser(withFreelancerUpdate(updateFreelancerByID()))))),
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		whenLoggedIn(whenFreelancer(withID(whenIDBelongsToUser(deleteFreelancerByID())))),
	},
	Route{
		"AddFreelancerReference",
		"PUT",
		"/freelancer/{id}/reference",
		whenLoggedIn(whenFreelancer(withID(withReference(addFreelancerReferenceByID())))),
	},
	Route{
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		whenLoggedIn(withID(withReview(addFreelancerReviewByID()))),
	},

	Route{
		"IndexProject",
		"GET",
		"/project",
		whenLoggedIn(getAllProjectsForUser()),
	},
	Route{
		"GetProject",
		"GET",
		"/project/{id}",
		whenLoggedIn(withID(whenProjectBelongsToUser(getProjectByID()))),
	},
	Route{
		"CreateProjectFromJobApplication",
		"POST",
		"/project/create_from_job_application/{id}",
		whenLoggedIn(whenClient(withID(whenJobApplicationBelongsToUser(createProjectFromJobApplication())))),
	},
	Route{
		"AddExtension",
		"POST",
		"/project/{id}/extension",
		whenLoggedIn(whenClient(whenProjectBelongsToUser(withID(withExtension(addExtensionToProjectContract()))))),
	},

	// Route{
	// 	"IndexClient",
	// 	"GET",
	// 	"/client",
	// 	whenLoggedIn(getAllClients()),
	// },
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
		whenLoggedIn(withID(getClientByID())),
	},
	Route{
		"UpdateClient",
		"POST",
		"/client/{id}",
		whenLoggedIn(whenClient(withID(whenIDBelongsToUser(updateClientByID())))),
	},

	Route{
		"IndexJob",
		"GET",
		"/job",
		whenLoggedIn(getAllJobsForUser()),
	},
	Route{
		"NewJob",
		"PUT",
		"/job/new",
		whenLoggedIn(whenClient(addJob())),
	},
	Route{
		"GetJob",
		"GET",
		"/job/{id}",
		whenLoggedIn(withID(getJob())),
	},
	Route{
		"ApplyForJob", // todo: prevent freelancers to apply twice
		"PUT",
		"/job/{id}/apply",
		whenLoggedIn(whenFreelancer(withID(withClientFromJobID(withJobApplication(addJobApplication()))))),
	},
	Route{
		"DeleteJobApplication",
		"DELETE",
		"/job_application/{id}",
		whenLoggedIn(whenFreelancer(whenJobApplicationBelongsToUser(withID(deleteJobApplicationByID())))),
	},
}
