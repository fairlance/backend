package application

import "github.com/fairlance/backend/middleware"

var privateRoutes = Routes{
	Route{
		"GetProjectPrivate",
		"GET",
		"/project/{id}",
		middleware.Chain(
			middleware.WithUINT("id"),
		)(getProjectByID()),
	},
	Route{
		"ProjectFundedPrivate",
		"GET",
		"/project/{id}/fund",
		middleware.Chain(
			middleware.WithUINT("id"),
			withProjectByID,
		)(projectFunded()),
	},
}
