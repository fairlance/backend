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

const (
	statusDeposited            = "deposited"
	statusInitiated            = "initiated"
	statusAwaitingConfirmation = "awaiting_confirmation"
	statusConfirmed            = "confirmed"
	statusError                = "error"
)

func newPayment(requester Requester, db DB, applicationDispatcher dispatcher.ApplicationDispatcher) *payment {
	return &payment{
		requester:             requester,
		receiversPercentile:   0.92,
		applicationDispatcher: applicationDispatcher,
		db: db,
	}
}

type payment struct {
	requester                       Requester
	receiversPercentile             float64
	paymentProviderChargePercentile float64
	paymentProviderChargeFixed      float64
	applicationDispatcher           dispatcher.ApplicationDispatcher
	db                              DB
}

func (p *payment) depositHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deposit, err := newDepositFromRequest(r)
		if err != nil {
			log.Printf("could not parse deposit request: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not parse deposit request (%+v): %v", deposit, err))
			return
		}
		project, err := p.getProject(deposit.projectID)
		if err != nil {
			log.Printf("could not get project %d: %v", deposit.projectID, err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not get project %d: %v", deposit.projectID, err))
			return
		}
		transactionReceivers := p.buildTransactionReceivers(project)
		t := &Transaction{
			TrackID:   deposit.trackID,
			Provider:  p.requester.ProviderID(),
			ProjectID: deposit.projectID,
			Amount:    fmt.Sprintf("%.2f", project.amount()),
			Status:    statusDeposited,
			Receivers: transactionReceivers,
		}
		if err := p.db.Insert(t); err != nil {
			log.Printf("could not save transaction: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not save transaction: %v", err))
			return
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}

func (p *payment) executeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execute, err := newExecuteFromRequest(r)
		if err != nil {
			log.Printf("could not parse execute request: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not parse execute request (%+v): %v", execute, err))
			return
		}
		t, err := p.db.Get(execute.projectID)
		if err != nil {
			log.Printf("could not get transaction: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not get transaction (%+v): %v", execute, err))
			return
		}
		t.Status = statusAwaitingConfirmation
		if err := p.db.Update(t); err != nil {
			log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s, status: %v", t.TrackID, err))
			return
		}
		var receivers []PayoutItem
		for _, r := range t.Receivers {
			receivers = append(receivers, PayoutItem{
				RecipientType: "EMAIL",
				Amount: PayoutItemAmount{
					Value:    r.Amount,
					Currency: "EUR",
				},
				Note:         fmt.Sprintf("Project %d", t.ProjectID),
				SenderItemID: fmt.Sprintf("%s_%d", t.TrackID, r.FairlanceID),
				Receiver:     r.Email,
			})
		}
		response, err := p.requester.Pay(PayoutRequest{
			SenderBatchHeader: PayoutSenderBatchHeader{
				SenderBatchID: t.TrackID,
				RecipientType: "EMAIL",
				EmailSubject:  fmt.Sprintf("Payment for project %d!", t.ProjectID),
			},
			Items: receivers,
		})
		if err != nil {
			log.Printf("could not execute a pay request: %v", err)
			respond.With(w, r, http.StatusFailedDependency, fmt.Errorf("could not execute a pay request: %v", err))
			return
		}
		t.Status = statusInitiated
		t.PaymentKey = response.PaymentKey
		t.ProviderStatus = response.Status
		if err := p.db.Update(t); err != nil {
			log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s: %v", t.TrackID, err))
			return
		}
		if !response.Success {
			t.Status = statusError
			t.ErrorMsg = "Bad response from PayPal on execute."
			if err := p.db.Update(t); err != nil {
				log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
				respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s: %v", t.TrackID, err))
				return
			}
			respond.With(w, r, http.StatusFailedDependency, response)
			return
		}
		respond.With(w, r, http.StatusOK, response)
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
		projectID := uint(0)
		// https://developer.paypal.com/docs/classic/ipn/ht_ipn/
		t, err := p.db.Get(projectID)
		if err != nil {
			log.Printf("could not get transaction: %v", err)
			return
		}
		t.Status = statusConfirmed
		if err := p.db.Update(t); err != nil {
			log.Printf("could not update transaction for project id %d, status: %v", projectID, err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction for project id %d: %v", projectID, err))
			return
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}

func (p *payment) buildTransactionReceivers(proj *Project) []TransactionReceiver {
	freelancerAmount := proj.amount() * p.receiversPercentile / float64(len(proj.Freelancers))
	var transactionReceivers []TransactionReceiver
	for _, freelancer := range proj.Freelancers {
		transactionReceivers = append(transactionReceivers, TransactionReceiver{
			FairlanceID: freelancer.ID,
			Email:       freelancer.Email,
			Amount:      fmt.Sprintf("%.2f", money(freelancerAmount)),
		})
	}
	return transactionReceivers
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
