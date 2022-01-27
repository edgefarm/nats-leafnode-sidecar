BIN_DIR ?= bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/edgefarm/pkg/version/version=$(VERSION)"

all: tidy test build ## default target: tidy, test, build

tidy: ## ensures go dependencies are up to date
	go mod tidy
	go mod vendor

client: tidy test ## build the client
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/client cmd/client/main.go

registry: tidy test ## build the registry
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/registry cmd/registry/main.go

build: test client registry ## build all


test-all: test integration-test ## run all tests

test: tidy ## run go tests
	go test ./...

integration-test: ## run integration test
	cd test/integration && ./test.sh

clean: ## cleanup
	rm -rf ${BIN_DIR}/client ${BIN_DIR}/registry

docker-amd64: client-docker-amd64 registry-docker-amd64 ## creates docker images for amd64 (prefix dev-)

client-docker-amd64: tidy test ## creates docker dev image for client (amd64 only)
	docker build -f build/client/Dockerfile -t ci4rail/dev-nats-leafnode-client:${VERSION} .
	docker push ci4rail/dev-nats-leafnode-client:${VERSION}

registry-docker-amd64: tidy test ## creates docker dev image for registry (amd64 only)
	docker build -f build/registry/Dockerfile -t ci4rail/dev-nats-leafnode-registry:${VERSION} .
	docker push ci4rail/dev-nats-leafnode-registry:${VERSION}

.PHONY: test clean build build-client all tidy test-all integration-test client-docker-amd64 registry-docker-amd64 docker-amd64

help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make [target]\033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m\t %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
