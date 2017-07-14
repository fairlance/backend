package dispatcher

import (
	"fmt"
	"net/http"
	"time"
)

type Application interface {
	GetProject(id uint) ([]byte, error)
	SetProjectFunded(id uint) error
}

func NewApplication(url string) Application {
	return &httpApplication{
		url: url,
		client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

type httpApplication struct {
	url    string
	client *http.Client
}

func (d *httpApplication) GetProject(id uint) ([]byte, error) {
	url := fmt.Sprintf("%s/private/project/%d", d.url, id)
	content, err := doGET(d.client, url)
	return content, err
}

func (d *httpApplication) SetProjectFunded(id uint) error {
	url := fmt.Sprintf("%s/private/project/%d/fund", d.url, id)
	_, err := doGET(d.client, url)
	return err
}
