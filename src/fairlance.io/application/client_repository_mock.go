package application

type ClientRepositoryMock struct {
	GetAllClientsCall struct {
		Returns struct {
			Clients []Client
			Error   error
		}
	}
	AddClientCall struct {
		Receives struct {
			Client *Client
		}
		Returns struct {
			Error error
		}
	}
	UpdateClientCall struct {
		Receives struct {
			Client *Client
		}
		Returns struct {
			Error error
		}
	}
	GetClientCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Client *Client
			Error  error
		}
	}
}

func (repo *ClientRepositoryMock) GetAllClients() ([]Client, error) {
	return repo.GetAllClientsCall.Returns.Clients,
		repo.GetAllClientsCall.Returns.Error
}

func (repo *ClientRepositoryMock) AddClient(client *Client) error {
	repo.AddClientCall.Receives.Client = client
	return repo.AddClientCall.Returns.Error
}

func (repo *ClientRepositoryMock) UpdateClient(client *Client) error {
	repo.UpdateClientCall.Receives.Client = client
	return repo.UpdateClientCall.Returns.Error
}

func (repo *ClientRepositoryMock) GetClient(id uint) (*Client, error) {
	repo.GetClientCall.Receives.ID = id
	return repo.GetClientCall.Returns.Client,
		repo.GetClientCall.Returns.Error
}
