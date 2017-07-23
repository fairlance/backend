package dispatcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Messaging interface {
	Send(m *Message) error
}

type MessageUser struct {
	ID       uint   `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

type Message struct {
	From      MessageUser            `json:"from,omitempty"`
	Type      string                 `json:"type,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	Read      bool                   `json:"read"`
	ProjectID uint                   `json:"projectId" bson:"projectId"`
}

type httpMessaging struct {
	url    string
	client *http.Client
}

func NewMessaging(messagingURL string) Messaging {
	return &httpMessaging{
		url: messagingURL,
		client: &http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

func (m *httpMessaging) Send(msg *Message) error {
	url := fmt.Sprintf("%s/private/%d/send", m.url, msg.ProjectID)
	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("could not mashall msg: %v", err)
		return err
	}
	return doPOST(m.client, url, b)
}
