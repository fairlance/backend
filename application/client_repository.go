package application

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type ClientRepository interface {
	GetAllClients() ([]Client, error)
	AddClient(client *Client) error
	UpdateClient(client *Client) error
	GetClient(id uint) (*Client, error)
}

type PostgreClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) (ClientRepository, error) {
	repo := &PostgreClientRepository{db}

	return repo, nil
}

func (repo *PostgreClientRepository) GetAllClients() ([]Client, error) {
	clients := []Client{}
	repo.db.Preload("Jobs").Preload("Projects").Preload("Reviews").Preload("Reviews.Freelancer").Find(&clients)
	return clients, nil
}

func (repo *PostgreClientRepository) AddClient(client *Client) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	client.Password = string(hashedPassword)
	return repo.db.Create(client).Error
}

func (repo *PostgreClientRepository) UpdateClient(client *Client) error {
	return repo.db.Save(client).Error
}

func (repo *PostgreClientRepository) GetClient(id uint) (*Client, error) {
	client := &Client{}
	if err := repo.db.Preload("Projects").Preload("Jobs").Preload("Reviews").Preload("Reviews.Freelancer").Find(client, id).Error; err != nil {
		return client, err
	}
	return client, nil
}
