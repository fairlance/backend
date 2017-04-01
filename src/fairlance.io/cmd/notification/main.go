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

	"fairlance.io/application"
	"fairlance.io/messaging"
	"fairlance.io/wsrouter"
	"github.com/dgrijalva/jwt-go"
)

// todo: addd getter/setter funcs with locking
var notification struct {
	Users map[string]wsrouter.User
}

var port int
var secret string
var db *mongoDB

// Examples:
// {"to":[{"type": "freelancer", "id": 1}],"from":{"type": "freelancer", "id": 1},"type":"notification","data":{"text":"hahahah", "projectId": 2}}
// {"type":"read", "from":{"type": "freelancer", "id": 1}, "to":[{"type": "freelancer", "id": 1}], "data": {"timestamp":1487717547735}}

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
			notification.Users[usr.UniqueID()] = usr
			var messages = []wsrouter.Message{}

			messages, err := db.LoadLastDocs(usr.UniqueID(), 20)
			if err != nil {
				log.Println(err)
				return messages
			}

			return messages
		},
		Unregister: func(usr wsrouter.User) {
			log.Println("unregistering", usr.Username)
			delete(notification.Users, usr.UniqueID())
		},
		BroadcastTo: func(msg *wsrouter.Message) []wsrouter.User {
			log.Printf("broadcast %v\n", msg)

			if len(msg.To) == 0 {
				log.Println("error: message not addressed to anyone")
				return []wsrouter.User{}
			}
			users := []wsrouter.User{}
			for _, userConf := range msg.To {
				user, ok := notification.Users[userConf.UniqueID()]
				if !ok {
					log.Printf("error: user not found [%v]", userConf)
				} else {
					users = append(users, user)
				}
			}
			return users
		},
		CreateUser: func(r *http.Request) *wsrouter.User {
			user := context.Get(r, "user").(*application.User)
			claims := context.Get(r, "claims").(jwt.MapClaims)
			userType := claims["userType"].(string)

			return &wsrouter.User{
				Username: user.FirstName + " " + user.LastName,
				Type:     userType,
				ID:       user.ID,
			}
		},
		BuildMessage: func(b []byte) *wsrouter.Message {
			var msg = &wsrouter.Message{}

			err := json.Unmarshal(b, msg)
			if err != nil {
				log.Println(err)
				return nil
			}

			msg.Timestamp = time.Now().Unix()

			switch msg.Type {
			case "read":
				uniqueIDFrom := msg.From.UniqueID()
				if _, ok := notification.Users[uniqueIDFrom]; !ok {
					log.Printf("error: user not found [%v]", msg.From)
					return nil
				}
				timestampString, ok := msg.Data["timestamp"].(string)
				if !ok {
					log.Printf("error: timestamp not provided [%s]", msg.Data["timestamp"])
					return nil
				}

				timestampInt, err := strconv.ParseInt(timestampString, 10, 0)
				if err != nil {
					log.Println(err)
					return nil
				}
				if err := db.MarkRead(uniqueIDFrom, timestampInt); err != nil {
					log.Println(err)
					return nil
				}
				return nil
			default:
				for _, userConf := range msg.To {
					db.Save(userConf.UniqueID(), *msg)
				}
			}

			return msg
		},
	}
	router := wsrouter.NewRouter(conf)
	http.Handle("/",
		application.RecoverHandler(
			application.LoggerHandler(
				messaging.WithTokenFromParams(
					application.AuthenticateTokenWithClaims(
						secret, application.WithUserFromClaims(
							router.ServeWS()))))))

	http.Handle("/send",
		application.RecoverHandler(
			application.LoggerHandler(
				application.CORSHandler(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.Method == "POST" {
							var msg wsrouter.Message

							decoder := json.NewDecoder(r.Body)
							if err := decoder.Decode(&msg); err != nil {
								w.WriteHeader(http.StatusBadRequest)
								w.Write([]byte(err.Error()))
								return
							}
							defer r.Body.Close()

							msg.Timestamp = time.Now().Unix()

							switch msg.Type {
							case "read":
								uniqueIDFrom := msg.From.UniqueID()
								if _, ok := notification.Users[uniqueIDFrom]; !ok {
									log.Printf("error: user not found [%v]", msg.From)
									return
								}
								timestampString, ok := msg.Data["timestamp"].(string)
								if !ok {
									log.Printf("error: timestamp not provided [%s]", msg.Data["timestamp"])
									return
								}

								timestampInt, err := strconv.ParseInt(timestampString, 10, 0)
								if err != nil {
									log.Println(err)
									return
								}
								if err := db.MarkRead(uniqueIDFrom, timestampInt); err != nil {
									log.Println(err)
									return
								}
								return
							default:
								for _, userConf := range msg.To {
									db.Save(userConf.UniqueID(), msg)
								}
							}

							router.BroadcastMessage(msg)
							w.WriteHeader(http.StatusOK)
							return
						}

						w.WriteHeader(http.StatusMethodNotAllowed)
						w.Write([]byte("method not allowed"))
					})))))

	go router.Run()

	log.Println("Started...")
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func newMongoDatabase(dbName string) *mongoDB {
	s, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("cannot connect to mongo:", err.Error())
	}
	s.DB(dbName).DropDatabase()

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
