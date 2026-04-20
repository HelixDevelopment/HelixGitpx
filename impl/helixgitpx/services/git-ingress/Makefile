SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/git-ingress ./cmd/git-ingress
image:
	docker build -t helixgitpx/git-ingress:dev -f deploy/Dockerfile ../..
