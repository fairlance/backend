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
		Type: "job_application_added",
		Data: map[string]interface{}{
			"jobApplication": jobApplication,
		},
	}
	if err := n.Notify(not); err != nil {
		log.Println(err)
	}
}

func notifyJobApplicationAccepted(n notifier.Notifier, jobApplication *JobApplication, project *Project) {
	var users []notifier.NotificationUser
	for _, freelancer := range project.Freelancers {
		users = append(users, notifier.NotificationUser{
			ID:   freelancer.ID,
			Type: "freelancer",
		})
	}

	not := &notifier.Notification{
		From: notifier.NotificationSystemUser,
		To:   users,
		Type: "job_application_accepted",
		Data: map[string]interface{}{
			"jobApplication": jobApplication,
			"project":        project,
		},
	}
	if err := n.Notify(not); err != nil {
		log.Println(err)
	}
}
