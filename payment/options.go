package payment

type Options struct {
	AdaptivePaymentsURL string
	AuthorizationURL    string
	ReturnURL           string
	CancelURL           string
	SecurityUserID      string
	SecurityPassword    string
	SecuritySignature   string
	ApplicationID       string
	PrimaryEmail        string
	ApplicationURL      string
}
