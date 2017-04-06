package application

import "github.com/jinzhu/gorm"

type ProjectRepository interface {
	getAllProjects() ([]Project, error)
	getByID(id uint) (Project, error)
	add(project *Project) error
	getAllProjectsForClient(id uint) ([]Project, error)
	getAllProjectsForFreelancer(id uint) ([]Project, error)
}

const (
	getAllProjectsForFreelancerQuery = "SELECT p.id, p.name, p.description, p.status, p.due_date, c.id AS client_id, c.first_name, c.last_name FROM projects p, clients c, project_freelancers pf WHERE pf.freelancer_id = ? AND p.client_id = c.id GROUP BY c.id, p.id"
)

type postgreProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) (ProjectRepository, error) {
	repo := &postgreProjectRepository{db}

	return repo, nil
}

func (repo *postgreProjectRepository) getAllProjects() ([]Project, error) {
	projects := []Project{}
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&projects).Error
	return projects, err
}

// todo: check if user has access to project
func (repo *postgreProjectRepository) getByID(id uint) (Project, error) {
	project := Project{}
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&project, id).Error
	return project, err
}

func (repo *postgreProjectRepository) add(project *Project) error {
	return repo.db.Create(project).Error
}
func (repo *postgreProjectRepository) getAllProjectsForClient(id uint) ([]Project, error) {
	projects := []Project{}
	err := repo.db.Preload("Client").Preload("Freelancers").Where("client_id = ?", id).Find(&projects).Error
	return projects, err
}

func (repo *postgreProjectRepository) getAllProjectsForFreelancer(id uint) ([]Project, error) {
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
