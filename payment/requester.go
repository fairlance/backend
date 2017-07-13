package payment

import "io"

type Requester interface {
	ProviderID() string
	Pay(request *PayRequest) (*PayResponse, error)
	VerifyPayment(reader io.Reader) (bool, error)
}

type PayRequest struct {
	ProjectID uint
	Receivers []PayRequestReceiver
}

type PayRequestReceiver struct {
	Amount string
	Email  string
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
