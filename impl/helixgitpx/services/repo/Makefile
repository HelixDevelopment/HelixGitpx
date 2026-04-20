SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/repo ./cmd/repo
image:
	docker build -t helixgitpx/repo:dev -f deploy/Dockerfile ../..
