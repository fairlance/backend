package mailer

import (
	"errors"
	"os"

	"github.com/mailgun/mailgun-go"
)

var (
	emailFrom  = "Fairlance <welcome@github.com/fairlance>"
	emailTitle = "Welcome"
)

type WelcomeMailgunMailer struct{}

// Send welcome message
func (m WelcomeMailgunMailer) SendWelcomeMessage(email string) (string, error) {
	var publicApiKey = os.Getenv("MAILGUN_PUBLIC_API_KEY")
	var apiKey = os.Getenv("MAILGUN_API_KEY")
	var domain = os.Getenv("MAILGUN_DOMAIN")

	if apiKey != "" && domain != "" && email != "" {
		mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
		m := mg.NewMessage(emailFrom, emailTitle, WelcomeMessage, email)
		m.SetHtml(HTMLWelcomeMessage)
		m.AddHeader("Content-Type", "text/html")
		_, id, err := mg.Send(m)
		return id, err
	}

	return "", errors.New("Set MAILGUN_API_KEY and MAILGUN_DOMAIN")
}
