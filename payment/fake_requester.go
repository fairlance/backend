package payment

import "io"

type FakeRequester struct{}

func (r *FakeRequester) ProviderID() string { return "fake" }

func (r *FakeRequester) Pay(request *PayRequest) (*PayResponse, error) {
	return &PayResponse{PaymentKey: "fakeKey", Success: true, Status: "fakeStatus"}, nil
}
func (r *FakeRequester) VerifyPayment(reader io.Reader) (bool, error) {
	return true, nil
}
