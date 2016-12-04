package application

import (
	"github.com/jinzhu/gorm"
)

type ReferenceRepository interface {
	AddReference(id uint, reference *Reference) error
}

type PostgreReferenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) (ReferenceRepository, error) {
	repo := &PostgreReferenceRepository{db}

	return repo, nil
}

func (repo *PostgreReferenceRepository) AddReference(id uint, reference *Reference) error {
	reference.FreelancerID = id
	return repo.db.Create(reference).Error
}
