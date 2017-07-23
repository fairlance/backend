package main

import (
	"log"
	"net/http"
	"os"

	"fmt"

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
	notifications := notification.New(secret, mongoHost)
	http.Handle("/", notifications.Handler())
	http.Handle("/send", notifications.SendHandler())
	go notifications.Run()
	log.Printf("Listening on: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
