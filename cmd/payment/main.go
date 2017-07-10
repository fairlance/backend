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
	paymentURL         = os.Getenv("PAYMENT_URL")
	applicationID      = os.Getenv("APPLICATION_ID")
	securityUserID     = os.Getenv("SECURITY_USER_ID")
	securityPassword   = os.Getenv("SECURITY_PASSWORD")
	securitySignature  = os.Getenv("SECURITY_SIGNATURE")
	ipnNotificationURL = os.Getenv("IPN_NOTIFICATION_URL")
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
		PaymentURL:         paymentURL,
		ApplicationID:      applicationID,
		SecurityUserID:     securityUserID,
		SecurityPassword:   securityPassword,
		SecuritySignature:  securitySignature,
		IPNNotificationURL: ipnNotificationURL,
	}
	paymentDB := payment.NewDB(db)
	paymentDB.Init()
	mux := payment.NewServeMux(
		&payment.PayPalRequester{Options: options},
		paymentDB,
		dispatcher.NewApplicationDispatcher(applicationURL),
	)

	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
