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
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"-"`
	Email     string `json:"email" valid:"required,email" sql:"index" gorm:"unique"`
	Image     string `json:"image"`
}

type Freelancer struct {
	User
	Rating          float64          `json:"rating"`
	Timezone        string           `json:"timezone"`
	Skills          stringList       `json:"skills" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	IsAvailable     bool             `json:"isAvailable"`
	HourlyRateFrom  uint             `json:"hourlyRateFrom"`
	HourlyRateTo    uint             `json:"hourlyRateTo"`
	Projects        []Project        `json:"projects" gorm:"many2many:project_freelancers;"`
	Reviews         []Review         `json:"reviews"`
	References      []Reference      `json:"references"`
	JobApplications []JobApplication `json:"jobApplications"`
	About           string           `json:"about"`
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
	Timezone string    `json:"timezone"`
	Payment  string    `json:"payment"`
	Industry string    `json:"industry"`
	Rating   float64   `json:"rating"`
	Jobs     []Job     `json:"jobs"`
	Projects []Project `json:"projects"`
	Reviews  []Review  `json:"reviews"`
}

type Project struct {
	Model
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	Freelancers         []Freelancer `json:"freelancers,omitempty" gorm:"many2many:project_freelancers;"`
	ClientID            uint         `json:"-"`
	Client              *Client      `json:"client,omitempty"`
	Status              string       `json:"status"`
	Deadline            time.Time    `json:"deadline,omitempty"`
	DeadlineFlexibility int          `json:"deadlineFlexibility"`
	Hours               int          `json:"hours"`
	PerHour             float64      `json:"perHour"`
	Contract            *Contract    `json:"contract,omitempty"`
	ContractID          uint         `json:"-"`
}

type Contract struct {
	Model
	Hours               int         `json:"hours,omitempty"`
	PerHour             float64     `json:"perHour,omitempty"`
	Deadline            time.Time   `json:"deadline,omitempty"`
	DeadlineFlexibility int         `json:"deadlineFlexibility,omitempty"`
	Extensions          []Extension `json:"extensions,omitempty"`
}

type Extension struct {
	Model
	ContractID          uint      `json:"-"`
	Hours               int       `json:"hours"`
	PerHour             float64   `json:"perHour"`
	Deadline            time.Time `json:"deadline"`
	DeadlineFlexibility int       `json:"deadlineFlexibility"`
}

type Job struct {
	Model
	Name            string           `json:"name"`
	Summary         string           `json:"summary"`
	Details         string           `json:"details"`
	ClientID        uint             `json:"-"`
	Client          *Client          `json:"client,omitempty"`
	IsActive        bool             `json:"isActive"`
	Price           int              `json:"price"`
	StartDate       time.Time        `json:"startDate"`
	Deadline        time.Time        `json:"deadline"`
	Tags            stringList       `json:"tags" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	JobApplications []JobApplication `json:"jobApplications"`
	Attachments     []Attachment     `json:"attachments" gorm:"polymorphic:Owner;"`
	Examples        []Example        `json:"examples" gorm:"polymorphic:Owner;"`
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
