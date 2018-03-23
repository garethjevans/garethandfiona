all: release

clean:
	rm -rf garethandfiona Godeps vendor

install: clean prepare build
	glide install

prepare: clean
	go get github.com/Masterminds/glide

build: clean prepare
	glide update
	go fmt

test: clean prepare build install
	go test ./... -cover
	go vet .

release: clean prepare build install test

.PHONY: clean install prepare build test release
