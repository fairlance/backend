package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	payEndpoint            = "Pay"
	paymentDetailsEndpont  = "PaymentDetails"
	executePaymentEndpoint = "ExecutePayment"
)

type payPalRequester struct {
	options *Options
}

func (p *payPalRequester) providerID() string { return "paypal" }

func (p *payPalRequester) pay(receivers []Receiver) (*payResponse, error) {
	payPrimaryRequest := &PayRequest{
		ActionType:   "PAY_PRIMARY",
		CurrencyCode: "EUR",
		ReceiverList: ReceiverList{
			Receiver: receivers,
		},
		ReturnURL: p.options.ReturnURL,
		CancelURL: p.options.CancelURL,
		RequestEnvelope: RequestEnvelope{
			ErrorLanguage: "en_US",
			DetailLevel:   "ReturnAll",
		},
		FeesPayer:          "PRIMARYRECEIVER",
		IPNNotificationURL: p.options.IPNNotificationURL,
	}
	req, err := p.newRequest(payPrimaryRequest, payEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	response := &PayResponse{}
	err = p.do(req, response)
	if err != nil {
		return nil, err
	}
	return &payResponse{
		paymentKey: response.PayKey,
		success:    response.ResponseEnvelope.Ack == "Success",
	}, err
}

func (p *payPalRequester) newRequest(request interface{}, apiEndpoint string) (*http.Request, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %v", err)
	}
	url := fmt.Sprintf("%s/%s", p.options.AdaptivePaymentsURL, apiEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("X-PAYPAL-SECURITY-USERID", p.options.SecurityUserID)
	req.Header.Add("X-PAYPAL-SECURITY-PASSWORD", p.options.SecurityPassword)
	req.Header.Add("X-PAYPAL-SECURITY-SIGNATURE", p.options.SecuritySignature)
	req.Header.Add("X-PAYPAL-APPLICATION-ID", p.options.ApplicationID)
	req.Header.Add("X-PAYPAL-REQUEST-DATA-FORMAT", "JSON")
	req.Header.Add("X-PAYPAL-RESPONSE-DATA-FORMAT", "JSON")
	return req, nil
}

func (p *payPalRequester) do(req *http.Request, response interface{}) error {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}
	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("could not unmarshal response: %v", err)
	}
	return nil
}

// func (p *payPalRequester) payPrimary(receivers []Receiver) (*payResponse, error) {
// 	payPrimaryRequest := &PayRequest{
// 		ActionType:   "PAY_PRIMARY",
// 		CurrencyCode: "EUR",
// 		ReceiverList: ReceiverList{
// 			Receiver: receivers,
// 		},
// 		ReturnURL: p.options.ReturnURL,
// 		CancelURL: p.options.CancelURL,
// 		RequestEnvelope: RequestEnvelope{
// 			ErrorLanguage: "en_US",
// 			DetailLevel:   "ReturnAll",
// 		},
// 		FeesPayer:          "PRIMARYRECEIVER",
// 		IPNNotificationURL: p.options.IPNNotificationURL,
// 	}
// 	req, err := p.newRequest(payPrimaryRequest, payEndpoint)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not create request: %v", err)
// 	}
// 	response := &PayResponse{}
// 	err = p.do(req, response)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &payResponse{
// 		paymentKey: response.PayKey,
// 		success:    response.ResponseEnvelope.Ack == "Success",
// 	}, err
// }

// func (p *payPalRequester) executePayment(payKey string) (*executeResponse, error) {
// 	executePaymentRequest := &ExecutePaymentRequest{
// 		PayKey: payKey,
// 		RequestEnvelope: RequestEnvelope{
// 			ErrorLanguage: "en_US",
// 			DetailLevel:   "ReturnAll",
// 		},
// 	}
// 	req, err := p.newRequest(executePaymentRequest, executePaymentEndpoint)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not create request: %v", err)
// 	}
// 	response := &ExecutePaymentResponse{}
// 	err = p.do(req, response)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &executeResponse{
// 		success: response.ResponseEnvelope.Ack == "Success",
// 	}, err
// }
