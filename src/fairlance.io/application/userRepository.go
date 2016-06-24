package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	repo := &UserRepository{db}

	return repo, nil
}

func (repo UserRepository) CheckCredentials(email string, password string) (User, error) {
	user, err := repo.GetUserByEmail(email)
	if err != nil {
		return user, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, errors.New("User not found (password is wrong)")
	}

	return user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (User, error) {
	user := User{}
	freelancer := Freelancer{}
	if repo.db.Where("email = ?", email).First(&freelancer).RecordNotFound() {
		client := Client{}
		if repo.db.Where("email = ?", email).First(&client).RecordNotFound() {
			return user, errors.New("User not found")
		}
		user = client.User
	} else {
		user = freelancer.User
	}

	return user, nil
}
