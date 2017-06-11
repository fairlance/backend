package payment

import (
	"net/http"

	"github.com/fairlance/backend/middleware"
)

func NewServeMux(options *Options) *http.ServeMux {
	payment := newPayment(options)
	mux := http.NewServeMux()
	mux.Handle("/deposit", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.CORSHandler,
		middleware.JSONEnvelope,
		// middleware.WithTokenFromHeader,
		// middleware.AuthenticateTokenWithClaims(options.Secret),
		middleware.WhenUserType("client"),
	)(payment.payPrimaryHandler()))
	mux.Handle("/check", middleware.Chain(middleware.JSONEnvelope)(payment.paymentDetailsHandler()))
	mux.Handle("/finalize", middleware.Chain(middleware.JSONEnvelope)(payment.executePaymentHandler()))

	return mux
}
