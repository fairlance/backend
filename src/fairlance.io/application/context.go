package application

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ApplicationContext struct {
	FreelancerRepository *FreelancerRepository
	ProjectRepository    *ProjectRepository
	ClientRepository     *ClientRepository
	JwtSecret            string
}

func NewContext(dbName string) (*ApplicationContext, error) {
	db, err := gorm.Open("postgres", "user=fairlance password=fairlance dbname=application sslmode=disable")
	if err != nil {
		return nil, err
	}
	prepareTables(db)

	freelancerRepository, _ := NewFreelancerRepository(db)
	projectRepository, _ := NewProjectRepository(db)
	clientRepository, _ := NewClientRepository(db)

	context := &ApplicationContext{
		FreelancerRepository: freelancerRepository,
		ProjectRepository:    projectRepository,
		ClientRepository:     clientRepository,
		JwtSecret:            "fairlance", //base64.StdEncoding.EncodeToString([]byte("fairlance")),
	}

	return context, nil
}

func prepareTables(db *gorm.DB) {
	db.DropTableIfExists(&Freelancer{})
	db.DropTableIfExists(&Project{})
	db.DropTableIfExists(&Client{})
	db.DropTableIfExists(&Job{})

	db.CreateTable(&Freelancer{})
	db.CreateTable(&Project{})
	db.CreateTable(&Client{})
	db.CreateTable(&Job{})

	db.Create(&Freelancer{
		FirstName: "First",
		LastName:  "Last",
		Title:     "Dev",
		Password:  "Pass",
		Email:     "first@mail.com",
	})

	db.Create(&Project{
		Name:        "Project",
		Description: "Description",
		ClientId:    1,
		IsActive:    true,
	})

	db.Create(&Client{
		Name:        "Client",
		Description: "Desc Client",
	})

	db.Create(&Job{
		Name:        "Client",
		Description: "Desc Client",
		ClientId:    1,
	})
}
