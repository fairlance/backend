package messaging

type Project struct {
	ID          uint
	Freelancers []Freelancer
	Client      *Client
}

type Client struct {
	ID        uint
	FirstName string
	LastName  string
}

type Freelancer struct {
	ID        uint
	FirstName string
	LastName  string
}
