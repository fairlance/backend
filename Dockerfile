FROM alpine

RUN apk update && apk add ca-certificates

COPY service /service

ENTRYPOINT ["/service"]
