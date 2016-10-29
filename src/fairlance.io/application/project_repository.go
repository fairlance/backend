package application

import (
	"github.com/jinzhu/gorm"
)

type ProjectRepository interface {
	GetAllProjects() ([]Project, error)
}

type PostgreProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) (ProjectRepository, error) {
	repo := &PostgreProjectRepository{db}

	return repo, nil
}

func (repo *PostgreProjectRepository) GetAllProjects() ([]Project, error) {
	projects := []Project{}
	repo.db.Preload("Freelancers").Find(&projects)
	return projects, nil
}
