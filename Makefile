BUILD_VERSION ?= 0.1.0
BUILD_DATE ?= $(shell date +%F\ %H:%M:%S)
BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)

help: ## display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

build: gophermart
.PHONY: build

gophermart: ## build gophermart
	go build \
		-ldflags "\
			-X 'main.buildVersion=$(BUILD_VERSION)' \
			-X 'main.buildDate=$(BUILD_DATE)' \
			-X 'main.buildCommit=$(BUILD_COMMIT)' \
		" \
		-o cmd/$@/$@ \
		cmd/$@/*.go
.PHONY: gophermart

clean: ## remove build artifacts
	rm -rf cmd/gophermart/gophermart
.PHONY: clean

unit-tests: ## run unit tests
	@go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v -E "(_mock|.pb).go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
.PHONY: unit-tests
