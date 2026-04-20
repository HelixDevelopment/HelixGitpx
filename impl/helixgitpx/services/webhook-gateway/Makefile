SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/webhook-gateway ./cmd/webhook-gateway
image:
	docker build -t helixgitpx/webhook-gateway:dev -f deploy/Dockerfile ../..
