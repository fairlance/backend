FROM alpine

COPY service /service

ENTRYPOINT ["/service"]