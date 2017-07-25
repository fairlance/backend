package notification

import "fmt"

type MessageUser struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

func (user *MessageUser) uniqueID() string {
	return fmt.Sprintf("%s_%d", user.Type, user.ID)
}

type Message struct {
	To        []MessageUser          `json:"to,omitempty"`
	From      MessageUser            `json:"from,omitempty"`
	Type      string                 `json:"type,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	Read      bool                   `json:"read"`
}
