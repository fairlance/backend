package application

import (
	"database/sql/driver"
	"time"
)

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
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Freelancers          []Freelancer `json:"freelancers,omitempty" gorm:"many2many:project_freelancers;"`
	ClientID             uint         `json:"-"`
	Client               *Client      `json:"client,omitempty"`
	Status               string       `json:"status"`
	Contract             *Contract    `json:"contract,omitempty"`
	ContractID           uint         `json:"-"`
	ClientAgreed         bool         `json:"clientAgreed"`
	FreelancersAgreed    uintList     `json:"freelancersAgreed" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	ClientConcluded      bool         `json:"clientConcluded"`
	FreelancersConcluded uintList     `json:"freelancersConcluded" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}

func NewProject(job *Job, jobApplication *JobApplication) *Project {
	return &Project{
		Name:        job.Name,
		Description: job.Details,
		ClientID:    job.ClientID,
		Status:      projectStatusFinilazingTerms,
		Freelancers: []Freelancer{
			*jobApplication.Freelancer,
		},
		Contract: &Contract{
			Deadline: job.Deadline,
			Hours:    jobApplication.Hours,
			PerHour:  jobApplication.HourPrice,
		},
		ClientAgreed:         false,
		FreelancersAgreed:    []uint{},
		ClientConcluded:      false,
		FreelancersConcluded: []uint{},
	}
}
func (p *Project) allUsersConcluded() bool {
	if p.ClientConcluded && len(p.FreelancersConcluded) == len(p.Freelancers) {
		return true
	}

	return false
}

func (p *Project) canBeStarted() bool {
	if p.ClientAgreed && len(p.FreelancersAgreed) == len(p.Freelancers) {
		return true
	}

	if p.Contract.Proposal != nil {
		freelancersNotAgree := uintList{}
		for _, f := range p.Freelancers {
			if !contains(p.FreelancersAgreed, f.ID) {
				freelancersNotAgree = append(freelancersNotAgree, f.ID)
			}
		}

		if p.ClientAgreed &&
			len(freelancersNotAgree) == 1 &&
			p.Contract.Proposal.UserType == "freelancer" &&
			p.Contract.Proposal.UserID == freelancersNotAgree[0] {
			return true
		}

		if len(freelancersNotAgree) == 0 &&
			!p.ClientAgreed &&
			p.Contract.Proposal.UserType == "client" {
			return true
		}
	}

	return false
}

func (p *Project) mergeProposalToContract() {
	p.ClientAgreed = true
	p.FreelancersAgreed = uintList{}
	for _, f := range p.Freelancers {
		p.FreelancersAgreed = append(p.FreelancersAgreed, f.ID)
	}
	p.Contract.mergeProposalToContract()
}

type Contract struct {
	Model
	Hours               int         `json:"hours"`
	PerHour             float64     `json:"perHour"`
	Deadline            time.Time   `json:"deadline"`
	DeadlineFlexibility int         `json:"deadlineFlexibility"`
	Extensions          []Extension `json:"extensions,omitempty"`
	Proposal            *Proposal   `json:"proposal,omitempty" sql:"type:JSONB"`
}

func (c *Contract) mergeProposalToContract() {
	if c.Proposal != nil {
		c.PerHour = c.Proposal.PerHour
		c.Hours = c.Proposal.Hours
		c.Deadline = c.Proposal.Deadline
		c.DeadlineFlexibility = c.Proposal.DeadlineFlexibility
		c.Proposal = nil
	}
}

type Proposal struct {
	UserType            string    `json:"userType"`
	UserID              uint      `json:"userId"`
	Deadline            time.Time `json:"deadline"`
	DeadlineFlexibility int       `json:"deadlineFlexibility"`
	Hours               int       `json:"hours"`
	PerHour             float64   `json:"perHour"`
	Time                time.Time `json:"time"`
}

func (p *Proposal) Value() (driver.Value, error) {
	return pValue(p)
}

func (p *Proposal) Scan(src interface{}) error {
	return pScan(p, src)
}

type Extension struct {
	Model
	ContractID          uint      `json:"-"`
	Hours               int       `json:"hours"`
	PerHour             float64   `json:"perHour"`
	Deadline            time.Time `json:"deadline"`
	DeadlineFlexibility int       `json:"deadlineFlexibility"`
	ClientAgreed        bool      `json:"clientAgreed""`
	FreelancersAgreed   uintList  `json:"freelancersAgreed" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
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
	Message      string       `json:"message"`
	Milestones   stringList   `json:"milestones" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	HourPrice    float64      `json:"hourPrice"`
	Hours        int          `json:"hours"`
	Freelancer   *Freelancer  `json:"freelancer,omitempty"`
	FreelancerID uint         `json:"-"`
	JobID        uint         `json:"-"`
	Attachments  []Attachment `json:"attachments" gorm:"polymorphic:Owner;"`
	Examples     []Example    `json:"examples" gorm:"polymorphic:Owner;"`
}

type Review struct {
	Model
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	Rating       float64     `json:"rating"`
	JobID        uint        `json:"jobId"`
	ClientID     uint        `json:"-"`
	Client       *Client     `json:"client,omitempty"`
	FreelancerID uint        `json:"-"`
	Freelancer   *Freelancer `json:"freelancer,omitempty"`
}

type Reference struct {
	Model
	Title        string `json:"title"`
	Content      string `json:"content"`
	Media        Media  `json:"media,omitempty"` // todo: should be a pointer
	FreelancerID uint   `json:"freelancerId"`
}

type Media struct {
	Model
	Image       string `json:"image"`
	Video       string `json:"video"`
	ReferenceID uint   `json:"referenceId"`
}

type Attachment struct {
	Model
	Name      string `json:"name"`
	URL       string `json:"url"`
	OwnerID   int    `json:"-"`
	OwnerType string `json:"-"`
}

type Example struct {
	Model
	URL         string `json:"url"`
	Description string `json:"description"`
	OwnerID     int    `json:"-"`
	OwnerType   string `json:"-"`
}
