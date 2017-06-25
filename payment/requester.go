package payment

type requester interface {
	providerID() string
	payPrimary(receivers []Receiver) (*payResponse, error)
	executePayment(paymentKey string) (*executeResponse, error)
}

type payResponse struct {
	success    bool
	paymentKey string
	data       map[string]string
}

type executeResponse struct {
	success bool
}
