package application

var privateRoutes = Routes{
	Route{
		"GetProjectPrivate",
		"GET",
		"/project/{id}",
		withID(getProjectByID()),
	},
}
