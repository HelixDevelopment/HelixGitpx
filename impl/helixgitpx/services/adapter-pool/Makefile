SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/adapter-pool ./cmd/adapter-pool
image:
	docker build -t helixgitpx/adapter-pool:dev -f deploy/Dockerfile ../..
