package payment

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fairlance/backend/dispatcher"

	"encoding/json"

	respond "gopkg.in/matryer/respond.v1"
)

type contextKey string

func newPayment(options *Options) *payment {
	return &payment{
		primaryEmail:          options.PrimaryEmail,
		authorizationURL:      options.AuthorizationURL,
		requester:             &payPalRequester{options},
		receiversPercentile:   0.92,
		applicationDispatcher: dispatcher.NewApplicationDispatcher(options.ApplicationURL),
	}
}

type payment struct {
	primaryEmail                    string
	authorizationURL                string
	requester                       *payPalRequester
	receiversPercentile             float64
	paymentProviderChargePercentile float64
	paymentProviderChargeFixed      float64
	applicationDispatcher           dispatcher.ApplicationDispatcher
}

func (p *payment) payPrimaryHandler() http.Handler {
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

func (p *payment) paymentDetailsHandler() http.Handler {
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

func (p *payment) executePaymentHandler() http.Handler {
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

func (p *payment) buildReceivers(r *http.Request) ([]Receiver, error) {
	amountParam := r.URL.Query().Get("amount")
	if amountParam == "" || !strings.HasSuffix(amountParam, ".00") || len(amountParam) > 8 {
		return nil, fmt.Errorf("amount wrong format: %s", amountParam)
	}
	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		return nil, err
	}
	projectID, err := strconv.Atoi(r.URL.Query().Get("project"))
	if err != nil {
		return nil, err
	}
	projectBytes, err := p.applicationDispatcher.GetProject(uint(projectID))
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(projectBytes, &project); err != nil {
		if err != nil {
			return nil, err
		}
	}
	receivers := []Receiver{
		Receiver{
			Amount:  fmt.Sprintf("%.2f", money(amount)),
			Email:   p.primaryEmail,
			Primary: true,
		},
	}
	for _, freelancer := range project.Freelancers {
		receivers = append(receivers, Receiver{
			Email:  freelancer.Email,
			Amount: fmt.Sprintf("%.2f", money(amount*p.receiversPercentile)),
		})
	}
	return receivers, nil
}

func money(amt float64) float64 {
	var intAmt int64 = int64(amt * 100)
	return float64(intAmt) / 100
}
