package payment

type TransactionReceiver struct {
	ID                     uint
	FairlanceID            uint
	ProviderIdentifier     string
	Amount                 string
	ProviderStatus         string
	ProviderTransactionKey string
}
