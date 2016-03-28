help:
	@echo "Available targets:"
	@echo "- install: install all services"
	@echo "- test: run tests"
	@echo "- installDependenciesAndTest: install dependencies declared in dependencies.txt and run tests"
	@echo "- installDependencies: installs dependencies declared in dependencies.txt"

installDependencies:
	cat dependencies.txt | xargs go get

test:
	go test fairlance.io/... -v

installDependenciesAndTest: installDependencies test

install:
	go install fairlance.io/services/...

installARM:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/registration_arm fairlance.io/services/registration
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/application_arm fairlance.io/services/application