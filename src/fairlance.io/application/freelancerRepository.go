package application

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type FreelancerRepository struct {
	db *gorm.DB
}

func NewFreelancerRepository(db *gorm.DB) (*FreelancerRepository, error) {
	repo := &FreelancerRepository{db}

	return repo, nil
}
func (repo *FreelancerRepository) GetAllFreelancers() ([]Freelancer, error) {
	freelancers := []Freelancer{}
	repo.db.Preload("Projects").Find(&freelancers)
	return freelancers, nil
}

func (repo *FreelancerRepository) GetFreelancerByEmail(email string) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.Where("email = ?", email).First(&freelancer).Error; err != nil {
		return freelancer, err
	}
	return freelancer, nil
}

func (repo *FreelancerRepository) GetFreelancer(id int) (Freelancer, error) {
	freelancer := Freelancer{}
	if err := repo.db.Preload("Projects").Where("id = ?", id).Find(&freelancer).Error; err != nil {
		return freelancer, err
	}
	return freelancer, nil
}

func (repo *FreelancerRepository) CheckCredentials(email string, password string) error {
	freelancer, err := repo.GetFreelancerByEmail(email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(freelancer.Password), []byte(password)); err != nil {
		return errors.New("Freelancer not found (password is wrong)")
	}

	return nil
}

func (repo *FreelancerRepository) AddFreelancer(freelancer *Freelancer) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(freelancer.Password), bcrypt.DefaultCost)
	freelancer.Password = string(hashedPassword)
	return repo.db.Create(freelancer).Error
}

func (repo *FreelancerRepository) DeleteFreelancer(id uint) error {
	freelancer := Freelancer{}
	freelancer.ID = id
	if repo.db.Find(&freelancer).RecordNotFound() {
		return errors.New("Freelancer not found")
	}
	return repo.db.Delete(&freelancer).Error
}

//func (repo *FreelancerRepository) addReference(freelancerId int, reference Reference) error {
//
//	//TODO: wrong
//	freelancer, _ := repo.GetFreelancer(freelancerId)
//	freelancer.Data.References = append(freelancer.Data.References, reference)
//	dataJson, err := json.Marshal(freelancer.Data)
//	if err != nil {
//		return err
//	}
//
//	_, err = repo.db.Exec(`
//            UPDATE freelancers SET data = $1 WHERE id = $2;`,
//		dataJson,
//		freelancerId,
//	)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
