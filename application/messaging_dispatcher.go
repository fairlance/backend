package application

import (
	"fmt"
	"time"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/models"
)

type MessagingDispatcher struct {
	messaging dispatcher.Messaging
}

func NewMessagingDispatcher(messaging dispatcher.Messaging) *MessagingDispatcher {
	return &MessagingDispatcher{messaging}
}

func (m *MessagingDispatcher) send(projectID uint, data map[string]interface{}) error {
	message := &dispatcher.Message{
		From: dispatcher.MessageUser{
			ID:       0,
			Type:     "system",
			Username: "system",
		},
		Data:      data,
		ProjectID: fmt.Sprint(projectID),
		Timestamp: timeToMillis(time.Now()),
	}
	return m.messaging.Send(message)
}

func (m *MessagingDispatcher) sendProjectContractProposalAdded(projectID uint, proposal *Proposal, user *models.User) error {
	return m.send(projectID, map[string]interface{}{
		"proposal": proposal,
		"user":     user,
		"type":     "project_contract_proposal",
	})
}

func (m *MessagingDispatcher) sendProjectContractExtensionProposalAdded(projectID uint, extension *Extension, proposal *Proposal) error {
	return m.send(projectID, map[string]interface{}{
		"proposal":  proposal,
		"extension": extension,
		"type":      "project_contract_extension_proposal",
	})
}

func (m *MessagingDispatcher) sendProjectStateChanged(project *Project) error {
	return m.send(project.ID, map[string]interface{}{
		"status": project.Status,
		"type":   "project_status_changed",
	})
}

func (m *MessagingDispatcher) sendContractAccepted(project *Project, user *models.User) error {
	return m.send(project.ID, map[string]interface{}{
		"user": user,
		"type": "project_contract_accepted",
	})
}

func (m *MessagingDispatcher) sendProjectFinishedByFreelancer(project *Project, user *models.User) error {
	return m.send(project.ID, map[string]interface{}{
		"user": user,
		"type": "project_finished_by_freelancer",
	})
}

func (m *MessagingDispatcher) sendProjectDone(project *Project, user *models.User) error {
	return m.send(project.ID, map[string]interface{}{
		"user": user,
		"type": "project_done",
	})
}

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
}
