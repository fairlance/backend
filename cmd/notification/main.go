package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/notification"
)

// Examples:
// {"to":[{"type": "freelancer", "id": 1}],"from":{"type": "freelancer", "id": 1},"type":"notification","data":{"text":"hahahah", "projectId": 2}}
// {"type":"read", "from":{"type": "freelancer", "id": 1}, "to":[{"type": "freelancer", "id": 1}], "data": {"timestamp":1487717547735}}

func main() {
	log.SetFlags(log.Lshortfile)
	var port = os.Getenv("PORT")
	var secret = os.Getenv("SECRET")
	var mongoHost = os.Getenv("MONGO_HOST")
	hub := notification.NewHub(mongoHost)
	router := notification.NewRouter(hub, secret)
	http.Handle("/", router)
	go hub.Run()
	log.Printf("Listening on: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
