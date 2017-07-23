package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/application"
	"github.com/fairlance/backend/mailer"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	log.SetFlags(log.Lshortfile)
	var port = os.Getenv("PORT")
	options := application.ContextOptions{
		DbHost:          os.Getenv("DB_HOST"),
		DbName:          os.Getenv("DB_NAME"),
		DbUser:          os.Getenv("DB_USER"),
		DbPass:          os.Getenv("DB_PASS"),
		Secret:          os.Getenv("SECRET"),
		NotificationURL: os.Getenv("NOTIFICATION_URL"),
		MessagingURL:    os.Getenv("MESSAGING_URL"),
		PaymentURL:      os.Getenv("PAYMENT_URL"),
		SearcherURL:     os.Getenv("SEARCHER_URL"),
		MailerOptions: mailer.Options{
			PublicApiKey: os.Getenv("MAILGUN_PUBLIC_API_KEY"),
			ApiKey:       os.Getenv("MAILGUN_API_KEY"),
			Domain:       os.Getenv("MAILGUN_DOMAIN"),
			Self:         os.Getenv("EMAIL_TO_SELF"),
		},
	}
	var appContext, err = application.NewContext(options)
	if err != nil {
		log.Fatal(err)
	}
	appContext.CreateTables()
	http.Handle("/", application.NewRouter(appContext))
	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
