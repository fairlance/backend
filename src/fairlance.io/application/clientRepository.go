package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) (*ClientRepository, error) {
	repo := &ClientRepository{db}

	return repo, nil
}

func (repo *ClientRepository) GetAllClients() ([]Client, error) {
	clients := []Client{}
	repo.db.Preload("Jobs").Preload("Projects").Preload("Reviews").Find(&clients)
	return clients, nil
}

func (repo *ClientRepository) AddClient(client *Client) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	client.Password = string(hashedPassword)
	return repo.db.Create(client).Error
}
