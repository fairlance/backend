package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type appContext struct {
	session *mgo.Session
	// ... and the rest of our globals.
}

func buildContext() *appContext {
	// Setup context
	context := &appContext{session: getMongoDBSession()}

	return context
}

func main() {
	context := buildContext()
	defer context.session.Close()

	// Instantiate handler
	indexHandler := &appHandler{context, IndexHandler}
	registerHandler := &appHandler{context, RegisterHandler}

	// Setup mux
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler)
	mux.Handle("/register", registerHandler)

	fmt.Println("Starting sever on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
