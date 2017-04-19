package application

import "fairlance.io/dispatcher"
import "encoding/json"
import "fmt"

type MessagingDispatcher struct {
	messaging dispatcher.Messaging
}

func NewMessagingDispatcher(messaging dispatcher.Messaging) *MessagingDispatcher {
	return &MessagingDispatcher{messaging}
}

func (m *MessagingDispatcher) send(projectID uint, textObject interface{}) error {
	textBytes, err := json.Marshal(textObject)
	if err != nil {
		return err
	}
	message := &dispatcher.Message{
		Text:      string(textBytes),
		UserID:    0,
		UserType:  "system",
		Username:  "system",
		ProjectID: fmt.Sprint(projectID),
	}
	return m.messaging.Send(message)
}

func (m *MessagingDispatcher) sendProjectContractProposalAdded(projectID uint, proposal *Proposal) error {
	msgTextObject := struct {
		Proposal *Proposal `json:"proposal"`
		Type     string    `json:"type"`
	}{
		proposal,
		"project_contract_proposal",
	}
	return m.send(projectID, msgTextObject)
}

func (m *MessagingDispatcher) sendProjectContractExtensionProposalAdded(projectID uint, extension *Extension, proposal *Proposal) error {
	msgTextObject := struct {
		Proposal  *Proposal  `json:"proposal"`
		Extension *Extension `json:"extension"`
		Type      string     `json:"type"`
	}{
		proposal,
		extension,
		"project_contract_extension_proposal",
	}
	return m.send(projectID, msgTextObject)
}

func (m *MessagingDispatcher) sendProjectStateChanged(project *Project) error {
	msgTextObject := struct {
		Status string `json:"status"`
		Type   string `json:"type"`
	}{
		project.Status,
		"project_status_changed",
	}
	return m.send(project.ID, msgTextObject)
}

func (m *MessagingDispatcher) sendContractAccepted(project *Project, userType string, user *User) error {
	msgTextObject := struct {
		User     *User  `json:"user"`
		UserType string `json:"userType"`
		Type     string `json:"type"`
	}{
		user,
		userType,
		"project_contract_accepted",
	}
	return m.send(project.ID, msgTextObject)
}

func (m *MessagingDispatcher) sendProjectConcludedByUser(project *Project, userType string, user *User) error {
	msgTextObject := struct {
		User     *User  `json:"user"`
		UserType string `json:"userType"`
		Type     string `json:"type"`
	}{
		user,
		userType,
		"project_concluded_by_user",
	}
	return m.send(project.ID, msgTextObject)
}
