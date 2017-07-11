package notification

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fairlance/backend/models"
	"github.com/fairlance/backend/notification/wsrouter"
	"github.com/gorilla/context"
)

func newRouterConf(users map[string]wsrouter.User, db *mongoDB) wsrouter.RouterConf {
	return wsrouter.RouterConf{
		Register: func(usr wsrouter.User) []wsrouter.Message {
			return register(usr, users, db)
		},
		Unregister: func(usr wsrouter.User) {
			unregister(usr, users)
		},
		BroadcastTo: func(msg *wsrouter.Message) []wsrouter.User {
			return broadcastTo(msg, users)
		},
		CreateUser: func(r *http.Request) *wsrouter.User {
			return createUser(r)
		},
		BuildMessage: func(b []byte) *wsrouter.Message {
			return buildMessage(b, users, db)
		},
	}
}

func register(usr wsrouter.User, users map[string]wsrouter.User, db *mongoDB) []wsrouter.Message {
	log.Println("registering", usr.UniqueID())
	users[usr.UniqueID()] = usr
	var messages = []wsrouter.Message{}

	messages, err := db.loadLastDocs(usr.UniqueID(), 20)
	if err != nil {
		log.Println(err)
		return messages
	}

	return messages
}

func unregister(usr wsrouter.User, users map[string]wsrouter.User) {
	log.Println("unregistering", usr.UniqueID())
	delete(users, usr.UniqueID())
}

func broadcastTo(msg *wsrouter.Message, users map[string]wsrouter.User) []wsrouter.User {
	log.Printf("broadcast %v", msg)
	if msg.Type != "read" && len(msg.To) == 0 {
		log.Println("error: message not addressed to anyone")
		return []wsrouter.User{}
	}
	u := []wsrouter.User{}
	for _, userConf := range msg.To {
		user, ok := users[userConf.UniqueID()]
		if !ok {
			log.Printf("user not found [%+v]", userConf)
		} else {
			u = append(u, user)
		}
	}
	return u
}

func createUser(r *http.Request) *wsrouter.User {
	user := context.Get(r, "user").(*models.User)
	return &wsrouter.User{
		Username: user.FirstName + " " + user.LastName,
		Type:     user.Type,
		ID:       user.ID,
	}
}

func buildMessage(b []byte, users map[string]wsrouter.User, db *mongoDB) *wsrouter.Message {
	var msg = &wsrouter.Message{}

	err := json.Unmarshal(b, msg)
	if err != nil {
		log.Printf("could not build message: %v", err)
		return nil
	}

	msg.Timestamp = timeToMillis(time.Now())

	switch msg.Type {
	case "read":
		uniqueIDFrom := msg.From.UniqueID()
		if _, ok := users[uniqueIDFrom]; !ok {
			log.Printf("could not find user: %+v", msg.From)
			return nil
		}
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
		if err := db.markRead(uniqueIDFrom, timestampInt); err != nil {
			log.Println(err)
			return nil
		}
	default:
		for _, userConf := range msg.To {
			db.save(userConf.UniqueID(), *msg)
		}
	}

	return msg
}

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
}
