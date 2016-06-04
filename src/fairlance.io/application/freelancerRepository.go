package main

import (
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
	if err := repo.db.Preload("Projects").Preload("References").Preload("Reviews").Find(&freelancers).Error; err != nil {
		return freelancers, err
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
	if err := repo.db.Preload("Projects").Preload("References").Preload("Reviews").Find(&freelancer, id).Error; err != nil {
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
	return repo.db.Save(freelancer).Error
}

func (repo *FreelancerRepository) DeleteFreelancer(id uint) error {
	freelancer := Freelancer{}
	if repo.db.Find(&freelancer, id).RecordNotFound() {
		return errors.New("Freelancer not found")
	}
	return repo.db.Delete(&freelancer).Error
}

func (repo *FreelancerRepository) AddReview(newReview *Review) error {
	freelancer := Freelancer{}
	err := repo.db.Preload("Reviews").Find(&freelancer, newReview.FreelancerId).Error
	if err != nil {
		return err
	}
	rating := newReview.Rating
	for _, review := range freelancer.Reviews {
		rating += review.Rating
	}
	freelancer.Rating = rating / float64(len(freelancer.Reviews)+1)
	err = repo.db.Save(newReview).Error
	if err != nil {
		return err
	}
	return repo.db.Save(&freelancer).Error
}
