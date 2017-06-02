package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fairlance/backend/middleware"
	"github.com/fairlance/backend/payment"
)

var (
	port                int
	secret              string
	authorizationURL    string
	adaptivePaymentsURL string
	returnURL           string
	cancelURL           string
	applicationID       string
	securityUserID      string
	securityPassword    string
	securitySignature   string
	primaryEmail        string
)

func init() {
	f, err := os.OpenFile("/var/log/fairlance/payment.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	flag.IntVar(&port, "port", 3008, "Port.")
	flag.StringVar(&secret, "secret", "", "Secret.")
	flag.StringVar(&authorizationURL, "authorizationURL", "", "authorizationURL")
	flag.StringVar(&adaptivePaymentsURL, "adaptivePaymentsURL", "", "adaptivePaymentsURL")
	flag.StringVar(&returnURL, "returnURL", "", "returnURL")
	flag.StringVar(&cancelURL, "cancelURL", "", "cancelURL")
	flag.StringVar(&applicationID, "applicationID", "", "applicationID")
	flag.StringVar(&securityUserID, "securityUserID", "", "securityUserID")
	flag.StringVar(&securityPassword, "securityPassword", "", "securityPassword")
	flag.StringVar(&securitySignature, "securitySignature", "", "securitySignature")
	flag.StringVar(&primaryEmail, "primaryEmail", "", "primaryEmail")
	flag.Parse()

	payment := payment.New(&payment.Options{
		AuthorizationURL:    authorizationURL,
		AdaptivePaymentsURL: adaptivePaymentsURL,
		ReturnURL:           returnURL,
		CancelURL:           cancelURL,
		ApplicationID:       applicationID,
		SecurityUserID:      securityUserID,
		SecurityPassword:    securityPassword,
		SecuritySignature:   securitySignature,
		PrimaryEmail:        primaryEmail,
	})

	http.Handle("/deposit", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.CORSHandler,
		middleware.JSONEnvelope,
	)(payment.PayPrimaryHandler()))
	http.Handle("/check", middleware.JSONEnvelope(payment.PaymentDetailsHandler()))
	http.Handle("/finalize", middleware.JSONEnvelope(payment.ExecutePaymentHandler()))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
