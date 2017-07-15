package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/messaging"

	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	port = os.Getenv("PORT")
	secret = os.Getenv("SECRET")
	mongoHost = os.Getenv("MONGO_HOST")
	notificationURL = os.Getenv("NOTIFICATION_URL")
	applicationURL := os.Getenv("APPLICATION_URL")
	hub := messaging.NewHub(
		messaging.NewMessageDB(mongoHost),
		dispatcher.NewNotifier(notificationURL),
		dispatcher.NewApplication(applicationURL),
	)
	go hub.Run()
	http.Handle("/", messaging.NewRouter(hub, secret))
	log.Printf("Listening on: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
