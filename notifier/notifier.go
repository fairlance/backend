package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Notifier interface {
	Notify(n *Notification) error
}

type NotificationUser struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

var NotificationSystemUser = NotificationUser{
	ID:   0,
	Type: "system",
}

type Notification struct {
	From NotificationUser       `json:"from,omitempty"`
	To   []NotificationUser     `json:"to,omitempty"`
	Type string                 `json:"type,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type HTTPNotifier struct {
	NotificationURL string
}

func NewHTTPNotifier(notificationURL string) *HTTPNotifier {
	return &HTTPNotifier{notificationURL}
}

func (notifier *HTTPNotifier) Notify(n *Notification) error {
	url := "http://" + notifier.NotificationURL + "/send"
	body, err := json.Marshal(&n)
	if err != nil {
		return err
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		err = fmt.Errorf(
			"bad status: %d\n body: %s\nfor request: %s\nrequest body: %s",
			response.StatusCode,
			contents,
			url,
			body)
		response.Body.Close()
		return err
	}

	return nil
}
