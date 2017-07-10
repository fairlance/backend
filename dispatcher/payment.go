package dispatcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"bytes"
	"log"
)

type Payment interface {
	// Deposit(amount float64, projectID uint) (string, string, error)
	Execute(trackID uint) error
}

func NewHTTPPayment(url string) Payment {
	return &httpPayment{
		url: url,
		client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

type executeRequest struct {
	ProjectID uint
}

type httpPayment struct {
	url    string
	client *http.Client
}

func (p *httpPayment) Execute(projectID uint) error {
	url := fmt.Sprintf("http://%s/private/execute", p.url)
	b, err := json.Marshal(executeRequest{
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	response, err := p.client.Do(request)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("status: %s, body: %s, url: %s", response.Status, content, url)
		log.Printf("could not execute payment: %v", err)
		return err
	}
	return nil
}

// type depositRequest struct {
// 	Project uint
// 	Amount  string
// }

// type depositResponse struct {
// 	TrackID     string
// 	RedirectURL string
// }

// func (d *paymentDispatcher) Deposit(amount float64, projectID uint) (string, string, error) {
// 	url := fmt.Sprintf("http://%s/private/deposit", d.url)
// 	b, err := json.Marshal(depositRequest{
// 		Amount:  fmt.Sprintf("%.2f", amount),
// 		Project: projectID,
// 	})
// 	if err != nil {
// 		return "", "", err
// 	}
// 	request, err := http.NewRequest("POST", url, bytes.NewReader(b))
// 	if err != nil {
// 		return "", "", err
// 	}
// 	response, err := d.client.Do(request)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	content, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	defer response.Body.Close()
// 	if response.StatusCode != http.StatusOK {
// 		err = fmt.Errorf("status: %s, body: %s, url: %s", response.Status, content, url)
// 		log.Printf("could not deposit payment: %v", err)
// 		return "", "", err
// 	}
// 	var responseData depositResponse
// 	if err := json.Unmarshal(content, &responseData); err != nil {
// 		return "", "", err
// 	}
// 	return responseData.TrackID, responseData.RedirectURL, nil
// }
