SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/sync-orchestrator ./cmd/sync-orchestrator
image:
	docker build -t helixgitpx/sync-orchestrator:dev -f deploy/Dockerfile ../..
