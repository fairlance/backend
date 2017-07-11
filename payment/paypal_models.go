package payment

// "sender_batch_header":{
//     "sender_batch_id":"2014021801",
//     "email_subject":"You have a payout!",
//     "recipient_type":"EMAIL"
//   },
//   "items":[
//     {
//       "recipient_type":"EMAIL",
//       "amount":{
//         "value":"1.0",
//         "currency":"USD"
//       },
//       "note":"Thanks for your patronage!",
//       "sender_item_id":"201403140001",
//       "receiver":"anybody01@gmail.com"
//     }
//   ]
// }'

type PayoutRequest struct {
	SenderBatchHeader PayoutSenderBatchHeader `json:"sender_batch_header"`
	Items             PayoutItems             `json:"items"`
}

type PayoutSenderBatchHeader struct {
	// SenderBatchID string `json:"sender_batch_id"`
	EmailSubject  string `json:"email_subject"`
	RecipientType string `json:"recipient_type"`
}

type PayoutItems []PayoutItem

type PayoutItem struct {
	RecipientType string           `json:"recipient_type"`
	Amount        PayoutItemAmount `json:"amount"`
	Note          string           `json:"note"`
	SenderItemID  string           `json:"sender_item_id"`
	Receiver      string           `json:"receiver"`
}

type PayoutItemAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// {
//   "batch_header": {
//     "sender_batch_header": {
//       "sender_batch_id": "2014021801",
//       "email_subject": "You have a payout!"
//     },
//     "payout_batch_id": "12345678",
//     "batch_status": "PENDING"
//   }
// }

type PayoutResponse struct {
	BatchHeader BatchHeader `json:"batch_header"`
}

type BatchHeader struct {
	SenderBatchHeader PayoutSenderBatchHeader `json:"sender_batch_header"`
	PayoutBatchID     string                  `json:"payout_batch_id"`
	BatchStatus       string                  `json:"batch_status"`
}

// {
// 	"name": "INSUFFICIENT_FUNDS",
// 	"message": "An internal service error has occurred.",
// 	"debug_id":"60adcac84df3",
// 	"information_link":"https://developer.paypal.com/docs/api/payments.payouts-batch/#errors"
// }

type PayoutErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// {
//   "scope": "https://uri.paypal.com/services/subscriptions https://api.paypal.com/v1/payments/.* https://api.paypal.com/v1/vault/credit-card https://uri.paypal.com/services/applications/webhooks openid https://uri.paypal.com/payments/payouts https://api.paypal.com/v1/vault/credit-card/.*",
//   "nonce": "2017-06-08T18:30:28ZCl54Q_OlDqP6-4D03sDT8wRiHjKrYlb5EH7Di0gRrds",
//   "access_token": "Access-Token",
//   "token_type": "Bearer",
//   "app_id": "APP-80W284485P519543T",
//   "expires_in": 32398
// }

type AuthTokenResponse struct {
	Scope       string `json:"scope"`
	Nonce       string `json:"nonce"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
}
