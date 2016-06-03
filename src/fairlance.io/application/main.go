package main

import (
	"net/http"
)

func main() {
	var appContext, err = NewContext("application")
	appContext.DropCreateFillTables()
	if err != nil {
		panic(err)
	}

	router := NewRouter(appContext)
	http.Handle("/", router)
	if err := http.ListenAndServe(":3001", nil); err != nil {
		panic(err)
	}
}
