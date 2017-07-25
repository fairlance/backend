package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var writeWait = 10 * time.Second
var readWait = 4 * time.Hour

type User struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	Type     string `json:"type"`
	send     chan Message
	conn     *websocket.Conn
	hub      *Hub
}

func (u *User) uniqueID() string {
	return fmt.Sprintf("%s_%d", u.Type, u.ID)
}

func (u *User) startReading() {
	defer func() {
		u.hub.unregister <- u
		u.conn.Close()
	}()
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
		msg := u.hub.buildAndSaveMessage(msgBytes)
		if msg != nil {
			u.hub.broadcast <- *msg
		}
	}
}

func (u *User) startWriting() {
	defer func() {
		u.conn.Close()
	}()
	for {
		select {
		case m, ok := <-u.send:
			if !ok {
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			sw, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			messages := []Message{m}
			n := len(u.hub.broadcast)
			for i := 0; i < n; i++ {
				messages = append(messages, <-u.hub.broadcast)
			}
			messagesAsBytes, err := json.Marshal(messages)
			if err != nil {
				log.Println(err.Error())
				return
			}
			sw.Write(messagesAsBytes)
			if err := sw.Close(); err != nil {
				return
			}
		}
	}
}
