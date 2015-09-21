package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type appContext struct {
	session *mgo.Session
	// ... and the rest of our globals.
}

type RegisteredUser struct {
	Email string
}

type appHandler struct {
	context *appContext
	handle  func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.handle(ah.context, w, r)
	// handle errors
	if err != nil {
		// todo: implement better error handling
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			log.Printf("HTTP %d: %q", status, err)
		default:
			log.Printf("HTTP %d: %q", status, err)
		}
	}
}

func IndexHandler(context *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "GET" {
		// todo: err is never used
		return http.StatusNotFound, errors.New("Bad method.")
	}

	var results []RegisteredUser

	session := context.session.Copy()
	defer session.Close()

	people := session.DB("registration").C("people")

	err := people.Find(nil).All(&results)
	if err != nil {
		return 500, err
	}

	response, err := json.Marshal(results)
	if err != nil {
		return 500, err
	}

	fmt.Fprintf(w, "{\"people\": %q}", response)
	return 200, nil
}

func RegisterHandler(context *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "POST" {
		// todo: err is never used
		return http.StatusNotFound, errors.New("Bad method.")
	}

	email := r.FormValue("email")

	if email != "" {
		session := context.session.Copy()
		defer session.Close()

		people := session.DB("registration").C("people")

		err := people.Insert(&RegisteredUser{email})
		if err != nil {
			if mgo.IsDup(err) {
				fmt.Fprintf(w, "{\"success\": \"false\", \"error\": \"Email %q exists!\"}", email)
				return 400, err
			}

			return 500, err
		}

		fmt.Fprintf(w, "{\"success\": \"true\", \"email\": \"%q\"}", email)
		return 200, nil
	}

	fmt.Fprintf(w, "{\"success\": \"false\", \"error\": \"Email missing!\"}")
	return 400, nil
}

func buildContext() *appContext {
	// Setup db connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	// Setup context
	context := &appContext{session: session}

	err = session.DB("registration").C("people").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	if err != nil {
		panic(err)
	}

	return context
}

func main() {
	context := buildContext()
	defer context.session.Close()

	// Instatiate handler
	indexHandler := &appHandler{context, IndexHandler}
	registerHandler := &appHandler{context, RegisterHandler}

	// Setup mux
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler)
	mux.Handle("/register", registerHandler)

	fmt.Println("Starting sever on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
