package application

type ApplicationContext struct {
    FreelancerRepository *FreelancerRepository
    JwtSecret            string
}

func NewContext(dbName string) *ApplicationContext {
    freelancerRepository, _ := NewFreelancerRepository(dbName)

    context := &ApplicationContext{
        FreelancerRepository:   freelancerRepository,
        JwtSecret:              "fairlance",//base64.StdEncoding.EncodeToString([]byte("fairlance")),
    }

    return context
}
