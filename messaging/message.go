package messaging

import (
	"bytes"
	"time"
)

func NewMessage(projectID, userID uint, userType, username string, text []byte) Message {
	return Message{
		From: MessageUser{
			ID:       userID,
			Type:     userType,
			Username: username,
		},
		Data: map[string]interface{}{
			"text": string(bytes.TrimSpace(text)),
		},
		Timestamp: timeToMillis(time.Now()),
		ProjectID: projectID,
	}
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

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
}
