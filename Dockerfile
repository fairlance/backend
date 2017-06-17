FROM alpine

RUN apk update && apk add ca-certificates

ENTRYPOINT ["/service"]

COPY service /service