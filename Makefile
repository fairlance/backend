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
