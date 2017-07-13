package payment

import "net/http"

type FakeRequester struct{}

func (r *FakeRequester) ProviderID() string { return "fake" }

func (r *FakeRequester) Pay(request *PayRequest) (*PayResponse, error) {
	return &PayResponse{PaymentKey: "fakeKey", Success: true, Status: "fakeStatus"}, nil
}
func (r *FakeRequester) VerifyPayment(request *http.Request) (bool, error) {
	return true, nil
}
