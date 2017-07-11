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
	AddReview(freelancerID uint, newReview *Review) error
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
	if err := repo.db.
		Preload("Projects").
		Preload("Projects.Client").
		Preload("References").
		Preload("References.Media").
		Preload("Reviews").
		Preload("Reviews.Client").
		Find(&freelancers).Error; err != nil {
		return freelancers, err
	}
	return freelancers, nil
}

func (repo *PostgreFreelancerRepository) GetFreelancer(freelancerID uint) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.
		Preload("Projects").
		Preload("Projects.Client").
		Preload("References").
		Preload("References.Media").
		Preload("Reviews").
		Preload("Reviews.Client").
		Preload("PortfolioItems", "type IN (?)", fileTypeFreelancerPortfolioItems).
		Preload("PortfolioLinks", "type IN (?)", fileTypeFreelancerPortfolioLinks).
		Preload("AdditionalFiles", "type IN (?)", fileTypeFreelancerAdditionalField).
		Find(&freelancer, freelancerID).Error; err != nil {
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

func (repo *PostgreFreelancerRepository) DeleteFreelancer(freelancerID uint) error {
	freelancer := Freelancer{}
	if repo.db.Preload("References").Find(&freelancer, freelancerID).RecordNotFound() {
		return errors.New("Freelancer not found")
	}
	tx := repo.db.Begin()
	for _, reference := range freelancer.References {
		if err := tx.Delete(&Media{}, "reference_id = ?", reference.ID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// probably faster than deleting one by one in loop above
	if err := tx.Delete(&Reference{}, "freelancer_id = ?", freelancerID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&Review{}, "freelancer_id = ?", freelancerID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("owner_type = 'freelancers' AND owner_id = ?", freelancerID).Delete(&File{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := repo.db.Delete(&freelancer).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *PostgreFreelancerRepository) AddReview(freelancerID uint, review *Review) error {
	freelancer := &Freelancer{}
	err := repo.db.Preload("Reviews").Find(freelancer, freelancerID).Error
	if err != nil {
		return err
	}
	freelancer.Reviews = append(freelancer.Reviews, *review)
	var rating float64
	for _, review := range freelancer.Reviews {
		rating += review.Rating
	}
	freelancer.Rating = round(rating/float64(len(freelancer.Reviews)), 0.5, 1)
	return repo.db.Save(freelancer).Error
}
