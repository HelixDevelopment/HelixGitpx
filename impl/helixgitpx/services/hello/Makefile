SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/hello ./cmd/hello
image:
	docker build -t helixgitpx/hello:dev -f deploy/Dockerfile ../..
