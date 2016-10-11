package mocks

import "fairlance.io/application"

type JobRepository struct {
	GettAllJobsCall struct {
		Returns struct {
			Jobs  []application.Job
			Error error
		}
	}
	AddJobCall struct {
		Receives struct {
			Job *application.Job
		}
		Returns struct {
			Error error
		}
	}
	GetJobCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Job   application.Job
			Error error
		}
	}
	AddJobApplicationCall struct {
		Receives struct {
			JobApplication *application.JobApplication
		}
		Returns struct {
			Error error
		}
	}
}

func (repo *JobRepository) GetAllJobs() ([]application.Job, error) {
	return repo.GettAllJobsCall.Returns.Jobs,
		repo.GettAllJobsCall.Returns.Error
}

func (repo *JobRepository) AddJob(job *application.Job) error {
	repo.AddJobCall.Receives.Job = job
	return repo.AddJobCall.Returns.Error
}

func (repo *JobRepository) GetJob(id uint) (application.Job, error) {
	repo.GetJobCall.Receives.ID = id
	return repo.GetJobCall.Returns.Job,
		repo.GetJobCall.Returns.Error
}

func (repo *JobRepository) AddJobApplication(jobApplication *application.JobApplication) error {
	repo.AddJobApplicationCall.Receives.JobApplication = jobApplication
	return repo.AddJobApplicationCall.Returns.Error
}
