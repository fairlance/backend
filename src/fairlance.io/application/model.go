package application

import (
    "gopkg.in/mgo.v2/bson"
    "time"
)

type Freelancer struct {
    Id        bson.ObjectId `bson:"_id,omitempty"`
    FirstName string        `bson:"firstName" valid:"required"`
    LastName  string        `bson:"lastName" valid:"required"`
    Password  string        `valid:"required"`
    Email     string        `valid:"required,email"`
    Created   time.Time     `valid:"required"`
}

func (freelancer *Freelancer) getRepresentationMap() map[string]string  {
    return map[string]string{
        "id":freelancer.Id.Hex(),
        "firstName":freelancer.FirstName,
        "lastName":freelancer.LastName,
        "email":freelancer.Email,
    }
}

type Client struct {
    Id          bson.ObjectId  `bson:"_id,omitempty"`
    Name        string
    Description string
    JobPostings []JobPosting   `bson:"-"`
    Projects    []Project      `bson:"-"`
    Created     time.Time
}

type Project struct {
    Id            bson.ObjectId     `bson:"_id,omitempty"`
    Name          string
    Description   string
    ClientId      bson.ObjectId     `bson:"clientId"`
    Client        Client            `bson:"-"`
    FreelancerIds []bson.ObjectId   `bson:"freelancerIds"`
    Freelancers   []Freelancer      `bson:"-"`
    IsActive      bool              `bson:"isActive"`
    Created       time.Time
}

type JobPosting struct {
    Id          bson.ObjectId  `bson:"_id,omitempty"`
    ClientId    bson.ObjectId  `bson:"clientId"`
    Client      Client         `bson:"-"`
    Name        string
    Description string
}