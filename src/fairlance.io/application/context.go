package application

import (
	"encoding/json"

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

	freelancerRepository, _ := NewFreelancerRepository(db)
	projectRepository, _ := NewProjectRepository(db)
	clientRepository, _ := NewClientRepository(db)

	context := &ApplicationContext{
		FreelancerRepository: freelancerRepository,
		ProjectRepository:    projectRepository,
		ClientRepository:     clientRepository,
		JwtSecret:            "fairlance", //base64.StdEncoding.EncodeToString([]byte("fairlance")),
	}

	context.prepareTables(db)

	return context, nil
}

func (ac *ApplicationContext) prepareTables(db *gorm.DB) {
	db.DropTableIfExists(&Freelancer{}, &Project{}, &Client{}, &Job{})
	db.CreateTable(&Freelancer{}, &Project{}, &Client{}, &Job{})

	db.Create(&Freelancer{
		FirstName:      "First",
		LastName:       "Last",
		Title:          "Dev",
		Password:       "Pass",
		Email:          "first@mail.com",
		JsonComments:   `[]`,
		JsonReferences: `[]`,
	})

	freelancer, _ := ac.FreelancerRepository.GetFreelancerByEmail("first@mail.com")
	js, _ := json.Marshal(append(freelancer.Comments, Comment{"text2", 1}))
	freelancer.JsonComments = string(js)
	js, _ = json.Marshal(append(freelancer.References, Reference{"title", "content", Media{"image", "video"}}))
	freelancer.JsonReferences = string(js)
	ac.FreelancerRepository.UpdateFreelancer(&freelancer)

	db.Create(&Freelancer{
		FirstName:      "Milos",
		LastName:       "Krsmanovic",
		Title:          "Dev",
		Password:       "$2a$10$VJ8H9EYOIj9mnyW5mUm/nOWUrz/Rkak4/Ov3Lnw1GsAm4gmYU6sQu",
		Email:          "milos@gmail.com",
		JsonComments:   `[]`,
		JsonReferences: `[]`,
	})

	milos, _ := ac.FreelancerRepository.GetFreelancerByEmail("milos@gmail.com")
	js, _ = json.Marshal(append(milos.Comments, Comment{"text2", 1}))
	milos.JsonComments = string(js)
	js, _ = json.Marshal(append(milos.References, Reference{"title", "content", Media{"image", "video"}}))
	milos.JsonReferences = string(js)
	ac.FreelancerRepository.UpdateFreelancer(&milos)

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

	ac.FreelancerRepository.DeleteFreelancer(1)
}
