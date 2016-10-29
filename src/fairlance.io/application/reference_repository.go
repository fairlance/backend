package application

import (
	"github.com/jinzhu/gorm"
)

type ReferenceRepository interface {
	AddReference(reference *Reference) error
}

type PostgreReferenceRepository struct {
	db *gorm.DB
}

func NewPostgreReferenceRepository(db *gorm.DB) (*PostgreReferenceRepository, error) {
	repo := &PostgreReferenceRepository{db}

	return repo, nil
}

func (repo *PostgreReferenceRepository) AddReference(reference *Reference) error {
	return repo.db.Create(reference).Error
}
