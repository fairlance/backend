package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fairlance/backend/dispatcher"
)

type Hub struct {
	projectRoomsMU         sync.RWMutex
	projectRooms           map[uint]*ProjectRoom
	broadcast              chan Message
	register               chan *userConn
	unregister             chan *User
	db                     messageDB
	notificationDispatcher dispatcher.Notifications
	applicationDispatcher  dispatcher.Application
}

func NewHub(db messageDB, notificationDispatcher dispatcher.Notifications, applicationDispatcher dispatcher.Application) *Hub {
	return &Hub{
		broadcast:    make(chan Message),
		register:     make(chan *userConn),
		unregister:   make(chan *User),
		projectRooms: make(map[uint]*ProjectRoom),
		db:           db,
		notificationDispatcher: notificationDispatcher,
		applicationDispatcher:  applicationDispatcher,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case newConnection := <-h.register:
			log.Println("registering", newConnection.userType, newConnection.userID, "to room:", newConnection.projectID)
			projectRoom, err := h.getProjectRoom(newConnection.projectID)
			if err != nil {
				log.Println(err)
				break
			}
			allowedUser := projectRoom.getAllowedUser(newConnection.userType, newConnection.userID)
			if allowedUser == nil {
				break
			}
			user := newUser(h, newConnection, allowedUser)
			projectRoom.addUser(user)
			go user.startWriting()
			go user.startReading()
			h.sendOldMessagesToUser(user)
		case user := <-h.unregister:
			log.Println("unregistering user", user.fairlanceType, user.fairlanceID, "from room:", user.projectID)
			h.removeUser(user)
		case msg := <-h.broadcast:
			log.Println("broadcasting message", msg)
			h.db.save(msg)
			projectRoom, err := h.getProjectRoom(msg.ProjectID)
			if err != nil {
				log.Printf("could get projectRoom %d: %v", msg.ProjectID, err)
				break
			}
			for usr := range projectRoom.users {
				select {
				case usr.send <- msg:
				default:
					h.removeUser(usr)
				}
			}
			for _, usr := range projectRoom.getAbsentUsers() {
				h.notifyUser(usr.fairlanceType, usr.fairlanceID, msg)
			}
		}
	}
}

func (h *Hub) getProjectRoom(projectID uint) (*ProjectRoom, error) {
	h.projectRoomsMU.RLock()
	projectRoom, ok := h.projectRooms[projectID]
	h.projectRoomsMU.RUnlock()
	if ok {
		return projectRoom, nil
	}
	projectBytes, err := h.applicationDispatcher.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.NewDecoder(bytes.NewReader(projectBytes)).Decode(&project); err != nil {
		return nil, err
	}
	users := buildAllowedUsersMap(project)
	projectRoom = newProjectRoom(projectID, users)
	h.projectRoomsMU.Lock()
	h.projectRooms[projectID] = projectRoom
	h.projectRoomsMU.Unlock()
	return projectRoom, nil
}

func buildAllowedUsersMap(project Project) map[*AllowedUser]bool {
	users := make(map[*AllowedUser]bool)
	client := &AllowedUser{
		fairlanceID:   project.Client.ID,
		fairlanceType: "client",
		firstName:     project.Client.FirstName,
		lastName:      project.Client.LastName,
	}
	users[client] = true
	for _, freelancer := range project.Freelancers {
		f := &AllowedUser{
			fairlanceID:   freelancer.ID,
			fairlanceType: "freelancer",
			firstName:     freelancer.FirstName,
			lastName:      freelancer.LastName,
		}
		users[f] = true
	}
	return users
}

func (h *Hub) removeUser(userToUnregister *User) {
	h.projectRoomsMU.Lock()
	close(userToUnregister.send)
	h.projectRooms[userToUnregister.projectID].removeUser(userToUnregister)
	h.projectRoomsMU.Unlock()
	if !h.projectRooms[userToUnregister.projectID].hasReasonToExist() {
		h.projectRoomsMU.Lock()
		delete(h.projectRooms, userToUnregister.projectID)
		h.projectRoomsMU.Unlock()
	}
}

func (h *Hub) SendMessage(room string, msg Message) {
	h.broadcast <- msg
}

func (h *Hub) notifyUser(userType string, userID uint, msg Message) {
	log.Printf("notifying %s %d with %+v", userType, userID, msg.Data)
	h.notificationDispatcher.Notify(&dispatcher.Notification{
		To: []dispatcher.NotificationUser{
			dispatcher.NotificationUser{
				ID:   userID,
				Type: userType,
			},
		},
		From: dispatcher.NotificationUser{
			ID:   msg.From.ID,
			Type: msg.From.Type,
		},
		Type: "new_message",
		Data: map[string]interface{}{
			"message":   msg,
			"username":  msg.From.Username,
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
	log.Println("no of rooms", len(h.projectRooms))
	for _, room := range h.projectRooms {
		log.Printf("in room %d", room.id)
		log.Println("allowed users")
		for usr := range room.allowedUsers {
			log.Printf("    %s %s", usr.firstName, usr.lastName)
		}
		log.Println("connected users")
		for usr := range room.users {
			log.Printf("    %s", usr.username)
		}
	}
}
