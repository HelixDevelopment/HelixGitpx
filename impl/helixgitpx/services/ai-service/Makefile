SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/ai-service ./cmd/ai-service
image:
	docker build -t helixgitpx/ai-service:dev -f deploy/Dockerfile ../..
