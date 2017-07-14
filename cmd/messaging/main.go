package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/messaging"
	"github.com/fairlance/backend/middleware"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	port            int
	secret          string
	dbHost          string
	dbName          string
	dbUser          string
	dbPass          string
	mongoHost       string
	notificationURL string
)

func init() {
	flag.IntVar(&port, "port", 3005, "Specify the port to listen to.")
	flag.StringVar(&dbHost, "dbHost", "localhost", "DB host.")
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.StringVar(&mongoHost, "mongoHost", "localhost", "Mongo host.")
	flag.StringVar(&notificationURL, "notificationUrl", "localhost:3007", "Notification service url.")
	flag.Parse()

	// f, err := os.OpenFile("/var/log/fairlance/messaging.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// log.SetOutput(f)
}

type Project struct {
	ID          uint
	Freelancers []Freelancer `json:"freelancers,omitempty" gorm:"many2many:project_freelancers;"`
	ClientID    uint         `json:"-"`
	Client      *Client      `json:"client,omitempty"`
}

type Client struct {
	ID        uint
	FirstName string
	LastName  string
}

type Freelancer struct {
	ID        uint
	FirstName string
	LastName  string
}

func main() {
	db, err := gorm.Open("postgres", "host="+dbHost+" user="+dbUser+" password="+dbPass+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	// todo: call application service
	getARoom := func(id string) (*messaging.Room, error) {
		users := make(map[string]*messaging.User)

		project := Project{}
		err := db.Preload("Client").Preload("Freelancers").Find(&project, id).Error
		if err != nil {
			return nil, err
		}

		client := messaging.NewUser(
			project.ClientID,
			project.Client.FirstName,
			project.Client.LastName,
			"client",
			id)
		users[client.UniqueID()] = client
		for _, freelancer := range project.Freelancers {
			f := messaging.NewUser(
				freelancer.ID,
				freelancer.FirstName,
				freelancer.LastName,
				"freelancer",
				id)
			users[f.UniqueID()] = f
		}

		return messaging.NewRoom(id, users), nil
	}

	hub := messaging.NewHub(messaging.NewMessageDB(mongoHost), dispatcher.NewNotifier(notificationURL), getARoom)
	go hub.Run()

	router.Handle("/private/{room}/send", middleware.Chain(
		middleware.RecoverHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			room := mux.Vars(r)["room"]
			var message messaging.Message
			if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("could not reade message from body"))
				return
			}
			r.Body.Close()
			hub.SendMessage(room, message)
		})))

	// requires a GET 'token' parameter
	router.Handle("/public/{room}/ws", middleware.Chain(
		middleware.RecoverHandler,
		messaging.WithRoom(hub),
		messaging.WithTokenFromParams,
		middleware.AuthenticateTokenWithClaims(secret),
		middleware.WithUserFromClaims,
		messaging.ValidateUser(hub),
	)(messaging.ServeWS(hub)))

	http.Handle("/", router)

	log.Printf("Listening on: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
