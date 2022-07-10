VERS ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

FOUND_GO_VERSION := $(shell go version)
EXPECTED_GO_VERSION = 1.17
.PHONY: check-go-version
check-go-version:
	@$(if $(findstring ${EXPECTED_GO_VERSION}, ${FOUND_GO_VERSION}),(exit 0),(echo Wrong go version! Please install ${EXPECTED_GO_VERSION}; exit 1))

test: check-go-version
	@echo "running all tests"
	@go install ./...
	@go fmt ./...
	go vet ./...
	go test ./...

run: test
	@echo "starting application"
	docker-compose up

stop:
	docker-compose down
