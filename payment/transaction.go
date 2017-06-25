package payment

import (
	"time"
)

type transaction struct {
	// transactions table
	id         uint
	trackID    string
	provider   string
	amount     string
	paymentKey string //payKey
	projectID  uint
	status     string
	createdAt  *time.Time
	updatedAt  *time.Time
	// receivers table
	receivers []paymentReceiver
}
