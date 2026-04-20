SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/collab-service ./cmd/collab-service
image:
	docker build -t helixgitpx/collab-service:dev -f deploy/Dockerfile ../..
