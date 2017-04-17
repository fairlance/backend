package messaging

import (
	"fmt"
	"log"
	"time"

	"fairlance.io/dispatcher"
)

// Hub ...
type Hub struct {
	// Registered users.
	rooms map[string]*Room

	// Inbound messages from the users.
	broadcast chan Message

	// Register requests from the users.
	register chan *userConn

	// Unregister requests from users.
	unregister chan *User

	db       messageDB
	notifier dispatcher.Notifier
	getARoom func(id string) (*Room, error)
}

// NewHub creates new Hub object
func NewHub(db messageDB, notifier dispatcher.Notifier, getARoomFunc func(id string) (*Room, error)) *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *userConn),
		unregister: make(chan *User),
		rooms:      make(map[string]*Room),
		db:         db,
		notifier:   notifier,
		getARoom:   getARoomFunc,
	}
}

// Run the Hub
func (h *Hub) Run() {
	for {
		select {
		case newConnection := <-h.register:
			if h.rooms[newConnection.room] == nil {
				log.Println("room does not exist", newConnection.room)
				break
			}
			user, err := h.rooms[newConnection.room].ActivateUser(newConnection)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("registering", user.username, "to room:", user.room)
			h.sendOldMessagesToUser(user)
		case user := <-h.unregister:
			log.Println("unregistering", user.username, "from room:", user.room)
			h.removeUser(user)
		case msg := <-h.broadcast:
			log.Println("broadcasting message", msg)
			h.db.save(msg)
			if h.rooms[msg.ProjectID] == nil {
				log.Println("sending to unknown room", msg.ProjectID)
				continue
			}
			h.printRooms()
			for _, usr := range h.rooms[msg.ProjectID].Users {
				if usr.online {
					select {
					case usr.send <- msg:
					default:
						usr.Close()
					}
				} else {
					h.notifyUser(usr, msg)
				}
			}
		}
	}
}

func (h *Hub) removeUser(userToUnregister *User) {
	for _, usr := range h.rooms[userToUnregister.room].Users {
		if userToUnregister.username == usr.username {
			usr.Close()
		}
	}

	if !h.rooms[userToUnregister.room].HasReasonToExist() {
		h.rooms[userToUnregister.room].Close()
		delete(h.rooms, userToUnregister.room)
	}
}

// SendMessage ...
func (h *Hub) SendMessage(room, name, msg string) {
	h.broadcast <- NewMessage(0, "system", name, []byte(msg), room)
}

func (h *Hub) notifyUser(u *User, msg Message) {
	log.Println("notifying", u.username, "with message", msg.Text)
	h.notifier.Notify(&dispatcher.Notification{
		To: []dispatcher.NotificationUser{
			dispatcher.NotificationUser{
				ID:   u.id,
				Type: u.userType,
			},
		},
		From: dispatcher.NotificationUser{
			ID:   msg.UserID,
			Type: msg.UserType,
		},
		Type: "new_message",
		Data: map[string]interface{}{
			"message":   msg.Text,
			"username":  msg.Username,
			"timestamp": fmt.Sprintf("%d", msg.Timestamp),
			"time":      time.Unix(0, msg.Timestamp*1000000),
			"projectId": msg.ProjectID,
		},
	})
}

func (h *Hub) sendOldMessagesToUser(u *User) {
	messages, err := h.db.loadLastMessagesForUser(u, 20)
	if err != nil {
		log.Println(err)
		return
	}

	for _, msg := range messages {
		u.send <- msg
	}
}

func (h *Hub) printRooms() {
	log.Println("no of rooms", len(h.rooms))
	for _, room := range h.rooms {
		log.Printf("in room %s", room.ID)
		for _, usr := range room.Users {
			log.Printf("    %s, online: %b ", usr.username, usr.online)
		}
	}
}
