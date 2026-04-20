SHELL := /usr/bin/env bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.PHONY: help lint test build gen fmt tidy cover

help:
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format code
	gofumpt -w .
	goimports -w .

lint: ## Run linters
	golangci-lint run ./...
	buf lint proto/

test: ## Run unit tests
	go test -race -shuffle=on -count=1 ./...

cover: ## Coverage report
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | tail -1

build: ## Build all binaries
	go build ./...

gen: ## Regenerate protobuf
	buf generate proto/

tidy: ## go work sync + tidy each module
	go work sync
	find . -name go.mod -not -path './gen/*' -execdir go mod tidy \;
