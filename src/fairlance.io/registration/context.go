package registration

import (
	"fairlance.io/mailer"
	"log"
	"os"
)

var context = NewContext("registration")

type RegistrationContext struct {
	userRepository *UserRepository
	mailer         mailer.Mailer
	Logger         *log.Logger
}

func NewContext(dbName string) *RegistrationContext {
	file, err := os.OpenFile("/var/log/fairlance.io/registration.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file: %v", err)
	}

	logger := log.New(file, "register: ", log.Ldate|log.Ltime|log.Lshortfile)

	userRepository, err := NewUserRepository(dbName)
	if err != nil {
		logger.Println("Failed to open user repository: %v", err)
	}

	// Setup context
	context := &RegistrationContext{
		userRepository: userRepository,
		mailer:         mailer.MailgunMailer{},
		Logger:         logger,
	}

	return context
}
