package main

import (
	"gopkg.in/mgo.v2"
)

func getMongoDBSession() *mgo.Session {
	// Setup db connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	err = session.DB("registration").C("users").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	if err != nil {
		panic(err)
	}

	return session
}

func getAllRegisteredUsers(context *appContext) ([]RegisteredUser, error) {
	registeredUsers := []RegisteredUser{}

	session := context.session.Copy()
	defer session.Close()

	users := session.DB("registration").C("users")

	err := users.Find(nil).All(&registeredUsers)
	if err != nil {
		return registeredUsers, err
	}

	return registeredUsers, nil
}

func addRegisteredUser(context *appContext, userEmail string) error {
	session := context.session.Copy()
	defer session.Close()

	users := session.DB("registration").C("users")

	err := users.Insert(&RegisteredUser{userEmail})
	if err != nil {
		if mgo.IsDup(err) {
			return err
		}

		return err
	}

	return nil
}
