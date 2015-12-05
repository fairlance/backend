package main

type Mailer interface {
	SendWelcomeMessage(string) (string, error)
}
