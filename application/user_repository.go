package application

import (
	"errors"
	"time"

	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckCredentials(email string, password string) (User, string, error)
	LoggedIn(ID uint, userType string) error
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

func (repo *PostgreUserRepository) LoggedIn(ID uint, userType string) error {
	user := User{Model: Model{ID: ID}}
	var db *gorm.DB
	if userType == "client" {
		db = repo.db.Model(&Client{User: user})
	} else if userType == "freelancer" {
		db = repo.db.Model(&Freelancer{User: user})
	} else {
		return fmt.Errorf("user type not regnized: %s", userType)
	}
	return db.Update("last_login", time.Now()).Error
}
