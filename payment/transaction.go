package payment

import (
	"database/sql/driver"
	"encoding/json"

	"time"

	"github.com/fairlance/backend/jsonb"
)

type transaction struct {
	ID          string
	trackID     string
	provider    string
	amount      string
	providerKey string //payKey
	projectID   uint
	status      string
	receivers   paymentReceivers // JSONB
	createdAt   *time.Time
	updatedAt   *time.Time
}

type paymentReceiver struct {
	email  string
	amount string
}

type paymentReceivers []paymentReceiver

func (r paymentReceivers) Value() (driver.Value, error) {
	return pValue(r)
}

func (r *paymentReceivers) Scan(src interface{}) error {
	return pScan(r, src)
}

func pValue(entity interface{}) (driver.Value, error) {
	var josnb jsonb.JSONB
	josnb, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	return josnb.Value()
}

func pScan(entity, src interface{}) error {
	var jsonb jsonb.JSONB
	err := (&jsonb).Scan(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, entity)
}
