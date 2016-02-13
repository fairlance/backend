package registration

type RegisteredUser struct {
    Email string `json:"email"`
}

type RegistrationError struct {
    Error string `json:"error"`
}
