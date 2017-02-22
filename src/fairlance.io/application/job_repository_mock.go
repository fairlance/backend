package application

type JobRepositoryMock struct {
	GettAllJobsCall struct {
		Returns struct {
			Jobs  []Job
			Error error
		}
	}
	AddJobCall struct {
		Receives struct {
			Job *Job
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
			Job   Job
			Error error
		}
	}
	AddJobApplicationCall struct {
		Receives struct {
			JobApplication *JobApplication
		}
		Returns struct {
			Error error
		}
	}
	GetJobApplicationCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			JobApplication *JobApplication
			Error          error
		}
	}
	DeleteJobApplicationCall struct {
		Receives struct {
			JobApplication *JobApplication
		}
		Returns struct {
			Error error
		}
	}
}

func (repo *JobRepositoryMock) GetAllJobs() ([]Job, error) {
	return repo.GettAllJobsCall.Returns.Jobs,
		repo.GettAllJobsCall.Returns.Error
}

func (repo *JobRepositoryMock) AddJob(job *Job) error {
	repo.AddJobCall.Receives.Job = job
	return repo.AddJobCall.Returns.Error
}

func (repo *JobRepositoryMock) GetJob(id uint) (Job, error) {
	repo.GetJobCall.Receives.ID = id
	return repo.GetJobCall.Returns.Job,
		repo.GetJobCall.Returns.Error
}

func (repo *JobRepositoryMock) AddJobApplication(jobApplication *JobApplication) error {
	repo.AddJobApplicationCall.Receives.JobApplication = jobApplication
	return repo.AddJobApplicationCall.Returns.Error
}

func (repo *JobRepositoryMock) GetJobApplication(id uint) (*JobApplication, error) {
	repo.GetJobApplicationCall.Receives.ID = id
	return repo.GetJobApplicationCall.Returns.JobApplication,
		repo.GetJobApplicationCall.Returns.Error
}

func (repo *JobRepositoryMock) DeleteJobApplication(jobApplication *JobApplication) error {
	repo.DeleteJobApplicationCall.Receives.JobApplication = jobApplication
	return repo.DeleteJobApplicationCall.Returns.Error
}
