package application

import "github.com/jinzhu/gorm"

type ProjectRepository interface {
	getAllProjects() ([]Project, error)
	getByID(id uint) (*Project, error)
	add(project *Project) error
	update(project *Project, fields map[string]interface{}) error
	getAllProjectsForClient(id uint) ([]Project, error)
	getAllProjectsForFreelancer(id uint) ([]Project, error)
	projectBelongsToUser(id uint, userType string, userID uint) (bool, error)
	addContract(contract *Contract) error
	getExtension(id uint) (*Extension, error)
	addExtension(extension *Extension) error
	updateContract(contract *Contract) error
	updateExtension(extension *Extension, fields map[string]interface{}) error
	setContractProposal(contract *Contract, proposal *Proposal) error
	setContractExtensionProposal(extension *Extension, proposal *Proposal) error
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
	err := repo.db.Preload("Contract").Preload("Contract.Extensions").Preload("Client").Preload("Freelancers").Find(&projects).Error
	return projects, err
}

func (repo *postgreProjectRepository) getByID(id uint) (*Project, error) {
	project := &Project{}
	err := repo.db.Preload("Contract").Preload("Contract.Extensions").Preload("Client").Preload("Freelancers").Find(project, id).Error
	return project, err
}

func (repo *postgreProjectRepository) add(project *Project) error {
	return repo.db.Create(project).Error
}

func (repo *postgreProjectRepository) update(project *Project, fields map[string]interface{}) error {
	return repo.db.Model(project).Update(fields).Error
}

func (repo *postgreProjectRepository) addContract(contract *Contract) error {
	return repo.db.Create(contract).Error
}

func (repo *postgreProjectRepository) addExtension(extension *Extension) error {
	return repo.db.Create(extension).Error
}

func (repo *postgreProjectRepository) projectBelongsToUser(id uint, userType string, userID uint) (bool, error) {
	var project Project
	err := repo.db.Preload("Client").Preload("Freelancers").Find(&project, id).Error
	if err != nil {
		return false, err
	}

	if userType == "client" && project.ClientID == userID {
		return true, nil
	}

	if userType == "freelancer" {
		for _, frelancer := range project.Freelancers {
			if frelancer.ID == userID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (repo *postgreProjectRepository) getAllProjectsForClient(id uint) ([]Project, error) {
	projects := []Project{}
	err := repo.db.Preload("Contract").Preload("Contract.Extensions").Preload("Client").Preload("Freelancers").Where("client_id = ?", id).Find(&projects).Error
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

	return projects, repo.db.Preload("Contract").Preload("Contract.Extensions").Preload("Client").Preload("Freelancers").Where("id IN (?)", projectIDs).Find(&projects).Error
}

func (repo *postgreProjectRepository) updateContract(contract *Contract) error {
	return repo.db.Save(contract).Error
}

func (repo *postgreProjectRepository) updateExtension(extension *Extension, fields map[string]interface{}) error {
	return repo.db.Model(extension).Update(fields).Error
}

func (repo *postgreProjectRepository) getExtension(id uint) (*Extension, error) {
	var extension Extension
	err := repo.db.Find(&extension, id).Error

	return &extension, err
}

func (repo *postgreProjectRepository) setContractProposal(contract *Contract, proposal *Proposal) error {
	return repo.db.Model(contract).Update(map[string]interface{}{
		"proposal": proposal,
	}).Error
}

func (repo *postgreProjectRepository) setContractExtensionProposal(extension *Extension, proposal *Proposal) error {
	return repo.db.Model(extension).Update(map[string]interface{}{
		"proposal": proposal,
	}).Error
}
