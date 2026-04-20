SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/live-events-service ./cmd/live-events-service
image:
	docker build -t helixgitpx/live-events-service:dev -f deploy/Dockerfile ../..
