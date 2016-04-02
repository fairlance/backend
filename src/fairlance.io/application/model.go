package application

import (
	"time"
)

type Freelancer struct {
	Id        int            `json:"id"`
	Title     string         `valid:"required" json:"title"`
	FirstName string         `valid:"required" json:"firstName"`
	LastName  string         `valid:"required" json:"lastName"`
	Password  string         `valid:"required" json:",omitempty"`
	Email     string         `valid:"required,email" json:"email"`
	_Data     string         `sql:"data" json:"-"`
	Data      FreelancerData `json:"data"`
	Projects  []Project      `json:"projects"`
	Created   time.Time      `valid:"required" json:"created"`
}

func (freelancer *Freelancer) getRepresentationMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        freelancer.Id,
		"firstName": freelancer.FirstName,
		"lastName":  freelancer.LastName,
		"email":     freelancer.Email,
	}
}

type ProjectFreelancers struct {
	FreelancerId int
	ProjectId    int
}

type Client struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Jobs        []Job     `json:"jobs"`
	Projects    []Project `json:"projects"`
	Created     time.Time `json:"created"`
}

type Project struct {
	Id          int
	Name        string
	Description string
	Freelancers []Freelancer `json:",omitempty"`
	ClientId    int          `json:"-"`
	Client      Client       `json:",omitempty"`
	IsActive    bool
	Created     time.Time
}

type Job struct {
	Id          int
	ClientId    int    `json:"-"`
	Client      Client `json:",omitempty"`
	Name        string
	Description string
	IsActive    bool
	Created     time.Time
}

type Review struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Rating   float32 `json:"rating"`
	Created  string  `json:"created"`
	ClientId int     `json:"clientId"`
	Client   Client  `json:"client"`
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

type FreelancerData struct {
	Reviews    []Review    `json:"reviews"`
	References []Reference `json:"references"`
}
