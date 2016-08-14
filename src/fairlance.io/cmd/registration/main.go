package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"fairlance.io/registration"
)

var port int
var user string
var pass string

func main() {
	flag.IntVar(&port, "port", 3000, "Specify the port to listen to.")
	flag.StringVar(&user, "user", "", "Auth user.")
	flag.StringVar(&pass, "pass", "", "Auth password.")
	flag.Parse()

	if user == "" || pass == "" {
		fmt.Println("User or pass empty!")
		return
	}

	context := registration.NewContext("registration")

	// Instantiate handler
	indexHandler := &registration.AppHandler{context, registration.IndexHandler}
	registerHandler := &registration.AppHandler{context, registration.RegisterHandler}

	// Setup mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if registration.Authenticated(w, r, user, pass) {
			indexHandler.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="FAIRLANCE"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized\n"))
	})
	mux.Handle("/register", registerHandler)

	panic(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
