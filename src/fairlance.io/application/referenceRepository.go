package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ReferenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) (*ReferenceRepository, error) {
	repo := &ReferenceRepository{db}

	return repo, nil
}

func (repo *ReferenceRepository) AddReference(reference *Reference) error {
	return repo.db.Create(reference).Error
}
