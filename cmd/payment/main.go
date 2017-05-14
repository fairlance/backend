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
	port   int
	secret string
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
	flag.Parse()

	payPalRequester := payment.NewPayPalRequester(&payment.Options{})

	http.Handle("/deposit", middleware.JSONEnvelope(payPalRequester.PayPrimaryHandler()))
	http.Handle("/check", middleware.JSONEnvelope(payPalRequester.PaymentDetailsHandler()))
	http.Handle("/finalize", middleware.JSONEnvelope(payPalRequester.ExecutePaymentHandler()))
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
