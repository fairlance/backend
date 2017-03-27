package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	app "fairlance.io/application"
	"fairlance.io/messaging"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	port       int
	secret     string
	dbName     string
	dbUser     string
	dbPass     string
	projectURL string
)

func init() {
	flag.IntVar(&port, "port", 3005, "Specify the port to listen to.")
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.StringVar(&projectURL, "projectURL", "", "Project url to check for user access rights.")
	flag.Parse()

	f, err := os.OpenFile("/var/log/fairlance/messaging.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	db, err := gorm.Open("postgres", "user="+dbUser+" password="+dbPass+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	getARoom := func(id string) (*messaging.Room, error) {
		users := make(map[string]*messaging.User)

		project := app.Project{}
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

	hub := messaging.NewHub(messaging.NewMessageDB(), getARoom)
	go hub.Run()

	// todo: make safe
	router.Handle("/{username}/{room}/send", app.RecoverHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			room := mux.Vars(r)["room"]
			username := mux.Vars(r)["username"]
			message := r.URL.Query().Get("message")
			hub.SendMessage(room, username, message)
		})))

	// requires a GET 'token' parameter
	router.Handle("/{room}/ws",
		app.RecoverHandler(
			messaging.WithRoom(
				hub, messaging.WithTokenFromParams(
					app.AuthenticateTokenWithClaims(
						secret, app.WithUserFromClaims(
							messaging.ValidateUser(
								hub, messaging.ServeWS(hub))))))))

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
