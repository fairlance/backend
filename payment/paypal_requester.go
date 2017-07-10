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
	payEndpoint = "payouts"
)

type PayPalRequester struct {
	Options *Options
}

func (p *PayPalRequester) ProviderID() string { return "paypal" }

func (p *PayPalRequester) Pay(request PayoutRequest) (*PayResponse, error) {
	req, err := p.newRequest(request, payEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	response := &PayoutResponse{}
	status, err := p.do(req, response)
	if err != nil {
		return nil, err
	}
	return &PayResponse{
		PaymentKey: response.PayoutBatchID,
		Success:    status == http.StatusCreated,
		Status:     response.BatchStatus,
	}, err
}

func (p *PayPalRequester) newRequest(request interface{}, apiEndpoint string) (*http.Request, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %v", err)
	}
	url := fmt.Sprintf("%s/%s", p.Options.PaymentURL, apiEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("X-PAYPAL-SECURITY-USERID", p.Options.SecurityUserID)
	req.Header.Add("X-PAYPAL-SECURITY-PASSWORD", p.Options.SecurityPassword)
	req.Header.Add("X-PAYPAL-SECURITY-SIGNATURE", p.Options.SecuritySignature)
	req.Header.Add("X-PAYPAL-APPLICATION-ID", p.Options.ApplicationID)
	req.Header.Add("X-PAYPAL-REQUEST-DATA-FORMAT", "JSON")
	req.Header.Add("X-PAYPAL-RESPONSE-DATA-FORMAT", "JSON")
	return req, nil
}

func (p *PayPalRequester) do(req *http.Request, response interface{}) (int, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("could not execute request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("could not read response: %v", err)
	}
	if err := json.Unmarshal(body, response); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("could not unmarshal response: %v", err)
	}
	return resp.StatusCode, nil
}
