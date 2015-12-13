package registration

import (
	"gopkg.in/mgo.v2"
)

type UserRepository struct {
	session *mgo.Session
	db      string
}

func NewUserRepository(db string) (*UserRepository, error) {
	// Setup db connection
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	repo := &UserRepository{session, db}
	err = repo.getUsers().EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo UserRepository) GetAllRegisteredUsers() ([]RegisteredUser, error) {
	registeredUsers := []RegisteredUser{}

	err := repo.getUsers().Find(nil).All(&registeredUsers)
	if err != nil {
		return registeredUsers, err
	}

	return registeredUsers, nil
}

func (repo UserRepository) AddRegisteredUser(userEmail string) error {
	err := repo.getUsers().Insert(&RegisteredUser{userEmail})
	if err != nil {
		if mgo.IsDup(err) {
			return err
		}

		return err
	}

	return nil
}

func (repo UserRepository) getUsers() *mgo.Collection {
	return repo.session.DB(repo.db).C("users")
}
