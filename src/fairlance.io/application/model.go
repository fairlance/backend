package application

import "time"

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" sql:"index"`
}

type User struct {
	Model
	FirstName string `json:"firstName,omitempty" valid:"required"`
	LastName  string `json:"lastName,omitempty" valid:"required"`
	Password  string `json:"-" valid:"required"`
	Email     string `json:"email,omitempty" valid:"required,email" sql:"index" gorm:"unique"`
}

type Freelancer struct {
	User
	Rating          float64          `json:"rating,omitempty"`
	Timezone        string           `json:"timezone,omitempty"`
	Skills          stringList       `json:"skills,omitempty" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	IsAvailable     bool             `json:"isAvailable,omitempty"`
	HourlyRateFrom  uint             `json:"hourlyRateFrom,omitempty"`
	HourlyRateTo    uint             `json:"hourlyRateTo,omitempty"`
	Projects        []Project        `json:"projects,omitempty" gorm:"many2many:project_freelancers;"`
	Reviews         []Review         `json:"reviews,omitempty"`
	References      []Reference      `json:"references,omitempty"`
	JobApplications []JobApplication `json:"jobApplications,omitempty"`
	About           string           `json:"about,omitempty"`
}

type FreelancerUpdate struct {
	Skills         stringList `json:"skills" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	Timezone       string     `json:"timezone" valid:"required"`
	IsAvailable    bool       `json:"isAvailable"`
	HourlyRateFrom uint       `json:"hourlyRateFrom" valid:"required"`
	HourlyRateTo   uint       `json:"hourlyRateTo" valid:"required"`
}

type Client struct {
	User
	Timezone string    `json:"timezone,omitempty"`
	Payment  string    `json:"payment,omitempty"`
	Industry string    `json:"industry,omitempty"`
	Rating   float64   `json:"rating,omitempty"`
	Jobs     []Job     `json:"jobs,omitempty"`
	Projects []Project `json:"projects,omitempty"`
	Reviews  []Review  `json:"reviews,omitempty"`
}

type Project struct {
	Model
	Name        string       `json:"name,omitempty" valid:"required"`
	Description string       `json:"description,omitempty" valid:"required"`
	Freelancers []Freelancer `json:"freelancers,omitempty" gorm:"many2many:project_freelancers;"`
	ClientID    uint         `json:"-" valid:"required"`
	Client      *Client      `json:"client,omitempty"`
	Status      string       `json:"status,omitempty"`
	DueDate     time.Time    `json:"dueDate,omitempty"`
}

type Job struct {
	Model
	Name            string           `json:"name,omitempty" valid:"required"`
	Summary         string           `json:"summary,omitempty" valid:"required"`
	Details         string           `json:"details,omitempty" valid:"required"`
	ClientID        uint             `json:"-" valid:"required"`
	Client          *Client          `json:"client,omitempty"`
	IsActive        bool             `json:"isActive,omitempty"`
	Price           int              `json:"price,omitempty"`
	StartDate       time.Time        `json:"startDate,omitempty"`
	Deadline        time.Time        `json:"deadline,omitempty"`
	Tags            stringList       `json:"tags,omitempty" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	JobApplications []JobApplication `json:"jobApplications,omitempty"`
	Attachments     []Attachment     `json:"attachments,omitempty" gorm:"polymorphic:Owner;"`
	Examples        []Example        `json:"examples,omitempty" gorm:"polymorphic:Owner;"`
}

type Review struct {
	Model
	Title        string  `json:"title,omitempty" valid:"required"`
	Content      string  `json:"content,omitempty"`
	Rating       float64 `json:"rating,omitempty" valid:"required"`
	JobID        uint    `json:"jobId,omitempty" valid:"required"`
	ClientID     uint    `json:"clientId,omitempty" valid:"required"`
	FreelancerID uint    `json:"freelancerId,omitempty"` //should be userID
}

type Reference struct {
	Model
	Title        string `json:"title,omitempty" valid:"required"`
	Content      string `json:"content,omitempty"`
	Media        Media  `json:"media,omitempty"`
	FreelancerID uint   `json:"freelancerId,omitempty"`
}

type Media struct {
	Model
	Image       string `json:"image,omitempty"`
	Video       string `json:"video,omitempty"`
	ReferenceID uint   `json:"referenceId,omitempty"`
}

type JobApplication struct {
	Model
	Message          string       `json:"message,omitempty" valid:"required"`
	Samples          uintList     `json:"samples,omitempty" valid:"required" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	DeliveryEstimate int          `json:"deliveryEstimate,omitempty" valid:"required"`
	Milestones       stringList   `json:"milestones,omitempty" valid:"required" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	HourPrice        float64      `json:"hourPrice,omitempty" valid:"required"`
	Hours            int          `json:"hours,omitempty" valid:"required"`
	FreelancerID     uint         `json:"freelancerId,omitempty" valid:"required"`
	JobID            uint         `json:"-"`
	Attachments      []Attachment `json:"attachments,omitempty" gorm:"polymorphic:Owner;"`
	Examples         []Example    `json:"examples,omitempty" gorm:"polymorphic:Owner;"`
}

type Attachment struct {
	Model
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	OwnerId   int    `json:"-"`
	OwnerType string `json:"-"`
}

type Example struct {
	Model
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerId     int    `json:"-"`
	OwnerType   string `json:"-"`
}
