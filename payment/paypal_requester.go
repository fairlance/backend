package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	payEndpoint = "payouts"
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

func (p *PayPalRequester) VerifyPayment(r *http.Request) (bool, error) {
	// if err := r.ParseForm(); err != nil {
	// 	log.Printf("could not parse IPN form: %v", err)
	// 	return false, fmt.Errorf("could not parse IPN form: %v", err)
	// }
	// notificationMap := make(map[string]string)
	// postStr := p.Options.IPNNotificationURL + "&cmd=_notify-validate&"
	// for key, v := range r.Form {
	// 	value := strings.Join(v, "")
	// 	log.Printf("key: %s, value: %s", key, value)
	// 	notificationMap[key] = value
	// 	postStr = postStr + key + "=" + url.QueryEscape(value) + "&"
	// }

	// To verify the message from PayPal, we must send
	// back the contents in the exact order they were received and precede it with
	// the command _notify-validate
	// PayPal will then send one single-word message, either VERIFIED,
	// if the message is valid, or INVALID if the messages is not valid.
	// See more at
	// https://developer.paypal.com/webapps/developer/docs/classic/ipn/integration-guide/IPNIntro/
	// post data back to PayPal
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %v", err)
		return false, fmt.Errorf("could not read request body: %v", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", p.Options.IPNNotificationURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("could not create verification POST request: %v", err)
		return false, fmt.Errorf("could not create verification POST request: %v", err)
	}
	req.Header.Add("Content-Type: ", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("could not send verification POST request: %v", err)
		return false, fmt.Errorf("could not send verification POST request: %v", err)
	}
	log.Println("Response:")
	log.Println(resp)
	log.Println("Status:")
	log.Println(resp.Status)
	// convert response to string
	respStr, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response String: ", string(respStr))
	if string(respStr) != "VERIFIED" {
		fmt.Println("IPN validation failed!")
		fmt.Println("Do not send the stuff out yet!")
		return false, nil
	}
	fmt.Println("IPN verified")
	fmt.Println("TODO : Email receipt, increase credit, etc")
	return true, nil
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
			Note:         fmt.Sprintf("Project %d", r.ProjectID),
			SenderItemID: time.Now().String(),
			Receiver:     receiver.Email,
		})
	}
	return &PayPalPayoutRequest{
		SenderBatchHeader: PayPalPayoutSenderBatchHeader{
			// SenderBatchID: t.TrackID,
			RecipientType: "EMAIL",
			EmailSubject:  fmt.Sprintf("Payment for project %d!", r.ProjectID),
		},
		Items: receivers,
	}
}
