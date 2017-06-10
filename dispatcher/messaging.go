package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	ProjectID string                 `json:"projectId" bson:"projectId"`
}

type HTTPMessaging struct {
	MessagingURL string
}

func NewHTTPMessaging(MessagingURL string) Messaging {
	return &HTTPMessaging{MessagingURL}
}

func (m *HTTPMessaging) Send(msg *Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("could not mashall msg: %v", err)
		return err
	}
	url := fmt.Sprintf("http://%s/%s/send", m.MessagingURL, msg.ProjectID)
	request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Println(err)
		return err
	}
	response, err := http.DefaultClient.Do(request)
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
		err = fmt.Errorf("\nStatus: %s\n Body: %s\nURL: %s", response.Status, contents, url)
		response.Body.Close()
		return err
	}

	return nil
}
