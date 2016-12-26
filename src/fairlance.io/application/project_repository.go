package application

import (
	"github.com/jinzhu/gorm"
)

type ProjectRepository interface {
	GetAllProjects() ([]Project, error)
	GetByID(id uint) (Project, error)
	Add(project *Project) error
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
	repo.db.Preload("Client").Preload("Freelancers").Find(&projects)
	return projects, nil
}

func (repo *PostgreProjectRepository) GetByID(id uint) (Project, error) {
	project := Project{}
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&project, id).Error
	return project, err
}

func (repo *PostgreProjectRepository) Add(project *Project) error {
	return repo.db.Create(project).Error
}
