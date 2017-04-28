package messaging

import (
	"fmt"
	"log"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	readWait = 4 * time.Hour

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

// user type
type User struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan Message
	// Room in which the client is participating
	room string
	// user id
	id       uint
	username string
	userType string
	online   bool
}

func NewUser(id uint, firstName, lastName, userType, room string) *User {
	return &User{
		id:       id,
		username: firstName + " " + lastName,
		userType: userType,
		room:     room,
	}
}

func (u *User) UniqueID() string {
	return fmt.Sprintf("%s.%d", u.userType, u.id)
}

func (u *User) Activate(conn *userConn) {
	u.hub = conn.hub
	u.conn = conn.conn
	u.send = make(chan Message, 256)
	u.online = true
}

func (u *User) Close() {
	if u.online {
		log.Println("close user", u.UniqueID(), u.online, u.send)
		u.hub = nil
		u.online = false
		close(u.send)
	}
}

func (u *User) startReading() {
	defer func() {
		u.hub.unregister <- u
		u.conn.Close()
	}()
	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(readWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(readWait)); return nil })
	for {
		_, msgBytes, err := u.conn.ReadMessage()
		if err != nil {
			log.Printf("could not read message: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		u.hub.broadcast <- NewMessage(u.id, u.userType, u.username, msgBytes, u.room)
	}
}

func (u *User) startWriting() {
	for {
		select {
		case msg, ok := <-u.send:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messages := []Message{msg}
			// Add queued chat messages to the current websocket message.
			n := len(u.send)
			for i := 0; i < n; i++ {
				messages = append(messages, <-u.send)
			}

			messagesAsBytes, err := json.Marshal(messages)
			if err != nil {
				log.Println(err.Error())
				return
			}

			w.Write(messagesAsBytes)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
