package application

import (
    "time"
)

type Freelancer struct {
    Id        int
    FirstName string        `valid:"required"`
    LastName  string        `valid:"required"`
    Password  string        `valid:"required" json:",omitempty"`
    Email     string        `valid:"required,email"`
    Projects  []Project     `json:",omitempty"`
    Created   time.Time     `valid:"required"`
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
    Client      *Client         `json:",omitempty"`
    IsActive    bool
    Created     time.Time
}

type Job struct {
    Id          int
    ClientId    int     `json:"-"`
    Client      *Client `json:",omitempty"`
    Name        string
    Description string
    IsActive    bool
    Created     time.Time
}