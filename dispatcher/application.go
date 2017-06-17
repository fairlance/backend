package dispatcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fairlance/backend/models"
)

type ApplicationDispatcher interface {
	GetProject(id uint) ([]byte, error)
}

func NewApplicationDispatcher(url string) ApplicationDispatcher {
	return &applicationDispatcher{url}
}

type applicationDispatcher struct {
	url string
}

func (d *applicationDispatcher) GetProject(id uint) ([]byte, error) {
	url := fmt.Sprintf("http://%s/private/project/%d", d.url, id)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("\nStatus: %s\n Body: %s\nURL: %s", response.Status, content, url)
		return nil, err
	}
	var envelope models.JSONEnvelope
	if err := json.Unmarshal(content, &envelope); err != nil {
		return nil, err
	}
	if envelope.Code != http.StatusOK {
		err = fmt.Errorf("\nEnvelope: %+v\nURL: %s", envelope, url)
		return nil, err
	}
	return json.Marshal(envelope.Data)
}
