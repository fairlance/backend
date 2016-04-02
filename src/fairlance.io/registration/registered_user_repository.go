package registration

import (
	"gopkg.in/mgo.v2"
)

type RegisteredUserRepository struct {
	session *mgo.Session
	db      string
}

func NewRegisteredUserRepository(db string) (*RegisteredUserRepository, error) {
	// Setup db connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	repo := &RegisteredUserRepository{session, db}
	err = repo.getUsersCollection().EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo RegisteredUserRepository) GetAllRegisteredUsers() ([]RegisteredUser, error) {
	registeredUsers := []RegisteredUser{}

	err := repo.getUsersCollection().Find(nil).All(&registeredUsers)
	if err != nil {
		return registeredUsers, err
	}

	return registeredUsers, nil
}

func (repo RegisteredUserRepository) AddRegisteredUser(user RegisteredUser) error {
	err := repo.getUsersCollection().Insert(&user)
	if err != nil {
		if mgo.IsDup(err) {
			return err
		}

		return err
	}

	return nil
}

func (repo RegisteredUserRepository) getUsersCollection() *mgo.Collection {
	return repo.session.DB(repo.db).C("users")
}
