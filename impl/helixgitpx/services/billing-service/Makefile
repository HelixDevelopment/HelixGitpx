SERVICE := billing-service
IMAGE   ?= ghcr.io/helixgitpx/$(SERVICE):dev

.PHONY: build test docker
build:
	GOTOOLCHAIN=go1.23.4 go build -o bin/$(SERVICE) ./cmd/$(SERVICE)
test:
	GOTOOLCHAIN=go1.23.4 go test ./...
docker:
	docker build -t $(IMAGE) -f deploy/Dockerfile .
