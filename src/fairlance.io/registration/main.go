package main

import (
	"net/http"
)

func main() {
	context := NewContext("registration")

	// Instantiate handler
	indexHandler := &AppHandler{context, IndexHandler}
	registerHandler := &AppHandler{context, RegisterHandler}

	// Setup mux
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler)
	mux.Handle("/register", registerHandler)

	panic(http.ListenAndServe(":3000", mux))
}
