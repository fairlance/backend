package application

type FreelancerRepositoryMock struct {
	GetAllFreelancersCall struct {
		Returns struct {
			Freelancers []Freelancer
			Error       error
		}
	}
	GetFreelancerCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Freelancer Freelancer
			Error      error
		}
	}
	AddFreelancerCall struct {
		Receives struct {
			Freelancer *Freelancer
		}
		Returns struct {
			Error error
		}
	}
	UpdateFreelancerCall struct {
		Receives struct {
			Freelancer *Freelancer
		}
		Returns struct {
			Error error
		}
	}
	DeleteFreelancerCall struct {
		Receives struct {
			ID uint
		}
		Returns struct {
			Error error
		}
	}
	AddReviewCall struct {
		Receives struct {
			Review *Review
		}
		Returns struct {
			Error error
		}
	}
}

func (repo *FreelancerRepositoryMock) GetAllFreelancers() ([]Freelancer, error) {
	return repo.GetAllFreelancersCall.Returns.Freelancers,
		repo.GetAllFreelancersCall.Returns.Error
}

func (repo *FreelancerRepositoryMock) GetFreelancer(id uint) (Freelancer, error) {
	repo.GetFreelancerCall.Receives.ID = id
	return repo.GetFreelancerCall.Returns.Freelancer,
		repo.GetFreelancerCall.Returns.Error
}

func (repo *FreelancerRepositoryMock) AddFreelancer(freelancer *Freelancer) error {
	repo.AddFreelancerCall.Receives.Freelancer = freelancer
	return repo.AddFreelancerCall.Returns.Error
}

func (repo *FreelancerRepositoryMock) UpdateFreelancer(freelancer *Freelancer) error {
	repo.UpdateFreelancerCall.Receives.Freelancer = freelancer
	return repo.UpdateFreelancerCall.Returns.Error
}

func (repo *FreelancerRepositoryMock) DeleteFreelancer(id uint) error {
	repo.DeleteFreelancerCall.Receives.ID = id
	return repo.DeleteFreelancerCall.Returns.Error
}

func (repo *FreelancerRepositoryMock) AddReview(newReview *Review) error {
	repo.AddReviewCall.Receives.Review = newReview
	return repo.AddReviewCall.Returns.Error
}
