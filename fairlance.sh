#!/bin/bash

CMD=$1
case $CMD in
    init )
        echo 'eval $("C:\Program Files\Docker Toolbox\docker-machine.exe" env --shell=bash)'
        ;;
    dependencies )
        echo "docker build -t fairlance/backend-dependencies -f dependencies.Dockerfile ."
        docker build -t fairlance/backend-dependencies -f dependencies.Dockerfile .
        ;;
    ssh )
        echo 'docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependencies sh'
        docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependencies sh
        ;;
    test )
        echo 'docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependencies sh -c "go list github.com/fairlance/backend/... | grep -v /cmd/ | xargs go test -v"'
        docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependencies sh -c "go list github.com/fairlance/backend/... | grep -v /cmd/ | xargs go test -v"
        ;;
    build )
        START=$(date +%s)
        docker run --rm -v "/$(pwd)":/go/src/github.com/fairlance/backend/ -it fairlance/backend-dependencies sh -c "GOOS=linux CGO_ENABLED=0 go build -o service ./cmd/$2" || exit 1
        docker build -t fairlance/$2 .
        rm -f service
        END=$(date +%s)
        DIFF=$(( $END - $START ))
        echo "It took $DIFF seconds"
        ;;
    buildAll )
        START_ALL=$(date +%s)
        ./fairlance.sh build application
        ./fairlance.sh build fileserver
        ./fairlance.sh build search
        ./fairlance.sh build searcher
        ./fairlance.sh build importer
        ./fairlance.sh build messaging
        ./fairlance.sh build notification
        ./fairlance.sh build payment
        END_ALL=$(date +%s)
        DIFF_ALL=$(( $END_ALL - $START_ALL ))
        echo "It took $DIFF_ALL seconds"
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
    ngrok )
        docker run -p 4040:4040 --net="host" fnichol/ngrok 8888
        ;;
     *)
        echo $"Usage: $0 {init|dependencies|ssh|build|buildAll|saveImages|ngrok}"
        exit 1
esac
