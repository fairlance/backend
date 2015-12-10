package registration

import "fairlance.io/mailer"

type RegistrationContext struct {
	userRepository *UserRepository
	mailer         mailer.Mailer
	// ... and the rest of our globals.
}

func NewContext(db string) *RegistrationContext {
	// Setup context
	context := &RegistrationContext{userRepository: NewUserRepository(db), mailer: mailer.MailgunMailer{}}

	return context
}
