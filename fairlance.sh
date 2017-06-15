#!/bin/bash

CMD=$1
case $CMD in
    init )
        echo 'eval $("C:\Program Files\Docker Toolbox\docker-machine.exe" env --shell=bash)'
        ;;
    dependancies )
        echo "docker build -t fairlance/backend-dependancies -f dependancies.Dockerfile ."
        docker build -t fairlance/backend-dependancies -f dependancies.Dockerfile .
        ;;
    ssh )
        echo 'docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependancies bash'
        docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependancies bash
        ;;
    build )
        START=$(date +%s)
        docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependancies bash -c "GOOS=linux CGO_ENABLED=0 go build -o service ./cmd/$2"
        docker build -t fairlance/$2 .
        rm -f service
        END=$(date +%s)
        DIFF=$(( $END - $START ))
        echo "It took $DIFF seconds"
        ;;
    buildAll )
        ./build.sh build application
        ./build.sh build fileserver
        ./build.sh build search
        ./build.sh build searcher
        ./build.sh build importer
        ./build.sh build messaging
        ./build.sh build notification
        ./build.sh build payment
        ;;
    saveImages )
        docker save -o images/application_image fairlance/application
        docker save -o images/fileserver_image fairlance/fileserver
        docker save -o images/search_image fairlance/search
        docker save -o images/searcher_image fairlance/searcher
        docker save -o images/importer_image fairlance/importer
        docker save -o images/messaging_image fairlance/messaging
        docker save -o images/notification_image fairlance/notification
        docker save -o images/payment_image fairlance/payment
        ;;
     *)
        echo $"Usage: $0 {init|dependancies|ssh|build|buildAll|saveImages}"
        exit 1
esac
