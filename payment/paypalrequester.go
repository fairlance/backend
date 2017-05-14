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

func (p *payPalRequester) payPrimary(receivers []Receiver) (*PayResponse, error) {
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
		FeesPayer: "PRIMARYRECEIVER",
	}
	req, err := p.newRequest(payPrimaryRequest, payEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	payResponse := &PayResponse{}
	err = p.do(req, payResponse)
	return payResponse, err
}

func (p *payPalRequester) paymentDetails(payKey string) (*PaymentDetailsResponse, error) {
	paymentDetailRequest := &PaymentDetailsRequest{
		PayKey: payKey,
		RequestEnvelope: RequestEnvelope{
			ErrorLanguage: "en_US",
			DetailLevel:   "ReturnAll",
		},
	}
	req, err := p.newRequest(paymentDetailRequest, paymentDetailsEndpont)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	paymentDetailsResponse := &PaymentDetailsResponse{}
	err = p.do(req, paymentDetailsResponse)
	return paymentDetailsResponse, err
}

func (p *payPalRequester) executePayment(payKey string) (*ExecutePaymentResponse, error) {
	executePaymentRequest := &ExecutePaymentRequest{
		PayKey: payKey,
		RequestEnvelope: RequestEnvelope{
			ErrorLanguage: "en_US",
			DetailLevel:   "ReturnAll",
		},
	}
	req, err := p.newRequest(executePaymentRequest, executePaymentEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	executePaymentResponse := &ExecutePaymentResponse{}
	err = p.do(req, executePaymentResponse)
	return executePaymentResponse, err
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
		Timeout: time.Duration(5 * time.Second),
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
