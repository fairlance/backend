package payment

type Project struct {
	ID          uint
	Freelancers []Freelancer `json:"freelancers,omitempty"`
}

type Freelancer struct {
	ID    uint
	Email string
}
