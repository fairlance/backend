package application

import (
    "golang.org/x/crypto/bcrypt"
    "errors"
    "database/sql"
    _ "github.com/lib/pq"
    "encoding/json"
)

type FreelancerRepository struct {
    db *sql.DB
}

func NewFreelancerRepository(db *sql.DB) (*FreelancerRepository, error) {
    repo := &FreelancerRepository{db}

    return repo, nil
}

func (repo *FreelancerRepository) GetAllFreelancers() ([]Freelancer, error) {
    freelancers := []Freelancer{}

    //todo: extract prepare statement outside of the function
    queryStmt, err := repo.db.Prepare("SELECT id,first_name,last_name,email,data,created FROM freelancers")
    if err != nil {
        return freelancers, err
    }

    rows, err := queryStmt.Query()
    defer rows.Close()
    if err != nil {
        return freelancers, err
    }

    for rows.Next() {
        freelancer := Freelancer{}

        if err := rows.Scan(
            &freelancer.Id,
            &freelancer.FirstName,
            &freelancer.LastName,
            &freelancer.Email,
            &freelancer._Data,
            &freelancer.Created,
        ); err != nil {
            return freelancers, err
        }

        if freelancer.Data, err = repo.getFreelancerData(freelancer._Data); err != nil {
            return freelancers, err
        }
        if freelancer.Projects, err = repo.getProjects(freelancer.Id); err != nil {
            return freelancers, err
        }

        freelancers = append(freelancers, freelancer)
    }

    return freelancers, nil
}

func (repo *FreelancerRepository) getProjects(freelancerId int) ([]Project, error) {
    projects := []Project{}

    queryStmt, err := repo.db.Prepare(`
        SELECT p.*, c.name, c.description, c.created
        FROM projects p
            LEFT JOIN clients c
                    ON p.client_id = c.id
            INNER JOIN project_freelancers
                ON project_freelancers.project_id = p.id
        WHERE project_freelancers.freelancer_id = $1`)

    if err != nil {
        return projects, err
    }
    rows, err := queryStmt.Query(freelancerId)
    defer rows.Close()

    for rows.Next() {
        project := Project{}
        client := Client{}
        client.Projects = []Project{}
        client.Jobs = []Job{}

        if err := queryStmt.QueryRow(freelancerId).Scan(
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

        project.Freelancers = []Freelancer{}
        project.Client = client

        projects = append(projects, project)
    }

    return projects, nil
}

func (repo *FreelancerRepository) GetFreelancer(id int) (Freelancer, error) {
    freelancer := Freelancer{}

    queryStmt, err := repo.db.Prepare(`
        SELECT id,first_name,last_name,email,data,created
        FROM freelancers
        WHERE id = $1`)
    if err != nil {
        return freelancer, err
    }

    if err := queryStmt.QueryRow(id).Scan(
        &freelancer.Id,
        &freelancer.FirstName,
        &freelancer.LastName,
        &freelancer.Email,
        &freelancer._Data,
        &freelancer.Created,
    ); err != nil {
        return freelancer, err
    }

    if freelancer.Data, err = repo.getFreelancerData(freelancer._Data); err != nil {
        return freelancer, err
    }
    if freelancer.Projects, err = repo.getProjects(freelancer.Id); err != nil {
        return freelancer, err
    }

    return freelancer, nil
}

func (repo *FreelancerRepository) GetFreelancerByEmail(email string) (Freelancer, error) {
    freelancer := Freelancer{}

    queryStmt, err := repo.db.Prepare(`
        SELECT id,first_name,last_name,email,data,created
        FROM freelancers
        WHERE email = $1`)
    if err != nil {
        return freelancer, err
    }

    if err := queryStmt.QueryRow(email).Scan(
        &freelancer.Id,
        &freelancer.FirstName,
        &freelancer.LastName,
        &freelancer.Email,
        &freelancer._Data,
        &freelancer.Created,
    ); err != nil {
        return freelancer, err
    }

    if freelancer.Data, err = repo.getFreelancerData(freelancer._Data); err != nil {
        return freelancer, err
    }
    if freelancer.Projects, err = repo.getProjects(freelancer.Id); err != nil {
        return freelancer, err
    }

    return freelancer, nil
}

func (repo *FreelancerRepository) AddFreelancer(freelancer Freelancer) error {
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(freelancer.Password), bcrypt.DefaultCost)
    freelancer.Password = string(hashedPassword)

    var insertId int
    err := repo.db.QueryRow(`
            INSERT INTO freelancers(first_name,last_name,email,password,created)
            VALUES($1,$2,$3,$4,$5) returning id;`,
        freelancer.FirstName,
        freelancer.LastName,
        freelancer.Email,
        freelancer.Password,
        freelancer.Created,
    ).Scan(&insertId)

    if err != nil {
        return err
    }

    return nil
}

func (repo *FreelancerRepository) DeleteFreelancer(id string) error {
    queryStmt, err := repo.db.Prepare("DELETE FROM freelancers WHERE id = $1")
    if err != nil {
        return err
    }

    _, err = queryStmt.Exec(id)
    if err != nil {
        return err
    }

    return nil
}

func (repo *FreelancerRepository) CheckCredentials(email string, password string) error {
    queryStmt, err := repo.db.Prepare("SELECT f.password FROM freelancers f WHERE email = $1")
    if err != nil {
        return err
    }

    var foundPassword string
    if err = queryStmt.QueryRow(email).Scan(&foundPassword); err != nil {
        return err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(password)); err != nil {
        return errors.New("Freelancer not found (password is wrong)")
    }

    return nil
}

func (repo *FreelancerRepository) getFreelancerData(data string) (FreelancerData, error) {
    var freelancerData = FreelancerData{}
    json.Unmarshal([]byte(data), &freelancerData)

    for i, review := range freelancerData.Reviews {
        client, err := repo.getClient(review.ClientId)
        if err != nil {
            return freelancerData, err
        }
        freelancerData.Reviews[i].Client = client
    }

    return freelancerData, nil
}

func (repo *FreelancerRepository) getClient(clientId int) (Client, error) {
    client := Client{}

    queryStmt, err := repo.db.Prepare("SELECT id,name,description,created FROM clients WHERE id = $1")
    if err != nil {
        return client, err
    }

    if err := queryStmt.QueryRow(clientId).Scan(
        &client.Id,
        &client.Name,
        &client.Description,
        &client.Created,
    ); err != nil {
        return client, err
    }

    client.Jobs = []Job{}
    client.Projects = []Project{}

    return client, nil
}