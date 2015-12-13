package mailer

import (
	"errors"
	"github.com/mailgun/mailgun-go"
	"os"
)

var (
	emailFrom    = "Fairlance <welcome@fairlance.io>"
	emailTitle   = "Welcome"
	emailContent = WelcomeMessage
)

type MailgunMailer struct{}

// Send welcome message
func (m MailgunMailer) SendWelcomeMessage(email string) (string, error) {
	var publicApiKey = os.Getenv("MAILGUN_PUBLIC_API_KEY")
	var apiKey = os.Getenv("MAILGUN_API_KEY")
	var domain = os.Getenv("MAILGUN_DOMAIN")

	if apiKey != "" && domain != "" && email != "" {
		mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
		m := mg.NewMessage(emailFrom, emailTitle, emailContent, email)
		_, id, err := mg.Send(m)
		return id, err
	}

	return "", errors.New("Set MAILGUN_API_KEY and MAILGUN_DOMAIN")
}
