package application

import (
    "net/http"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{"Index", "GET", "/", Index},
    Route{"IndexFreelancer", "GET", "/freelancer/", IndexFreelancer},
    Route{"NewFreelancer", "PUT", "/freelancer/new", NewFreelancer},
    Route{"GetFreelancer", "GET", "/freelancer/{id}", GetFreelancer},
    Route{"DeleteFreelancer", "DELETE", "/freelancer/{id}", DeleteFreelancer},
}
