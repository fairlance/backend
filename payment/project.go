package payment

type Project struct {
	ID          uint         `json:"id"`
	Freelancers []Freelancer `json:"freelancers"`
}

type Freelancer struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}
