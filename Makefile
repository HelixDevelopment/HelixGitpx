# HelixGitpx root Makefile — thin orchestrator delegating to impl/<subdir>/
SHELL := /usr/bin/env bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory

COMPOSE := impl/helixgitpx-platform/compose/bin/compose
COMPOSE_FILE := impl/helixgitpx-platform/compose/compose.yml

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ {printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: bootstrap
bootstrap: ## Fetch deps for every sub-project
	cd impl/helixgitpx && go work sync && go mod download
	cd impl/helixgitpx-web && npm install
	cd impl/helixgitpx-clients && ./gradlew --refresh-dependencies help
	cd impl/helixgitpx-docs && npm install

.PHONY: gen
gen: ## Regenerate protobuf/OpenAPI code
	cd impl/helixgitpx && buf generate

.PHONY: dev
dev: ## Bring up compose stack + hello service
	$(COMPOSE) --profile all up -d --build
	@echo "hello REST:  http://localhost:8001/v1/hello?name=world"
	@echo "hello gRPC:  localhost:9001"
	@echo "Grafana:     http://localhost:3000 (admin/admin)"
	@echo "Jaeger:      http://localhost:16686"

.PHONY: dev-down
dev-down: ## Tear down compose stack
	$(COMPOSE) --profile all down -v

.PHONY: test
test: ## Run tests across all sub-projects
	cd impl/helixgitpx && go test -race ./...
	cd impl/helixgitpx-web && npx nx run-many -t test
	cd impl/helixgitpx-clients && ./gradlew check
	cd impl/helixgitpx-docs && npm test --if-present

.PHONY: lint
lint: ## Lint across all sub-projects
	cd impl/helixgitpx && golangci-lint run ./... && buf lint proto/
	cd impl/helixgitpx-web && npx nx run-many -t lint
	cd impl/helixgitpx-clients && ./gradlew detekt ktlintCheck
	cd impl/helixgitpx-docs && npm run lint --if-present

.PHONY: build
build: ## Build all sub-projects
	cd impl/helixgitpx && go build ./...
	cd impl/helixgitpx-web && npx nx run-many -t build
	cd impl/helixgitpx-clients && ./gradlew assemble
	cd impl/helixgitpx-docs && npm run build

.PHONY: docs
docs: ## Build and serve the Docusaurus documentation site
	cd impl/helixgitpx-docs && node sync-docs.mjs && npm run start -- --port 3001

.PHONY: docs-build
docs-build: ## Build Docusaurus site (no serve)
	cd impl/helixgitpx-docs && node sync-docs.mjs && npm run build

.PHONY: runbook-lint
runbook-lint: ## Check runbooks conform to the template
	bash docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_lint.sh

.PHONY: clean
clean: ## Remove build artifacts (keeps compose volumes)
	cd impl/helixgitpx && go clean -cache -testcache
	cd impl/helixgitpx-web && rm -rf dist .nx
	cd impl/helixgitpx-clients && ./gradlew clean || true
	cd impl/helixgitpx-docs && rm -rf build .docusaurus
