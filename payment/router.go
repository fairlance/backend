package payment

import (
	"net/http"

	"github.com/fairlance/backend/middleware"
)

// NewServeMux creates an http.ServeMux with all the routes configured and handeled
func NewServeMux(options *Options) *http.ServeMux {
	payment := newPayment(options)
	mux := http.NewServeMux()
	mux.Handle("/public/deposit", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.CORSHandler,
		middleware.JSONEnvelope,
		middleware.HTTPMethod(http.MethodPost),
		// middleware.WithTokenFromHeader,
		// middleware.AuthenticateTokenWithClaims(options.Secret),
		// middleware.WhenUserType("client"),
	)(payment.payPrimaryHandler()))
	// temp endpoint fo testing
	mux.Handle("/public/check", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.JSONEnvelope,
		middleware.HTTPMethod(http.MethodGet),
	)(payment.paymentDetailsHandler()))
	// private
	mux.Handle("/public/finalize", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.JSONEnvelope,
		middleware.HTTPMethod(http.MethodGet), // POST
	)(payment.executePaymentHandler()))
	mux.Handle("/public/notification", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodGet),
	)(payment.ipnNotificationHandler()))

	return mux
}
