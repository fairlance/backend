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
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Password  string `json:"-"`
	Email     string `json:"email,omitempty" valid:"required,email" sql:"index" gorm:"unique"`
	Image     string `json:"image,omitempty"`
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
	Timezone       string     `json:"timezone"`
	IsAvailable    bool       `json:"isAvailable"`
	HourlyRateFrom uint       `json:"hourlyRateFrom"`
	HourlyRateTo   uint       `json:"hourlyRateTo"`
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
	Name                string       `json:"name,omitempty"`
	Description         string       `json:"description,omitempty"`
	Freelancers         []Freelancer `json:"freelancers,omitempty" gorm:"many2many:project_freelancers;"`
	ClientID            uint         `json:"-"`
	Client              *Client      `json:"client,omitempty"`
	Status              string       `json:"status,omitempty"`
	Deadline            time.Time    `json:"deadline,omitempty"`
	DeadlineFlexibility int          `json:"deadlineFlexibility,omitempty"`
	WorkhoursPerDay     int          `json:"workhoursPerDay,omitempty"`
	PerHour             float64      `json:"perHour,omitempty"`
	Contract            *Contract    `json:"contract,omitempty"`
	ContractID          uint         `json:"-"`
}

type Contract struct {
	Model
	WorkhoursPerDay     int         `json:"workhoursPerDay,omitempty"`
	PerHour             float64     `json:"perHour,omitempty"`
	Deadline            time.Time   `json:"deadline,omitempty"`
	DeadlineFlexibility int         `json:"deadlineFlexibility,omitempty"`
	Extensions          []Extension `json:"extensions,omitempty"`
}

type Extension struct {
	Model
	ContractID          uint      `json:"-"`
	WorkhoursPerDay     int       `json:"workhoursPerDay,omitempty"`
	PerHour             float64   `json:"perHour,omitempty"`
	Deadline            time.Time `json:"deadline,omitempty"`
	DeadlineFlexibility int       `json:"deadlineFlexibility,omitempty"`
}

type Job struct {
	Model
	Name            string           `json:"name,omitempty"`
	Summary         string           `json:"summary,omitempty"`
	Details         string           `json:"details,omitempty"`
	ClientID        uint             `json:"-"`
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

type JobApplication struct {
	Model
	Message          string       `json:"message,omitempty"`
	DeliveryEstimate int          `json:"deliveryEstimate,omitempty"`
	Milestones       stringList   `json:"milestones,omitempty" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	HourPrice        float64      `json:"hourPrice,omitempty"`
	Hours            int          `json:"hours,omitempty"`
	Freelancer       *Freelancer  `json:"freelancer,omitempty"`
	FreelancerID     uint         `json:"-"`
	JobID            uint         `json:"-"`
	Attachments      []Attachment `json:"attachments,omitempty" gorm:"polymorphic:Owner;"`
	Examples         []Example    `json:"examples,omitempty" gorm:"polymorphic:Owner;"`
}

type Review struct {
	Model
	Title        string  `json:"title,omitempty"`
	Content      string  `json:"content,omitempty"`
	Rating       float64 `json:"rating,omitempty"`
	JobID        uint    `json:"jobId,omitempty"`
	ClientID     uint    `json:"clientId,omitempty"`
	FreelancerID uint    `json:"freelancerId,omitempty"`
}

type Reference struct {
	Model
	Title        string `json:"title,omitempty"`
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

type Attachment struct {
	Model
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	OwnerID   int    `json:"-"`
	OwnerType string `json:"-"`
}

type Example struct {
	Model
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerID     int    `json:"-"`
	OwnerType   string `json:"-"`
}
