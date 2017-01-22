package application

import "github.com/jinzhu/gorm"

type ProjectRepository interface {
	GetAllProjects() ([]Project, error)
	GetByID(id uint) (Project, error)
	Add(project *Project) error
	GetAllProjectsForClient(id uint) ([]Project, error)
	GetAllProjectsForFreelancer(id uint) ([]Project, error)
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
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&projects).Error
	return projects, err
}

func (repo *PostgreProjectRepository) GetByID(id uint) (Project, error) {
	project := Project{}
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&project, id).Error
	return project, err
}

func (repo *PostgreProjectRepository) Add(project *Project) error {
	return repo.db.Create(project).Error
}

func (repo *PostgreProjectRepository) GetAllProjectsForClient(id uint) ([]Project, error) {
	projects := []Project{}
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&projects).Where("client_id = ?", id).Error
	return projects, err
}

func (repo *PostgreProjectRepository) GetAllProjectsForFreelancer(id uint) ([]Project, error) {
	projects := []Project{}
	var projectIDs []uint
	rows, err := repo.db.Table("project_freelancers").Select("project_id").Where("freelancer_id = ?", id).Rows()
	if err != nil {
		return projects, err
	}
	defer rows.Close()
	for rows.Next() {
		var projectID uint
		rows.Scan(&projectID)
		projectIDs = append(projectIDs, projectID)
	}

	return projects, repo.db.Preload("Client").Preload("Freelancers").Where("id IN (?)", projectIDs).Find(&projects).Error
}
