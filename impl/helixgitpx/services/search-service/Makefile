SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/search-service ./cmd/search-service
image:
	docker build -t helixgitpx/search-service:dev -f deploy/Dockerfile ../..
