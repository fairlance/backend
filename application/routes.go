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
	// 	"AllFreelancers",
	// 	"GET",
	// 	"/freelancer",
	// 	whenLoggedIn(getAllFreelancers()),
	// },
	Route{
		"RegisterFreelancer",
		"PUT",
		"/freelancer/new",
		middleware.Chain(
			withUserToAdd,
		)(addFreelancer()),
	},
	Route{
		"GetFreelancer",
		"GET",
		"/freelancer/{id}",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
		)(getFreelancerByID()),
	},
	Route{
		"UpdateFreelancer",
		"POST",
		"/freelancer/{id}/complete_profile",
		middleware.Chain(
			whenLoggedIn,
			whenFreelancer,
			withUINT("id"),
			whenIDBelongsToUser,
			withFreelancerUpdateFromRequest,
		)(updateFreelancerByID()),
	},
	Route{
		"DeleteFreelancer",
		"DELETE",
		"/freelancer/{id}",
		middleware.Chain(
			whenLoggedIn,
			whenFreelancer,
			withUINT("id"),
			whenIDBelongsToUser,
		)(deleteFreelancerByID()),
	},
	Route{
		"AddFreelancerReference",
		"PUT",
		"/freelancer/{id}/reference",
		middleware.Chain(
			whenLoggedIn,
			whenFreelancer,
			withUINT("id"),
			withReference,
		)(addFreelancerReferenceByID()),
	},
	Route{ // todo: reviews should be related to a project?
		"AddFreelancerReview",
		"PUT",
		"/freelancer/{id}/review",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
			withReview,
		)(addFreelancerReviewByID()),
	},

	Route{
		"AllProjects",
		"GET",
		"/project",
		middleware.Chain(
			whenLoggedIn,
		)(basedOnUserType(getAllProjectsForClient(), getAllProjectsForFreelancer())),
	},
	Route{
		"GetProject",
		"GET",
		"/project/{id}",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
			whenBasedOnUserType(
				whenProjectBelongsToClientByID,
				whenProjectBelongsToFreelancerByID,
			),
		)(getProjectByID()),
	},
	Route{
		"CreateProjectFromJobApplication",
		"POST",
		"/project/create_from_job_application/{id}",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			withUINT("id"),
			whenJobApplicationBelongsToClientByID,
		)(createProjectFromJobApplication()),
	},
	Route{
		"AddProjectContractProposal",
		"POST",
		"/project/{id}/contract/proposal",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
			whenBasedOnUserType(
				whenProjectBelongsToClientByID,
				whenProjectBelongsToFreelancerByID,
			),
			withProjectByID,
			whenCurrentProjectStatus(projectStatusFinalizingTerms),
			withProposal,
		)(setProposalToProjectContract()),
	},
	Route{
		"AgreeToContractTerms",
		"POST",
		"/project/{id}/contract/agree",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
			whenBasedOnUserType(
				whenProjectBelongsToClientByID,
				whenProjectBelongsToFreelancerByID,
			),
			withProjectByID,
			whenCurrentProjectStatus(projectStatusFinalizingTerms),
		)(agreeToContractTerms()),
	},
	Route{
		"AddExtension",
		"POST",
		"/project/{id}/extension",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
			whenBasedOnUserType(
				whenProjectBelongsToClientByID,
				whenProjectBelongsToFreelancerByID,
			),
			withExtension,
		)(addExtensionToProjectContract()),
	},
	Route{
		"AgreeToExtensionTerms",
		"POST",
		"/project/{id}/extension/{extension_id}/agree",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
			whenBasedOnUserType(
				whenProjectBelongsToClientByID,
				whenProjectBelongsToFreelancerByID,
			),
			withProjectByID,
			whenCurrentProjectStatus(projectStatusInProgress),
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
			withUINT("id"),
			whenBasedOnUserType(
				whenProjectBelongsToClientByID,
				whenProjectBelongsToFreelancerByID,
			),
			withProjectByID,
			whenCurrentProjectStatus(projectStatusInProgress),
			withProposal,
			withUINT("extension_id"),
			withExtensionWhenBelongsToProject,
		)(setProposalToProjectContractExtension()),
	},
	Route{
		"FundProject",
		"POST",
		"/project/{id}/funded",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			withUINT("id"),
			whenProjectBelongsToClientByID,
			withProjectByID,
			whenCurrentProjectStatus(projectStatusPendingFunds),
		)(fundedProject()),
	},
	Route{
		"FinishProjectByFreelancer",
		"POST",
		"/project/{id}/finish",
		middleware.Chain(
			whenLoggedIn,
			whenFreelancer,
			withUINT("id"),
			whenProjectBelongsToFreelancerByID,
			withProjectByID,
			whenCurrentProjectStatus(projectStatusInProgress),
		)(freelancerFinishProject()),
	},
	Route{
		"ProjectDone",
		"POST",
		"/project/{id}/done",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			withUINT("id"),
			whenProjectBelongsToClientByID,
			withProjectByID,
			whenCurrentProjectStatus(projectStatusPendingFinished),
		)(projectDone()),
	},

	// Route{
	// 	"AllClients",
	// 	"GET",
	// 	"/client",
	// 	whenLoggedIn(getAllClients()),
	// },
	Route{
		"RegisterClient",
		"PUT",
		"/client/new",
		middleware.Chain(
			withUserToAdd,
		)(addClient()),
	},
	Route{
		"GetClient",
		"GET",
		"/client/{id}",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
		)(getClientByID()),
	},
	Route{
		"UpdateClient",
		"POST",
		"/client/{id}/complete_profile",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			withUINT("id"),
			whenIDBelongsToUser,
			withClientUpdateFromRequest,
		)(updateClientByID()),
	},

	Route{
		"AllJobs",
		"GET",
		"/job",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
		)(getAllJobsForClient()),
	},
	Route{
		"NewJob",
		"POST",
		"/job/new",
		middleware.Chain(
			whenLoggedIn,
			whenClient,
			whenProfileCompleted,
			withJobFromRequest,
		)(addJob()),
	},
	Route{
		"GetJob",
		"GET",
		"/job/{id}",
		middleware.Chain(
			whenLoggedIn,
			withUINT("id"),
		)(getJob()),
	},
	Route{
		"ApplyForJob", // todo: prevent freelancers to apply twice
		"PUT",
		"/job/{id}/apply",
		middleware.Chain(
			whenLoggedIn,
			whenFreelancer,
			whenProfileCompleted,
			withUINT("id"),
			whenFreelancerHasNotAppliedBeforeByID,
			withClientFromJobID,
			withJobApplicationFromRequest,
		)(addJobApplication()),
	},
	Route{
		"DeleteJobApplication",
		"DELETE",
		"/job_application/{id}",
		middleware.Chain(
			whenLoggedIn,
			whenFreelancer,
			withUINT("id"),
			whenJobApplicationBelongsToFreelancerByID,
		)(deleteJobApplicationByID()),
	},
}
