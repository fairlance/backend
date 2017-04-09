package application

type JobRepositoryMock struct {
	GettAllJobsCall struct {
		Returns struct {
			Jobs  []Job
			Error error
		}
	}
	GetAllJobsForClientCall struct {
		Receives struct {
			ID uint
		}
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
			Job   *Job
			Error error
		}
	}
	GetJobForClientCall struct {
		Receives struct {
			ID       uint
			ClientID uint
		}
		Returns struct {
			Job   *Job
			Error error
		}
	}
	GetJobForFreelancerCall struct {
		Receives struct {
			ID           uint
			FreelancerID uint
		}
		Returns struct {
			Job   *Job
			Error error
		}
	}
	DeleteJobCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Error error
		}
	}
	DeactivateJobCall struct {
		Receives struct {
			Job *Job
		}
		Returns struct {
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
			ID uint
		}
		Returns struct {
			Error error
		}
	}
	JobApplicationBelongsToClientCall struct {
		Receives struct {
			ID       uint
			ClientID uint
		}
		Returns struct {
			OK    bool
			Error error
		}
	}
	JobApplicationBelongsToFreelancerCall struct {
		Receives struct {
			ID           uint
			FreelancerID uint
		}
		Returns struct {
			OK    bool
			Error error
		}
	}
}

func (repo *JobRepositoryMock) GetAllJobs() ([]Job, error) {
	return repo.GettAllJobsCall.Returns.Jobs,
		repo.GettAllJobsCall.Returns.Error
}

func (repo *JobRepositoryMock) GetAllJobsForClient(id uint) ([]Job, error) {
	repo.GetAllJobsForClientCall.Receives.ID = id
	return repo.GetAllJobsForClientCall.Returns.Jobs,
		repo.GetAllJobsForClientCall.Returns.Error
}

func (repo *JobRepositoryMock) AddJob(job *Job) error {
	repo.AddJobCall.Receives.Job = job
	return repo.AddJobCall.Returns.Error
}

func (repo *JobRepositoryMock) GetJob(id uint) (*Job, error) {
	repo.GetJobCall.Receives.ID = id
	return repo.GetJobCall.Returns.Job,
		repo.GetJobCall.Returns.Error
}

func (repo *JobRepositoryMock) GetJobForClient(id, clientID uint) (*Job, error) {
	repo.GetJobForClientCall.Receives.ID = id
	repo.GetJobForClientCall.Receives.ClientID = clientID
	return repo.GetJobForClientCall.Returns.Job,
		repo.GetJobForClientCall.Returns.Error
}

func (repo *JobRepositoryMock) GetJobForFreelancer(id, freelancerID uint) (*Job, error) {
	repo.GetJobForFreelancerCall.Receives.ID = id
	repo.GetJobForFreelancerCall.Receives.FreelancerID = freelancerID
	return repo.GetJobForFreelancerCall.Returns.Job,
		repo.GetJobForFreelancerCall.Returns.Error
}

func (repo *JobRepositoryMock) DeleteJob(id uint) error {
	repo.DeleteJobCall.Receives.ID = id
	return repo.DeleteJobCall.Returns.Error
}

func (repo *JobRepositoryMock) DeactivateJob(job *Job) error {
	repo.DeactivateJobCall.Receives.Job = job
	return repo.DeactivateJobCall.Returns.Error
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

func (repo *JobRepositoryMock) DeleteJobApplication(id uint) error {
	repo.DeleteJobApplicationCall.Receives.ID = id
	return repo.DeleteJobApplicationCall.Returns.Error
}

func (repo *JobRepositoryMock) jobApplicationBelongsToClient(id uint, clientID uint) (bool, error) {
	repo.JobApplicationBelongsToClientCall.Receives.ID = id
	repo.JobApplicationBelongsToClientCall.Receives.ClientID = clientID
	return repo.JobApplicationBelongsToClientCall.Returns.OK,
		repo.JobApplicationBelongsToClientCall.Returns.Error
}

func (repo *JobRepositoryMock) jobApplicationBelongsToFreelancer(id uint, freelancerID uint) (bool, error) {
	repo.JobApplicationBelongsToFreelancerCall.Receives.ID = id
	repo.JobApplicationBelongsToFreelancerCall.Receives.FreelancerID = freelancerID
	return repo.JobApplicationBelongsToFreelancerCall.Returns.OK,
		repo.JobApplicationBelongsToFreelancerCall.Returns.Error
}
