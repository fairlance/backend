package application

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
	Model
	FirstName string `json:"firstName" valid:"required"`
	LastName  string `json:"lastName" valid:"required"`
	Password  string `json:"-" valid:"required"`
	Email     string `json:"email" valid:"required,email" sql:"index" gorm:"unique"`
}

type Freelancer struct {
	User
	Rating         float64     `json:"rating"`
	Timezone       string      `json:"timezone"`
	Skills         []Tag       `json:"skills" gorm:"polymorphic:Owner;"`
	IsAvailable    bool        `json:"isAvailable"`
	HourlyRateFrom uint        `json:"hourlyRateFrom"`
	HourlyRateTo   uint        `json:"hourlyRateTo"`
	Projects       []Project   `json:"projects" gorm:"many2many:project_freelancers;"`
	Reviews        []Review    `json:"reviews"`
	References     []Reference `json:"references"`
}

type FreelancerUpdate struct {
	Skills         []Tag  `json:"skills"`
	Timezone       string `json:"timezone" valid:"required"`
	IsAvailable    bool   `json:"isAvailable"`
	HourlyRateFrom uint   `json:"hourlyRateFrom" valid:"required"`
	HourlyRateTo   uint   `json:"hourlyRateTo" valid:"required"`
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
	Name        string       `json:"name" valid:"required"`
	Description string       `json:"description" valid:"required"`
	Freelancers []Freelancer `json:"freelancers" gorm:"many2many:project_freelancers;"`
	ClientId    uint         `json:"-" valid:"required"`
	Client      Client       `json:"client"`
	IsActive    bool         `json:"isActive"`
}

type Job struct {
	Model
	Name        string `json:"name" valid:"required"`
	Description string `json:"description" valid:"required"`
	ClientId    uint   `json:"-" valid:"required"`
	Client      Client `json:"client"`
	IsActive    bool   `json:"isActive"`
	Tags        []Tag  `json:"tags" gorm:"polymorphic:Owner;"`
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

type Tag struct {
	ID        uint   `json:"-" gorm:"primary_key"`
	Name      string `json:"name" valid:"required"`
	OwnerId   uint   `json:"-"`
	OwnerType string `json:"-"`
}
