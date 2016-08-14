package application

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) (*JobRepository, error) {
	repo := &JobRepository{db}

	return repo, nil
}

func (repo *JobRepository) GetAllJobs() ([]Job, error) {
	jobs := []Job{}
	repo.db.Preload("Client").Find(&jobs)
	return jobs, nil
}

func (repo *JobRepository) AddJob(job *Job) error {
	return repo.db.Create(job).Error
}

func (repo *JobRepository) GetJob(id uint) (Job, error) {
	job := Job{}
	if err := repo.db.Preload("Client").Find(&job, id).Error; err != nil {
		return job, err
	}
	return job, nil
}
