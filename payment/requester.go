package payment

type Requester interface {
	ProviderID() string
	Pay(request PayoutRequest) (*PayResponse, error)
}

type PayResponse struct {
	Success    bool
	PaymentKey string
	Status     string
}

type ExecuteResponse struct {
	Success bool
}
