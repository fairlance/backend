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
	users      map[string]*User
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
		users:      make(map[string]*User),
		db:         newMongoDatabase(mongoHost, "notification"),
	}
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case usr := <-h.register:
			messages := h.registerUser(usr)
			for _, message := range messages {
				usr.send <- message
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
		Username: user.FirstName + " " + user.LastName,
		Type:     user.Type,
		ID:       user.ID,
		conn:     conn,
		send:     make(chan Message),
		hub:      h,
	}
	h.register <- usr
	go usr.startWriting()
	go usr.startReading()
}

func (h *Hub) sendMessage(b []byte) {
	msg := h.buildAndSaveMessage(b)
	h.broadcast <- *msg
}

func (h *Hub) registerUser(usr *User) []Message {
	h.usersMU.Lock()
	h.users[usr.uniqueID()] = usr
	h.usersMU.Unlock()
	var messages = []Message{}
	messages, err := h.db.loadLastDocs(usr.uniqueID(), 20)
	if err != nil {
		log.Println(err)
		return messages
	}
	return messages
}

func (h *Hub) unregisterUser(usr *User) {
	h.usersMU.Lock()
	delete(h.users, usr.uniqueID())
	h.usersMU.Unlock()
}

func (h *Hub) getUsersForMessage(msg *Message) []*User {
	if msg.Type != "read" && len(msg.To) == 0 {
		log.Println("error: message not addressed to anyone")
		return []*User{}
	}
	u := []*User{}
	for _, messageUser := range msg.To {
		h.usersMU.RLock()
		user, ok := h.users[messageUser.uniqueID()]
		h.usersMU.RUnlock()
		if !ok {
			log.Printf("user not found [%+v]", messageUser)
			continue
		}
		u = append(u, user)
	}
	return u
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
		if _, ok := h.users[uniqueIDFrom]; !ok {
			log.Printf("could not find user: %+v", msg.From)
			h.usersMU.RUnlock()
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
