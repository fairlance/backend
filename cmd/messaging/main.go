package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

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

	f, err := os.OpenFile("/var/log/fairlance/messaging.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
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

	hub := messaging.NewHub(messaging.NewMessageDB(mongoHost), dispatcher.NewHTTPNotifier(notificationURL), getARoom)
	go hub.Run()

	// todo: make safe
	router.Handle("/{username}/{room}/send", middleware.RecoverHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			room := mux.Vars(r)["room"]
			username := mux.Vars(r)["username"]
			message := r.URL.Query().Get("message")
			hub.SendMessage(room, username, message)
		})))

	// requires a GET 'token' parameter
	router.Handle("/{room}/ws", middleware.Chain(
		middleware.RecoverHandler,
		messaging.WithRoom(hub),
		messaging.WithTokenFromParams,
		middleware.AuthenticateTokenWithClaims(secret),
		middleware.WithUserFromClaims,
		messaging.ValidateUser(hub),
	)(messaging.ServeWS(hub)))

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
