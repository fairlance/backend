package payment

type Project struct {
	ID          uint         `json:"id"`
	Freelancers []Freelancer `json:"freelancers"`
	Contract    Contract     `json:"contract"`
	Name        string       `json:"name"`
}

func (p *Project) amount() float64 {
	return float64(p.Contract.Hours) * p.Contract.PerHour
}

type Freelancer struct {
	ID    uint   `json:"id"`
	Email string `json:"payPalEmail"`
}

type Contract struct {
	Hours   int     `json:"hours"`
	PerHour float64 `json:"perHour"`
}
