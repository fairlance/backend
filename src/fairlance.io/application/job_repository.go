package application

import (
	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	GetAllJobs() ([]Job, error)
	GetAllJobsForClient(id uint) ([]Job, error)
	AddJob(job *Job) error
	GetJob(id uint) (*Job, error)
	GetJobForClient(id, clientID uint) (*Job, error)
	GetJobForFreelancer(id, freelancerID uint) (*Job, error)
	DeactivateJob(job *Job) error
	DeleteJob(id uint) error
	GetJobApplication(id uint) (*JobApplication, error)
	AddJobApplication(jobApplication *JobApplication) error
	DeleteJobApplication(id uint) error
	jobApplicationBelongsToClient(id uint, clientID uint) (bool, error)
	jobApplicationBelongsToFreelancer(id uint, freelancerID uint) (bool, error)
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
	err := repo.db.Preload("JobApplications").Preload("JobApplications.Freelancer").Preload("JobApplications.Examples").Preload("JobApplications.Attachments").Preload("Examples").Preload("Attachments").Preload("Client").Find(&jobs).Error
	return jobs, err
}

func (repo *PostgreJobRepository) GetAllJobsForClient(id uint) ([]Job, error) {
	jobs := []Job{}
	err := repo.db.Preload("JobApplications").Preload("JobApplications.Freelancer").Preload("JobApplications.Examples").Preload("JobApplications.Attachments").Preload("Examples").Preload("Attachments").Preload("Client").Find(&jobs).Where("client_id = ?", id).Error
	return jobs, err
}

func (repo *PostgreJobRepository) AddJob(job *Job) error {
	return repo.db.Create(job).Error
}

func (repo *PostgreJobRepository) DeleteJob(id uint) error {
	return repo.db.Where("id = ?", id).Delete(&Job{}).Error
}

func (repo *PostgreJobRepository) DeactivateJob(job *Job) error {
	return repo.db.Model(job).Update("is_active", false).Error
}

func (repo *PostgreJobRepository) GetJob(id uint) (*Job, error) {
	job := &Job{}
	if err := repo.db.Preload("JobApplications").Preload("JobApplications.Freelancer").Preload("JobApplications.Examples").Preload("JobApplications.Attachments").Preload("Examples").Preload("Attachments").Preload("Client").Find(job, id).Error; err != nil {
		return job, err
	}
	return job, nil
}

func (repo *PostgreJobRepository) GetJobForClient(id, clientID uint) (*Job, error) {
	job := &Job{}
	if err := repo.db.Preload("JobApplications").Preload("JobApplications.Freelancer").Preload("JobApplications.Examples").Preload("JobApplications.Attachments").Preload("Examples").Preload("Attachments").Preload("Client").Where("client_id = ?", clientID).Find(job, id).Error; err != nil {
		return job, err
	}
	return job, nil
}

func (repo *PostgreJobRepository) GetJobForFreelancer(id, freelancerID uint) (*Job, error) {
	job := &Job{}
	if err := repo.db.Preload("Examples").Preload("Attachments").Preload("Client").Find(job, id).Error; err != nil {
		return job, err
	}

	var jobApplications []JobApplication
	if err := repo.db.Preload("Examples").Preload("Attachments").Where("freelancer_id = ?", freelancerID).Model(job).Related(&jobApplications).Error; err != nil {
		return job, err
	}

	job.JobApplications = jobApplications

	return job, nil
}

func (repo *PostgreJobRepository) GetJobApplication(id uint) (*JobApplication, error) {
	var jobApplication JobApplication
	err := repo.db.Preload("Freelancer").Find(&jobApplication, id).Error

	return &jobApplication, err
}

func (repo *PostgreJobRepository) AddJobApplication(jobApplication *JobApplication) error {
	return repo.db.Save(jobApplication).Error
}

func (repo *PostgreJobRepository) DeleteJobApplication(id uint) error {
	return repo.db.Where("id = ?", id).Delete(&JobApplication{}).Error
}

func (repo *PostgreJobRepository) jobApplicationBelongsToClient(id uint, clientID uint) (bool, error) {
	var jobApplication JobApplication
	if err := repo.db.Find(&jobApplication, id).Error; err != nil {
		return false, err
	}

	job, err := repo.GetJob(jobApplication.JobID)
	if err != nil {
		return false, err
	}

	if job.ClientID == clientID {
		return true, nil
	}

	return false, nil
}

func (repo *PostgreJobRepository) jobApplicationBelongsToFreelancer(id uint, freelancerID uint) (bool, error) {
	var jobApplication JobApplication
	if err := repo.db.Find(&jobApplication, id).Error; err != nil {
		return false, err
	}

	if jobApplication.FreelancerID == freelancerID {
		return true, nil
	}

	return false, nil
}
