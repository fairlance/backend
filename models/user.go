package models

type User struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" valid:"required,email"`
	Image     string `json:"image"`
}
