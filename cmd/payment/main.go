package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/payment"
	_ "github.com/lib/pq"
)

var (
	port               = os.Getenv("PORT")
	secret             = os.Getenv("SECRET")
	payPalPaymentURL   = os.Getenv("PAYPAL_PAYMENT_URL")
	payPalClientID     = os.Getenv("PAYPAL_CLIENT_ID")
	payPalSecret       = os.Getenv("PAYPAL_SECRET")
	payPalOAuth2URL    = os.Getenv("PAYPAL_OAUTH2_URL")
	ipnNotificationURL = os.Getenv("PAYPAL_IPN_NOTIFICATION_URL")
	applicationURL     = os.Getenv("APPLICATION_URL")
	dbHost             = os.Getenv("DB_HOST")
	dbUser             = os.Getenv("DB_USER")
	dbPass             = os.Getenv("DB_PASS")
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=payment sslmode=disable", dbHost, dbUser, dbPass))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	options := &payment.Options{
		PaymentURL:         payPalPaymentURL,
		ClientID:           payPalClientID,
		Secret:             payPalSecret,
		OAuth2URL:          payPalOAuth2URL,
		IPNNotificationURL: ipnNotificationURL,
	}
	paymentDB := payment.NewDB(db)
	paymentDB.Init()
	mux := payment.NewServeMux(
		&payment.PayPalRequester{Options: options},
		paymentDB,
		dispatcher.NewApplication(applicationURL),
	)

	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
