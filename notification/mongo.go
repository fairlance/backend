package notification

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func newMongoDatabase(mongoHost, dbName string) *mongoDB {
	s, err := mgo.Dial(mongoHost)
	if err != nil {
		log.Fatal("cannot connect to mongo:", err.Error())
	}
	return &mongoDB{s, dbName}
}

type mongoDB struct {
	s      *mgo.Session
	dbName string
}

func (m mongoDB) save(collection string, doc Message) error {
	session := m.s.Copy()
	defer session.Close()

	return session.DB(m.dbName).C(collection).Insert(doc)
}

func (m mongoDB) markRead(collection string, timestamp int64) error {
	session := m.s.Copy()
	defer session.Close()

	var msg Message
	if err := session.DB(m.dbName).C(collection).Find(bson.M{"timestamp": timestamp}).One(&msg); err != nil {
		return err
	}
	msg.Read = true

	return session.DB(m.dbName).C(collection).Update(bson.M{"timestamp": msg.Timestamp}, msg)
}

func (m mongoDB) loadLastDocs(collection string, num int) ([]Message, error) {
	session := m.s.Copy()
	defer session.Close()

	documents := []Message{}

	count, err := session.DB(m.dbName).C(collection).Count()
	if err != nil {
		return documents, err
	}

	query := session.DB(m.dbName).C(collection).Find(bson.M{"read": false})
	if count > num {
		query = query.Skip(count - num)
	}

	if err := query.All(&documents); err != nil {
		return documents, err
	}

	return documents, nil
}
