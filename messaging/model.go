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
		UserID:    userID,
		UserType:  userType,
		Username:  username,
		Text:      string(bytes.TrimSpace(text)),
		Timestamp: timeToMillis(time.Now()),
		ProjectID: projectID,
	}
}

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

type Message struct {
	UserID    uint   `json:"userId" bson:"userId"`
	UserType  string `json:"userType" bson:"userType"`
	Username  string `json:"username" bson:"username"`
	Text      string `json:"text" bson:"text"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	ProjectID string `json:"projectId" bson:"projectId"`
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

	return nil, fmt.Errorf("user %d not found", conn.id)
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
