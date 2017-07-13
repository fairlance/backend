package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		transactionReceivers, err := p.buildTransactionReceivers(project)
		if err != nil {
			log.Printf("could not build transaction receivers for project %d: %v", deposit.projectID, err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not build transaction receivers for project %d: %v", deposit.projectID, err))
			return
		}
		t := &Transaction{
			TrackID:   deposit.trackID,
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
		t, err := p.db.GetByProjectID(execute.projectID)
		if err != nil {
			log.Printf("could not get transaction: %v", err)
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not get transaction (%+v): %v", execute, err))
			return
		}
		t.Status = statusAwaitingConfirmation
		t.Provider = p.requester.ProviderID()
		if err := p.db.Update(t); err != nil {
			log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s, status: %v", t.TrackID, err))
			return
		}
		response, err := p.requester.Pay(p.buildPayRequest(t))
		if err != nil {
			t.Status = statusError
			if err := p.db.Update(t); err != nil {
				log.Printf("could not update transaction %s: %v", t.TrackID, err)
				respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s: %v", t.TrackID, err))
				return
			}
			log.Printf("could not execute a pay request: %v", err)
			respond.With(w, r, http.StatusFailedDependency, fmt.Errorf("could not execute a pay request: %v", err))
			return
		}
		if !response.Success {
			t.ProviderStatus = response.Status
			t.ErrorMsg = response.ErrorMessage
			t.Status = statusError
			if err := p.db.Update(t); err != nil {
				log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
				respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s: %v", t.TrackID, err))
				return
			}
			respond.With(w, r, http.StatusFailedDependency, response)
			return
		}
		t.ProviderStatus = response.Status
		t.PaymentKey = response.PaymentKey
		t.Status = statusInitiated
		if err := p.db.Update(t); err != nil {
			log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not update transaction %s: %v", t.TrackID, err))
			return
		}
		respond.With(w, r, http.StatusOK, response)
	})
}

// https://developer.paypal.com/docs/classic/ipn/integration-guide/IPNIntro/#id08CKFJ00JYK
// https://developer.paypal.com/docs/classic/ipn/ht_ipn/
func (p *payment) notificationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("IPN notificcation received")
		var buff bytes.Buffer
		tee := io.TeeReader(r.Body, &buff)
		verified, err := p.requester.VerifyPayment(tee)
		if err != nil {
			log.Printf("could not verify payment: %v", err)
			return
		}
		if !verified {
			log.Println("notification request is not verified")
			return
		}
		var notification PayPalPaymentPayoutsBaseNotification
		if err := json.NewDecoder(&buff).Decode(&notification); err != nil {
			log.Printf("could not decode notification: %v", err)
			return
		}
		defer r.Body.Close()
		log.Printf("IPN notificcation: %+v", notification)
		var batchID string
		switch notification.ResourceType {
		case "payouts":
			var batchResource PayPalPaymentPayoutsBatchNotificationResource
			if err := json.Unmarshal(notification.Resource, &batchResource); err != nil {
				log.Printf("could not decode the resource: %v", err)
				return
			}
			batchID = batchResource.BatchHeader.PayoutBatchID
			t, err := p.db.GetByProviderTransactionKey(batchID)
			if err != nil {
				log.Printf("could not get transaction with provider transaction key: %s: %v", batchID, err)
				return
			}
			t.ProviderStatus = batchResource.BatchHeader.BatchStatus
			if err := p.db.Update(t); err != nil {
				log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
				return
			}
			log.Printf("batch resource: %+v", batchResource)
		case "payouts_item":
			var itemResource PayPalPaymentPayoutsItemNotificationResource
			if err := json.Unmarshal(notification.Resource, &itemResource); err != nil {
				log.Printf("could not decode the resource: %v", err)
				return
			}
			batchID = itemResource.PayoutBatchID
			t, err := p.db.GetByProviderTransactionKey(batchID)
			if err != nil {
				log.Printf("could not get transaction with provider transaction key: %s: %v", batchID, err)
				return
			}
			for _, receiver := range t.Receivers {
				if receiver.ProviderIdentifier == itemResource.PayoutItem.Receiver {
					receiver.Status = itemResource.TransactionStatus
				}
			}
			if err := p.db.Update(t); err != nil {
				log.Printf("could not update transaction %s, status: %v", t.TrackID, err)
				return
			}
			log.Printf("item resource: %+v", itemResource)
		}
		respond.With(w, r, http.StatusOK, nil)
	})
}

func (p *payment) buildTransactionReceivers(proj *Project) ([]TransactionReceiver, error) {
	freelancerAmount := proj.amount() * p.receiversPercentile / float64(len(proj.Freelancers))
	var transactionReceivers []TransactionReceiver
	for _, freelancer := range proj.Freelancers {
		if freelancer.ID == 0 || freelancer.Email == "" {
			return nil, fmt.Errorf("could not build TransactionReceiver, email or id is not provided: id=%d, email=%s", freelancer.ID, freelancer.Email)
		}
		transactionReceivers = append(transactionReceivers, TransactionReceiver{
			FairlanceID:        freelancer.ID,
			ProviderIdentifier: freelancer.Email,
			Amount:             fmt.Sprintf("%.2f", money(freelancerAmount)),
		})
	}
	return transactionReceivers, nil
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

func (p *payment) buildPayRequest(t *Transaction) *PayRequest {
	var receivers []PayRequestReceiver
	for _, r := range t.Receivers {
		receivers = append(receivers, PayRequestReceiver{
			Email:  r.ProviderIdentifier,
			Amount: r.Amount,
		})
	}
	return &PayRequest{
		ProjectID: t.ProjectID,
		Receivers: receivers,
	}
}

func money(amt float64) float64 {
	var intAmt = int64(amt * 100)
	return float64(intAmt) / 100
}
