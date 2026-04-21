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

# ---------------------------------------------------------------------------
# Test matrix (Constitution Article II).
# All seven types are mandatory. Mocks are allowed ONLY in test-unit.
# ---------------------------------------------------------------------------

.PHONY: test-unit test-integration test-e2e test-security test-stress test-ddos test-benchmark test-all
test-unit: ## Run unit tests (mocks allowed)
	cd impl/helixgitpx && GOTOOLCHAIN=go1.23.4 go test -race -cover ./...
	cd impl/helixgitpx-web && npx nx test web --watchAll=false
	cd impl/helixgitpx-clients && gradle :shared:jvmTest || true

test-integration: ## Run integration tests (compose must be up)
	cd test/integration && GOTOOLCHAIN=go1.23.4 go test -tags=integration -v ./...

test-e2e: ## Run end-to-end tests (k3d must be up)
	cd impl/helixgitpx-web && npx playwright test

test-security: ## Run security scans
	bash test/security/run.sh

test-stress: ## Run stress scenarios (k6 required)
	cd tools/perf && make all

test-ddos: ## Run ddos / rate-limit scenarios
	bash test/ddos/run.sh

test-benchmark: ## Run Go micro-benchmarks and k6 budget scenarios
	cd impl/helixgitpx && GOTOOLCHAIN=go1.23.4 go test -run='^$$' -bench=. -benchmem ./...
	cd tools/perf && python3 check_budgets.py /tmp/k6-out-*.json budgets.json || true

test-all: test-unit test-integration test-e2e test-security test-stress test-ddos test-benchmark ## Run all seven mandatory test types
	@echo "All seven required test types executed."

.PHONY: coverage-audit
coverage-audit: ## Per-package coverage audit; fails below threshold
	bash tools/coverage-audit/audit.sh

.PHONY: ci-local
ci-local: ## Run every check that the ci-verifiers workflow runs (green-suite)
	bash scripts/verify-everything.sh

.PHONY: verify-proto-gen verify-secrets
verify-proto-gen: ## Detect drift between proto/ and committed gen/
	bash scripts/verify-proto-gen.sh

verify-secrets: ## Scan the repo for committed secrets (gitleaks)
	bash scripts/verify-secrets.sh

# ---------------------------------------------------------------------------
# Upstream federation (Constitution Article IV §2).
# ---------------------------------------------------------------------------

.PHONY: upstream-push upstream-status
upstream-push: ## Push main + tags to ALL configured upstreams
	bash scripts/push-to-all-upstreams.sh

upstream-status: ## Show divergence for each configured upstream
	bash scripts/upstream-status.sh
