package notification

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/fairlance/backend/models"
	"github.com/gorilla/websocket"
)

type Hub struct {
	usersMU    sync.RWMutex
	users      map[*User]bool
	broadcast  chan Message
	register   chan *User
	unregister chan *User
	db         *mongoDB
}

func NewHub(mongoHost string) *Hub {
	hub := &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *User),
		unregister: make(chan *User),
		users:      make(map[*User]bool),
		db:         newMongoDatabase(mongoHost, "notification"),
	}
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case usr := <-h.register:
			h.registerUser(usr)
			if err := h.sendOldMessagesToUser(usr); err != nil {
				log.Printf("could not send old messages to user %s: %v", usr.uniqueID(), err)
			}
		case usr := <-h.unregister:
			h.unregisterUser(usr)
		case msg := <-h.broadcast:
			users := h.getUsersForMessage(&msg)
			for _, user := range users {
				select {
				case user.send <- msg:
				default:
					h.unregisterUser(user)
				}
			}
		}
	}
}

func (h *Hub) addUser(user *models.User, conn *websocket.Conn) {
	usr := &User{
		Type: user.Type,
		ID:   user.ID,
		conn: conn,
		send: make(chan Message),
		hub:  h,
	}
	h.register <- usr
	go usr.startWriting()
	go usr.startReading()
}

func (h *Hub) sendMessage(b []byte) {
	msg := h.buildAndSaveMessage(b)
	h.broadcast <- *msg
}

func (h *Hub) registerUser(usr *User) {
	h.usersMU.Lock()
	h.users[usr] = true
	h.usersMU.Unlock()
}

func (h *Hub) sendOldMessagesToUser(usr *User) error {
	messages, err := h.db.loadLastDocs(usr.uniqueID(), 20)
	if err != nil {
		return err
	}
	for _, message := range messages {
		usr.send <- message
	}
	return nil
}

func (h *Hub) unregisterUser(usr *User) {
	h.usersMU.Lock()
	delete(h.users, usr)
	h.usersMU.Unlock()
}

func (h *Hub) getUsersForMessage(msg *Message) []*User {
	if msg.Type != "read" && len(msg.To) == 0 {
		log.Printf("message %d not addressed to anyone", msg.Timestamp)
		return []*User{}
	}
	u := []*User{}
	for _, messageUser := range msg.To {
		h.usersMU.RLock()
		for user := range h.users {
			if user.uniqueID() == messageUser.uniqueID() {
				u = append(u, user)
			}
		}
		h.usersMU.RUnlock()
	}
	return u
}

func (h *Hub) userExists(uniqueID string) bool {
	h.usersMU.RLock()
	defer h.usersMU.RUnlock()
	for user := range h.users {
		if user.uniqueID() == uniqueID {
			return true
		}
	}
	return false
}

func (h *Hub) buildAndSaveMessage(b []byte) *Message {
	var msg = &Message{}
	if err := json.Unmarshal(b, msg); err != nil {
		log.Printf("could not build message: %v", err)
		return nil
	}
	msg.Timestamp = timeToMillis(time.Now())
	switch msg.Type {
	case "read":
		h.usersMU.RLock()
		uniqueIDFrom := msg.From.uniqueID()
		if !h.userExists(uniqueIDFrom) {
			h.usersMU.RUnlock()
			log.Printf("could not find user: %+v", msg.From)
			return nil
		}
		h.usersMU.RUnlock()
		timestampString, ok := msg.Data["timestamp"].(string)
		if !ok {
			log.Printf("could not parse timestamp: %s", msg.Data["timestamp"])
			return nil
		}
		timestampInt, err := strconv.ParseInt(timestampString, 10, 64)
		if err != nil {
			log.Println(err)
			return nil
		}
		if err := h.db.markRead(uniqueIDFrom, timestampInt); err != nil {
			log.Println(err)
			return nil
		}
	default:
		for _, messageUser := range msg.To {
			h.db.save(messageUser.uniqueID(), *msg)
		}
	}
	return msg
}

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func (h *Hub) printUsers() {
	log.Printf("num of users %d", len(h.users))
	log.Println("connected users")
	for usr := range h.users {
		log.Printf("    %s", usr.uniqueID())
	}
}
