package application

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) (*ProjectRepository, error) {
	repo := &ProjectRepository{db}

	return repo, nil
}

func (repo *ProjectRepository) GetAllProjects() ([]Project, error) {
	projects := []Project{}

	//todo: extract prepare statement outside of the function
	queryStmt, err := repo.db.Prepare(`
        SELECT p.*, c.name, c.description, c.created
        FROM projects p
            LEFT JOIN clients c
                ON p.client_id = c.id`)

	if err != nil {
		return projects, err
	}

	rows, err := queryStmt.Query()
	defer rows.Close()

	for rows.Next() {
		project := Project{}
		client := Client{}
		client.Projects = []Project{}
		client.Jobs = []Job{}

		if err := rows.Scan(
			&project.Id,
			&project.Name,
			&project.Description,
			&client.Id,
			&project.IsActive,
			&project.Created,
			&client.Name,
			&client.Description,
			&client.Created,
		); err != nil {
			return projects, err
		}

		freelancers, err := repo.getFreelancers(project.Id)
		if err != nil {
			return projects, err
		}
		project.Freelancers = freelancers
		project.Client = client

		projects = append(projects, project)
	}

	return projects, nil
}

func (repo *ProjectRepository) getFreelancers(projectId int) ([]Freelancer, error) {
	freelancers := []Freelancer{}

	queryStmt, err := repo.db.Prepare(`
        SELECT f.id, f.first_name, f.last_name, f.email, f.created
        FROM freelancers f
            INNER JOIN project_freelancers
                ON project_freelancers.freelancer_id = f.id
        WHERE project_freelancers.project_id = $1`)

	if err != nil {
		return freelancers, err
	}

	rows, err := queryStmt.Query(projectId)
	defer rows.Close()
	if err != nil {
		return freelancers, err
	}

	for rows.Next() {
		freelancer := Freelancer{}
		freelancer.Projects = []Project{}
		freelancer.Password = ""

		if err := rows.Scan(
			&freelancer.Id,
			&freelancer.FirstName,
			&freelancer.LastName,
			&freelancer.Email,
			&freelancer.Created,
		); err != nil {
			return freelancers, err
		}
		freelancers = append(freelancers, freelancer)
	}

	return freelancers, nil
}
