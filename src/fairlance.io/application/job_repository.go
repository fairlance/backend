package application

import (
	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	GetAllJobs() ([]Job, error)
	AddJob(job *Job) error
	GetJob(id uint) (Job, error)
	GetJobApplication(id uint) (*JobApplication, error)
	AddJobApplication(jobApplication *JobApplication) error
	DeleteJobApplication(jobApplication *JobApplication) error
}

type PostgreJobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) (JobRepository, error) {
	repo := &PostgreJobRepository{db}

	return repo, nil
}

func (repo *PostgreJobRepository) GetAllJobs() ([]Job, error) {
	jobs := []Job{}
	repo.db.Preload("JobApplications").Preload("JobApplications.Freelancer").Preload("JobApplications.Examples").Preload("JobApplications.Attachments").Preload("Examples").Preload("Attachments").Preload("Client").Find(&jobs)
	return jobs, nil
}

func (repo *PostgreJobRepository) AddJob(job *Job) error {
	return repo.db.Create(job).Error
}

func (repo *PostgreJobRepository) GetJob(id uint) (Job, error) {
	job := Job{}
	if err := repo.db.Preload("JobApplications").Preload("JobApplications.Freelancer").Preload("JobApplications.Examples").Preload("JobApplications.Attachments").Preload("Examples").Preload("Attachments").Preload("Client").Find(&job, id).Error; err != nil {
		return job, err
	}
	return job, nil
}

func (repo *PostgreJobRepository) GetJobApplication(id uint) (*JobApplication, error) {
	var jobApplication JobApplication
	err := repo.db.Find(&jobApplication, id).Error

	return &jobApplication, err
}

func (repo *PostgreJobRepository) AddJobApplication(jobApplication *JobApplication) error {
	return repo.db.Save(jobApplication).Error
}

func (repo *PostgreJobRepository) DeleteJobApplication(jobApplication *JobApplication) error {
	return repo.db.Delete(jobApplication).Error
}
