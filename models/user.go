package models

type User struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Image      string `json:"image"`
	Salutation string `json:"salutation"`
	IsCompany  bool   `json:"isCompany"`
}
