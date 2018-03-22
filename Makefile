all: release

clean:
	rm -rf garethandfiona Godeps vendor

install: clean prepare build
	godep go install

prepare: clean
	go get github.com/tools/godep
	go get github.com/gorilla/mux
	go get github.com/newrelic/go-agent
	go get github.com/magiconair/properties
	go get github.com/gorilla/schema
	go get github.com/go-sql-driver/mysql
	go get github.com/simplereach/timeutils

build: clean prepare
	godep save ./...
	godep go build

test: clean prepare build install
	go test ./... -cover
	go vet .

release: clean prepare build install test

.PHONY: clean install prepare build test release
