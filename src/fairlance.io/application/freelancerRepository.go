package application

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
	if err := repo.db.Preload("Skills", "owner_type = ?", "freelancers").Preload("Projects").Preload("References").Preload("References.Media").Preload("Reviews").Find(&freelancers).Error; err != nil {
		return freelancers, err
	}
	return freelancers, nil
}

func (repo *FreelancerRepository) GetFreelancer(id uint) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.Preload("Skills", "owner_type = ?", "freelancers").Preload("Projects").Preload("References").Preload("References.Media").Preload("Reviews").Find(&freelancer, id).Error; err != nil {
		return freelancer, err
	}
	return freelancer, nil
}

func (repo *FreelancerRepository) AddFreelancer(freelancer *Freelancer) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(freelancer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	freelancer.Password = string(hashedPassword)
	return repo.db.Create(freelancer).Error
}

func (repo *FreelancerRepository) ClearSkills(freelancer *Freelancer) error {
	return repo.db.Delete(&Tag{}, "owner_id = ? AND owner_type = ?", freelancer.ID, "freelancers").Error
}

func (repo *FreelancerRepository) UpdateFreelancer(freelancer *Freelancer) error {
	return repo.db.Save(freelancer).Error
}

func (repo *FreelancerRepository) DeleteFreelancer(id uint) error {
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
