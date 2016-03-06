package application

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "errors"
    "golang.org/x/crypto/bcrypt"
)

var collectionName = "freelancers"

type FreelancerRepository struct {
    session *mgo.Session
    db      string
}

func NewFreelancerRepository(db string) (*FreelancerRepository, error) {
    session, err := mgo.Dial("localhost")
    if err != nil {
        return nil, err
    }

    repo := &FreelancerRepository{session, db}
    if err != nil {
        return nil, err
    }

    return repo, nil
}

func (repo FreelancerRepository) GetAllFreelancers() ([]Freelancer, error) {
    session := repo.session.Copy()
    defer session.Close()

    freelancers := []Freelancer{}

    collection := session.DB(repo.db).C(collectionName)
    err := collection.Find(nil).All(&freelancers)
    if err != nil {
        return freelancers, err
    }

    return freelancers, nil
}

func (repo FreelancerRepository) AddFreelancer(freelancer Freelancer) error {
    session := repo.session.Copy()
    defer session.Close()

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(freelancer.Password), bcrypt.DefaultCost)
    freelancer.Password = string(hashedPassword)

    collection := session.DB(repo.db).C(collectionName)
    err := collection.Insert(&freelancer)
    if err != nil {
        if mgo.IsDup(err) {
            return err
        }

        return err
    }

    return nil
}

func (repo FreelancerRepository) DeleteFreelancer(id string) error {
    if !bson.IsObjectIdHex(id) {
        return errors.New("Invalid id provided")
    }

    session := repo.session.Copy()
    defer session.Close()

    collection := session.DB(repo.db).C(collectionName)
    err := collection.RemoveId(bson.ObjectIdHex(id))
    if err != nil {
        return err
    }

    return nil
}

func (repo FreelancerRepository) GetFreelancer(id string) (Freelancer, error) {
    freelancer := Freelancer{}
    if !bson.IsObjectIdHex(id) {
        return freelancer, errors.New("Invalid id provided")
    }

    session := repo.session.Copy()
    defer session.Close()

    collection := session.DB(repo.db).C(collectionName)
    if err := collection.FindId(bson.ObjectIdHex(id)).One(&freelancer); err != nil {
        return freelancer, err
    }

    return freelancer, nil
}

func (repo FreelancerRepository) CheckCredentials(email string, password string) (bool, error) {
    session := repo.session.Copy()
    defer session.Close()

    freelancer := Freelancer{}
    collection := session.DB(repo.db).C(collectionName)
    if err := collection.Find(bson.M{"email": email}).One(&freelancer); err != nil {
        if err == mgo.ErrNotFound {
            return false, errors.New("Freelancer not found")
        }
        return false, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(freelancer.Password), []byte(password)); err != nil {
        return false, errors.New("Freelancer not found (password is wrong)")
    }

    return true, nil
}

func (repo FreelancerRepository) UpdateFreelancer(id string, data Freelancer) error {
    if !bson.IsObjectIdHex(id) {
        return errors.New("Invalid id provided")
    }

    change := bson.M{"name": data.FirstName, "email": data.Email}
    if data.Password != "" {
        change["password"] = data.Password
    }
    update := bson.M{"$set": change}

    session := repo.session.Copy()
    defer session.Close()

    collection := session.DB(repo.db).C(collectionName)
    if err := collection.UpdateId(bson.ObjectIdHex(id), update); err != nil {
        return err
    }

    return nil
}
