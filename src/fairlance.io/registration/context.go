package registration

import (
	"log"
	"os"

	"fairlance.io/mailer"
)

type RegistrationContext struct {
	RegisteredUserRepository *RegisteredUserRepository
	Mailer                   mailer.Mailer
	Logger                   *log.Logger
}

func NewContext(dbName string) *RegistrationContext {
	file, err := os.OpenFile("/var/log/fairlance.io/registration.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file: %v", err)
	}

	logger := log.New(file, "register: ", log.Ldate|log.Ltime|log.Lshortfile)

	registeredUserRepository, err := NewRegisteredUserRepository(dbName)
	if err != nil {
		logger.Fatalf("Failed to open user repository: %q", err.Error())
	}

	// Setup context
	context := &RegistrationContext{
		RegisteredUserRepository: registeredUserRepository,
		Mailer: mailer.MailgunMailer{},
		Logger: logger,
	}

	return context
}
