package registration

import (
    "fairlance.io/mailer"
    "log"
    "os"
)

type RegistrationContext struct {
    registeredUserRepository    *RegisteredUserRepository
    mailer                      mailer.Mailer
    Logger                      *log.Logger
}

func NewContext(dbName string) *RegistrationContext {
    file, err := os.OpenFile("/var/log/fairlance.io/registration.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file: %v", err)
    }

    logger := log.New(file, "register: ", log.Ldate | log.Ltime | log.Lshortfile)

    registeredUserRepository, err := NewRegisteredUserRepository(dbName)
    if err != nil {
        logger.Println("Failed to open user repository: %v", err)
    }

    // Setup context
    context := &RegistrationContext{
        registeredUserRepository: registeredUserRepository,
        mailer:         mailer.MailgunMailer{},
        Logger:         logger,
    }

    return context
}
