package wsrouter

import (
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

var writeWait = 10 * time.Second
var readWait = 30 * time.Minute

func (router *Router) StartReading(u User, conn *websocket.Conn) {
	defer func() {
		router.unregister <- u
		conn.Close()
	}()
	conn.SetReadDeadline(time.Now().Add(readWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(readWait)); return nil })
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := router.conf.BuildMessage(msgBytes)
		if msg != nil {
			router.broadcast <- *msg
		}
	}
}

func (router *Router) StartWriting(u User, conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()
	for {
		select {
		case m, ok := <-u.send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			sw, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messages := []Message{m}

			n := len(router.broadcast)
			for i := 0; i < n; i++ {
				messages = append(messages, <-router.broadcast)
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

func (router *Router) ServeWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		router.Handle(r, conn)
	})
}
