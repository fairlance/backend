package messaging

import (
	"gopkg.in/mgo.v2"
	"log"
)

type messageDB interface {
	save(msg message) error
	loadLastMessagesForUser(user *user) ([]message, error)
}

func NewMessageDB() messageDB {
	// Setup mongo db connection
	s, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf(err.Error())
	}

	s.DB("messaging").DropDatabase()

	return &mongoDB{s}
}

type mongoDB struct {
	s *mgo.Session
}

func (m mongoDB) save(msg message) error {
	session := m.s.Copy()
	defer session.Close()
	return session.DB("messaging").C("project_" + msg.ProjectID).Insert(&msg)
}

func (m mongoDB) loadLastMessagesForUser(user *user) ([]message, error) {
	session := m.s.Copy()
	defer session.Close()

	roomName := "project_" + user.projectID
	messages := []message{}

	count, err := session.DB("messaging").C(roomName).Count()
	if err != nil {
		return messages, err
	}

	query := session.DB("messaging").C(roomName).Find(nil)
	if count > 20 {
		query = query.Skip(count - 20)
	}

	if err := query.All(&messages); err != nil {
		return messages, err
	}

	return messages, nil
}
