package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func NewHTTPNotifier(notificationURL string) Notifier {
	return &HTTPNotifier{notificationURL}
}

func (notifier *HTTPNotifier) Notify(n *Notification) error {
	url := "http://" + notifier.NotificationURL + "/send"
	body, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	if response.StatusCode != http.StatusOK {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			return err
		}
		err = fmt.Errorf(
			"bad status: %d\n body: %s\nfor request: %s\nrequest body: %s",
			response.StatusCode,
			contents,
			url,
			body)
		response.Body.Close()
		log.Println(err)
		return err
	}

	return nil
}
