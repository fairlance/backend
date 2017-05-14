package payment

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	respond "gopkg.in/matryer/respond.v1"
)

func New(options *Options) *Payment {
	return &Payment{
		primaryEmail:        options.PrimaryEmail,
		authorizationURL:    options.AuthorizationURL,
		requester:           &payPalRequester{options},
		receiversPercentile: 0.92,
	}
}

type Payment struct {
	primaryEmail                    string
	authorizationURL                string
	requester                       *payPalRequester
	receiversPercentile             float64
	paymentProviderChargePercentile float64
	paymentProviderChargeFixed      float64
}

func (p *Payment) PayPrimaryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivers, err := p.buildReceivers(r)
		if err != nil {
			log.Printf("could not build receivers: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not build receivers: %v", err))
			return
		}
		response, err := p.requester.payPrimary(receivers)
		if err != nil {
			log.Printf("could not execute a payPrimary request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a payPrimary request: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, struct {
				RedirectURL string
				Response    *PayResponse
			}{
				RedirectURL: fmt.Sprintf("%s%s", p.authorizationURL, response.PayKey),
				Response:    response,
			})
			return
		}
		respond.With(w, r, http.StatusInternalServerError, response)
	})
}

func (p *Payment) PaymentDetailsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("payKey") == "" { // should be project, and pay key is stored localy
			respond.With(w, r, http.StatusBadRequest, "payKey missing")
			return
		}
		payKey := r.URL.Query().Get("payKey")
		response, err := p.requester.paymentDetails(payKey)
		if err != nil {
			log.Printf("could not execute a paymentDetails request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a paymentDetails request: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, response)
			return
		}
		respond.With(w, r, http.StatusInternalServerError, response)
	})
}

func (p *Payment) ExecutePaymentHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("payKey") == "" { // should be project, and pay key is stored localy
			respond.With(w, r, http.StatusBadRequest, "payKey missing")
			return
		}
		payKey := r.URL.Query().Get("payKey")
		response, err := p.requester.executePayment(payKey)
		if err != nil {
			log.Printf("could not execute a executePayment request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a executePayment request: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, response)
			return
		}
		respond.With(w, r, http.StatusInternalServerError, response)
	})
}

func (p *Payment) buildReceivers(r *http.Request) ([]Receiver, error) {
	amountParam := r.URL.Query().Get("amount")
	if amountParam == "" || !strings.HasSuffix(amountParam, ".00") || len(amountParam) > 8 {
		return []Receiver{}, fmt.Errorf("amount invalid")
	}
	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		return []Receiver{}, err
	}

	receivers := []Receiver{
		Receiver{
			Email:  r.URL.Query().Get("email"),
			Amount: fmt.Sprintf("%.2f", money(amount*p.receiversPercentile)),
		},
		Receiver{
			Amount:  fmt.Sprintf("%.2f", money(amount)),
			Email:   p.primaryEmail,
			Primary: true,
		},
	}
	return receivers, nil
}

func money(amt float64) float64 {
	var intAmt int64 = int64(amt * 100)
	return float64(intAmt) / 100
}
