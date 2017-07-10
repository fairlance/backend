package payment

type FakeRequester struct{}

func (r *FakeRequester) ProviderID() string { return "fake" }

func (r *FakeRequester) Pay(request PayoutRequest) (*PayResponse, error) {
	return &PayResponse{PaymentKey: "fakeKey", Success: true, Status: "fakeStatus"}, nil
}
