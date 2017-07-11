package payment_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"encoding/json"

	"github.com/fairlance/backend/payment"
)

type dbMock struct {
	getCall struct {
		receives struct{ projectID uint }
		returns  struct {
			transaction *payment.Transaction
			err         error
		}
	}
	insertCall struct {
		receives struct{ transaction *payment.Transaction }
		returns  struct{ err error }
	}
	updateCall struct {
		receives struct{ transactions []payment.Transaction }
		returns  struct{ err error }
	}
}

func (db *dbMock) Init() {}
func (db *dbMock) Insert(t *payment.Transaction) error {
	db.insertCall.receives.transaction = t
	return db.insertCall.returns.err
}
func (db *dbMock) Update(t *payment.Transaction) error {
	db.updateCall.receives.transactions = append(db.updateCall.receives.transactions, *t)
	return db.updateCall.returns.err
}
func (db *dbMock) Get(projectID uint) (*payment.Transaction, error) {
	db.getCall.receives.projectID = projectID
	return db.getCall.returns.transaction, db.getCall.returns.err
}

type applicationDispatcherMock struct {
	getProjectCall struct {
		receives struct{ id uint }
		returns  struct {
			project payment.Project
			err     error
		}
	}
}

func (d *applicationDispatcherMock) GetProject(id uint) ([]byte, error) {
	d.getProjectCall.receives.id = id
	b, _ := json.Marshal(d.getProjectCall.returns.project)
	return b, d.getProjectCall.returns.err
}
func (d *applicationDispatcherMock) SetProjectFunded(id uint) error { return nil }

func TestDepositHandler(t *testing.T) {
	db := &dbMock{}
	applicationDispatcher := &applicationDispatcherMock{}
	applicationDispatcher.getProjectCall.returns.project = payment.Project{
		ID: 1,
		Freelancers: []payment.Freelancer{
			{
				ID:    1,
				Email: "freelancer@email.com",
			},
		},
		Contract: payment.Contract{
			Hours:   2,
			PerHour: 8,
		},
	}
	router := payment.NewServeMux(&payment.FakeRequester{}, db, applicationDispatcher)
	projectID := uint(7)
	respRec := httptest.NewRecorder()
	body := fmt.Sprintf(`{ "projectID": %d }`, projectID)
	req, err := http.NewRequest(http.MethodPost, "/private/deposit", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatal("Creating '/private/deposit' request failed!")
	}
	router.ServeHTTP(respRec, req)
	// rb, err := ioutil.ReadAll(respRec.Body)
	// log.Printf("%s %v", rb, err)
	if respRec.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusBadRequest)
	}
	transaction := db.insertCall.receives.transaction
	if transaction.ProjectID != projectID {
		t.Fatal("Error: Returned transaction ProjectID", transaction.ProjectID, "instead of", projectID)
	}
	if transaction.Status != "deposited" {
		t.Fatal("Error: Returned transaction Status", transaction.Status, "instead of", "deposited")
	}
	if transaction.Amount != "16.00" {
		t.Fatal("Error: Returned transaction Amount", transaction.Amount, "instead of", "16.00")
	}
	if len(transaction.Receivers) != 1 {
		t.Fatal("Error: Returned transaction len(transaction.Receivers)", len(transaction.Receivers), "instead of", 1)
	}
	if transaction.Receivers[0].Email != "freelancer@email.com" {
		t.Fatal("Error: Returned transaction transaction.Receivers[0].Email", transaction.Receivers[0].Email, "instead of", "freelancer@email.com")
	}
	if transaction.Receivers[0].Amount != "14.72" {
		t.Fatal("Error: Returned transaction transaction.Receivers[0].Amount", transaction.Receivers[0].Amount, "instead of", "14.72")
	}
}

func TestExecuteHandler(t *testing.T) {
	db := &dbMock{}
	db.getCall.returns.transaction = &payment.Transaction{
		ID:        1,
		ProjectID: uint(7),
		TrackID:   "trackID",
		Receivers: []payment.TransactionReceiver{
			{
				ID:          1,
				Email:       "receiver@mail.com",
				Amount:      "14.72",
				FairlanceID: uint(1),
			},
		},
	}
	applicationDispatcher := &applicationDispatcherMock{}
	router := payment.NewServeMux(&payment.FakeRequester{}, db, applicationDispatcher)
	projectID := uint(7)
	respRec := httptest.NewRecorder()
	body := fmt.Sprintf(`{ "projectID": %d }`, projectID)
	req, err := http.NewRequest(http.MethodPost, "/private/execute", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatal("Creating '/private/execute' request failed!")
	}
	router.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusBadRequest)
	}
	if len(db.updateCall.receives.transactions) != 2 {
		t.Fatal("Error: Update called", len(db.updateCall.receives.transactions), "times instead of", 2)
	}
	firstTransactionUpdate := db.updateCall.receives.transactions[0]
	secondTransactionUpdate := db.updateCall.receives.transactions[1]
	if firstTransactionUpdate.Status != "awaiting_confirmation" {
		t.Fatal("Error: Status is", firstTransactionUpdate.Status, "instead of", "awaiting_confirmation")
	}
	if firstTransactionUpdate.PaymentKey != "" {
		t.Fatal("Error: PaymentKey is", firstTransactionUpdate.PaymentKey, "instead of", "")
	}
	if firstTransactionUpdate.ProviderStatus != "" {
		t.Fatal("Error: ProviderStatus is", firstTransactionUpdate.ProviderStatus, "instead of", "")
	}
	if firstTransactionUpdate.Provider != "fake" {
		t.Fatal("Error: Returned transaction Provider", firstTransactionUpdate.Provider, "instead of", "fake")
	}
	if secondTransactionUpdate.Status != "initiated" {
		t.Fatal("Error: Status is", secondTransactionUpdate.Status, "instead of", "initiated")
	}
	if secondTransactionUpdate.PaymentKey != "fakeKey" {
		t.Fatal("Error: PaymentKey is", secondTransactionUpdate.PaymentKey, "instead of", "fakeKey")
	}
	if secondTransactionUpdate.ProviderStatus != "fakeStatus" {
		t.Fatal("Error: ProviderStatus is", secondTransactionUpdate.ProviderStatus, "instead of", "fakeStatus")
	}
}
