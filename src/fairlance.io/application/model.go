package application

import (
	"encoding/json"
	"time"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
}

type Freelancer struct {
	Model
	Title          string    `json:"title" valid:"required"`
	FirstName      string    `json:"firstName" valid:"required"`
	LastName       string    `json:"lastName" valid:"required"`
	Password       string    `json:"-" valid:"required"`
	Email          string    `json:"email" valid:"required,email"`
	TimeZone       string    `json:"timeZone"`
	Rating         float64   `json:"rating"`
	HourlyRateFrom float64   `json:"hourlyRateFrom"`
	HourlyRateTo   float64   `json:"hourlyRateTo"`
	Projects       []Project `json:"projects" gorm:"many2many:project_freelancers;"`
	Reviews        []Review  `json:"reviews"`

	JsonReferences string      `json:"-" sql:"type:JSONB NOT NULL DEFAULT '[]'::JSONB"`
	References     []Reference `json:"references" sql:"-"`
}

func NewFreelancer(title string, firstName string, lastName string, password string, email string) *Freelancer {
	return &Freelancer{
		Title:          title,
		FirstName:      firstName,
		LastName:       lastName,
		Password:       password,
		Email:          email,
		TimeZone:       "UTC",
		HourlyRateFrom: 2,
		HourlyRateTo:   52,
		JsonReferences: `[]`,
	}
}

func (freelancer *Freelancer) getRepresentationMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         freelancer.ID,
		"firstName":  freelancer.FirstName,
		"lastName":   freelancer.LastName,
		"email":      freelancer.Email,
		"title":      freelancer.Title,
		"timeZone":   freelancer.TimeZone,
		"hourlyRate": []float64{freelancer.HourlyRateFrom, freelancer.HourlyRateTo},
	}
}

func (freelancer *Freelancer) AfterFind() (err error) {
	if err := json.Unmarshal([]byte(freelancer.JsonReferences), &freelancer.References); err != nil {
		return err
	}
	return nil
}

type Client struct {
	Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Jobs        []Job     `json:"jobs"`
	Projects    []Project `json:"projects"`
	Reviews     []Review  `json:"reviews"`
}

type Project struct {
	Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Freelancers []Freelancer `json:"freelancers" gorm:"many2many:project_freelancers;"`
	ClientId    uint         `json:"clientId"`
	IsActive    bool         `json:"isActive"`
}

type Job struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	ClientId    uint   `json:"clientId"`
	IsActive    bool   `json:"isActive"`
}

type Review struct {
	Model
	Title        string  `json:"title"`
	Content      string  `json:"content"`
	Rating       float64 `json:"rating"`
	ClientId     uint    `json:"clientId"`
	FreelancerId uint    `json:"freelancerId"`
}

type Reference struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Media   Media  `json:"media"`
}

type Media struct {
	Image string `json:"image"`
	Video string `json:"video"`
}
