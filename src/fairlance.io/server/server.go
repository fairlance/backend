package main

import (
	"log"
	"net/http"
)

type appContext struct {
	userRepository *UserRepository
	// ... and the rest of our globals.
}

func buildContext(db string) *appContext {
	// Setup context
	context := &appContext{userRepository: NewUserRepository(db)}

	return context
}

func main() {
	context := buildContext("registration")

	// Instantiate handler
	indexHandler := &appHandler{context, IndexHandler}
	registerHandler := &appHandler{context, RegisterHandler}

	// Setup mux
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler)
	mux.Handle("/register", registerHandler)

	log.Fatal(http.ListenAndServe(":3000", mux))
}
