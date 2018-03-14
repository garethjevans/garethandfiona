FROM golang:1.9

ADD . /go/src/garethandfiona

RUN cd /go/src/garethandfiona && \
   make

ENTRYPOINT cd /go/src/garethandfiona/;go run main.go

EXPOSE 8080
