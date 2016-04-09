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
	db.DropTableIfExists(&Freelancer{}, &Project{}, &Client{}, &Job{}, &Review{})
	db.CreateTable(&Freelancer{}, &Project{}, &Client{}, &Job{}, &Review{})

	ac.FreelancerRepository.AddFreelancer(NewFreelancer("First", "Last", "Dev", "Pass", "first@mail.com", 3, 55, "UTC"))

	ac.FreelancerRepository.AddReview(Review{
		Title:        "text2",
		Content:      "content",
		Rating:       4.1,
		ClientId:     1,
		FreelancerId: 1,
	})
	ac.FreelancerRepository.AddReference(1, Reference{"title", "content", Media{"image", "video"}})
	ac.FreelancerRepository.AddFreelancer(NewFreelancer(
		"Milos",
		"Krsmanovic",
		"Dev",
		"$2a$10$VJ8H9EYOIj9mnyW5mUm/nOWUrz/Rkak4/Ov3Lnw1GsAm4gmYU6sQu",
		"milos@gmail.com",
		12,
		22,
		"CET",
	))

	ac.FreelancerRepository.AddReview(Review{
		Title:        "text2",
		Content:      "content",
		Rating:       4.1,
		ClientId:     1,
		FreelancerId: 2,
	})

	ac.FreelancerRepository.AddReview(Review{
		Title:        "text2",
		Content:      "content",
		Rating:       2.4,
		ClientId:     2,
		FreelancerId: 2,
	})
	ac.FreelancerRepository.AddReference(2, Reference{"title", "content", Media{"image", "video"}})

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

	ac.FreelancerRepository.AddFreelancer(NewFreelancer("Third", "Last", "Dev", "Pass", "third@mail.com", 3, 55, "UTC"))
	ac.FreelancerRepository.DeleteFreelancer(3)
}
