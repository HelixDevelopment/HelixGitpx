SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/audit ./cmd/audit
image:
	docker build -t helixgitpx/audit:dev -f deploy/Dockerfile ../..
