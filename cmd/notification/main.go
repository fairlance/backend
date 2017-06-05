package main

import (
	"flag"
	"log"
	"net/http"

	"fmt"

	"github.com/fairlance/backend/notification"
)

var port int
var secret string
var mongoHost string

// Examples:
// {"to":[{"type": "freelancer", "id": 1}],"from":{"type": "freelancer", "id": 1},"type":"notification","data":{"text":"hahahah", "projectId": 2}}
// {"type":"read", "from":{"type": "freelancer", "id": 1}, "to":[{"type": "freelancer", "id": 1}], "data": {"timestamp":1487717547735}}

func init() {
	flag.IntVar(&port, "port", 3007, "Specify the port to listen on.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.StringVar(&mongoHost, "mongoHost", "localhost", "Mongo host.")
	flag.Parse()

	// f, err := os.OpenFile("/var/log/fairlance/notification.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// log.SetOutput(f)
}

func main() {
	notifications := notification.New(secret, mongoHost)
	http.Handle("/", notifications.Handler())
	http.Handle("/send", notifications.SendHandler())
	go notifications.Run()

	log.Printf("Listening on: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
