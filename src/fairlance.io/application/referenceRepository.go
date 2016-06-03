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

func (repo *ReferenceRepository) GetReferences(freelancerId int) ([]Reference, error) {
	references := []Reference{}
	repo.db.Find(&references).Where("freelancerId = ?", freelancerId)
	return references, nil
}
