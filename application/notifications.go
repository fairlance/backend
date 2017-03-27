package application

import (
	"log"

	"fairlance.io/notifier"
)

func notifyJobApplicationAdded(n notifier.Notifier, jobApplication *JobApplication, clientID uint) {
	not := &notifier.Notification{
		From: notifier.NotificationSystemUser,
		To: []notifier.NotificationUser{
			notifier.NotificationUser{
				ID:   clientID,
				Type: "client",
			},
		},
		Type: "jobApplicationAdded",
		Data: map[string]interface{}{
			"jobApplication": jobApplication,
		},
	}
	if err := n.Notify(not); err != nil {
		log.Println(err)
	}
}
