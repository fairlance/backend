package application

import "github.com/fairlance/backend/middleware"

var privateRoutes = Routes{
	Route{
		"GetProjectPrivate",
		"GET",
		"/project/{id}",
		middleware.Chain(
			withUINT("id"),
		)(getProjectByID()),
	},
	Route{
		"ProjectFundedPrivate",
		"GET",
		"/project/{id}/fund",
		middleware.Chain(
			withUINT("id"),
			withProjectByID,
		)(projectFunded()),
	},
}
