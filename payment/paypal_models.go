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
	SenderBatchID string `json:"sender_batch_id"`
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
	SenderBatchHeader PayoutSenderBatchHeader `json:"sender_batch_header"`
	PayoutBatchID     string                  `json:"payout_batch_id"`
	BatchStatus       string                  `json:"batch_status"`
}
