package registration

type RegisteredUser struct {
	Email string `json:"email"`
}

type RegisteredError struct {
	Error string `json:"error"`
}
