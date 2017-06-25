package payment

type PayRequest struct {
	ActionType         string          `json:"actionType"`
	CurrencyCode       string          `json:"currencyCode"`
	ReceiverList       ReceiverList    `json:"receiverList"`
	ReturnURL          string          `json:"returnUrl"`
	CancelURL          string          `json:"cancelUrl"`
	RequestEnvelope    RequestEnvelope `json:"requestEnvelope"`
	FeesPayer          string          `json:"feesPayer"`
	IPNNotificationURL string          `json:"ipnNotificationUrl,omitempty"`
}

type RequestEnvelope struct {
	ErrorLanguage string `json:"errorLanguage"`
	DetailLevel   string `json:"detailLevel"`
}

type ReceiverList struct {
	Receiver []Receiver `json:"receiver"`
}

type Receiver struct {
	Amount      string `json:"amount"`
	Email       string `json:"email"`
	Primary     bool   `json:"primary"`
	fairlanceID uint   // fairlance specific ID, not sent to PAYPAL
}

type PayResponse struct {
	ResponseEnvelope  ResponseEnvelope
	Error             []Error
	PayKey            string
	PaymentExecStatus string
}

type ResponseEnvelope struct {
	Timestamp     string
	Ack           string
	CorrelationID string
	Build         string
}

type Error struct {
	ErrorID   string
	Domain    string
	Subdomain string
	Severity  string
	Category  string
	Message   string
}

type PaymentDetailsRequest struct {
	PayKey          string          `json:"payKey"`
	RequestEnvelope RequestEnvelope `json:"requestEnvelope"`
}

type PaymentInfoList struct {
	PaymentInfo []PaymentInfo
}

type PaymentInfoReceiver struct {
	Amount string
	Email  string
}

type PaymentInfo struct {
	PendingReason           string
	PendingRefund           string
	Receiver                PaymentInfoReceiver
	RefundedAmount          string
	SenderTransactionId     string
	SenderTransactionStatus string
	TransactionId           string
	TransactionStatus       string
}

type PaymentDetailsResponse struct {
	ResponseEnvelope ResponseEnvelope
	Status           string
	Error            []Error
	PaymentInfoList  PaymentInfoList
}

type ExecutePaymentRequest struct {
	PayKey          string          `json:"payKey"`
	RequestEnvelope RequestEnvelope `json:"requestEnvelope"`
}

type ExecutePaymentResponse struct {
	ResponseEnvelope ResponseEnvelope
	Error            []Error
}
