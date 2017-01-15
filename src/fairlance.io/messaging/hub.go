package messaging

import (
	"container/list"
	"log"

	"time"
)

// Hub ...
type Hub struct {
	// Registered users.
	rooms map[string]*list.List

	// Inbound messages from the users.
	broadcast chan message

	// Register requests from the users.
	register chan *user

	// Unregister requests from users.
	unregister chan *user

	db messageDB
}

// NewHub creates new Hub object
func NewHub(db messageDB) *Hub {
	return &Hub{
		broadcast:                 make(chan message),
		register:                  make(chan *user),
		unregister:                make(chan *user),
		rooms:                     make(map[string]*list.List),
		db:                        db,
	}
}

// Run the Hub
func (h *Hub) Run() {
	for {
		select {
		case userToRegister := <-h.register:
			log.Println("registering", userToRegister.username, "to room:", userToRegister.projectID)
			if h.rooms[userToRegister.projectID] == nil {
				h.rooms[userToRegister.projectID] = list.New()
			}
			h.sendOldMessagesToUser(userToRegister)
			h.rooms[userToRegister.projectID].PushFront(userToRegister)
		case user := <-h.unregister:
			log.Println("unregistering", user.username, "from room:", user.projectID)
			h.removeUser(user)
		case message := <-h.broadcast:
			h.db.save(message)
			if h.rooms[message.ProjectID] == nil {
				continue
			}

			log.Println("broadcasting message", message)
			for el := h.rooms[message.ProjectID].Front(); el != nil; el = el.Next() {
				userInRoom := el.Value.(*user)
				select {
				case userInRoom.send <- message:
				default:
					h.rooms[userInRoom.projectID].Remove(el)
					close(userInRoom.send)
				}
			}
		}
	}
}

func (h *Hub) removeUser(userToUnregister *user) {
	for el := h.rooms[userToUnregister.projectID].Front(); el != nil; el = el.Next() {
		u := el.Value.(*user)
		if userToUnregister.username == u.username {
			h.rooms[userToUnregister.projectID].Remove(el)
		}
	}

	if h.rooms[userToUnregister.projectID].Len() == 0 {
		delete(h.rooms, userToUnregister.projectID)
	}
}

// SendMessage ...
func (h *Hub) SendMessage(room, name, msg string) {
	h.broadcast <- newMessage(0, name, []byte(msg), time.Now().Unix(), room)
}

func (h *Hub) sendOldMessagesToUser(u *user) {
	messages, err := h.db.loadLastMessagesForUser(u)
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
		log.Printf("in room %s", room.Front().Value.(*user).projectID)
		for el := room.Front(); el != nil; el = el.Next() {
			c := el.Value.(*user)
			log.Printf("    %s", c.username)
		}
	}
}
