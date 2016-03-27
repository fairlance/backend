package application

import (
    "time"
)

type Freelancer struct {
    Id        int             `json:"id"`
    FirstName string          `valid:"required" json:"firstName"`
    LastName  string          `valid:"required" json:"lastName"`
    Password  string          `valid:"required" json:",omitempty"`
    Email     string          `valid:"required,email" json:"email"`
    _Data     string          `sql:"data" json:"-"`
    Data      FreelancerData  `json:"data"`
    Projects  []Project       `json:"projects,omitempty"`
    Created   time.Time       `valid:"required" json:"created"`
}

func (freelancer *Freelancer) getRepresentationMap() map[string]interface{} {
    return map[string]interface{}{
        "id":freelancer.Id,
        "firstName":freelancer.FirstName,
        "lastName":freelancer.LastName,
        "email":freelancer.Email,
    }
}

type ProjectFreelancers struct {
    FreelancerId int
    ProjectId    int
}

type Client struct {
    Id          int
    Name        string
    Description string
    Jobs        []Job
    Projects    []Project
    Created     time.Time
}

type Project struct {
    Id          int
    Name        string
    Description string
    Freelancers []Freelancer    `json:",omitempty"`
    ClientId    int             `json:"-"`
    Client      Client          `json:",omitempty"`
    IsActive    bool
    Created     time.Time
}

type Job struct {
    Id          int
    ClientId    int     `json:"-"`
    Client      Client  `json:",omitempty"`
    Name        string
    Description string
    IsActive    bool
    Created     time.Time
}

type Review struct {
    Title    string     `json:"title"`
    Content  string     `json:"content"`
    Rating   float32    `json:"rating"`
    Created  string     `json:"created"`
    ClientId int        `json:"clientId"`
    Client   Client     `json:"client"`
}

type FreelancerData struct {
    Reviews []Review `json:"reviews"`
}