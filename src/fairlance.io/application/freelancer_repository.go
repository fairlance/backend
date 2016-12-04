package application

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type FreelancerRepository interface {
	GetAllFreelancers() ([]Freelancer, error)
	GetFreelancer(id uint) (Freelancer, error)
	AddFreelancer(freelancer *Freelancer) error
	UpdateFreelancer(freelancer *Freelancer) error
	DeleteFreelancer(id uint) error
	AddReview(id uint, newReview *Review) error
}

type PostgreFreelancerRepository struct {
	db *gorm.DB
}

func NewFreelancerRepository(db *gorm.DB) (FreelancerRepository, error) {
	repo := &PostgreFreelancerRepository{db}

	return repo, nil
}

func (repo *PostgreFreelancerRepository) GetAllFreelancers() ([]Freelancer, error) {
	freelancers := []Freelancer{}
	if err := repo.db.Preload("Projects").Preload("References").Preload("References.Media").Preload("Reviews").Find(&freelancers).Error; err != nil {
		return freelancers, err
	}
	return freelancers, nil
}

func (repo *PostgreFreelancerRepository) GetFreelancer(id uint) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.Preload("Projects").Preload("References").Preload("References.Media").Preload("Reviews").Find(&freelancer, id).Error; err != nil {
		return freelancer, err
	}
	return freelancer, nil
}

func (repo *PostgreFreelancerRepository) AddFreelancer(freelancer *Freelancer) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(freelancer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	freelancer.Password = string(hashedPassword)
	return repo.db.Create(freelancer).Error
}

func (repo *PostgreFreelancerRepository) UpdateFreelancer(freelancer *Freelancer) error {
	return repo.db.Save(freelancer).Error
}

func (repo *PostgreFreelancerRepository) DeleteFreelancer(id uint) error {
	freelancer := Freelancer{}
	if repo.db.Preload("References").Find(&freelancer, id).RecordNotFound() {
		return errors.New("Freelancer not found")
	}

	for _, reference := range freelancer.References {
		if err := repo.db.Delete(&Media{}, "reference_id = ?", reference.ID).Error; err != nil {
			return err
		}
	}

	// probably faster than deleting one by one in loop above
	if err := repo.db.Delete(&Reference{}, "freelancer_id = ?", id).Error; err != nil {
		return err
	}

	if err := repo.db.Delete(&Review{}, "freelancer_id = ?", id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&freelancer).Error
}

func (repo *PostgreFreelancerRepository) AddReview(freelancerID uint, review *Review) error {
	freelancer := Freelancer{}

	err := repo.db.Preload("Reviews").Find(&freelancer, freelancerID).Error
	if err != nil {
		return err
	}

	freelancer.Reviews = append(freelancer.Reviews, *review)

	rating := review.Rating
	for _, review := range freelancer.Reviews {
		rating += review.Rating
	}
	freelancer.Rating = rating / float64(len(freelancer.Reviews)+1)

	return repo.db.Save(&freelancer).Error
}
