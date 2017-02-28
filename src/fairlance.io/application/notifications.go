package application

import "log"

func notifyJobApplicationAdded(n notifier, jobApplication *JobApplication, clientID uint) {
	not := &notification{
		From: notificationSystemUser,
		To: []notificationUser{
			notificationUser{
				ID:   clientID,
				Type: "client",
			},
		},
		Type: "jobApplicationAdded",
		Data: map[string]interface{}{
			"freelancer": jobApplication.Freelancer,
		},
	}
	if err := n.Notify(not); err != nil {
		log.Println(err)
	}
}
