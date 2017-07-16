package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/messaging"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var port = os.Getenv("PORT")
	var secret = os.Getenv("SECRET")
	var mongoHost = os.Getenv("MONGO_HOST")
	var notificationURL = os.Getenv("NOTIFICATION_URL")
	var applicationURL = os.Getenv("APPLICATION_URL")
	hub := messaging.NewHub(
		messaging.NewMessageDB(mongoHost),
		dispatcher.NewNotifier(notificationURL),
		&fakeDispatcher{applicationURL}, //dispatcher.NewApplication(applicationURL),
	)
	go hub.Run()
	http.Handle("/", messaging.NewRouter(hub, secret))
	log.Printf("Listening on: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

type fakeDispatcher struct{ applicationURL string }

func (d *fakeDispatcher) GetProject(id uint) ([]byte, error) {
	project := messaging.Project{
		Freelancers: []messaging.Freelancer{
			{
				ID:        1,
				FirstName: "First",
				LastName:  "First",
			},
		},
		Client: &messaging.Client{
			ID:        1,
			FirstName: "Client",
			LastName:  "Client",
		},
	}
	content, err := json.Marshal(project)
	return content, err
}

func (d *fakeDispatcher) SetProjectFunded(id uint) error { return nil }
