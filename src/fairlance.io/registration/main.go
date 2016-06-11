package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
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

	context := NewContext("registration")

	// Instantiate handler
	indexHandler := &AppHandler{context, IndexHandler}
	registerHandler := &AppHandler{context, RegisterHandler}

	// Setup mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if authenticated(w, r, user, pass) {
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
