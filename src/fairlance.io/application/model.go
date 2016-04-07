package application

import "github.com/jinzhu/gorm"

type Freelancer struct {
	gorm.Model
	Title          string      `valid:"required"`
	FirstName      string      `valid:"required"`
	LastName       string      `valid:"required"`
	Password       string      `valid:"required"`
	Email          string      `valid:"required,email"`
	Projects       []Project   `gorm:"many2many:project_freelancers;"`
	JsonComments   string      `json:"-" sql:"type:JSONB NOT NULL"`
	Comments       []Comment   `sql:"-"`
	JsonReferences string      `json:"-" sql:"type:JSONB NOT NULL"`
	References     []Reference `sql:"-"`
}

func (freelancer *Freelancer) getRepresentationMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        freelancer.ID,
		"firstName": freelancer.FirstName,
		"lastName":  freelancer.LastName,
		"email":     freelancer.Email,
		"title":     freelancer.Title,
	}
}

type Client struct {
	gorm.Model
	Name        string
	Description string
	Jobs        []Job
	Projects    []Project
}

type Project struct {
	gorm.Model
	Name        string
	Description string
	Freelancers []Freelancer `gorm:"many2many:project_freelancers;"`
	ClientId    uint
	IsActive    bool
}

type Job struct {
	gorm.Model
	Name        string
	Description string
	ClientId    uint
	IsActive    bool
}

// type Review struct {
// 	Title    string  `json:"title"`
// 	Content  string  `json:"content"`
// 	Rating   float32 `json:"rating"`
// 	Created  string  `json:"created"`
// 	ClientId int     `json:"clientId"`
// 	Client   Client  `json:"client"`
// }

type Reference struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Media   Media  `json:"media"`
}

type Media struct {
	Image string `json:"image"`
	Video string `json:"video"`
}

type Comment struct {
	Text     string `json:"text"`
	ClientId uint   `json:"clientId"`
}
