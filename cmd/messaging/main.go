package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/messaging"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var port = os.Getenv("PORT")
	var secret = os.Getenv("SECRET")
	var mongoHost = os.Getenv("MONGO_HOST")
	var notificationURL = os.Getenv("NOTIFICATION_URL")
	var applicationURL = os.Getenv("APPLICATION_URL")
	hub := messaging.NewHub(
		messaging.NewMessageDB(mongoHost),
		dispatcher.NewNotifications(notificationURL),
		dispatcher.NewApplication(applicationURL),
	)
	go hub.Run()
	http.Handle("/", messaging.NewRouter(hub, secret))
	log.Printf("Listening on: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
