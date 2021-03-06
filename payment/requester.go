package payment

type Requester interface {
	ProviderID() string
	Pay(request *PayRequest) (*PayResponse, error)
}

type PayRequest struct {
	TrackID     string
	ProjectID   uint
	ProjectName string
	Receivers   []PayRequestReceiver
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
