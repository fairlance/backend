package application

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) (*ClientRepository, error) {
	repo := &ClientRepository{db}

	return repo, nil
}

func (repo *ClientRepository) GetAllClients() ([]Client, error) {
	clients := []Client{}

	//todo: extract prepare statement outside of the function
	queryStmt, err := repo.db.Prepare("SELECT * FROM clients")
	if err != nil {
		return clients, err
	}

	rows, err := queryStmt.Query()
	defer rows.Close()
	if err != nil {
		return clients, err
	}

	for rows.Next() {
		client := Client{}

		if err := rows.Scan(
			&client.Id,
			&client.Name,
			&client.Description,
			&client.Created,
		); err != nil {
			return clients, err
		}

		err := repo.hydrate(&client)
		if err != nil {
			return clients, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func (repo *ClientRepository) hydrate(client *Client) error {
	projects, err := repo.getProjects(client.Id)
	if err != nil {
		return err
	}
	client.Projects = projects

	jobs, err := repo.getJobs(client.Id)
	if err != nil {
		return err
	}
	client.Jobs = jobs

	return nil
}

func (repo *ClientRepository) getProjects(clientId int) ([]Project, error) {
	projects := []Project{}

	queryStmt, err := repo.db.Prepare(`
        SELECT p.id, p.name, p.description, p.is_active, p.created
        FROM projects p
        WHERE p.client_id = $1`)

	if err != nil {
		return projects, err
	}

	rows, err := queryStmt.Query(clientId)
	defer rows.Close()
	if err != nil {
		return projects, err
	}

	for rows.Next() {
		project := Project{}

		if err := rows.Scan(
			&project.Id,
			&project.Name,
			&project.Description,
			&project.IsActive,
			&project.Created,
		); err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (repo *ClientRepository) getJobs(clientId int) ([]Job, error) {
	jobs := []Job{}

	queryStmt, err := repo.db.Prepare(`
        SELECT j.id, j.name, j.description, j.is_active, j.created
        FROM jobs j
        WHERE j.client_id = $1`)

	if err != nil {
		return jobs, err
	}

	rows, err := queryStmt.Query(clientId)
	defer rows.Close()
	if err != nil {
		return jobs, err
	}

	for rows.Next() {
		job := Job{}

		if err := rows.Scan(
			&job.Id,
			&job.Name,
			&job.Description,
			&job.IsActive,
			&job.Created,
		); err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
