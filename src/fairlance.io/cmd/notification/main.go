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

func newUser(r *http.Request) *wsrouter.User {
	user := context.Get(r, "user").(*application.User)
	claims := context.Get(r, "claims").(jwt.MapClaims)
	userType := claims["userType"].(string)

	id := strconv.Itoa(int(user.ID))

	return &wsrouter.User{
		Username: user.FirstName + " " + user.LastName,
		ID:       userType + "." + id,
	}
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
	db = NewMongoDatabase("notification")
	conf := wsrouter.RouterConf{
		Register: func(usr wsrouter.User) []wsrouter.Message {
			log.Println("registering", usr.Username)
			notification.Users[usr.ID] = usr
			var messages = []wsrouter.Message{}

			messages, err := db.LoadLastDocs(usr.ID, 5)
			if err != nil {
				log.Println(err)
				return messages
			}

			return messages
		},
		Unregister: func(usr wsrouter.User) {
			log.Println("unregistering", usr.Username)
			delete(notification.Users, usr.ID)
		},
		BroadcastTo: func(msg wsrouter.Message) []wsrouter.User {
			log.Println("broadcast", msg)
			user, ok := notification.Users[msg.To[0]]
			if !ok {
				return []wsrouter.User{}
			}

			return []wsrouter.User{user}
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

			usr := notification.Users[msg.To[0]]
			db.Save(usr.ID, *msg)

			return msg
		},
	}
	router := wsrouter.NewRouter(conf)
	http.Handle("/", messaging.WithTokenFromParams(
		application.AuthenticateTokenWithClaims(
			secret,
			application.WithUserFromClaims(router.ServeWS()))))

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := wsrouter.Message{}
		router.BroadcastMessage(msg)
	})

	go router.Run()

	log.Println("Started...")
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func NewMongoDatabase(dbName string) *mongoDB {
	// Setup mongo db connection
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

func (m mongoDB) Update(collection string, doc wsrouter.Message) error {
	session := m.s.Copy()
	defer session.Close()

	return session.DB(m.dbName).C(collection).Update(bson.M{"timestamp": doc.Timestamp}, doc)
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
