package main

import (
    "fairlance.io/registration"
    "net/http"
)

func main() {
    context := registration.NewContext("registration")

    // Instantiate handler
    indexHandler := &registration.AppHandler{context, registration.IndexHandler}
    registerHandler := &registration.AppHandler{context, registration.RegisterHandler}

    // Setup mux
    mux := http.NewServeMux()
    mux.Handle("/", indexHandler)
    mux.Handle("/register", registerHandler)

    panic(http.ListenAndServe(":3000", mux))
}
