# Fairlance backend

![CircleCI](https://circleci.com/gh/fairlance/backend.svg?style=shield&circle-token=274b1fc821de530df06b3cc3e99b599c12abfaab
 "")


#### Run tests with:
```bash
make test
```

#### See available *make* commands
```bash
make
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
        ├── mailer                  <---- Utility package
        │   ├── mailer.go
        │   └── mailgun.go
        ├── registration            <---- Package that contains all service relevant functionality
        │   ├── context.go
        │   ├── handlers.go
        │   ├── handlers_test.go
        │   ├── model.go
        │   └── user_repository.go
        └── services                <---- Folder that contains all runnable services
            └── registration        <---- Every service is located in a folder with the same name
                └── registration.go <---- Service with package "main"
```
