# Fairlance backend

![CircleCI](https://circleci.com/gh/fairlance/backend.svg?style=shield&circle-token=274b1fc821de530df06b3cc3e99b599c12abfaab
 "")


#### Available *make* commands
```bash
Available targets:
- install: install all services
- test: run tests
- installDependenciesAndTest: install dependencies declared in dependencies.txt and run tests
- installDependencies: installs dependencies declared in dependencies.txt

```

#### Run tests with:
```bash
make test
```

#### Structure:
```
├── bin
│   └── registration                <---- Service executable
├── pkg
│   └── linux_amd64
│       ├── fairlance.io
│       │   ├── mailer.a
│       │   └── registration.a
│       └── ...
└── src
    └── fairlance.io/
        ├── mailer                  <---- Utility package; package mailer
        │   ├── mailer.go
        │   └── mailgun.go
        │
        ├── registration            <---- Package that contains a service; generates a binary; package main
        │   ├── main.go             <---- contains main function
        │   ├── context.go
        │   ├── handlers.go
        │   ├── handlers_test.go
        │   ├── model.go
        │   └── user_repository.go
        │
        └── application             <---- Second package that also generates a binary; package main
            ├── main.go
            ├── context.go
            ├── handlers.go
            ├── handlers_test.go
            ├── model.go
            └── repository.go
```
