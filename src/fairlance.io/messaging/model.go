package messaging

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// NewMessage ...
func NewMessage(userID uint, userType string, username string, text []byte, projectID string) Message {
	return Message{
		UserID:    userID,
		UserType:  userType,
		Username:  username,
		Text:      string(bytes.TrimSpace(text)),
		Timestamp: time.Now().Unix(),
		ProjectID: projectID,
	}
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
	ID    string
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
		user.hub = conn.hub
		user.conn = conn.conn
		user.send = make(chan Message, 256)
		user.online = true

		go user.startWriting()
		go user.startReading()

		return user, nil
	}

	return nil, fmt.Errorf("user %d not found", conn.id)
}

func (r *Room) Close() {
	for _, user := range r.Users {
		r.CloseUser(user)
	}
}

func (r *Room) CloseUser(u *User) {
	user, ok := r.Users[fmt.Sprintf("%s.%d", u.userType, u.id)]
	if ok && user.online == true {
		user.hub = nil
		if user.send != nil {
			close(user.send)
		}
		user.online = false
	}
}

func (r *Room) HasUser(id uint, userType string) bool {
	for _, user := range r.Users {
		if id == user.id && userType == user.userType {
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
