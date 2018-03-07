FROM golang:1.8

ADD . /go/src/github.com/cagiti/kirstenandchris

RUN cd /go/src/github.com/cagiti/kirstenandchris;make

ENTRYPOINT cd /go/src/github.com/cagiti/kirstenandchris/;go run main.go

EXPOSE 8080
