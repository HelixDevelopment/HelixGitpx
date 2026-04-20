SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/auth ./cmd/auth
image:
	docker build -t helixgitpx/auth:dev -f deploy/Dockerfile ../..
