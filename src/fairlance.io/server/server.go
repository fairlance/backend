package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type RegisteredUser struct {
	Name  string
	Email string
}

func Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")

	if name != "" && email != "" {
		session, err := mgo.Dial("localhost")
		if err != nil {
			fmt.Printf("Can't connect to mongo, go error %v\n", err)
			panic(err)
		}
		defer session.Close()
		c := session.DB("registration").C("people")
		err = c.Insert(&RegisteredUser{name, email})
		if err != nil {
			fmt.Printf("Can't save to mongo, go error %v\n", err)
			panic(err)
		} else {
			fmt.Fprintf(w, "Registered user %q with email %q", name, email)
		}
	}

}

func Index(w http.ResponseWriter, r *http.Request) {
	var results []RegisteredUser

	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer session.Close()
	c := session.DB("registration").C("people")
	err = c.Find(nil).All(&results)
	if err != nil {
		fmt.Printf("Can't query mongo, go error %v\n", err)
		panic(err)
	}
	fmt.Fprintf(w, "Results All: %q", results)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/register", Register).Methods("POST")

	fmt.Println("Starting sever on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
