# services/repo-service/Makefile
# Canonical service Makefile — every HelixGitpx Go service ships with this shape.
# Many targets delegate to shared scripts in tools/make/*.mk.

SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

SERVICE      := repo-service
MODULE       := github.com/vasic-digital/helixgitpx
REGISTRY     := ghcr.io/vasic-digital/helixgitpx
IMAGE        := $(REGISTRY)/$(SERVICE)
VERSION      ?= dev
COMMIT       := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO           ?= go
GOFLAGS      ?= -trimpath
LDFLAGS      := -s -w \
    -X '$(MODULE)/pkg/buildinfo.Version=$(VERSION)' \
    -X '$(MODULE)/pkg/buildinfo.Commit=$(COMMIT)' \
    -X '$(MODULE)/pkg/buildinfo.BuildTime=$(BUILD_TIME)'

PKG          := ./services/$(SERVICE)/...
COVERAGE_OUT := coverage.out
GOLANGCI     := $(shell command -v golangci-lint 2>/dev/null || echo "./bin/golangci-lint")

## --- Help ----------------------------------------------------
.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## --- Generate ------------------------------------------------
.PHONY: generate
generate: generate-proto generate-sql ## Run all code generators.

.PHONY: generate-proto
generate-proto: ## Regenerate protobuf Go + grpc-gateway + docs.
	buf generate

.PHONY: generate-sql
generate-sql: ## Run sqlc to regenerate typed DB code from .sql queries.
	sqlc generate -f services/$(SERVICE)/db/sqlc.yaml

## --- Build / Run ---------------------------------------------
.PHONY: build
build: ## Build the binary under ./bin/.
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(SERVICE) ./services/$(SERVICE)/cmd/$(SERVICE)

.PHONY: run
run: ## Run the service locally with devbox-provided deps.
	$(GO) run ./services/$(SERVICE)/cmd/$(SERVICE) serve --config config/dev.yaml

.PHONY: docker
docker: ## Build the container image.
	docker buildx build --load \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(IMAGE):$(VERSION) \
		-t $(IMAGE):$(COMMIT) \
		-f services/$(SERVICE)/Dockerfile .

.PHONY: docker-push
docker-push: ## Push the image + create immutable digest ref.
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(COMMIT)

## --- Lint / Format -------------------------------------------
.PHONY: fmt
fmt: ## Format code (gofumpt + goimports).
	gofumpt -w services/$(SERVICE) pkg internal
	goimports -w -local $(MODULE) services/$(SERVICE) pkg internal

.PHONY: lint
lint: ## Run golangci-lint.
	$(GOLANGCI) run $(PKG) ./pkg/... ./internal/...

.PHONY: vet
vet: ## Run go vet.
	$(GO) vet $(PKG)

.PHONY: check
check: fmt vet lint ## Format + vet + lint.

## --- Tests ---------------------------------------------------
.PHONY: test
test: ## Run unit tests with coverage + race detector.
	$(GO) test -race -count=1 -timeout=120s -coverprofile=$(COVERAGE_OUT) -covermode=atomic $(PKG)

.PHONY: test-integration
test-integration: ## Run integration tests (Testcontainers).
	$(GO) test -tags=integration -race -count=1 -timeout=600s $(PKG)

.PHONY: test-e2e
test-e2e: ## Run E2E tests against staging.
	E2E_TARGET?=https://api.staging.helixgitpx.example.com \
	$(GO) test -tags=e2e -count=1 -timeout=900s $(PKG)

.PHONY: bench
bench: ## Run Go benchmarks.
	$(GO) test -run=NONE -bench=. -benchmem -count=5 $(PKG) | tee bench-$(VERSION).txt

.PHONY: mutation
mutation: ## Mutation tests (gremlins) on changed files.
	gremlins unleash --tags-filter='!integration,!e2e' $(PKG)

.PHONY: coverage
coverage: test ## Show coverage report.
	$(GO) tool cover -html=$(COVERAGE_OUT) -o coverage.html
	$(GO) tool cover -func=$(COVERAGE_OUT)

## --- Security ------------------------------------------------
.PHONY: sec
sec: ## Run govulncheck + gosec + gitleaks.
	govulncheck $(PKG)
	gosec -severity=medium $(PKG)
	gitleaks detect --source=. --redact --no-git

## --- SBOM / Sign ---------------------------------------------
.PHONY: sbom
sbom: ## Generate a CycloneDX SBOM for the image.
	syft $(IMAGE):$(VERSION) -o cyclonedx-json > sbom.cdx.json

.PHONY: sign
sign: ## Cosign keyless sign the image.
	COSIGN_EXPERIMENTAL=1 cosign sign --yes $(IMAGE):$(VERSION)

.PHONY: verify
verify: ## Cosign keyless verify the image.
	COSIGN_EXPERIMENTAL=1 cosign verify $(IMAGE):$(VERSION) \
		--certificate-identity-regexp='https://github.com/vasic-digital/helixgitpx/.*' \
		--certificate-oidc-issuer='https://token.actions.githubusercontent.com'

## --- Deploy --------------------------------------------------
.PHONY: helm-lint
helm-lint: ## Lint the Helm chart.
	helm lint deploy/helm/$(SERVICE)

.PHONY: deploy-dev
deploy-dev: docker ## Deploy to local Kind cluster via Tilt.
	tilt up

## --- Database migrations -------------------------------------
.PHONY: migrate-up
migrate-up: ## Apply pending migrations to the local DB.
	goose -dir services/$(SERVICE)/db/migrations postgres "$$DATABASE_URL" up

.PHONY: migrate-down
migrate-down: ## Roll back one migration in the local DB.
	goose -dir services/$(SERVICE)/db/migrations postgres "$$DATABASE_URL" down

.PHONY: migrate-status
migrate-status:
	goose -dir services/$(SERVICE)/db/migrations postgres "$$DATABASE_URL" status

## --- Clean ---------------------------------------------------
.PHONY: clean
clean:
	rm -rf bin/ dist/ $(COVERAGE_OUT) coverage.html

## --- CI convenience -----------------------------------------
.PHONY: ci
ci: check test sec ## CI pipeline target.
