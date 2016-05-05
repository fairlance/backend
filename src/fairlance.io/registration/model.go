package main

import "time"

type RegisteredUser struct {
	Email   string    `json:"email"`
	Created time.Time `json:"created,omitempty"`
}

type RegistrationError struct {
	Error string `json:"error"`
}
