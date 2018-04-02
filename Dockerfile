FROM golang:1.9

ADD . /go/src/garethandfiona

RUN cd /go/src/garethandfiona && \
   make build

EXPOSE 8080

WORKDIR /go/src/garethandfiona/
ENTRYPOINT go run main.go app.go rsvp.go wedding_database.go wedding_database_mysql.go
