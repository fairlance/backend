package main

import (
	"time"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
}

type User struct {
	FirstName string `json:"firstName" valid:"required"`
	LastName  string `json:"lastName" valid:"required"`
	Password  string `json:"-" valid:"required"`
	Email     string `json:"email" valid:"required,email" sql:"index"`
}

type Freelancer struct {
	Model
	User
	Title          string      `json:"title" valid:"required"`
	TimeZone       string      `json:"timeZone"`
	Rating         float64     `json:"rating"`
	HourlyRateFrom float64     `json:"hourlyRateFrom"`
	HourlyRateTo   float64     `json:"hourlyRateTo"`
	Projects       []Project   `json:"projects" gorm:"many2many:project_freelancers;"`
	Reviews        []Review    `json:"reviews"`
	References     []Reference `json:"references"`
}

func NewFreelancer(
	firstName string,
	lastName string,
	title string,
	password string,
	email string,
	hourlyRateFrom float64,
	hourlyRateTo float64,
	timeZone string,
) *Freelancer {
	return &Freelancer{
		User: User{
			FirstName: firstName,
			LastName:  lastName,
			Password:  password,
			Email:     email,
		},
		Title:          title,
		HourlyRateFrom: hourlyRateFrom,
		HourlyRateTo:   hourlyRateTo,
		TimeZone:       timeZone,
		// Reviews:        []Review{},
		// Projects:       []Project{},
		// References:     []Reference{},
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

type Client struct {
	Model
	User
	Jobs     []Job     `json:"jobs"`
	Projects []Project `json:"projects"`
	Reviews  []Review  `json:"reviews"`
}

type Project struct {
	Model
	Name        string       `json:"name" valid:"required"`
	Description string       `json:"description" valid:"required"`
	Freelancers []Freelancer `json:"freelancers" gorm:"many2many:project_freelancers;"`
	ClientId    uint         `json:"clientId" valid:"required"`
	IsActive    bool         `json:"isActive"`
}

type Job struct {
	Model
	Name        string `json:"name" valid:"required"`
	Description string `json:"description" valid:"required"`
	ClientId    uint   `json:"clientId" valid:"required"`
	IsActive    bool   `json:"isActive"`
}

type Review struct {
	Model
	Title        string  `json:"title" valid:"required"`
	Content      string  `json:"content"`
	Rating       float64 `json:"rating" valid:"required"`
	JobId        uint    `json:"jobId" valid:"required"`
	ClientId     uint    `json:"clientId" valid:"required"`
	FreelancerId uint    `json:"freelancerId" valid:"required"`
}

type Reference struct {
	Model
	Title        string `json:"title" valid:"required"`
	Content      string `json:"content"`
	Media        Media  `json:"media"`
	FreelancerId uint   `json:"freelancerId" valid:"required"`
}

type Media struct {
	Model
	Image       string `json:"image"`
	Video       string `json:"video"`
	ReferenceId uint   `json:"referenceId"`
}
