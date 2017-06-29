package application

import "net/http"
import "github.com/fairlance/backend/middleware"

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
		whenLoggedIn(whenFreelancer(withID(
			whenIDBelongsToUser(withFreelancerUpdate(
				updateFreelancerByID()))))),
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		whenLoggedIn(whenFreelancer(withID(
			whenIDBelongsToUser(deleteFreelancerByID())))),
	},
	Route{
		"AddFreelancerReference",
		"PUT",
		"/freelancer/{id}/reference",
		whenLoggedIn(whenFreelancer(withID(
			withReference(addFreelancerReferenceByID())))),
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
		middleware.Chain(
			whenLoggedIn,
			withID,
			whenProjectBelongsToUserByID,
		)(getProjectByID()),
	},
	Route{
		"CreateProjectFromJobApplication",
		"POST",
		"/project/create_from_job_application/{id}",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			withID,
			whenJobApplicationBelongsToUser,
		)(createProjectFromJobApplication()),
	},
	Route{
		"AddProjectContractProposal",
		"POST",
		"/project/{id}/contract/proposal",
		middleware.Chain(
			whenLoggedIn,
			withID,
			whenProjectBelongsToUserByID,
			withProjectByID,
			withProposal,
		)(setProposalToProjectContract()),
	},
	Route{
		"AgreeToContractTerms",
		"POST",
		"/project/{id}/contract/agree",
		middleware.Chain(
			whenLoggedIn,
			withID,
			whenProjectBelongsToUserByID,
			withProjectByID,
		)(agreeToContractTerms()),
	},
	Route{
		"AddExtension",
		"POST",
		"/project/{id}/extension",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			withID,
			whenProjectBelongsToUserByID,
			withExtension,
		)(addExtensionToProjectContract()),
	},
	Route{
		"AgreeToExtensionTerms",
		"POST",
		"/project/{id}/extension/{extension_id}/agree",
		middleware.Chain(
			whenLoggedIn,
			whenProjectBelongsToUserByID,
			withProjectByID,
			withUINT("extension_id"),
			withExtensionWhenBelongsToProject,
		)(agreeToExtensionTerms()),
	},
	Route{
		"AddProjectContractExtensionProposal",
		"POST",
		"/project/{id}/extension/{extension_id}/proposal",
		middleware.Chain(
			whenLoggedIn,
			withID,
			whenProjectBelongsToUserByID,
			withProjectByID,
			withProposal,
			withUINT("extension_id"),
			withExtensionWhenBelongsToProject,
		)(setProposalToProjectContractExtension()),
	},
	Route{
		"FinishProject",
		"POST",
		"/project/{id}/finish",
		middleware.Chain(
			whenLoggedIn,
			withID,
			whenProjectBelongsToUserByID,
			withProjectByID,
		)(finishProject()),
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
		"POST",
		"/job/new",
		whenLoggedIn(whenClient(withJob(addJob()))),
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
		whenLoggedIn(whenFreelancer(withID(
			withClientFromJobID(withJobApplication(
				addJobApplication()))))),
	},
	Route{
		"DeleteJobApplication",
		"DELETE",
		"/job_application/{id}",
		whenLoggedIn(whenFreelancer(
			whenJobApplicationBelongsToUser(withID(
				deleteJobApplicationByID())))),
	},
}
