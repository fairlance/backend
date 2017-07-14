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
	update(job *Job) error
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
	err := repo.db.
		Preload("JobApplications").
		Preload("JobApplications.Freelancer").
		Preload("Examples", "type IN (?)", fileTypeJobExample).
		Preload("Attachments", "type IN (?)", fileTypeJobAttachment).
		Preload("Client").
		Find(&jobs).Error
	return jobs, err
}

func (repo *PostgreJobRepository) GetAllJobsForClient(clientID uint) ([]Job, error) {
	jobs := []Job{}
	err := repo.db.
		Preload("JobApplications").
		Preload("JobApplications.Freelancer").
		Preload("Examples", "type IN (?)", fileTypeJobExample).
		Preload("Attachments", "type IN (?)", fileTypeJobAttachment).
		Preload("Client").
		Where("client_id = ?", clientID).
		Find(&jobs).Error
	return jobs, err
}

func (repo *PostgreJobRepository) AddJob(job *Job) error {
	return repo.db.Create(job).Error
}

func (repo *PostgreJobRepository) DeleteJob(jobID uint) error {
	tx := repo.db.Begin()
	if err := tx.Where("owner_type = 'jobs' AND owner_id = ?", jobID).Delete(&File{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := repo.db.Where("job_id = ?", jobID).Delete(&JobApplication{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := repo.db.Where("id = ?", jobID).Delete(&Job{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *PostgreJobRepository) update(job *Job) error {
	return repo.db.Save(job).Error
}

func (repo *PostgreJobRepository) GetJob(jobID uint) (*Job, error) {
	job := &Job{}
	if err := repo.db.
		Preload("JobApplications").
		Preload("JobApplications.Freelancer").
		Preload("Examples", "type IN (?)", fileTypeJobExample).
		Preload("Attachments", "type IN (?)", fileTypeJobAttachment).
		Preload("Client").
		Find(job, jobID).Error; err != nil {
		return job, err
	}
	return job, nil
}

func (repo *PostgreJobRepository) GetJobForClient(jobID, clientID uint) (*Job, error) {
	job := &Job{}
	if err := repo.db.
		Preload("JobApplications").
		Preload("JobApplications.Freelancer").
		Preload("JobApplications.Examples", "type IN (?)", fileTypeJobApplicationExample).
		Preload("JobApplications.Attachments", "type IN (?)", fileTypeJobApplicationAttachment).
		Preload("Examples", "type IN (?)", fileTypeJobExample).
		Preload("Attachments", "type IN (?)", fileTypeJobAttachment).
		Preload("Client").
		Where("client_id = ?", clientID).
		Find(job, jobID).Error; err != nil {
		return job, err
	}
	return job, nil
}

func (repo *PostgreJobRepository) GetJobForFreelancer(jobID, freelancerID uint) (*Job, error) {
	job := &Job{}
	if err := repo.db.
		Preload("Examples", "type IN (?)", fileTypeJobExample).
		Preload("Attachments", "type IN (?)", fileTypeJobAttachment).
		Preload("Client").
		Find(job, jobID).Error; err != nil {
		return job, err
	}
	var jobApplications []JobApplication
	if err := repo.db.
		Model(job).
		Related(&jobApplications).
		Preload("Examples", "type IN (?)", fileTypeJobApplicationExample).
		Preload("Attachments", "type IN (?)", fileTypeJobApplicationAttachment).
		Where("freelancer_id = ?", freelancerID).Error; err != nil {
		return job, err
	}
	for i := range jobApplications {
		if err := repo.db.Model(&Freelancer{}).Where("id = ?", jobApplications[i].FreelancerID).Count(&jobApplications[i].FreelancerNumProjects).Error; err != nil {
			return nil, err
		}
	}
	return job, nil
}

func (repo *PostgreJobRepository) GetJobApplication(jobApplicationID uint) (*JobApplication, error) {
	var jobApplication JobApplication
	if err := repo.db.
		Preload("Freelancer").
		Preload("Examples", "type IN (?)", fileTypeJobApplicationExample).
		Preload("Attachments", "type IN (?)", fileTypeJobApplicationAttachment).
		Find(&jobApplication, jobApplicationID).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Model(&Freelancer{}).Where("id = ?", jobApplication.FreelancerID).Count(&jobApplication.FreelancerNumProjects).Error; err != nil {
		return nil, err
	}
	return &jobApplication, nil
}

func (repo *PostgreJobRepository) AddJobApplication(jobApplication *JobApplication) error {
	return repo.db.Save(jobApplication).Error
}

func (repo *PostgreJobRepository) DeleteJobApplication(jobApplicationID uint) error {
	tx := repo.db.Begin()
	if err := tx.Where("id = ?", jobApplicationID).Delete(&JobApplication{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("owner_type = job_applications AND owner_id = ?", jobApplicationID).Delete(&File{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *PostgreJobRepository) jobApplicationBelongsToClient(jobApplicationID uint, clientID uint) (bool, error) {
	var jobApplication JobApplication
	if err := repo.db.Find(&jobApplication, jobApplicationID).Error; err != nil {
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

func (repo *PostgreJobRepository) jobApplicationBelongsToFreelancer(jobApplicationID uint, freelancerID uint) (bool, error) {
	var jobApplication JobApplication
	if err := repo.db.Find(&jobApplication, jobApplicationID).Error; err != nil {
		return false, err
	}

	if jobApplication.FreelancerID == freelancerID {
		return true, nil
	}

	return false, nil
}
