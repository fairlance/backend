package messaging

import "bytes"

// NewMessage ...
func newMessage(userID uint, username string, text []byte, timestamp int64, projectID string) message {
	return message{
		UserID:    userID,
		Username:  username,
		Text:      string(bytes.TrimSpace(text)),
		Timestamp: timestamp,
		ProjectID: projectID,
	}
}

type message struct {
	UserID    uint   `json:"userId" bson:"userId"`
	Username  string `json:"username" bson:"username"`
	Text      string `json:"text" bson:"text"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	ProjectID string `json:"projectId" bson:"projectId"`
}
