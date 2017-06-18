package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fairlance/backend/dispatcher"
	respond "gopkg.in/matryer/respond.v1"
)

type contextKey string

func newPayment(options *Options, db db) *payment {
	return &payment{
		primaryEmail:          options.PrimaryEmail,
		authorizationURL:      options.AuthorizationURL,
		requester:             &payPalRequester{options},
		receiversPercentile:   0.92,
		applicationDispatcher: dispatcher.NewApplicationDispatcher(options.ApplicationURL),
		db: db,
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
	db                              db
}

func (p *payment) payPrimaryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deposit, err := newDepositFromRequest(r)
		if err != nil {
			log.Printf("could not parse request: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not parse request: %v", err))
			return
		}
		receivers, err := p.buildReceivers(deposit)
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
		var paymentReceivers paymentReceivers
		for _, r := range receivers {
			paymentReceivers = append(paymentReceivers, paymentReceiver{email: r.Email, amount: r.Amount})
		}
		if err := p.db.insert(transaction{
			trackID:     deposit.trackID,
			provider:    "paypal",
			providerKey: response.PayKey,
			projectID:   deposit.project,
			amount:      fmt.Sprintf("%.2f", deposit.amount),
			status:      response.ResponseEnvelope.Ack, // todo: not good enough
			receivers:   paymentReceivers,              // todo: this is not saved to db
		}); err != nil {
			log.Printf("could not save transaction: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not save transaction: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, struct {
				RedirectURL string
				TrackID     string
				Response    *PayResponse
			}{
				RedirectURL: fmt.Sprintf("%s%s", p.authorizationURL, response.PayKey),
				TrackID:     deposit.trackID,
				Response:    response,
			})
			return
		}
		respond.With(w, r, http.StatusFailedDependency, response)
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
		respond.With(w, r, http.StatusFailedDependency, response)
	})
}

func (p *payment) executePaymentHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execute, err := newExecuteFromRequest(r)
		if err != nil {
			log.Printf("could not parse request: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not parse request: %v", err))
			return
		}
		t, err := p.db.get(execute.TrackID)
		if err != nil {
			log.Printf("could not find transaction: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not find transaction: %v", err))
			return
		}
		response, err := p.requester.executePayment(t.providerKey)
		if err != nil {
			log.Printf("could not execute a executePayment request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a executePayment request: %v", err))
			return
		}
		t.status = response.ResponseEnvelope.Ack // todo: not good enough
		if err := p.db.update(t); err != nil {
			log.Printf("could not update transaction: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, response)
			return
		}
		respond.With(w, r, http.StatusFailedDependency, response)
	})
}

// https://developer.paypal.com/docs/classic/ipn/integration-guide/IPNIntro/#id08CKFJ00JYK
func (p *payment) ipnNotificationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("could not receive ipn notification: %v", err)
			return
		}
		defer r.Body.Close()
		log.Printf("IPN notificcation received: %v", string(body))
		// https://developer.paypal.com/docs/classic/ipn/ht_ipn/
	})
}

func (p *payment) buildReceivers(d deposit) ([]Receiver, error) {
	project, err := p.getProject(d.project)
	if err != nil {
		return nil, err
	}
	receivers := []Receiver{
		Receiver{
			Amount:  fmt.Sprintf("%.2f", money(d.amount)),
			Email:   p.primaryEmail,
			Primary: true,
		},
	}
	freelancerAmount := d.amount * p.receiversPercentile / float64(len(project.Freelancers))
	for _, freelancer := range project.Freelancers {
		receivers = append(receivers, Receiver{
			Email:  freelancer.Email,
			Amount: fmt.Sprintf("%.2f", money(freelancerAmount)),
		})
	}
	return receivers, nil
}

func (p *payment) getProject(projectID uint) (*Project, error) {
	projectBytes, err := p.applicationDispatcher.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(projectBytes, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func money(amt float64) float64 {
	var intAmt = int64(amt * 100)
	return float64(intAmt) / 100
}
