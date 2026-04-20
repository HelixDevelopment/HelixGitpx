SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/upstream ./cmd/upstream
image:
	docker build -t helixgitpx/upstream:dev -f deploy/Dockerfile ../..
