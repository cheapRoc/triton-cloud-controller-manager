SHELL := /bin/bash

## Display this help message
help:
		@awk '/^##.*$$/,/[a-zA-Z_-]+:/' $(MAKEFILE_LIST) | awk '!(NR%2){print $$0p}{p=$$0}' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' | sort

## Install dev/test CLI tooling
tools:
		@go version | grep 1.8 || (echo 'go1.8 not installed'; exit 1)
		@echo Installing dev/test tools/dependencies
		go get -u github.com/golang/dep/cmd/dep
		go get -u github.com/golang/lint/golint

## Cleanly install all dependencies
deps: clean tools
		dep ensure

## Update all dependencies if possible
deps-update:
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
