package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
)

type RegisteredUser struct {
	Name  string
	Email string
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	fmt.Println("Starting sever on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", router))

}
