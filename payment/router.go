package payment

import (
	"net/http"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/middleware"
)

// NewServeMux creates an http.ServeMux with all the routes configured and handeled
func NewServeMux(requester Requester, paymentDB DB, applicationDispatcher dispatcher.Application) *http.ServeMux {
	payment := newPayment(requester, paymentDB, applicationDispatcher)
	mux := http.NewServeMux()
	mux.Handle("/private/deposit", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(payment.depositHandler()))
	mux.Handle("/private/execute", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(payment.executeHandler()))
	mux.Handle("/public/webhook/paypal", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(payment.notificationHandler()))
	return mux
}
