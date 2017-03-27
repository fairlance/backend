package application

import (
	"time"

	"fairlance.io/notifier"
	"github.com/jinzhu/gorm"
)

type ApplicationContext struct {
	db                   *gorm.DB
	FreelancerRepository FreelancerRepository
	ProjectRepository    ProjectRepository
	ClientRepository     ClientRepository
	ReferenceRepository  ReferenceRepository
	JobRepository        JobRepository
	UserRepository       UserRepository
	JwtSecret            string
	Notifier             notifier.Notifier
}

type ContextOptions struct {
	DbName          string
	DbUser          string
	DbPass          string
	Secret          string
	NotificationURL string
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
		Notifier:             notifier.NewHTTPNotifier(options.NotificationURL),
	}

	return context, nil
}

func (ac *ApplicationContext) DropCreateFillTables() {
	ac.DropTables()
	ac.CreateTables()
	ac.FillTables()
}

func (ac *ApplicationContext) DropTables() {
	ac.db.DropTableIfExists(&Freelancer{}, &Project{}, &Client{}, &Job{}, &Review{}, &Reference{}, &Media{}, &JobApplication{}, &Attachment{}, &Example{})
	ac.db.DropTable("project_freelancers", "job_applications")
}

func (ac *ApplicationContext) CreateTables() {
	ac.db.CreateTable(&Freelancer{}, &Project{}, &Client{}, &Job{}, &Review{}, &Reference{}, &Media{}, &JobApplication{}, &Attachment{}, &Example{})
}

func (ac *ApplicationContext) FillTables() {
	f1 := &Freelancer{
		User: User{
			FirstName: "First",
			LastName:  "Last",
			Password:  "Pass",
			Email:     "first@mail.com",
		},
		HourlyRateFrom: 3,
		HourlyRateTo:   55,
		Skills:         stringList{"man", "dude", "boyyyy"},
		Timezone:       "UTC",
		About:          "Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.",
	}
	ac.FreelancerRepository.AddFreelancer(f1)

	ac.FreelancerRepository.AddReview(1, &Review{
		Title:    "text2",
		Content:  "content",
		Rating:   4.3,
		ClientID: 1,
		JobID:    4,
	})

	ac.ReferenceRepository.AddReference(1, &Reference{
		Title:   "title",
		Content: "content",
		Media:   Media{Image: "image", Video: "video"},
	})

	f2 := &Freelancer{
		User: User{
			FirstName: "Pera",
			LastName:  "Peric",
			Password:  "123456",
			Email:     "second@mail.com",
		},
		HourlyRateFrom: 12,
		HourlyRateTo:   22,
		Timezone:       "CET",
		Skills:         stringList{"good", "bad", "ugly"},
		About:          "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
	}
	ac.FreelancerRepository.AddFreelancer(f2)

	ac.FreelancerRepository.AddReview(2, &Review{
		Title:    "text2",
		Content:  "content",
		Rating:   4.1,
		JobID:    1,
		ClientID: 1,
	})

	ac.FreelancerRepository.AddReview(2, &Review{
		Title:    "text2",
		Content:  "content",
		Rating:   2.7,
		JobID:    2,
		ClientID: 1,
	})
	ac.ReferenceRepository.AddReference(2, &Reference{
		Title:   "title",
		Content: "content",
		Media:   Media{Image: "image", Video: "video"},
	})

	ac.db.Create(&Project{
		Name:        "Project Pending",
		Description: "Description Pending",
		ClientID:    1,
		Status:      projectStatusPending,
		Deadline:    time.Now().Add(time.Duration(21) * time.Hour),
	}).Association("Freelancers").Append([]Freelancer{*f1, *f2})

	ac.db.Create(&Project{
		Name:            "Project Finilazing Terms",
		Description:     "Description Finilazing Terms",
		ClientID:        2,
		Status:          projectStatusFinilazingTerms,
		Deadline:        time.Now().Add(time.Duration(5*24+1) * time.Hour),
		WorkhoursPerDay: 6,
		PerHour:         8,
	}).Association("Freelancers").Replace([]Freelancer{*f1}) // just to see that there is a replace as well

	ac.db.Create(&Project{
		Name:                "Project Working",
		Description:         "Description Working",
		ClientID:            1,
		Status:              projectStatusWorking,
		Deadline:            time.Now().Add(time.Hour),
		DeadlineFlexibility: 1,
		WorkhoursPerDay:     4,
		PerHour:             9.5,
	}).Association("Freelancers").Append([]Freelancer{*f2})

	ac.db.Create(&Project{
		Name:                "Project Archived",
		Description:         "Description Archived",
		ClientID:            1,
		Status:              projectStatusArchived,
		Deadline:            time.Now().Add(time.Hour),
		DeadlineFlexibility: 4,
		WorkhoursPerDay:     2,
		PerHour:             10,
	}).Association("Freelancers").Append([]Freelancer{*f2})

	ac.db.Create(&Project{
		Name:            "Project Canceled",
		Description:     "Description Canceled",
		ClientID:        2,
		Status:          projectStatusCanceled,
		Deadline:        time.Now().Add(time.Hour),
		WorkhoursPerDay: 8,
		PerHour:         7.85,
	}).Association("Freelancers").Append([]Freelancer{*f1})

	jobApplication := ac.db.Create(&JobApplication{
		Message:          "I apply",
		JobID:            1,
		FreelancerID:     1,
		Milestones:       stringList{"Milestone1", "Milestone2"},
		DeliveryEstimate: 15,
		HourPrice:        9.5,
		Hours:            6,
	})
	jobApplication.Association("Attachments").Replace([]Attachment{{Name: "job application attachment", URL: "www.google.com"}})
	jobApplication.Association("Examples").Replace([]Example{{Description: "Some job application example", URL: "www.google.com"}})

	ac.ClientRepository.AddClient(&Client{
		User: User{
			FirstName: "Clint",
			LastName:  "Clienter",
			Password:  "123456",
			Email:     "client@mail.com",
		},
		Timezone: "UTC",
		Payment:  "PayPal",
		Industry: "Dildos",
		Rating:   4.6,
	})

	ac.ClientRepository.AddClient(&Client{
		User: User{
			FirstName: "Clientoni",
			LastName:  "Clientello",
			Password:  "654321",
			Email:     "clientoni@mail.com",
		},
		Timezone: "UTC",
		Payment:  "Skrill",
		Industry: "Handcufs",
		Rating:   2.4,
	})

	job := ac.db.Create(&Job{
		Name:     "Job",
		Summary:  "Summary Job",
		Details:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		ClientID: 1,
		Tags:     stringList{"tag"},
		Deadline: time.Now().Add(time.Hour * 24 * 5),
	})
	job.Association("Attachments").Replace([]Attachment{{Name: "job attachment", URL: "www.google.com"}})
	job.Association("Examples").Replace([]Example{{Description: "Some job example", URL: "www.google.com"}})

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

	ac.FreelancerRepository.AddReview(3, &Review{
		Title:    "deleted text2",
		Content:  "deleted content",
		Rating:   2.4,
		JobID:    2,
		ClientID: 1,
	})
	ac.ReferenceRepository.AddReference(3, &Reference{
		Title:   "deleted title",
		Content: "deleted content",
		Media:   Media{Image: "deleted media image", Video: "deleted media video"},
	})
	ac.FreelancerRepository.DeleteFreelancer(3)
}
