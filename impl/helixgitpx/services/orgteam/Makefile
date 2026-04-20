SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/orgteam ./cmd/orgteam
image:
	docker build -t helixgitpx/orgteam:dev -f deploy/Dockerfile ../..
