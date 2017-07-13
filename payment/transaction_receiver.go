package payment

type TransactionReceiver struct {
	ID                 uint
	FairlanceID        uint
	ProviderIdentifier string
	Amount             string
	Status             string
}
