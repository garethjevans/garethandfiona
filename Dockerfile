FROM golang:1.8

ADD . /go/src/github.com/garethjevans/garethandfiona

RUN cd /go/src/github.com/garethjevans/garethandfiona;make

ENTRYPOINT cd /go/src/github.com/garethjevans/garethandfiona/;go run main.go

EXPOSE 8080
