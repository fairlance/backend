package application

type UserRepositoryMock struct {
	CheckCredentialsCall struct {
		Receives struct {
			Email    string
			Password string
		}
		Returns struct {
			User     User
			UserType string
			Error    error
		}
	}
	GetUserByEmailCall struct {
		Receives struct {
			Email string
		}
		Returns struct {
			User     User
			UserType string
			Error    error
		}
	}
}

func (repo *UserRepositoryMock) CheckCredentials(email string, password string) (User, string, error) {
	repo.CheckCredentialsCall.Receives.Email = email
	repo.CheckCredentialsCall.Receives.Password = password

	return repo.CheckCredentialsCall.Returns.User,
		repo.CheckCredentialsCall.Returns.UserType,
		repo.CheckCredentialsCall.Returns.Error
}

func (repo *UserRepositoryMock) GetUserByEmail(email string) (User, string, error) {
	repo.GetUserByEmailCall.Receives.Email = email

	return repo.GetUserByEmailCall.Returns.User,
		repo.GetUserByEmailCall.Returns.UserType,
		repo.GetUserByEmailCall.Returns.Error
}
