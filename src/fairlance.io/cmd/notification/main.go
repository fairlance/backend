package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/context"

	"log"

	"strconv"

	"encoding/json"
	"time"

	"flag"
	"os"

	"fmt"

	"fairlance.io/application"
	"fairlance.io/messaging"
	"fairlance.io/wsrouter"
	"github.com/dgrijalva/jwt-go"
)

func newUser(r *http.Request) *wsrouter.User {
	user := context.Get(r, "user").(*application.User)
	claims := context.Get(r, "claims").(jwt.MapClaims)
	userType := claims["userType"].(string)

	return &wsrouter.User{
		Username: user.FirstName + " " + user.LastName,
		Type:     userType,
		ID:       user.ID,
	}
}

func userUniqueID(user wsrouter.User) string {
	return fmt.Sprintf("%s.%d", user.Type, user.ID)
}

func uniqueID(id uint, userType string) string {
	return fmt.Sprintf("%s.%d", userType, id)
}

var notification struct {
	Users map[string]wsrouter.User
}

var port int
var secret string
var db *mongoDB

func init() {
	flag.IntVar(&port, "port", 3007, "Specify the port to listen on.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.Parse()

	f, err := os.OpenFile("/var/log/fairlance/notification.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	notification.Users = make(map[string]wsrouter.User)
	db = newMongoDatabase("notification")
	conf := wsrouter.RouterConf{
		Register: func(usr wsrouter.User) []wsrouter.Message {
			log.Println("registering", usr.Username)
			notification.Users[userUniqueID(usr)] = usr
			var messages = []wsrouter.Message{}

			messages, err := db.LoadLastDocs(userUniqueID(usr), 5)
			if err != nil {
				log.Println(err)
				return messages
			}

			return messages
		},
		Unregister: func(usr wsrouter.User) {
			log.Println("unregistering", usr.Username)
			delete(notification.Users, userUniqueID(usr))
		},
		BroadcastTo: func(msg *wsrouter.Message) []wsrouter.User {
			log.Printf("broadcast %v\n", msg)
			switch msg.Type {
			case "read":
				uniqueIDFrom := uniqueID(msg.From.ID, msg.From.Type)
				if _, ok := notification.Users[uniqueIDFrom]; !ok {
					log.Printf("error: user not found [%v]", msg.From)
					return []wsrouter.User{}
				}
				timestampFloat, ok := msg.Data["timestamp"].(float64)
				if !ok {
					log.Printf("error: timestamp not provided [%s]", msg.Data["timestamp"])
					return []wsrouter.User{}
				}

				if err := db.MarkRead(uniqueIDFrom, int64(timestampFloat)); err != nil {
					log.Println(err)
					return []wsrouter.User{}
				}
				return []wsrouter.User{}
			default:
				if len(msg.To) == 0 {
					log.Println("error: message not addressed to anyone")
					return []wsrouter.User{}
				}
				users := []wsrouter.User{}
				for _, userConf := range msg.To {
					user, ok := notification.Users[uniqueID(userConf.ID, userConf.Type)]
					if !ok {
						log.Printf("error: user not found [%v]", userConf)
					} else {
						users = append(users, user)
					}
				}
				return users
			}
		},
		CreateUser: func(r *http.Request) *wsrouter.User {
			return newUser(r)
		},
		BuildMessage: func(b []byte) *wsrouter.Message {
			var msg = &wsrouter.Message{}

			err := json.Unmarshal(b, msg)
			if err != nil {
				panic(err)
			}

			msg.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)

			if msg.Type != "read" {
				for _, userConf := range msg.To {
					db.Save(uniqueID(userConf.ID, userConf.Type), *msg)
				}
			}

			return msg
		},
	}
	router := wsrouter.NewRouter(conf)
	http.Handle("/", messaging.WithTokenFromParams(
		application.AuthenticateTokenWithClaims(
			secret,
			application.WithUserFromClaims(router.ServeWS()))))

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			var msg wsrouter.Message

			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&msg); err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			msg.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
			if msg.Type != "read" {
				for _, userConf := range msg.To {
					db.Save(uniqueID(userConf.ID, userConf.Type), msg)
				}
			}

			router.BroadcastMessage(msg)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Write([]byte("method not allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	go router.Run()

	log.Println("Started...")
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func newMongoDatabase(dbName string) *mongoDB {
	s, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("cannot connect to mongo:", err.Error())
	}

	return &mongoDB{s, dbName}
}

type mongoDB struct {
	s      *mgo.Session
	dbName string
}

func (m mongoDB) Save(collection string, doc wsrouter.Message) error {
	session := m.s.Copy()
	defer session.Close()

	return session.DB(m.dbName).C(collection).Insert(doc)
}

func (m mongoDB) MarkRead(collection string, timestamp int64) error {
	session := m.s.Copy()
	defer session.Close()

	var msg wsrouter.Message
	if err := session.DB(m.dbName).C(collection).Find(bson.M{"timestamp": timestamp}).One(&msg); err != nil {
		return err
	}
	msg.Read = true

	return session.DB(m.dbName).C(collection).Update(bson.M{"timestamp": msg.Timestamp}, msg)
}

func (m mongoDB) LoadLastDocs(collection string, num int) ([]wsrouter.Message, error) {
	session := m.s.Copy()
	defer session.Close()

	documents := []wsrouter.Message{}

	count, err := session.DB(m.dbName).C(collection).Count()
	if err != nil {
		return documents, err
	}

	query := session.DB(m.dbName).C(collection).Find(nil)
	if count > num {
		query = query.Skip(count - num)
	}

	if err := query.All(&documents); err != nil {
		return documents, err
	}

	return documents, nil
}
