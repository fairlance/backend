package messaging

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

type messageDB interface {
	save(msg Message) error
	loadLastMessagesForUser(user *User, num int) ([]Message, error)
}

func NewMessageDB(host string) messageDB {
	// Setup mongo db connection
	s, err := mgo.Dial(host)
	if err != nil {
		log.Fatalf("open mongo db: %v", err)
	}
	return &mongoDB{s}
}

type mongoDB struct {
	s *mgo.Session
}

func (m *mongoDB) save(msg Message) error {
	session := m.s.Copy()
	defer session.Close()
	return session.DB("messaging").C("project_" + msg.ProjectID).Insert(&msg)
}

func (m *mongoDB) loadLastMessagesForUser(user *User, num int) ([]Message, error) {
	session := m.s.Copy()
	defer session.Close()

	roomName := "project_" + user.room
	messages := []Message{}

	count, err := session.DB("messaging").C(roomName).Count()
	if err != nil {
		return messages, err
	}

	query := session.DB("messaging").C(roomName).Find(nil)
	if count > num {
		query = query.Skip(count - num)
	}

	if err := query.All(&messages); err != nil {
		return messages, err
	}

	return messages, nil
}
