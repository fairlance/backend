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
			Project *Project
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
	AddContractCall struct {
		Receives struct {
			Contract *Contract
		}
		Returns struct {
			Error error
		}
	}
	AddExtensionCall struct {
		Receives struct {
			Extension *Extension
		}
		Returns struct {
			Error error
		}
	}
	ProjectBelongsToUserCall struct {
		Receives struct {
			ID       uint
			UserType string
			UserID   uint
		}
		Returns struct {
			OK    bool
			Error error
		}
	}
}

func (repo *projectRepositoryMock) getAllProjects() ([]Project, error) {
	return repo.GetAllProjectsCall.Returns.Projects,
		repo.GetAllProjectsCall.Returns.Error
}

func (repo *projectRepositoryMock) getByID(id uint) (*Project, error) {
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

func (repo *projectRepositoryMock) addContract(contract *Contract) error {
	repo.AddContractCall.Receives.Contract = contract
	return repo.AddContractCall.Returns.Error
}

func (repo *projectRepositoryMock) addExtension(extension *Extension) error {
	repo.AddExtensionCall.Receives.Extension = extension
	return repo.AddExtensionCall.Returns.Error
}

func (repo *projectRepositoryMock) projectBelongsToUser(id uint, userType string, userID uint) (bool, error) {
	repo.ProjectBelongsToUserCall.Receives.ID = id
	repo.ProjectBelongsToUserCall.Receives.UserType = userType
	repo.ProjectBelongsToUserCall.Receives.UserID = userID
	return repo.ProjectBelongsToUserCall.Returns.OK,
		repo.ProjectBelongsToUserCall.Returns.Error
}
