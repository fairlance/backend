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

func (m *MessagingDispatcher) sendProjectContractProposalAdded(projectID uint, proposal *Proposal) error {
	msgTextObject := struct {
		Proposal *Proposal `json:"proposal"`
		Type     string    `json:"type"`
	}{
		proposal,
		"project_contract_proposal",
	}

	text, err := json.Marshal(msgTextObject)
	if err != nil {
		return err
	}

	message := &dispatcher.Message{
		Text:      string(text),
		UserID:    0,
		UserType:  "system",
		Username:  "system",
		ProjectID: fmt.Sprint(projectID),
	}
	return m.messaging.Send(message)
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

	text, err := json.Marshal(msgTextObject)
	if err != nil {
		return err
	}

	message := &dispatcher.Message{
		Text:      string(text),
		UserID:    0,
		UserType:  "system",
		Username:  "system",
		ProjectID: fmt.Sprint(projectID),
	}
	return m.messaging.Send(message)
}
