package messaging

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/fairlance/backend/dispatcher"
)

type Hub struct {
	projectRooms           map[uint]*ProjectRoom
	broadcast              chan Message
	register               chan *userConn
	unregister             chan *User
	db                     messageDB
	notificationDispatcher dispatcher.Notifier
	applicationDispatcher  dispatcher.Application
}

func NewHub(db messageDB, notificationDispatcher dispatcher.Notifier, applicationDispatcher dispatcher.Application) *Hub {
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
			projectRoom, err := h.getProjectRoom(newConnection.projectID)
			if err != nil {
				log.Println(err)
				break
			}
			// user := projectRoom.getUser("", newConnection.id)
			// projectRoom.register(*user, newConnection)
			// log.Println("registering", newConnection.id, "to room:", user.projectID)
			// h.sendOldMessagesToUser(user)
		case user := <-h.unregister:
			// log.Println("unregistering user", user.id, "of type", user.userType, "from room:", user.projectID)
			// h.removeUser(user)
		case msg := <-h.broadcast:
			// log.Println("broadcasting message", msg)
			// h.db.save(msg)
			// room, err := h.getRoom(msg.ProjectID)
			// if err != nil {
			// 	log.Printf("could get room %d: %v", msg.ProjectID, err)
			// 	break
			// }
			// for _, usr := range room.Users {
			// if usr.online {
			// 	select {
			// 	case usr.send <- msg:
			// 	default:
			// 		usr.Close()
			// 	}
			// } else {
			// 	h.notifyUser(usr, msg)
			// }
			// }
		}
	}
}

func (h *Hub) getProjectRoom(projectID uint) (*ProjectRoom, error) {
	projectRoom, ok := h.projectRooms[projectID]
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
	// users := make(map[string]*User)
	// client := &User{
	// 	id:        project.Client.ID,
	// 	userType:  "client",
	// 	username:  fmt.Sprintf("%s %s", project.Client.FirstName, project.Client.LastName),
	// 	projectID: projectID,
	// }
	// users[fmt.Sprintf("%d_%s", timeToMillis(time.Now()), client.UniqueID())] = client
	// for _, freelancer := range project.Freelancers {
	// 	f := &User{
	// 		id:       freelancer.ID,
	// 		username: fmt.Sprintf("%s %s", freelancer.FirstName, freelancer.LastName),
	// 		userType: "freelancer",
	// 		room:     projectID,
	// 	}
	// 	users[f.UniqueID()] = f
	// }

	projectRoom = &ProjectRoom{
		id:      projectID,
		project: &project,
	}
	h.projectRooms[projectID] = projectRoom
	return projectRoom, nil
}

// func (h *Hub) removeUser(userToUnregister *User) {
// 	for _, usr := range h.rooms[userToUnregister.room].Users {
// 		if userToUnregister.username == usr.username {
// 			usr.Close()
// 		}
// 	}

// 	if !h.rooms[userToUnregister.room].HasReasonToExist() {
// 		h.rooms[userToUnregister.room].Close()
// 		delete(h.rooms, userToUnregister.room)
// 	}
// }

func (h *Hub) SendMessage(room string, msg Message) {
	h.broadcast <- msg
}

// func (h *Hub) notifyUser(u *User, msg Message) {
// 	log.Printf("notifying %s with %+v", u.username, msg.Data)
// 	h.notifier.Notify(&dispatcher.Notification{
// 		To: []dispatcher.NotificationUser{
// 			dispatcher.NotificationUser{
// 				ID:   u.id,
// 				Type: u.userType,
// 			},
// 		},
// 		From: dispatcher.NotificationUser{
// 			ID:   msg.From.ID,
// 			Type: msg.From.Type,
// 		},
// 		Type: "new_message",
// 		Data: map[string]interface{}{
// 			"message":   msg,
// 			"username":  msg.From.Username,
// 			"timestamp": fmt.Sprintf("%d", msg.Timestamp),
// 			"time":      time.Unix(0, msg.Timestamp*1000000),
// 			"projectId": msg.ProjectID,
// 		},
// 	})
// }

// func (h *Hub) sendOldMessagesToUser(u *User) {
// 	messages, err := h.db.loadLastMessagesForUser(u, 20)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	for _, msg := range messages {
// 		u.send <- msg
// 	}
// }

// func (h *Hub) printRooms() {
// 	log.Println("no of rooms", len(h.rooms))
// 	for _, room := range h.rooms {
// 		log.Printf("in room %s", room.ID)
// 		for _, usr := range room.Users {
// 			log.Printf("    %s, online: %t ", usr.username, usr.online)
// 		}
// 	}
// }
