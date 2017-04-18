package dispatcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Messaging interface {
	Send(m *Message) error
}

type Message struct {
	UserID    uint   `json:"userId,omitempty"`
	UserType  string `json:"userType,omitempty"`
	Username  string `json:"username,omitempty"`
	Text      string `json:"text,omitempty"`
	ProjectID string `json:"projectId,omitempty"`
}

type HTTPMessaging struct {
	MessagingURL string
}

func NewHTTPMessaging(MessagingURL string) *HTTPMessaging {
	return &HTTPMessaging{MessagingURL}
}

func (m *HTTPMessaging) Send(msg *Message) error {
	url := fmt.Sprintf("http://%s/%s/%s/send?message=%s", m.MessagingURL, msg.Username, msg.ProjectID, msg.Text)
	request, err := http.NewRequest("GET", url, nil)
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
