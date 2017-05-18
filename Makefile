SHELL := /bin/bash
IMPORT_PATH := github.com/cheapRoc/triton-cloud-controller-manager
VERSION ?= dev-build-not-for-release
LDFLAGS := -X ${IMPORT_PATH}/core.GitHash='$(shell git rev-parse --short HEAD)' -X ${IMPORT_PATH}/core.Version='${VERSION}'

## Display this help message
help:
		@awk '/^##.*$$/,/[a-zA-Z_-]+:/' $(MAKEFILE_LIST) | awk '!(NR%2){print $$0p}{p=$$0}' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' | sort

## Build the controller binary for the local OS only
build: build/binary
build/binary: */*.go *.go
		go build -o build/triton-cloud-controller-manager -ldflags "$(LDFLAGS)"

## Install dev/test CLI tooling
tools:
		@go version | grep 1.7 || (echo 'go1.7 not installed'; exit 1)
		@echo Installing tools and dependencies
		go get -u github.com/golang/dep/cmd/dep
		go get -u github.com/golang/lint/golint

## Cleanly install all dependencies
deps: tools
		dep ensure

## Update all dependencies if possible
deps-update: clean tools
		deps ensure -update

## Clean out all dependencies from vendor/
clean:
		rm -rf vendor

## Execute golint against the project
lint:
		@golint

## Execute all unit tests
test: tools
		@go test -v
