help:
	@echo "Available targets:"
	@echo "- install: install all services"
	@echo "- installARM: install all services with GOARCH=arm GOARM=7"
	@echo "- test: run tests"
	@echo "- testShort: run short tests, without using the db"
	@echo "- installDependencies: installs dependencies declared in dependencies.txt"

installDependencies:
	cat dependencies.txt | xargs go get

test:
	 go list fairlance.io/... | grep -v /cmd/ | xargs go test -v

install:
	go install fairlance.io/cmd/...

installARM:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/registration_arm fairlance.io/cmd/registration
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/application_arm fairlance.io/cmd/application
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/search_arm fairlance.io/cmd/search
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/importer_arm fairlance.io/cmd/importer
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/searcher_arm fairlance.io/cmd/searcher
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/messaging_arm fairlance.io/cmd/messaging
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/fileserver_arm fairlance.io/cmd/fileserver
