package application

import (
    "net/http"
)

type Route struct {
    Name    string
    Method  string
    Pattern string
    Handler http.Handler
}

type Routes []Route

var routes = Routes{
    Route{"Index", "GET", "/", http.HandlerFunc(Index)},
    Route{"IndexFreelancer", "GET", "/freelancer/", AuthHandler(http.HandlerFunc(IndexFreelancer))},
    Route{"NewFreelancer", "PUT", "/freelancer/new", AuthHandler(http.HandlerFunc(NewFreelancer))},
    Route{"GetFreelancer", "GET", "/freelancer/{id}", AuthHandler(http.HandlerFunc(GetFreelancer))},
    Route{"DeleteFreelancer", "DELETE", "/freelancer/{id}", AuthHandler(http.HandlerFunc(DeleteFreelancer))},
}
