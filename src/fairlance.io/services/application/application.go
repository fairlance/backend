package main

import (
	"net/http"

	app "fairlance.io/application"
)

func main() {
	var appContext, err = app.NewContext("application")
	appContext.PrepareTables()
	if err != nil {
		panic(err)
	}
	router := app.NewRouter(appContext)
	http.Handle("/", router)
	if err := http.ListenAndServe(":3001", nil); err != nil {
		panic(err)
	}
}
