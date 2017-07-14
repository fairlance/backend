package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	payEndpoint           = "payouts"
	payPalEmailNoteFormat = "Thank you for finishing project: %s"
	payPalEmailSubject    = "Fairlance and PayPal just gave you money!"
)

type PayPalRequester struct {
	Options *Options
}

func (p *PayPalRequester) ProviderID() string { return "paypal" }

func (p *PayPalRequester) Pay(r *PayRequest) (*PayResponse, error) {
	token, err := p.getToken()
	if err != nil {
		return nil, fmt.Errorf("could not get Auth token: %v", err)
	}
	req, err := p.newHTTPRequest(token, payEndpoint, p.buildPayoutRequest(r))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	resp, err := p.do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do the request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		payoutErrorResponse := &PayPalPayoutErrorResponse{}
		if err := json.Unmarshal(body, payoutErrorResponse); err != nil {
			return nil, fmt.Errorf("could not unmarshal error response: %v", err)
		}
		return &PayResponse{
			Success:      false,
			Status:       payoutErrorResponse.Name,
			ErrorMessage: payoutErrorResponse.Message,
		}, nil
	}
	payoutResponse := &PayPalPayoutResponse{}
	if err := json.Unmarshal(body, payoutResponse); err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %v", err)
	}
	return &PayResponse{
		Success:    resp.StatusCode == http.StatusCreated,
		PaymentKey: payoutResponse.BatchHeader.PayoutBatchID,
		Status:     payoutResponse.BatchHeader.BatchStatus,
	}, err
}

func (p *PayPalRequester) getToken() (string, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, p.Options.OAuth2URL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(p.Options.ClientID, p.Options.Secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Language", "en_US")
	resp, err := p.do(req)
	if err != nil {
		return "", fmt.Errorf("could not do the request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unespected status (%d) while getting token: %v", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response: %v", err)
	}
	authTokenResponse := &PayPalAuthTokenResponse{}
	if err := json.Unmarshal(body, authTokenResponse); err != nil {
		return "", fmt.Errorf("could not unmarshal response: %v", err)
	}
	return authTokenResponse.AccessToken, nil
}

func (p *PayPalRequester) newHTTPRequest(token, apiEndpoint string, request interface{}) (*http.Request, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %v", err)
	}
	url := fmt.Sprintf("%s/%s", p.Options.PaymentURL, apiEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

func (p *PayPalRequester) do(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %v", err)
	}
	return resp, nil
}

func (p *PayPalRequester) buildPayoutRequest(r *PayRequest) *PayPalPayoutRequest {
	var receivers []PayPalPayoutItem
	for _, receiver := range r.Receivers {
		receivers = append(receivers, PayPalPayoutItem{
			RecipientType: "EMAIL",
			Amount: PayPalPayoutItemAmount{
				Value:    receiver.Amount,
				Currency: "EUR",
			},
			Note:         fmt.Sprintf(payPalEmailNoteFormat, r.ProjectName),
			SenderItemID: time.Now().String(),
			Receiver:     receiver.Email,
		})
	}
	return &PayPalPayoutRequest{
		SenderBatchHeader: PayPalPayoutSenderBatchHeader{
			SenderBatchID: r.TrackID,
			RecipientType: "EMAIL",
			EmailSubject:  fmt.Sprint(payPalEmailSubject),
		},
		Items: receivers,
	}
}
