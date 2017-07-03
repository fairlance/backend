package application

import "github.com/fairlance/backend/dispatcher"

type PaymentDispatcher struct {
	payment dispatcher.Payment
}

func NewPaymentDispatcher(payment dispatcher.Payment) *PaymentDispatcher {
	return &PaymentDispatcher{payment}
}

func (d *PaymentDispatcher) execute(projectID uint) error {
	return d.payment.Execute(projectID)
}
