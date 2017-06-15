FROM golang:1.8.3-alpine

RUN apk update && apk add git

ADD dependencies.txt .
RUN cat dependencies.txt | xargs go get

RUN mkdir -p /go/src/github.com/fairlance/backend/
WORKDIR /go/src/github.com/fairlance/backend/
