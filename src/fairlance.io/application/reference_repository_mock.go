package application

type ReferenceRepositoryMock struct {
	AddReferenceCall struct {
		Receives struct {
			Reference *Reference
			ID        uint
		}
		Returns struct {
			Error error
		}
	}
}

func (repo *ReferenceRepositoryMock) AddReference(id uint, reference *Reference) error {
	repo.AddReferenceCall.Receives.Reference = reference
	repo.AddReferenceCall.Receives.ID = id
	return repo.AddReferenceCall.Returns.Error
}
