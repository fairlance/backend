package dispatcher

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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

type httpNotifier struct {
	url    string
	client *http.Client
}

func NewNotifier(notificationURL string) Notifier {
	return &httpNotifier{
		url: notificationURL,
		client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

func (n *httpNotifier) Notify(notification *Notification) error {
	url := n.url + "/send"
	body, err := json.Marshal(notification)
	if err != nil {
		log.Println(err)
		return err
	}
	return doPOST(n.client, url, body)
}
