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
}

func (repo *UserRepositoryMock) CheckCredentials(email string, password string) (User, string, error) {
	repo.CheckCredentialsCall.Receives.Email = email
	repo.CheckCredentialsCall.Receives.Password = password

	return repo.CheckCredentialsCall.Returns.User,
		repo.CheckCredentialsCall.Returns.UserType,
		repo.CheckCredentialsCall.Returns.Error
}
