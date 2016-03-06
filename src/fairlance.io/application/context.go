package application

type ApplicationContext struct {
    FreelancerRepository *FreelancerRepository
}

func NewContext(dbName string) *ApplicationContext {
    freelancerRepository, _ := NewFreelancerRepository(dbName)

    context := &ApplicationContext{
        FreelancerRepository: freelancerRepository,
    }

    return context
}
