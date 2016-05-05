package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) (*ProjectRepository, error) {
	repo := &ProjectRepository{db}

	return repo, nil
}

func (repo *ProjectRepository) GetAllProjects() ([]Project, error) {
	projects := []Project{}
	repo.db.Preload("Freelancers").Find(&projects)
	return projects, nil
}
