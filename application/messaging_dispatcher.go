package application

import "github.com/fairlance/backend/dispatcher"

import "fmt"

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
	}
	return m.messaging.Send(message)
}

func (m *MessagingDispatcher) sendProjectContractProposalAdded(projectID uint, proposal *Proposal) error {
	return m.send(projectID, map[string]interface{}{
		"proposal": proposal,
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

func (m *MessagingDispatcher) sendContractAccepted(project *Project, userType string, user *User) error {
	return m.send(project.ID, map[string]interface{}{
		"user":     user,
		"userType": userType,
		"type":     "project_contract_accepted",
	})
}

func (m *MessagingDispatcher) sendProjectConcludedByUser(project *Project, userType string, user *User) error {
	return m.send(project.ID, map[string]interface{}{
		"user":     user,
		"userType": userType,
		"type":     "project_contract_accepted",
	})
}
