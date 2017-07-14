package dispatcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Payment interface {
	Deposit(projectID uint) error
	Execute(projectID uint) error
}

func NewPayment(url string) Payment {
	return &httpPayment{
		url: url,
		client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

type depositRequest struct {
	ProjectID uint
}

type executeRequest struct {
	ProjectID uint
}

type httpPayment struct {
	url    string
	client *http.Client
}

func (p *httpPayment) Deposit(projectID uint) error {
	url := fmt.Sprintf("%s/private/deposit", p.url)
	b, err := json.Marshal(depositRequest{
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	return doPOST(p.client, url, b)
}

func (p *httpPayment) Execute(projectID uint) error {
	url := fmt.Sprintf("%s/private/execute", p.url)
	b, err := json.Marshal(executeRequest{
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	return doPOST(p.client, url, b)
}
