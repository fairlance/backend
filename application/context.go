package application

import (
	"time"

	"github.com/fairlance/backend/dispatcher"
	"github.com/fairlance/backend/mailer"
	"github.com/jinzhu/gorm"
)

type ApplicationContext struct {
	db                     *gorm.DB
	FreelancerRepository   FreelancerRepository
	ProjectRepository      ProjectRepository
	ClientRepository       ClientRepository
	ReferenceRepository    ReferenceRepository
	JobRepository          JobRepository
	UserRepository         UserRepository
	JwtSecret              string
	NotificationDispatcher *NotificationDispatcher
	MessagingDispatcher    *MessagingDispatcher
	PaymentDispatcher      *PaymentDispatcher
	Indexer                Indexer
	Mailer                 mailer.Mailer
}

type ContextOptions struct {
	DbHost          string
	DbName          string
	DbUser          string
	DbPass          string
	Secret          string
	NotificationURL string
	MessagingURL    string
	PaymentURL      string
	SearcherURL     string
	MailerOptions   mailer.Options
}

func NewContext(options ContextOptions) (*ApplicationContext, error) {
	db, err := gorm.Open("postgres", "host="+options.DbHost+" user="+options.DbUser+" password="+options.DbPass+" dbname="+options.DbName+" sslmode=disable")
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
		db:                     db,
		UserRepository:         userRepository,
		FreelancerRepository:   freelancerRepository,
		ClientRepository:       clientRepository,
		JobRepository:          jobRepository,
		ProjectRepository:      projectRepository,
		ReferenceRepository:    referenceRepository,
		JwtSecret:              options.Secret, //base64.StdEncoding.EncodeToString([]byte(options.Secret)),
		NotificationDispatcher: NewNotificationDispatcher(dispatcher.NewNotifications(options.NotificationURL)),
		MessagingDispatcher:    NewMessagingDispatcher(dispatcher.NewMessaging(options.MessagingURL)),
		PaymentDispatcher:      NewPaymentDispatcher(dispatcher.NewPayment(options.PaymentURL)),
		Indexer:                NewHTTPIndexer(options.SearcherURL),
		Mailer:                 mailer.NewMailgun(options.MailerOptions),
	}

	return context, nil
}

func (ac *ApplicationContext) DropCreateFillTables() {
	ac.DropTables()
	ac.CreateTables()
	ac.FillTables()
}

func (ac *ApplicationContext) DropTables() {
	ac.db.DropTableIfExists(&Freelancer{}, &Extension{}, &Contract{}, &Project{}, &Client{}, &Job{}, &Review{}, &Reference{}, &Media{}, &JobApplication{}, &File{})
	ac.db.DropTableIfExists("project_freelancers") //, "job_applications")
}

func (ac *ApplicationContext) CreateTables() {
	ac.db.AutoMigrate(&Freelancer{}, &Extension{}, &Contract{}, &Project{}, &Client{}, &Job{}, &Review{}, &Reference{}, &Media{}, &JobApplication{}, &File{})
}

func (ac *ApplicationContext) FillTables() {
	f1 := &Freelancer{
		User: User{
			FirstName: "First",
			LastName:  "Last",
			Password:  "Pass",
			Email:     "first@mail.com",
		},
		Skills:   stringList{"man", "dude", "boyyyy"},
		Timezone: "UTC",
		About:    "Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.",
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
		Timezone: "CET",
		Skills:   stringList{"good", "bad", "ugly"},
		About:    "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
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
		Name:        "Project Done",
		Description: "Description Done",
		ClientID:    1,
		Status:      projectStatusDone,
		Contract: &Contract{
			Deadline:            time.Now().Add(24 * time.Hour * 2),
			PerHour:             5,
			Hours:               6,
			DeadlineFlexibility: 1,
		},
		ClientAgreed:      false,
		FreelancersAgreed: []uint{},
	}).Association("Freelancers").Append([]Freelancer{*f1, *f2})

	ac.db.Create(&Project{
		Name:        "Project Finalizing Terms",
		Description: "Description Finalizing Terms",
		ClientID:    2,
		Status:      projectStatusFinalizingTerms,
		Contract: &Contract{
			Deadline:            time.Now().Add(24 * time.Hour * 5),
			PerHour:             15,
			Hours:               6,
			DeadlineFlexibility: 1,
			Extensions: []Extension{Extension{
				Deadline:            time.Now().Add(24 * time.Hour * 6),
				PerHour:             15,
				Hours:               4,
				DeadlineFlexibility: 0,
				ClientAgreed:        true,
				FreelancersAgreed:   []uint{},
			}},
		},
		FreelancersAgreed: []uint{},
	}).Association("Freelancers").Replace([]Freelancer{*f1})

	ac.db.Create(&Project{
		Name:        "Project In Progress",
		Description: "Description In Progress",
		ClientID:    1,
		Status:      projectStatusInProgress,
		Contract: &Contract{
			Deadline:            time.Now().Add(time.Hour),
			DeadlineFlexibility: 1,
			Hours:               4,
			PerHour:             9.5,
		},
	}).Association("Freelancers").Append([]Freelancer{*f2})

	jobApplication := ac.db.Create(&JobApplication{
		Message:      "I apply",
		JobID:        1,
		FreelancerID: 1,
		HourPrice:    9.5,
		Hours:        6,
	})
	jobApplication.Association("Attachments").Replace([]File{{Name: "job application attachment", URL: "www.google.com"}})
	jobApplication.Association("Examples").Replace([]File{{Name: "Some job application example", URL: "www.google.com"}})

	ac.ClientRepository.AddClient(&Client{
		User: User{
			FirstName: "Clint",
			LastName:  "Clienter",
			Password:  "123456",
			Email:     "client@mail.com",
		},
		Timezone: "UTC",
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
	job.Association("Attachments").Replace([]File{{Name: "job attachment", URL: "www.google.com"}})
	job.Association("Examples").Replace([]File{{Name: "Some job example", URL: "www.google.com"}})

	ac.FreelancerRepository.AddFreelancer(&Freelancer{
		User: User{
			FirstName: "Third",
			LastName:  "Last",
			Password:  "Pass",
			Email:     "third@mail.com",
		},
		Timezone: "UTC",
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
