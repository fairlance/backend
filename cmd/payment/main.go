package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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
	applicationUrl      string
)

func init() {
	// f, err := os.OpenFile("/var/log/fairlance/payment.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// log.SetOutput(f)
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
	flag.StringVar(&applicationUrl, "applicationUrl", "localhost:3001", "applicationUrl")
	flag.Parse()

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
		ApplicationURL:      applicationUrl,
	})

	log.Printf("Listening on: %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
