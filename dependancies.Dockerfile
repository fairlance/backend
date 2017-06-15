FROM golang:1.8

RUN go get gopkg.in/mgo.v2
RUN go get github.com/asaskevich/govalidator
RUN go get github.com/mailgun/mailgun-go
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/context
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/lib/pq
RUN go get github.com/jinzhu/gorm
RUN go get gopkg.in/matryer/respond.v1
RUN go get github.com/cheekybits/is
RUN go get github.com/blevesearch/bleve/...
RUN go get github.com/gorilla/websocket
RUN go get github.com/PuerkitoBio/goquery

RUN mkdir -p /go/src/github.com/fairlance/backend/
WORKDIR /go/src/github.com/fairlance/backend/
