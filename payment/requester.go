package payment

type Requester interface {
	ProviderID() string
	// should be something else other than Transaction, but yeah...
	Pay(t *Transaction) (*PayResponse, error)
}

type PayResponse struct {
	Success      bool
	PaymentKey   string
	Status       string
	ErrorMessage string
}

type ExecuteResponse struct {
	Success bool
}
