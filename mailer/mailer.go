package mailer

import mailgun "github.com/mailgun/mailgun-go"
import "log"
import "fmt"

// Mailer send emails
type Mailer interface {
	SendProjectFunded(id uint, name string, clientID uint, clientName string) error
}

// Options holds option needed for the mailgun api
type Options struct {
	PublicApiKey string
	ApiKey       string
	Domain       string
	Self         string
}

// NewMailgunMailer creates a new mailer
func NewMailgun(o Options) Mailer {
	return &mailer{
		mailgun: mailgun.NewMailgun(
			o.Domain,
			o.ApiKey,
			o.PublicApiKey,
		),
		Self: o.Self,
	}
}

type mailer struct {
	mailgun mailgun.Mailgun
	Self    string
}

func (m *mailer) SendProjectFunded(projectID uint, projectName string, clientID uint, clientName string) error {
	title := fmt.Sprintf(projectFundedTitle, clientName, clientID, projectName, projectID)
	mail := m.mailgun.NewMessage("office@fairlance.io", title, " ", m.Self)
	mail.AddHeader("Content-Type", "text/html")
	resp, msgID, err := m.mailgun.Send(mail)
	if err != nil {
		log.Printf("could not send email: %v", err)
		return err
	}
	log.Printf("project funded email sent; id: %s;resp: %s => client (%d) has funded project (%d)", msgID, resp, clientID, projectID)
	return nil
}
