package application

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	repo.db.Preload("Projects").Find(&clients)
	return clients, nil
}
