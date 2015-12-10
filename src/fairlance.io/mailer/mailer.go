package mailer

type Mailer interface {
	SendWelcomeMessage(string) (string, error)
}
