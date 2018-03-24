all: release

clean:
	rm -rf garethandfiona Godeps vendor

install: clean prepare build
	glide install

prepare: clean
	go get github.com/Masterminds/glide

build: clean prepare
	glide update
	go build
	go fmt
	go vet .

test: clean prepare build install
	go test ./... -cover

release: clean prepare build install test

.PHONY: clean install prepare build test release
