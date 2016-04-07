package application

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type FreelancerRepository struct {
	db *gorm.DB
}

func NewFreelancerRepository(db *gorm.DB) (*FreelancerRepository, error) {
	repo := &FreelancerRepository{db}

	return repo, nil
}
func (repo *FreelancerRepository) GetAllFreelancers() ([]Freelancer, error) {
	freelancers := []Freelancer{}
	if err := repo.db.Preload("Projects").Find(&freelancers).Limit(100).Error; err != nil {
		return freelancers, err
	}
	for i := 0; i < len(freelancers); i++ {
		if err := repo.hydrate(&freelancers[i]); err != nil {
			return freelancers, err
		}
	}
	return freelancers, nil
}

func (repo *FreelancerRepository) GetFreelancerByEmail(email string) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.Where("email = ?", email).First(&freelancer).Error; err != nil {
		return freelancer, err
	}
	return freelancer, nil
}

func (repo *FreelancerRepository) GetFreelancer(id uint) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.Preload("Projects").Where("id = ?", id).Find(&freelancer).Error; err != nil {
		return freelancer, err
	}
	if err := repo.hydrate(&freelancer); err != nil {
		return freelancer, err
	}
	return freelancer, nil
}

func (repo *FreelancerRepository) CheckCredentials(email string, password string) error {
	freelancer, err := repo.GetFreelancerByEmail(email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(freelancer.Password), []byte(password)); err != nil {
		return errors.New("Freelancer not found (password is wrong)")
	}

	return nil
}

func (repo *FreelancerRepository) AddFreelancer(freelancer *Freelancer) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(freelancer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	freelancer.Password = string(hashedPassword)
	return repo.db.Create(freelancer).Error
}

func (repo *FreelancerRepository) UpdateFreelancer(freelancer *Freelancer) error {
	if freelancer.ID == 0 {
		return errors.New("Can't update entity without id")
	}
	return repo.db.Save(freelancer).Error
}

func (repo *FreelancerRepository) DeleteFreelancer(id uint) error {
	freelancer := Freelancer{}
	freelancer.ID = id
	if repo.db.Find(&freelancer).RecordNotFound() {
		return errors.New("Freelancer not found")
	}
	return repo.db.Delete(&freelancer).Error
}

func (repo *FreelancerRepository) hydrate(freelancer *Freelancer) error {
	if err := json.Unmarshal([]byte(freelancer.JsonComments), &freelancer.Comments); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(freelancer.JsonReferences), &freelancer.References); err != nil {
		return err
	}
	return nil
}

func (repo *FreelancerRepository) AddReference(freelancerId uint, reference Reference) error {

	freelancer, err := repo.GetFreelancer(freelancerId)
	if err != nil {
		return err
	}
	references, err := json.Marshal(append(freelancer.References, reference))
	if err != nil {
		return err
	}
	freelancer.JsonReferences = string(references)

	return nil
}
