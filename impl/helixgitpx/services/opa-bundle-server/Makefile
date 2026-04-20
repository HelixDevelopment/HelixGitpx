SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/opa-bundle-server ./cmd/opa-bundle-server
image:
	docker build -t helixgitpx/opa-bundle-server:dev -f deploy/Dockerfile ../..
