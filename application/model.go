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
	LastLogin        *time.Time `json:"lastLogin"`
	ProfileCompleted bool       `json:"profileCompleted"`
	FirstName        string     `json:"firstName"`
	LastName         string     `json:"lastName"`
	Password         string     `json:"-"`
	Email            string     `json:"email" valid:"required,email" sql:"index" gorm:"unique"`
	Image            string     `json:"image"`
	Salutation       string     `json:"salutation"`
	IsCompany        bool       `json:"isCompany"`
	CompanyName      string     `json:"companyName"`
}

type Freelancer struct {
	User
	Rating          float64          `json:"rating"`
	Timezone        string           `json:"timezone"`
	Skills          stringList       `json:"skills" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	Projects        []Project        `json:"projects" gorm:"many2many:project_freelancers;"`
	Reviews         []Review         `json:"reviews"`
	References      []Reference      `json:"references"`
	JobApplications []JobApplication `json:"jobApplications"`
	About           string           `json:"about"`
	PayPalEmail     string           `json:"payPalEmail"`
	Phone           string           `json:"phone"`
	AdditionalFiles []File           `json:"additionalFiles" gorm:"polymorphic:Owner;"`
	PortfolioItems  []File           `json:"portfolioItems" gorm:"polymorphic:Owner;"`
	PortfolioLinks  []File           `json:"portfolioLinks" gorm:"polymorphic:Owner;"`
	Birthdate       string           `json:"birthdate"`
}

// BeforeSave updates file types before saving
func (freelancer *Freelancer) BeforeSave() error {
	for i := range freelancer.AdditionalFiles {
		freelancer.AdditionalFiles[i].Type = fileTypeFreelancerAdditionalField
	}
	for i := range freelancer.PortfolioItems {
		freelancer.PortfolioItems[i].Type = fileTypeFreelancerPortfolioItems
	}
	for i := range freelancer.PortfolioLinks {
		freelancer.PortfolioLinks[i].Type = fileTypeFreelancerPortfolioLinks
	}
	return nil
}

type FreelancerUpdate struct {
	Image           string     `json:"image" valid:"required"`
	About           string     `json:"about" valid:"required"`
	Timezone        string     `json:"timezone" valid:"required"`
	PayPalEmail     string     `json:"payPalEmail" valid:"required"`
	Phone           string     `json:"phone"`
	Skills          stringList `json:"skills" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	AdditionalFiles []File     `json:"additionalFiles"`
	PortfolioItems  []File     `json:"portfolioItems"`
	PortfolioLinks  []File     `json:"portfolioLinks"`
	Birthdate       string     `json:"birthdate"`
}

type Client struct {
	User
	Timezone  string    `json:"timezone"`
	About     string    `json:"about"`
	Birthdate string    `json:"birthdate"`
	Rating    float64   `json:"rating"`
	Jobs      []Job     `json:"jobs"`
	Projects  []Project `json:"projects"`
	Reviews   []Review  `json:"reviews"`
}

type ClientUpdate struct {
	Image     string `json:"image" valid:"required"`
	About     string `json:"about" valid:"required"`
	Timezone  string `json:"timezone" valid:"required"`
	Birthdate string `json:"birthdate"`
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
		Status:      projectStatusFinalizingTerms,
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

func (p *Project) allFreelancersConcluded() bool {
	if len(p.FreelancersConcluded) == len(p.Freelancers) {
		return true
	}
	return false
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
	ClientAgreed        bool      `json:"clientAgreed"`
	FreelancersAgreed   uintList  `json:"freelancersAgreed" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}

type Job struct {
	Model
	Name                string           `json:"name"`
	Summary             string           `json:"summary"`
	PriceFrom           int              `json:"priceFrom"`
	PriceTo             int              `json:"priceTo"`
	Tags                stringList       `json:"tags" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	Details             string           `json:"details"`
	ClientID            uint             `json:"-"`
	Client              *Client          `json:"client,omitempty"`
	Deadline            time.Time        `json:"deadline"`
	DeadlineFlexibility int              `json:"flexibility"`
	JobApplications     []JobApplication `json:"jobApplications"`
	Attachments         []File           `json:"attachments" gorm:"polymorphic:Owner;"`
	Examples            []File           `json:"examples" gorm:"polymorphic:Owner;"`
}

// BeforeSave updates file types before saving
func (job *Job) BeforeSave() error {
	for i := range job.Attachments {
		job.Attachments[i].Type = fileTypeJobAttachment
	}
	for i := range job.Examples {
		job.Examples[i].Type = fileTypeJobExample
	}
	return nil
}

type JobApplication struct {
	Model
	Message               string      `json:"message"`
	HourPrice             float64     `json:"hourPrice"`
	Hours                 int         `json:"hours"`
	Freelancer            *Freelancer `json:"freelancer,omitempty"`
	FreelancerID          uint        `json:"-"`
	FreelancerNumProjects int         `json:"freelancer_num_projects" sql:"-"`
	JobID                 uint        `json:"-"`
	Attachments           []File      `json:"attachments" gorm:"polymorphic:Owner;"`
	Examples              []File      `json:"examples" gorm:"polymorphic:Owner;"`
}

// BeforeSave updates file types before saving
func (jobApplication *JobApplication) BeforeSave() error {
	for i := range jobApplication.Attachments {
		jobApplication.Attachments[i].Type = fileTypeJobApplicationAttachment
	}
	for i := range jobApplication.Examples {
		jobApplication.Examples[i].Type = fileTypeJobApplicationExample
	}
	return nil
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

const (
	fileTypeFreelancerAdditionalField = "freelancer_additional_field"
	fileTypeFreelancerPortfolioItems  = "freelancer_portfolio_items"
	fileTypeFreelancerPortfolioLinks  = "freelancer_portfolio_links"
	fileTypeJobApplicationAttachment  = "job_application_attachment"
	fileTypeJobApplicationExample     = "job_application_example"
	fileTypeJobAttachment             = "job_attachment"
	fileTypeJobExample                = "job_example"
)

type File struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Type        string `json:"-"`
	OwnerID     int    `json:"-"`
	OwnerType   string `json:"-"`
}
