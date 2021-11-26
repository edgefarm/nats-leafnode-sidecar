BIN_DIR ?= bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/edgefarm/pkg/version/version=$(VERSION)"

all: tidy test build

tidy:
	go mod tidy
	go mod vendor

client: test tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/client cmd/client/main.go

registry: test tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/registry cmd/registry/main.go

build: client registry

test: tidy
	go test ./...

clean:
	rm -rf ${BIN_DIR}/client ${BIN_DIR}/registry

.PHONY: test clean build build-client all tidy
