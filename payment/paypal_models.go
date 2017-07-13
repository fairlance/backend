package payment

import "time"
import "encoding/json"

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

type PayPalPayoutRequest struct {
	SenderBatchHeader PayPalPayoutSenderBatchHeader `json:"sender_batch_header"`
	Items             PayPalPayoutItems             `json:"items"`
}

type PayPalPayoutSenderBatchHeader struct {
	// SenderBatchID string `json:"sender_batch_id"`
	EmailSubject  string `json:"email_subject"`
	RecipientType string `json:"recipient_type"`
}

type PayPalPayoutItems []PayPalPayoutItem

type PayPalPayoutItem struct {
	RecipientType string                 `json:"recipient_type"`
	Amount        PayPalPayoutItemAmount `json:"amount"`
	Note          string                 `json:"note"`
	SenderItemID  string                 `json:"sender_item_id"`
	Receiver      string                 `json:"receiver"`
}

type PayPalPayoutItemAmount struct {
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

type PayPalPayoutResponse struct {
	BatchHeader PayPalBatchHeader `json:"batch_header"`
}

type PayPalBatchHeader struct {
	SenderBatchHeader PayPalPayoutSenderBatchHeader `json:"sender_batch_header"`
	PayoutBatchID     string                        `json:"payout_batch_id"`
	BatchStatus       string                        `json:"batch_status"`
}

// {
// 	"name": "INSUFFICIENT_FUNDS",
// 	"message": "An internal service error has occurred.",
// 	"debug_id":"60adcac84df3",
// 	"information_link":"https://developer.paypal.com/docs/api/payments.payouts-batch/#errors"
// }

type PayPalPayoutErrorResponse struct {
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

type PayPalAuthTokenResponse struct {
	Scope       string `json:"scope"`
	Nonce       string `json:"nonce"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
}

type PayPalPaymentPayoutsBaseNotification struct {
	ID           string    `json:"id"`
	CreateTime   time.Time `json:"create_time"`
	ResourceType string    `json:"resource_type"`
	EventType    string    `json:"event_type"`
	Summary      string    `json:"summary"`
	Resource     json.RawMessage
}

// {
// 	"id": "WH-68L35871MU875914X-8DM79009TB591331B",
// 	"create_time": "2016-06-06T03:44:51Z",
// 	"resource_type": "payouts_item",
// 	"event_type": "PAYMENT.PAYOUTS-ITEM.SUCCEEDED",
// 	"summary": "A payout item has succeeded",
// 	"resource": {
// 		"transaction_status": "SUCCESS",
// 		"payout_item_fee": {
// 			"currency": "USD",
// 			"value": "0.25"
// 		},
// 		"payout_batch_id": "Q8B9WFS7ZZJ4Q",
// 		"payout_item": {
// 			"recipient_type": "EMAIL",
// 			"amount": {
// 				"currency": "USD",
// 				"value": "1.0"
// 			},
// 			"note": "First payout",
// 			"receiver": "beamdaddy@paypal.com",
// 			"sender_item_id": "Item1"
// 		},
// 		"links": [
// 			{
// 				"href": "https://api.paypal.com/v1/payments/payouts-item/AYNYWNCHBD8KS",
// 				"rel": "self",
// 				"method": "GET"
// 			},
// 			{
// 				"href": "https://api.paypal.com/v1/payments/payouts/Q8B9WFS7ZZJ4Q",
// 				"rel": "batch",
// 				"method": "GET"
// 			}
// 		],
// 		"payout_item_id": "AYNYWNCHBD8KS",
// 		"time_processed": "2016-06-06T03:44:51Z",
// 		"transaction_id": "57J64166G9424913F"
// 	},
// 	"links": [
// 		{
// 			"href": "https://api.paypal.com/v1/notifications/webhooks-events/WH-68L35871MU875914X-8DM79009TB591331B",
// 			"rel": "self",
// 			"method": "GET",
// 			"encType": "application/json"
// 		},
// 		{
// 			"href": "https://api.paypal.com/v1/notifications/webhooks-events/WH-68L35871MU875914X-8DM79009TB591331B/resend",
// 			"rel": "resend",
// 			"method": "POST",
// 			"encType": "application/json"
// 		}
// 	],
// 	"event_version": "1.0"
// }

type PayPalPaymentPayoutsItemNotificationResource struct {
	TransactionStatus string `json:"transaction_status"`
	PayoutItemFee     struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"payout_item_fee"`
	PayoutBatchID string `json:"payout_batch_id"`
	PayoutItem    struct {
		RecipientType string `json:"recipient_type"`
		Amount        struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"amount"`
		Note         string `json:"note"`
		Receiver     string `json:"receiver"`
		SenderItemID string `json:"sender_item_id"`
	} `json:"payout_item"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
	PayoutItemID  string    `json:"payout_item_id"`
	TimeProcessed time.Time `json:"time_processed"`
	TransactionID string    `json:"transaction_id"`
}

// {
// 	"id": "WH-83C777576Y8332450-2L845887S3616745G",
// 	"create_time": "2015-09-07T12:46:46Z",
// 	"resource_type": "payouts",
// 	"event_type": "PAYMENT.PAYOUTSBATCH.SUCCESS",
// 	"summary": "Payouts batch completed successfully.",
// 	"resource": {
// 		"batch_header": {
// 			"payout_batch_id": "CQGA9SFAU8WSN",
// 			"batch_status": "SUCCESS",
// 			"time_created": "2015-09-07T12:46:41Z",
// 			"time_completed": "2015-09-07T12:46:45Z",
// 			"sender_batch_header": {
// 				"sender_batch_id": "REL1"
// 			},
// 			"amount": {
// 				"currency": "CAD",
// 				"value": "25.0"
// 			},
// 			"fees": {
// 				"currency": "CAD",
// 				"value": "0.5"
// 			},
// 			"payments": 1
// 		},
// 		"links": [
// 			{
// 				"href": "https://api.paypal.com/v1/payments/payouts/CQGA9SFAU8WSN",
// 				"rel": "self",
// 				"method": "GET"
// 			}
// 		]
// 	},
// 	"links": [
// 		{
// 			"href": "https://api.paypal.com/v1/notifications/webhooks-events/WH-83C777576Y8332450-2L845887S3616745G",
// 			"rel": "self",
// 			"method": "GET",
// 			"encType": "application/json"
// 		},
// 		{
// 			"href": "https://api.paypal.com/v1/notifications/webhooks-events/WH-83C777576Y8332450-2L845887S3616745G/resend",
// 			"rel": "resend",
// 			"method": "POST",
// 			"encType": "application/json"
// 		}
// 	],
// 	"event_version": "1.0"
// }

type PayPalPaymentPayoutsBatchNotificationResource struct {
	BatchHeader struct {
		PayoutBatchID     string    `json:"payout_batch_id"`
		BatchStatus       string    `json:"batch_status"`
		TimeCreated       time.Time `json:"time_created"`
		TimeCompleted     time.Time `json:"time_completed"`
		SenderBatchHeader struct {
			SenderBatchID string `json:"sender_batch_id"`
		} `json:"sender_batch_header"`
		Amount struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"amount"`
		Fees struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"fees"`
		Payments int `json:"payments"`
	} `json:"batch_header"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}
