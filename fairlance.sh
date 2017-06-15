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
        ./fairlance.sh build application
        ./fairlance.sh build fileserver
        ./fairlance.sh build search
        ./fairlance.sh build searcher
        ./fairlance.sh build importer
        ./fairlance.sh build messaging
        ./fairlance.sh build notification
        ./fairlance.sh build payment
        ;;
    saveImages )
        docker save -o application_image fairlance/application
        docker save -o fileserver_image fairlance/fileserver
        docker save -o search_image fairlance/search
        docker save -o searcher_image fairlance/searcher
        docker save -o importer_image fairlance/importer
        docker save -o messaging_image fairlance/messaging
        docker save -o notification_image fairlance/notification
        docker save -o payment_image fairlance/payment
        ;;
     *)
        echo $"Usage: $0 {init|dependancies|ssh|build|buildAll|saveImages}"
        exit 1
esac
