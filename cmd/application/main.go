package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/application"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var port = os.Getenv("PORT")
	var dbHost = os.Getenv("DB_HOST")
	var dbName = os.Getenv("DB_NAME")
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var secret = os.Getenv("SECRET")
	var notificationURL = os.Getenv("NOTIFICATION_URL")
	var messagingURL = os.Getenv("MESSAGING_URL")
	var searcherURL = os.Getenv("SEARCHER_URL")
	options := application.ContextOptions{
		DbHost:          dbHost,
		DbName:          dbName,
		DbUser:          dbUser,
		DbPass:          dbPass,
		Secret:          secret,
		NotificationURL: notificationURL,
		MessagingURL:    messagingURL,
		SearcherURL:     searcherURL,
	}
	var appContext, err = application.NewContext(options)
	if err != nil {
		log.Fatal(err)
	}
	appContext.DropCreateFillTables()
	http.Handle("/", application.NewRouter(appContext))
	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
