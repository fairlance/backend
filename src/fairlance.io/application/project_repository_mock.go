package application

type ProjectRepositoryMock struct {
	GetAllProjectsCall struct {
		Returns struct {
			Projects []Project
			Error    error
		}
	}
	GetByIDCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Project Project
			Error   error
		}
	}
	AddCall struct {
		Receives struct {
			Project *Project
		}
		Returns struct {
			Error error
		}
	}
	GetAllProjectsForClientCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Projects []Project
			Error    error
		}
	}
	GetAllProjectsForFreelancerCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Projects []Project
			Error    error
		}
	}
}

func (repo *ProjectRepositoryMock) GetAllProjects() ([]Project, error) {
	return repo.GetAllProjectsCall.Returns.Projects,
		repo.GetAllProjectsCall.Returns.Error
}

func (repo *ProjectRepositoryMock) GetByID(id uint) (Project, error) {
	repo.GetByIDCall.Receives.ID = id
	return repo.GetByIDCall.Returns.Project,
		repo.GetByIDCall.Returns.Error
}

func (repo *ProjectRepositoryMock) Add(project *Project) error {
	repo.AddCall.Receives.Project = project
	return repo.AddCall.Returns.Error
}

func (repo *ProjectRepositoryMock) GetAllProjectsForClient(id uint) ([]Project, error) {
	repo.GetAllProjectsForClientCall.Receives.ID = id
	return repo.GetAllProjectsForClientCall.Returns.Projects,
		repo.GetAllProjectsForClientCall.Returns.Error
}

func (repo *ProjectRepositoryMock) GetAllProjectsForFreelancer(id uint) ([]Project, error) {
	repo.GetAllProjectsForFreelancerCall.Receives.ID = id
	return repo.GetAllProjectsForFreelancerCall.Returns.Projects,
		repo.GetAllProjectsForFreelancerCall.Returns.Error
}
