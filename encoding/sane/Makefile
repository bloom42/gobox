.PHONY: all test release build

VERSION := $(shell cat version.go| grep "\sVersion\s=" | cut -d '"' -f2)

build:
	go build ./...

all: test build

test:
	go vet ./...
	go test -v -race ./...

release:
	git tag v$(VERSION)
	git push origin v$(VERSION)
