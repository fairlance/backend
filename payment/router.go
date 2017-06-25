package payment

import (
	"net/http"

	"database/sql"

	"github.com/fairlance/backend/middleware"
)

// NewServeMux creates an http.ServeMux with all the routes configured and handeled
func NewServeMux(options *Options, db *sql.DB) *http.ServeMux {
	paymentDB := newDB(db)
	paymentDB.init()
	payment := newPayment(options, paymentDB)
	mux := http.NewServeMux()
	mux.Handle("/private/deposit", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(payment.payPrimaryHandler()))
	mux.Handle("/private/execute", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(payment.executePaymentHandler()))
	mux.Handle("/public/notification", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodGet),
	)(payment.ipnNotificationHandler()))
	return mux
}
