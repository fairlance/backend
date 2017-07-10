package payment

type Options struct {
	PaymentURL         string
	SecurityUserID     string
	SecurityPassword   string
	SecuritySignature  string
	ApplicationID      string
	IPNNotificationURL string
}
