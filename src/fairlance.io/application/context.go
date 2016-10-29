package application

import (
	"github.com/jinzhu/gorm"
)

type ApplicationContext struct {
	db                   *gorm.DB
	FreelancerRepository *FreelancerRepository
	ProjectRepository    ProjectRepository
	ClientRepository     *ClientRepository
	ReferenceRepository  ReferenceRepository
	JobRepository        JobRepository // todo: think of a better name
	UserRepository       *UserRepository
	JwtSecret            string
}

type ContextOptions struct {
	DbName string
	DbUser string
	DbPass string
	Secret string
}

func NewContext(options ContextOptions) (*ApplicationContext, error) {
	db, err := gorm.Open("postgres", "user="+options.DbUser+" password="+options.DbPass+" dbname="+options.DbName+" sslmode=disable")
	if err != nil {
		return nil, err
	}

	userRepository, _ := NewUserRepository(db)
	freelancerRepository, _ := NewFreelancerRepository(db)
	clientRepository, _ := NewClientRepository(db)
	jobRepository, _ := NewJobRepository(db)
	projectRepository, _ := NewProjectRepository(db)
	referenceRepository, _ := NewReferenceRepository(db)

	context := &ApplicationContext{
		db:                   db,
		UserRepository:       userRepository,
		FreelancerRepository: freelancerRepository,
		ClientRepository:     clientRepository,
		JobRepository:        jobRepository,
		ProjectRepository:    projectRepository,
		ReferenceRepository:  referenceRepository,
		JwtSecret:            options.Secret, //base64.StdEncoding.EncodeToString([]byte(options.Secret)),
	}

	return context, nil
}

func (ac *ApplicationContext) DropCreateFillTables() {
	ac.DropTables()
	ac.CreateTables()
	ac.FillTables()
}

func (ac *ApplicationContext) DropTables() {
	ac.db.DropTableIfExists(&Freelancer{}, &Project{}, &Client{}, &Job{}, &Review{}, &Reference{}, &Media{}, &JobApplication{})
}

func (ac *ApplicationContext) CreateTables() {
	ac.db.CreateTable(&Freelancer{}, &Project{}, &Client{}, &Job{}, &Review{}, &Reference{}, &Media{}, &JobApplication{})
}

func (ac *ApplicationContext) FillTables() {
	ac.FreelancerRepository.AddFreelancer(&Freelancer{
		User: User{
			FirstName: "First",
			LastName:  "Last",
			Password:  "Pass",
			Email:     "first@mail.com",
		},
		HourlyRateFrom: 3,
		HourlyRateTo:   55,
		Timezone:       "UTC",
	})

	ac.FreelancerRepository.AddReview(&Review{
		Title:        "text2",
		Content:      "content",
		Rating:       4.1,
		ClientID:     1,
		FreelancerID: 1,
	})

	ac.ReferenceRepository.AddReference(&Reference{
		Title:        "title",
		Content:      "content",
		Media:        Media{Image: "image", Video: "video"},
		FreelancerID: 1,
	})

	ac.FreelancerRepository.AddFreelancer(&Freelancer{
		User: User{
			FirstName: "Pera",
			LastName:  "Peric",
			Password:  "123456",
			Email:     "second@mail.com",
		},
		HourlyRateFrom: 12,
		HourlyRateTo:   22,
		Timezone:       "CET",
		Skills:         strings{"good", "bad", "ugly"},
	})

	ac.FreelancerRepository.AddReview(&Review{
		Title:        "text2",
		Content:      "content",
		Rating:       4.1,
		JobID:        1,
		ClientID:     1,
		FreelancerID: 2,
	})

	ac.FreelancerRepository.AddReview(&Review{
		Title:        "text2",
		Content:      "content",
		Rating:       2.4,
		JobID:        2,
		ClientID:     1,
		FreelancerID: 2,
	})
	ac.ReferenceRepository.AddReference(&Reference{
		Title:        "title",
		Content:      "content",
		Media:        Media{Image: "image", Video: "video"},
		FreelancerID: 2,
	})

	ac.db.Create(&Project{
		Name:        "Project",
		Description: "Description",
		ClientID:    1,
		IsActive:    true,
	})

	ac.db.Create(&JobApplication{
		Message:          "I apply",
		JobID:            1,
		FreelancerID:     1,
		Milestones:       strings{"Milestone1", "Milestone2"},
		Samples:          uints{1, 2},
		DeliveryEstimate: 15,
	})

	ac.ClientRepository.AddClient(&Client{
		User: User{
			FirstName: "ClientName",
			LastName:  "ClientLast",
			Password:  "123456",
			Email:     "client@mail.com",
		},
		Timezone: "UTC",
		Payment:  "PayPal",
		Industry: "Dildos",
		Rating:   4.6,
	})

	ac.db.Create(&Job{
		Name:     "Job",
		Summary:  "Summary Job",
		Details:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		ClientID: 1,
		Tags:     strings{"tag"},
		Links:    strings{"http://www.google.com/"},
	})

	ac.FreelancerRepository.AddFreelancer(&Freelancer{
		User: User{
			FirstName: "Third",
			LastName:  "Last",
			Password:  "Pass",
			Email:     "third@mail.com",
		},
		HourlyRateFrom: 3,
		HourlyRateTo:   55,
		Timezone:       "UTC",
	})

	ac.FreelancerRepository.AddReview(&Review{
		Title:        "deleted text2",
		Content:      "deleted content",
		Rating:       2.4,
		JobID:        2,
		ClientID:     1,
		FreelancerID: 3,
	})
	ac.ReferenceRepository.AddReference(&Reference{
		Title:        "deleted title",
		Content:      "deleted content",
		Media:        Media{Image: "deleted media image", Video: "deleted media video"},
		FreelancerID: 3,
	})
	ac.FreelancerRepository.DeleteFreelancer(3)
}
