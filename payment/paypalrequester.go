package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	respond "gopkg.in/matryer/respond.v1"
)

const (
	payEndpoint            = "Pay"
	executePaymentEndpoint = "ExecutePayment"
	paymentDetailsEndpont  = "PaymentDetails"
)

type Options struct {
	AdaptivePaymentsURL string
	AuthorizationURL    string
	ReturnURL           string
	CancelURL           string
	SecurityUserID      string
	SecurityPassword    string
	SecuritySignature   string
	ApplicationID       string
}

func NewPayPalRequester(options *Options) *PayPalRequester {
	return &PayPalRequester{options}
}

type PayPalRequester struct {
	options *Options
}

func (p *PayPalRequester) PayPrimaryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivers := []Receiver{}
		response, err := p.payPrimary(receivers)
		if err != nil {
			log.Printf("could not execute a payPrimary request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a payPrimary request: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, struct {
				RedirectURL string
				Response    *PayResponse
			}{
				RedirectURL: fmt.Sprintf("%s%s", p.AuthorizationURL, response.PayKey),
				Response:    response,
			})
			return
		}
		respond.With(w, r, http.StatusInternalServerError, response)
	})
}

func (p *PayPalRequester) PaymentDetailsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("payKey") == "" { // should be project, and pay key is stored localy
			respond.With(w, r, http.StatusBadRequest, "payKey missing")
			return
		}
		payKey := r.URL.Query().Get("payKey")
		response, err := p.paymentDetails(payKey)
		if err != nil {
			log.Printf("could not execute a paymentDetails request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a paymentDetails request: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, response)
			return
		}
		respond.With(w, r, http.StatusInternalServerError, response)
	})
}

func (p *PayPalRequester) ExecutePaymentHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("payKey") == "" { // should be project, and pay key is stored localy
			respond.With(w, r, http.StatusBadRequest, "payKey missing")
			return
		}
		payKey := r.URL.Query().Get("payKey")
		response, err := p.executePayment(payKey)
		if err != nil {
			log.Printf("could not execute a executePayment request: %v", err)
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not execute a executePayment request: %v", err))
			return
		}
		if response.ResponseEnvelope.Ack == "Success" {
			respond.With(w, r, http.StatusOK, response)
			return
		}
		respond.With(w, r, http.StatusInternalServerError, response)
	})
}

func (p *PayPalRequester) payPrimary(receivers []Receiver) (*PayResponse, error) {
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
	}
	req, err := p.newRequest(payPrimaryRequest, payEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	payResponse := &PayResponse{}
	err = p.do(req, payResponse)
	return payResponse, err
}

func (p *PayPalRequester) paymentDetails(payKey string) (*PaymentDetailsResponse, error) {
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

func (p *PayPalRequester) executePayment(payKey string) (*ExecutePaymentResponse, error) {
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

func (p *PayPalRequester) newRequest(request interface{}, apiEndpoint string) (*http.Request, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %v", err)
	}
	url := fmt.Sprintf("%s/%s", p.options.AdaptivePaymentsURL, apiEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	// todo move to PayPalRequster struct
	req.Header.Add("X-PAYPAL-SECURITY-USERID", p.options.SecurityUserID)
	req.Header.Add("X-PAYPAL-SECURITY-PASSWORD", p.options.SecurityPassword)
	req.Header.Add("X-PAYPAL-SECURITY-SIGNATURE", p.options.SecuritySignature)
	req.Header.Add("X-PAYPAL-APPLICATION-ID", p.options.ApplicationID)
	req.Header.Add("X-PAYPAL-REQUEST-DATA-FORMAT", "JSON")
	req.Header.Add("X-PAYPAL-RESPONSE-DATA-FORMAT", "JSON")
	return req, nil
}

func (p *PayPalRequester) do(req *http.Request, response interface{}) error {
	client := &http.Client{
		Timeout: time.Duration(2 * time.Second),
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
	log.Println(string(body))
	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("could not unmarshal response: %v", err)
	}
	return nil
}
