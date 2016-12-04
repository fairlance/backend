package application

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckCredentials(email string, password string) (User, string, error)
}

type PostgreUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) (UserRepository, error) {
	repo := &PostgreUserRepository{db}

	return repo, nil
}

func (repo PostgreUserRepository) CheckCredentials(email string, password string) (User, string, error) {
	user, userType, err := repo.getUserByEmail(email)
	if err != nil {
		return user, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, "", errors.New("User not found (password is wrong)")
	}

	return user, userType, nil
}

func (repo *PostgreUserRepository) getUserByEmail(email string) (User, string, error) {
	user := User{}
	var userType string
	freelancer := Freelancer{}
	if repo.db.Where("email = ?", email).First(&freelancer).RecordNotFound() {
		client := Client{}
		if repo.db.Where("email = ?", email).First(&client).RecordNotFound() {
			return user, userType, errors.New("User not found")
		}
		userType = "client"
		user = client.User
	} else {
		userType = "freelancer"
		user = freelancer.User
	}

	return user, userType, nil
}
