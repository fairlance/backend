package messaging

import (
	"log"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	readWait = 15 * time.Minute

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

// user type
type user struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan message

	// client name
	username string

	// Room in which the client is participating
	projectID string

	// user id
	id uint
}

func (u *user) startReading() {
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		u.hub.broadcast <- newMessage(u.id, u.username, msgBytes, time.Now().Unix(), u.projectID)
	}
}

func (u *user) startWriting() {
	defer func() {
		u.conn.Close()
	}()
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

			messages := []message{msg}
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
