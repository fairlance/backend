# Fairlance backend

![CircleCI](https://circleci.com/gh/fairlance/backend.svg?style=shield&circle-token=274b1fc821de530df06b3cc3e99b599c12abfaab
 "")


#### Available *make* commands
```bash
Available targets:
- install: install all services
- test: run tests
- testShort: run short tests, without using the db
- installDependencies: installs dependencies declared in dependencies.txt
- installARM: install all services with GOARCH=arm GOARM=7
```

#### Run tests with:
```bash
make test
```

#### Import test data to DB and index in Bleve
```bash
run fairlance_importer service (sudo service fairlance_importer start)
go to http://localhost:3004
Use the interface to import
```

#### Structure:
```
├── bin
│   ├── registration                <---- Service executable
│   └── application
├── pkg
│   └── linux_amd64
│       ├── github.com/fairlance
│       │   ├── mailer.a
│       │   └── registration.a
│       └── ...
└── src
    └── github.com/fairlance/
        ├── cmd
        │   ├── application             <---- contains main function, used to build an executable; package main
        │   │   └── main.go
        │   │
        │   └── registration
        │       └── main.go
        │
        ├── mailer                  <---- Utility package; package mailer
        │   ├── mailer.go
        │   └── mailgun.go
        │
        ├── registration            <---- Package that contains a service; package registration
        │   ├── context.go
        │   ├── handlers.go
        │   ├── handlers_test.go
        │   ├── model.go
        │   └── user_repository.go
        │
        └── application             <---- Another package; package application
            ├── context.go
            ├── handlers.go
            ├── handlers_test.go
            ├── model.go
            └── repository.go
```
