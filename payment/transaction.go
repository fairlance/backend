package payment

import (
	"time"
)

type Transaction struct {
	// transactions table
	ID                     uint
	TrackID                string
	Provider               string
	ProviderStatus         string
	ProviderTransactionKey string
	Amount                 string
	ProjectID              uint
	Status                 string
	ErrorMsg               string
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
	// receivers table
	Receivers []TransactionReceiver
}
