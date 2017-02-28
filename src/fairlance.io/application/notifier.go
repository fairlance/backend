package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type notifier interface {
	Notify(n *notification) error
}

type notificationUser struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

var notificationSystemUser = notificationUser{
	ID:   0,
	Type: "system",
}

type notification struct {
	From notificationUser       `json:"from,omitempty"`
	To   []notificationUser     `json:"to,omitempty"`
	Type string                 `json:"type,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type httpNotifier struct {
	NotificationURL string
}

func newHTTPNotifier(notificationURL string) *httpNotifier {
	return &httpNotifier{notificationURL}
}

func (notifier *httpNotifier) Notify(n *notification) error {
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
