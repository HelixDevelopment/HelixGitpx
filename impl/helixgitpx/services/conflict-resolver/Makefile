SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/conflict-resolver ./cmd/conflict-resolver
image:
	docker build -t helixgitpx/conflict-resolver:dev -f deploy/Dockerfile ../..
