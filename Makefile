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

.PHONY: build buildAll saveImages
build:
	GOOS=linux go build -o service ./cmd/${service}
	docker build -t fairlance/${service} .
	rm -f service
buildAll:
	make service=application build
	make service=fileserver build
	make service=search build
	make service=searcher build
	make service=importer build
	make service=messaging build
	make service=notification build
	make service=payment build
saveImages:
	docker save -o images/application_image fairlance/application
	docker save -o images/fileserver_image fairlance/fileserver
	docker save -o images/search_image fairlance/search
	docker save -o images/searcher_image fairlance/searcher
	docker save -o images/importer_image fairlance/importer
	docker save -o images/messaging_image fairlance/messaging
	docker save -o images/notification_image fairlance/notification
	docker save -o images/payment_image fairlance/payment
