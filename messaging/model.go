package messaging

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fairlance/backend/models"
	"github.com/gorilla/websocket"
)

// NewMessage ...
func NewMessage(userID uint, userType string, username string, text []byte, projectID string) Message {
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

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
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

func NewRoom(id string, users map[string]*User) *Room {
	return &Room{
		ID:    id,
		Users: users,
	}
}

type Room struct {
	ID string
	// todo: allow same user twice
	Users map[string]*User
}

func (r *Room) HasReasonToExist() bool {
	for _, user := range r.Users {
		if user.online {
			return true
		}
	}
	return false
}

func (r *Room) ActivateUser(conn *userConn) (*User, error) {
	user, ok := r.Users[conn.id]
	if ok {
		user.Activate(conn)
		go user.startWriting()
		go user.startReading()
		return user, nil
	}

	return nil, fmt.Errorf("user %s not found", conn.id)
}

func (r *Room) Close() {
	for _, user := range r.Users {
		user.Close()
	}
}

func (r *Room) HasUser(user *models.User) bool {
	for _, u := range r.Users {
		if user.ID == u.id && user.Type == u.userType {
			return true
		}
	}
	return false
}

type userConn struct {
	id   string
	room string
	conn *websocket.Conn
	hub  *Hub
}
