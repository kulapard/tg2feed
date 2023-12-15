all: lint test build

build:
	docker build .

test:
	go clean -testcache
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	rm coverage.out

lint:
	golangci-lint run

.PHONY: build test lint