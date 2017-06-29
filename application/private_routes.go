package application

import "github.com/fairlance/backend/middleware"

var privateRoutes = Routes{
	Route{
		"GetProjectPrivate",
		"GET",
		"/project/{id}",
		withID(getProjectByID()),
	},
	Route{
		"ProjectFundedPrivate",
		"GET",
		"/project/{id}/fund",
		middleware.Chain(
			withID,
			withProjectByID,
		)(projectFunded()),
	},
}
