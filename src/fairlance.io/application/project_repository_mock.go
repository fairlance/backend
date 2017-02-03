package application

type projectRepositoryMock struct {
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

func (repo *projectRepositoryMock) getAllProjects() ([]Project, error) {
	return repo.GetAllProjectsCall.Returns.Projects,
		repo.GetAllProjectsCall.Returns.Error
}

func (repo *projectRepositoryMock) getByID(id uint) (Project, error) {
	repo.GetByIDCall.Receives.ID = id
	return repo.GetByIDCall.Returns.Project,
		repo.GetByIDCall.Returns.Error
}

func (repo *projectRepositoryMock) add(project *Project) error {
	repo.AddCall.Receives.Project = project
	return repo.AddCall.Returns.Error
}

func (repo *projectRepositoryMock) getAllProjectsForClient(id uint) ([]Project, error) {
	repo.GetAllProjectsForClientCall.Receives.ID = id
	return repo.GetAllProjectsForClientCall.Returns.Projects,
		repo.GetAllProjectsForClientCall.Returns.Error
}

func (repo *projectRepositoryMock) getAllProjectsForFreelancer(id uint) ([]Project, error) {
	repo.GetAllProjectsForFreelancerCall.Receives.ID = id
	return repo.GetAllProjectsForFreelancerCall.Returns.Projects,
		repo.GetAllProjectsForFreelancerCall.Returns.Error
}
