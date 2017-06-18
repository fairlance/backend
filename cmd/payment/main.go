package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/payment"
	_ "github.com/lib/pq"
)

var (
	port                = os.Getenv("PORT")
	secret              = os.Getenv("SECRET")
	authorizationURL    = os.Getenv("AUTHORIZATION_URL")
	adaptivePaymentsURL = os.Getenv("ADAPTIVE_PAYMENTS_URL")
	returnURL           = os.Getenv("RETURN_URL")
	cancelURL           = os.Getenv("CANCEL_URL")
	applicationID       = os.Getenv("APPLICATION_ID")
	securityUserID      = os.Getenv("SECURITY_USER_ID")
	securityPassword    = os.Getenv("SECURITY_PASSWORD")
	securitySignature   = os.Getenv("SECURITY_SIGNATURE")
	primaryEmail        = os.Getenv("PRIMARY_EMAIL")
	ipnNotificationURL  = os.Getenv("IPN_NOTIFICATION_URL")
	applicationURL      = os.Getenv("APPLICATION_URL")
	dbHost              = os.Getenv("DB_HOST")
	dbUser              = os.Getenv("DB_USER")
	dbPass              = os.Getenv("DB_PASS")
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=payment sslmode=disable", dbHost, dbUser, dbPass))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := payment.NewServeMux(&payment.Options{
		AuthorizationURL:    authorizationURL,
		AdaptivePaymentsURL: adaptivePaymentsURL,
		ReturnURL:           returnURL,
		CancelURL:           cancelURL,
		ApplicationID:       applicationID,
		SecurityUserID:      securityUserID,
		SecurityPassword:    securityPassword,
		SecuritySignature:   securitySignature,
		PrimaryEmail:        primaryEmail,
		ApplicationURL:      applicationURL,
		IPNNotificationURL:  ipnNotificationURL,
	}, db)

	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
