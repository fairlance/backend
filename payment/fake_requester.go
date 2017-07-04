package payment

type fakeRequester struct{}

func (r *fakeRequester) providerID() string { return "fake" }

// func (r *fakeRequester) payPrimary(receivers []Receiver) (*payResponse, error) {
// 	return &payResponse{paymentKey: "fakeKey", success: true}, nil
// }

// func (r *fakeRequester) executePayment(paymentKey string) (*executeResponse, error) {
// 	return &executeResponse{success: true}, nil
// }

func (r *fakeRequester) pay(receivers []Receiver) (*payResponse, error) {
	return &payResponse{paymentKey: "fakeKey", success: true}, nil
}
