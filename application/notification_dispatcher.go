package application

import "github.com/fairlance/dispatcher"

type NotificationDispatcher struct {
	notifier dispatcher.Notifier
}

func NewNotificationDispatcher(notifier dispatcher.Notifier) *NotificationDispatcher {
	return &NotificationDispatcher{notifier}
}

func (n *NotificationDispatcher) notifyJobApplicationAdded(jobApplication *JobApplication, clientID uint) error {
	notification := &dispatcher.Notification{
		From: dispatcher.NotificationSystemUser,
		To: []dispatcher.NotificationUser{
			dispatcher.NotificationUser{
				ID:   clientID,
				Type: "client",
			},
		},
		Type: "job_application_added",
		Data: map[string]interface{}{
			"jobApplication": jobApplication,
			"jobId":          jobApplication.JobID,
		},
	}
	return n.notifier.Notify(notification)
}

func (n *NotificationDispatcher) notifyJobApplicationAccepted(jobApplication *JobApplication, project *Project) error {
	var users []dispatcher.NotificationUser
	for _, freelancer := range project.Freelancers {
		users = append(users, dispatcher.NotificationUser{
			ID:   freelancer.ID,
			Type: "freelancer",
		})
	}
	notification := &dispatcher.Notification{
		From: dispatcher.NotificationSystemUser,
		To:   users,
		Type: "job_application_accepted",
		Data: map[string]interface{}{
			"jobApplication": jobApplication,
			"project":        project,
			"jobId":          jobApplication.JobID,
		},
	}
	return n.notifier.Notify(notification)
}
