.PHONY: help installDependencies test install installARM
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
	 go list github.com/fairlance/backend/... | grep -v /cmd/ | xargs go test -v

install:
	go install github.com/fairlance/backend/cmd/...

installARM:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/registration_arm github.com/fairlance/backend/cmd/registration
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/application_arm github.com/fairlance/backend/cmd/application
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/search_arm github.com/fairlance/backend/cmd/search
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/importer_arm github.com/fairlance/backend/cmd/importer
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/searcher_arm github.com/fairlance/backend/cmd/searcher
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/messaging_arm github.com/fairlance/backend/cmd/messaging
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/fileserver_arm github.com/fairlance/backend/cmd/fileserver
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/notification_arm github.com/fairlance/backend/cmd/notification
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ./../bin/payment_arm github.com/fairlance/backend/cmd/payment
