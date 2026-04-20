# M1 Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver the full 18-item M1 Foundation from `docs/specifications/.../HelixGitpx/13-roadmap/17-milestones.md` §§1.1–1.4 — working monorepo skeleton, `hello` service reachable via gRPC+REST through a runtime-portable compose stack, and every CI/infra/docs artifact the roadmap calls for (those needing infra ship as config that activates in M2+).

**Architecture:** Five logical sub-projects under `impl/` in one git history: Go monorepo (`helixgitpx`), Nx web (`helixgitpx-web`), KMP clients (`helixgitpx-clients`), GitOps/infra (`helixgitpx-platform`), Docusaurus docs (`helixgitpx-docs`). Shared Go `platform` module (14 packages). Spine-first: a vertical slice (hello service + compose + CI) is wired before breadth-fill (scanners, Kyverno, Kata/Vault config, design system scaffolds, ADRs, runbook template, Docusaurus).

**Tech Stack:** Go 1.23, Angular 19 + Nx, Kotlin Multiplatform + Compose, `buf` + protobuf + Connect, `mise` for toolchain pins, compose via docker-or-podman wrapper, GitHub Actions (workflow_dispatch only), Kyverno, Checkov, Docusaurus.

**Locked constraints (from design spec §4):**
- C-1 — code under `impl/<repo>/`, single git history
- C-2 — GitHub Actions only, `workflow_dispatch`-only (mandatory)
- C-3 — compose-based local dev, portable docker/podman
- C-4 — full 18 items (config-only where infra absent)
- C-5 — `mise.toml` toolchain
- C-6 — spine-first sequencing

**Phases:**
- **Phase A — Spine** (Tasks 1–14): monorepo layout, shared Go libs, proto root, hello service, compose stack, end-to-end greeting works.
- **Phase B — Breadth** (Tasks 15–32): CI workflows, Kyverno/Checkov, Kata runner config, Vault/OIDC, Nx web workspace, KMP clients, Docusaurus, ADRs, runbook template, completion-matrix verification.

**Conventions every task follows:**
- Commits are Conventional Commits with `-s` (DCO sign-off).
- Branch: keep working on `main` (solo mode; SOLO-NOTES.md documents deviation from 2-approver rule).
- Go packages: TDD. Write failing test → run → implement → run → commit.
- Config files: lint-first policy — write, then run the linter/validator shown in the step.
- Every task ends with a `git commit` step.

---

## File Structure (all created by this plan)

```
HelixGitpx/
├── CHANGELOG.md                       (new)
├── README.md                          (new — replaces 1-line existing)
├── SOLO-NOTES.md                      (new)
├── Makefile                           (new — root orchestrator)
├── mise.toml                          (new)
├── .gitattributes                     (new — mark gen/ as linguist-generated)
├── .github/
│   ├── CODEOWNERS                     (new)
│   ├── branch-protection.json         (new — config artifact)
│   └── workflows/
│       ├── ci-go.yml                  (new)
│       ├── ci-web.yml                 (new)
│       ├── ci-clients.yml             (new)
│       ├── ci-docs.yml                (new)
│       ├── ci-platform.yml            (new)
│       ├── security-scan.yml          (new)
│       ├── supply-chain.yml           (new)
│       ├── release.yml                (new)
│       ├── deploy.yml                 (new)
│       └── _reusable/
│           ├── _setup-go.yml          (new)
│           ├── _setup-node.yml        (new)
│           ├── _setup-jvm.yml         (new)
│           ├── _scan.yml              (new)
│           ├── _build-image.yml       (new)
│           ├── _push-image.yml        (new)
│           ├── _cosign-sign.yml       (new)
│           ├── _sbom.yml              (new)
│           └── _vault-creds.yml       (new)
└── impl/
    ├── helixgitpx/
    │   ├── go.work                    (new)
    │   ├── .golangci.yml              (new)
    │   ├── Makefile                   (new)
    │   ├── README.md                  (new)
    │   ├── .github/workflows/ci.yml   (new — calls reusables)
    │   ├── platform/
    │   │   ├── go.mod                 (new)
    │   │   ├── doc.go                 (new)
    │   │   └── {log,telemetry,errors,config,grpc,gin,kafka,pg,redis,temporal,spire,opa,health,testkit}/
    │   │       ├── doc.go             (new, per package)
    │   │       ├── <pkg>.go           (new, per package)
    │   │       └── <pkg>_test.go      (new, per package)
    │   ├── proto/
    │   │   ├── buf.yaml               (new)
    │   │   ├── buf.gen.yaml           (new)
    │   │   ├── buf.work.yaml          (new)
    │   │   ├── buf.lock               (new — generated)
    │   │   ├── README.md              (new)
    │   │   └── helixgitpx/
    │   │       ├── common/v1/{errors,pagination,time}.proto (new)
    │   │       ├── hello/v1/hello.proto (new)
    │   │       └── {auth,repo,sync,conflict,upstream,collab,events,platform}/v1/stubs.proto (new)
    │   ├── gen/go/                    (generated, committed)
    │   ├── api/openapi/               (generated, committed)
    │   ├── tools/scaffold/
    │   │   ├── go.mod                 (new)
    │   │   ├── main.go                (new)
    │   │   ├── template.go            (new)
    │   │   ├── main_test.go           (new)
    │   │   └── templates/             (new — the service template)
    │   └── services/
    │       └── hello/
    │           ├── go.mod             (new)
    │           ├── cmd/hello/main.go  (new)
    │           ├── internal/
    │           │   ├── app/app.go     (new)
    │           │   ├── domain/greeter.go (new)
    │           │   ├── domain/greeter_test.go (new)
    │           │   ├── handler/grpc/server.go (new)
    │           │   ├── handler/http/router.go (new)
    │           │   ├── repo/counter_pg.go (new)
    │           │   ├── repo/cache_redis.go (new)
    │           │   ├── repo/event_kafka.go (new)
    │           │   └── wire.go        (new)
    │           ├── migrations/001_init.sql (new)
    │           ├── config/{config.yaml,config.schema.json} (new)
    │           ├── deploy/
    │           │   ├── Dockerfile     (new)
    │           │   ├── Dockerfile.dev (new)
    │           │   ├── helm/          (new — chart + values + templates)
    │           │   └── skaffold.yaml  (new)
    │           ├── test/integration/hello_e2e_test.go (new)
    │           ├── .golangci.yml      (symlink or include)
    │           ├── Makefile           (new)
    │           ├── README.md          (new)
    │           └── CHANGELOG.md       (new)
    ├── helixgitpx-web/
    │   ├── nx.json                    (new)
    │   ├── package.json               (new)
    │   ├── tsconfig.base.json         (new)
    │   ├── README.md                  (new)
    │   ├── .github/workflows/ci.yml   (new)
    │   ├── apps/web/                  (new — minimal Angular 19 app)
    │   └── libs/proto/                (generated, committed)
    ├── helixgitpx-clients/
    │   ├── settings.gradle.kts        (new)
    │   ├── build.gradle.kts           (new)
    │   ├── gradle.properties          (new)
    │   ├── README.md                  (new)
    │   ├── .github/workflows/ci.yml   (new)
    │   ├── buildSrc/                  (new — convention plugins)
    │   └── shared/                    (new — KMP shared module skeleton)
    ├── helixgitpx-platform/
    │   ├── README.md                  (new)
    │   ├── .github/workflows/ci.yml   (new)
    │   ├── compose/
    │   │   ├── compose.yml            (new)
    │   │   ├── bin/compose            (new — runtime wrapper)
    │   │   ├── observability/
    │   │   │   ├── prometheus.yml     (new)
    │   │   │   └── grafana/provisioning/ (new)
    │   │   └── postgres/init/         (new — bootstrap SQL)
    │   ├── k8s-local/
    │   │   ├── up.sh                  (new)
    │   │   ├── down.sh                (new)
    │   │   ├── load-images.sh         (new)
    │   │   ├── k3d-config.yaml        (new)
    │   │   └── kind-config.yaml       (new)
    │   ├── kyverno/policies/
    │   │   ├── disallow-privileged.yaml (new)
    │   │   ├── require-labels.yaml    (new)
    │   │   ├── require-signed-images.yaml (new)
    │   │   └── enforce-resource-limits.yaml (new)
    │   ├── kyverno/tests/             (new)
    │   ├── checkov/.checkov.yml       (new)
    │   ├── github-actions-runner-controller/
    │   │   ├── values.yaml            (new)
    │   │   ├── runner-scale-set.yaml  (new)
    │   │   └── kata-runtimeclass.yaml (new)
    │   ├── vault/
    │   │   ├── terraform/main.tf      (new)
    │   │   ├── terraform/variables.tf (new)
    │   │   ├── policies/*.hcl         (new)
    │   │   └── oidc-role.json         (new)
    │   ├── argocd/                    (new — app-of-apps skeleton)
    │   ├── helm/                      (new — umbrella chart skeleton)
    │   ├── kustomize/                 (new)
    │   ├── terraform/                 (new — root module skeleton)
    │   └── Tiltfile                   (new)
    └── helixgitpx-docs/
        ├── package.json               (new)
        ├── tsconfig.json              (new)
        ├── docusaurus.config.ts       (new)
        ├── sidebars.ts                (new)
        ├── README.md                  (new)
        ├── .github/workflows/ci.yml   (new)
        ├── sync-docs.mjs              (new)
        ├── src/                       (new)
        ├── static/                    (new)
        └── docs/                      (sync'd, gitignored)

docs/specifications/main/main_implementation_material/HelixGitpx/
├── 12-operations/runbooks/_TEMPLATE.md   (new — added in Phase B)
├── 15-reference/adr/
│   ├── 0001-github-actions-workflow-dispatch-only.md (new)
│   ├── 0002-portable-container-runtime.md (new)
│   ├── 0003-single-git-history-impl-subdirs.md (new)
│   ├── 0004-mise-toolchain.md             (new)
│   └── 0005-spine-first-sequencing.md     (new)
```

---

## Phase A — Spine

### Task 1: Root scaffolding — `mise.toml`, `.gitattributes`, root `README.md`, `CHANGELOG.md`, `SOLO-NOTES.md`

**Files:**
- Create: `mise.toml`
- Create: `.gitattributes`
- Create: `README.md` (replaces existing 1-line)
- Create: `CHANGELOG.md`
- Create: `SOLO-NOTES.md`

- [ ] **Step 1: Write `mise.toml` with pinned toolchain**

Write `mise.toml`:

```toml
[tools]
go = "1.23.4"
node = "20.18.1"
java = "temurin-21.0.5+11"
gradle = "8.10.2"
protoc = "28.3"
"go:github.com/bufbuild/buf/cmd/buf" = "1.47.2"
"go:github.com/sqlc-dev/sqlc/cmd/sqlc" = "1.27.0"
kubectl = "1.31.3"
helm = "3.16.3"
kind = "0.25.0"
k3d = "5.7.5"
skaffold = "2.13.2"
tilt = "0.33.21"
cosign = "2.4.1"
syft = "1.18.0"
grype = "0.85.0"
kyverno-cli = "1.13.2"
"pipx:checkov" = "3.2.334"
"go:github.com/go-delve/delve/cmd/dlv" = "1.23.1"
"go:github.com/pressly/goose/v3/cmd/goose" = "3.22.1"
"go:github.com/go-task/task/v3/cmd/task" = "3.40.0"
"go:mvdan.cc/gofumpt" = "0.7.0"
"go:github.com/golangci/golangci-lint/cmd/golangci-lint" = "1.62.2"
"go:github.com/go-gremlins/gremlins/cmd/gremlins" = "0.5.0"

[env]
GOFLAGS = "-mod=mod"
KUBECONFIG = "{{config_root}}/.kube/config"
```

- [ ] **Step 2: Write `.gitattributes`**

Write `.gitattributes`:

```
# Codegen output — mark as generated to suppress diffs in UIs that honor linguist
impl/helixgitpx/gen/** linguist-generated=true
impl/helixgitpx/api/openapi/** linguist-generated=true
impl/helixgitpx-web/libs/proto/** linguist-generated=true
impl/helixgitpx-clients/shared/src/commonMain/kotlin/gen/** linguist-generated=true
impl/helixgitpx-clients/iosApp/Gen/** linguist-generated=true

# Normalize line endings
* text=auto eol=lf
*.sh text eol=lf
*.bat text eol=crlf
*.png binary
*.jpg binary
*.ico binary
*.pdf binary
```

- [ ] **Step 3: Replace root `README.md`**

Write `README.md`:

```markdown
# HelixGitpx

**Helix Git Proxy eXtended** — a federated, privacy-preserving Git proxy that mirrors a single source of truth across multiple upstream Git hosts (GitHub, GitLab, GitFlic, GitVerse, …) and resolves the inevitable conflicts.

## Where to go next

| If you want to… | Read |
|---|---|
| Understand what HelixGitpx is | [`docs/specifications/main/main_implementation_material/HelixGitpx/00-core/01-vision-scope-constraints.md`](docs/specifications/main/main_implementation_material/HelixGitpx/00-core/01-vision-scope-constraints.md) |
| Browse the full spec suite | [`docs/specifications/main/main_implementation_material/HelixGitpx/README.md`](docs/specifications/main/main_implementation_material/HelixGitpx/README.md) |
| Read it as a website | `make docs && open http://localhost:3001` |
| Start hacking | [`impl/helixgitpx/README.md`](impl/helixgitpx/README.md) |
| Understand the roadmap | [`docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md`](docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md) |
| See current milestone plans | [`docs/superpowers/plans/`](docs/superpowers/plans/) |

## Quick start

```sh
mise install                    # pin toolchain per mise.toml
make bootstrap                  # fetch deps for every sub-project
make dev                        # bring up compose stack + hello service
curl "http://localhost:8001/v1/hello?name=world"
```

## Repository layout

- `docs/` — authoritative specifications (immutable input; edit via PR with spec review)
- `impl/` — implementation sub-projects (Go monorepo, Angular web, KMP clients, GitOps, Docusaurus)
- `Upstreams/` — scripts configuring this repo's own federation across GitHub/GitLab/GitFlic/GitVerse
- `.github/workflows/` — CI pipelines (all manually-triggered; see ADR-0001)

## License

Apache-2.0 (code) / CC-BY-SA-4.0 (documentation). See `LICENSE`.
```

- [ ] **Step 4: Write `CHANGELOG.md`**

Write `CHANGELOG.md`:

```markdown
# Changelog

All notable changes to this repository are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning: per-artifact semver. This file tracks repo-level milestones.

## [Unreleased]

### Added
- Authoritative specification suite v1.0.0 (superseded Git_Proxy_Master_Specification v4.0.0).
- Milestone M1 Foundation design spec and implementation plan.
- Upstream mirror scripts for GitHub, GitLab, GitFlic, GitVerse.

## [0.0.0] — 2026-04-19

### Added
- Initial repository seeded with Git_Proxy_Master_Specification v4.0.0 (now superseded).
```

- [ ] **Step 5: Write `SOLO-NOTES.md`**

Write `SOLO-NOTES.md`:

```markdown
# Solo-Operation Notes

This file documents deviations from `CONTRIBUTING.md` that apply while HelixGitpx is maintained by a single engineer. Each deviation is temporary; re-enable the rule once the constraint no longer applies.

## Active deviations

| Rule | Deviation | Re-enable when |
|---|---|---|
| 2 approvers required for PRs to `main` | 1 self-review | team size ≥ 2 |
| CODEOWNERS-enforced reviews | configured but not enforced in GitHub branch protection | team size ≥ 2 |
| DCO `Signed-off-by` | enforced | always |
| Conventional Commits | enforced | always |
| Signed commits (GPG/SSH) | enforced | always |

## How this file is used

- New engineers read this first after `README.md`.
- Every deviation lists the exact condition under which it is lifted.
- When a deviation is lifted, this file is updated in the same PR that enables the enforcement.
```

- [ ] **Step 6: Commit**

```sh
git add mise.toml .gitattributes README.md CHANGELOG.md SOLO-NOTES.md
git commit -s -m "chore(m1): seed root scaffolding (mise, gitattributes, README, CHANGELOG, SOLO-NOTES)"
```

---

### Task 2: Root `Makefile` orchestrator

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Write the Makefile**

Write `Makefile`:

```makefile
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
```

- [ ] **Step 2: Verify Makefile parses**

Run: `make -n help`
Expected: no errors; prints the help output.

- [ ] **Step 3: Commit**

```sh
git add Makefile
git commit -s -m "chore(m1): add root Makefile orchestrator"
```

---

### Task 3: Compose runtime wrapper (portable docker/podman)

**Files:**
- Create: `impl/helixgitpx-platform/compose/bin/compose`
- Create: `impl/helixgitpx-platform/README.md`

- [ ] **Step 1: Create directory and write the wrapper**

```sh
mkdir -p impl/helixgitpx-platform/compose/bin
```

Write `impl/helixgitpx-platform/compose/bin/compose`:

```sh
#!/usr/bin/env sh
# Compose runtime wrapper — auto-detects docker or podman.
# Every Makefile/script must call this wrapper, never docker/podman directly.
# See ADR-0002.
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
COMPOSE_FILE="${COMPOSE_FILE:-${SCRIPT_DIR}/../compose.yml}"

if [ ! -f "${COMPOSE_FILE}" ]; then
    printf 'compose: cannot find compose file at %s\n' "${COMPOSE_FILE}" >&2
    exit 2
fi

# Resolution order: docker compose → podman compose → podman-compose.
# docker first because it is the ecosystem default; podman is first-class via compose v2.
if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    exec docker compose -f "${COMPOSE_FILE}" "$@"
elif command -v podman >/dev/null 2>&1 && podman compose version >/dev/null 2>&1; then
    # podman 4.7+ supports `podman compose` via compose plugin.
    exec podman compose -f "${COMPOSE_FILE}" "$@"
elif command -v podman-compose >/dev/null 2>&1; then
    exec podman-compose -f "${COMPOSE_FILE}" "$@"
else
    cat >&2 <<'EOF'
compose: no container runtime found. Install one of:
  - Docker Desktop (macOS/Windows) or docker-ce + compose v2 plugin (Linux)
  - Podman 4.7+ (with the compose plugin: `podman compose --help` must succeed)
  - podman-compose (fallback for older podman)
EOF
    exit 3
fi
```

- [ ] **Step 2: Make executable and shellcheck it**

```sh
chmod +x impl/helixgitpx-platform/compose/bin/compose
shellcheck impl/helixgitpx-platform/compose/bin/compose
```

Expected: no shellcheck output (clean).

- [ ] **Step 3: Write the platform README**

Write `impl/helixgitpx-platform/README.md`:

```markdown
# helixgitpx-platform

GitOps and infrastructure artifacts for HelixGitpx.

## Contents

| Path | Purpose |
|---|---|
| `compose/` | Local development stack (Postgres, Kafka, Redis, observability, hello) — see [ADR-0002] for the runtime wrapper |
| `k8s-local/` | `kind`/`k3d` cluster scripts + Tiltfile (full-platform rehearsal path) |
| `kyverno/policies/` | Cluster admission policies (enforced in M2+) |
| `checkov/` | IaC static analysis config |
| `github-actions-runner-controller/` | Self-hosted Kata runner configuration (M2 activation) |
| `vault/` | Vault policies + Terraform module + OIDC role (M2 activation) |
| `argocd/` | Argo CD Application definitions (M2 activation) |
| `helm/` | Umbrella Helm chart (M2+) |
| `kustomize/` | Kustomize overlays per environment |
| `terraform/` | Root Terraform modules |
| `Tiltfile` | Alternative local dev path on Kubernetes (vs. compose) |

## Local stack

```sh
make dev          # bring up everything
make dev-down     # tear down (removes volumes)
```

The compose file lives at `compose/compose.yml`. The `compose/bin/compose` wrapper auto-detects `docker`, `podman`, or `podman-compose` — never invoke them directly.

[ADR-0002]: ../../docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0002-portable-container-runtime.md
```

- [ ] **Step 4: Verify wrapper errors cleanly when no file present**

```sh
COMPOSE_FILE=/tmp/does-not-exist.yml impl/helixgitpx-platform/compose/bin/compose ps; echo "exit=$?"
```

Expected: "compose: cannot find compose file at /tmp/does-not-exist.yml", exit 2.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx-platform/compose/bin/compose impl/helixgitpx-platform/README.md
git commit -s -m "feat(platform): add portable compose wrapper (docker/podman)"
```

---

### Task 4: Minimal `compose.yml` — Postgres + Kafka + Redis

**Files:**
- Create: `impl/helixgitpx-platform/compose/compose.yml`
- Create: `impl/helixgitpx-platform/compose/postgres/init/00-schemas.sql`

- [ ] **Step 1: Write Postgres init script**

```sh
mkdir -p impl/helixgitpx-platform/compose/postgres/init
```

Write `impl/helixgitpx-platform/compose/postgres/init/00-schemas.sql`:

```sql
-- Per-service schemas for the hello service (M1 spine).
-- Later milestones add more services; each gets its own schema.
CREATE SCHEMA IF NOT EXISTS hello;
CREATE USER hello_svc WITH PASSWORD 'hello_svc';
GRANT USAGE, CREATE ON SCHEMA hello TO hello_svc;
GRANT ALL ON ALL TABLES IN SCHEMA hello TO hello_svc;
ALTER DEFAULT PRIVILEGES IN SCHEMA hello GRANT ALL ON TABLES TO hello_svc;
ALTER DEFAULT PRIVILEGES IN SCHEMA hello GRANT ALL ON SEQUENCES TO hello_svc;
```

- [ ] **Step 2: Write the compose file (core + hello; observability added in a later task)**

Write `impl/helixgitpx-platform/compose/compose.yml`:

```yaml
# Local development stack for HelixGitpx M1 Foundation.
# Invoked via impl/helixgitpx-platform/compose/bin/compose (see ADR-0002).
name: helixgitpx

x-common-env: &common-env
  TZ: UTC

networks:
  helix:
    driver: bridge

volumes:
  pg-data:
  kafka-data:
  redis-data:
  prom-data:
  grafana-data:

services:
  postgres:
    image: postgres:16-alpine
    profiles: ["core", "all"]
    environment:
      <<: *common-env
      POSTGRES_USER: helix
      POSTGRES_PASSWORD: helix
      POSTGRES_DB: helixgitpx
    ports: ["5432:5432"]
    volumes:
      - pg-data:/var/lib/postgresql/data
      - ./postgres/init:/docker-entrypoint-initdb.d:ro
    networks: [helix]
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U helix -d helixgitpx"]
      interval: 5s
      timeout: 5s
      retries: 10

  kafka:
    image: apache/kafka:3.8.1
    profiles: ["core", "all"]
    environment:
      <<: *common-env
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: "broker,controller"
      KAFKA_LISTENERS: "PLAINTEXT://:9092,CONTROLLER://:9093"
      KAFKA_ADVERTISED_LISTENERS: "PLAINTEXT://kafka:9092"
      KAFKA_CONTROLLER_LISTENER_NAMES: "CONTROLLER"
      KAFKA_CONTROLLER_QUORUM_VOTERS: "1@kafka:9093"
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
      CLUSTER_ID: "K6VU7KM6S2KjVSdQ7Y1-DA"
    ports: ["9092:9092"]
    volumes:
      - kafka-data:/var/lib/kafka/data
    networks: [helix]
    healthcheck:
      test: ["CMD-SHELL", "/opt/kafka/bin/kafka-cluster.sh cluster-id --bootstrap-server localhost:9092 || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 12

  redis:
    image: redis:7-alpine
    profiles: ["core", "all"]
    command: ["redis-server", "--appendonly", "yes"]
    ports: ["6379:6379"]
    volumes:
      - redis-data:/data
    networks: [helix]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 10

  karapace:
    image: ghcr.io/aiven-open/karapace:latest
    profiles: ["registry", "all"]
    environment:
      <<: *common-env
      KARAPACE_ADVERTISED_HOSTNAME: karapace
      KARAPACE_BOOTSTRAP_URI: kafka:9092
      KARAPACE_PORT: 8081
      KARAPACE_HOST: 0.0.0.0
      KARAPACE_REGISTRY_HOST: 0.0.0.0
      KARAPACE_LOG_LEVEL: WARNING
    command: ["karapace", "/opt/karapace/karapace.config.json"]
    ports: ["8081:8081"]
    networks: [helix]
    depends_on:
      kafka:
        condition: service_healthy

  jaeger:
    image: jaegertracing/all-in-one:1.62
    profiles: ["observability", "all"]
    environment:
      <<: *common-env
      COLLECTOR_OTLP_ENABLED: "true"
    ports:
      - "16686:16686"     # UI
      - "4317:4317"       # OTLP gRPC
      - "4318:4318"       # OTLP HTTP
    networks: [helix]

  prometheus:
    image: prom/prometheus:v3.0.1
    profiles: ["observability", "all"]
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--storage.tsdb.path=/prometheus"
    ports: ["9090:9090"]
    volumes:
      - ./observability/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prom-data:/prometheus
    networks: [helix]

  grafana:
    image: grafana/grafana:11.4.0
    profiles: ["observability", "all"]
    environment:
      <<: *common-env
      GF_SECURITY_ADMIN_USER: admin
      GF_SECURITY_ADMIN_PASSWORD: admin
      GF_USERS_ALLOW_SIGN_UP: "false"
    ports: ["3000:3000"]
    volumes:
      - grafana-data:/var/lib/grafana
      - ./observability/grafana/provisioning:/etc/grafana/provisioning:ro
    networks: [helix]
    depends_on:
      - prometheus

  bufstream:
    image: bufbuild/bufstream:latest
    profiles: ["registry"]
    environment:
      <<: *common-env
    ports: ["8082:8082"]
    networks: [helix]

  hello:
    build:
      context: ../../helixgitpx
      dockerfile: services/hello/deploy/Dockerfile
    profiles: ["all"]
    environment:
      <<: *common-env
      HELLO_HTTP_ADDR: ":8001"
      HELLO_GRPC_ADDR: ":9001"
      HELLO_HEALTH_ADDR: ":8081"
      HELLO_POSTGRES_DSN: "postgres://hello_svc:hello_svc@postgres:5432/helixgitpx?search_path=hello&sslmode=disable"
      HELLO_REDIS_ADDR: "redis:6379"
      HELLO_KAFKA_BROKERS: "kafka:9092"
      HELLO_KAFKA_TOPIC: "hello.said"
      OTEL_EXPORTER_OTLP_ENDPOINT: "http://jaeger:4317"
      OTEL_SERVICE_NAME: "hello"
    ports:
      - "8001:8001"
      - "9001:9001"
      - "8081:8081"
    networks: [helix]
    depends_on:
      postgres:
        condition: service_healthy
      kafka:
        condition: service_healthy
      redis:
        condition: service_healthy
```

- [ ] **Step 3: Write Prometheus config stub**

```sh
mkdir -p impl/helixgitpx-platform/compose/observability/grafana/provisioning
```

Write `impl/helixgitpx-platform/compose/observability/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ["localhost:9090"]

  - job_name: hello
    static_configs:
      - targets: ["hello:8081"]
    metrics_path: /metrics
```

- [ ] **Step 4: Validate compose file**

```sh
impl/helixgitpx-platform/compose/bin/compose --profile core config --quiet
```

Expected: exit 0, no output.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx-platform/compose
git commit -s -m "feat(platform): add compose stack (pg, kafka, redis, observability, registry)"
```

---

### Task 5: Go workspace + `platform` module skeleton + `errors` package (TDD)

**Files:**
- Create: `impl/helixgitpx/go.work`
- Create: `impl/helixgitpx/.golangci.yml`
- Create: `impl/helixgitpx/Makefile`
- Create: `impl/helixgitpx/README.md`
- Create: `impl/helixgitpx/platform/go.mod`
- Create: `impl/helixgitpx/platform/doc.go`
- Create: `impl/helixgitpx/platform/errors/doc.go`
- Create: `impl/helixgitpx/platform/errors/errors.go`
- Create: `impl/helixgitpx/platform/errors/errors_test.go`
- Create: `impl/helixgitpx/platform/errors/problem.go`

- [ ] **Step 1: Initialize Go workspace**

```sh
mkdir -p impl/helixgitpx/platform/errors
cd impl/helixgitpx
go work init
cd platform && go mod init github.com/helixgitpx/platform
cd ../.. ; cd impl/helixgitpx && go work use ./platform
```

- [ ] **Step 2: Write root Go linter config**

Write `impl/helixgitpx/.golangci.yml`:

```yaml
run:
  timeout: 5m
  go: "1.23"

linters:
  disable-all: true
  enable:
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - revive
    - gosec
    - gofumpt
    - goimports
    - misspell
    - unparam
    - prealloc
    - gocritic
    - bodyclose
    - contextcheck
    - errorlint
    - nilerr
    - noctx
    - thelper
    - tparallel
    - unconvert

linters-settings:
  gofumpt:
    extra-rules: true
  revive:
    rules:
      - name: exported
        arguments: [checkPrivateReceivers, sayRepetitiveInsteadOfStutters]
  errcheck:
    check-type-assertions: true

issues:
  exclude-rules:
    - path: _test\.go
      linters: [gosec, errcheck, unparam]
```

- [ ] **Step 3: Write the helixgitpx Go monorepo Makefile**

Write `impl/helixgitpx/Makefile`:

```makefile
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
```

- [ ] **Step 4: Write `impl/helixgitpx/README.md`**

Write `impl/helixgitpx/README.md`:

```markdown
# helixgitpx (Go monorepo)

Contains the Go implementation of HelixGitpx: shared `platform` libraries, service binaries, the `scaffold` tool, and the proto root.

## Layout

```
├── go.work                Go 1.23 workspace
├── platform/              shared libraries (14 packages)
├── services/              service binaries (hello in M1; more in M3+)
├── tools/scaffold/        service-template renderer (Go binary, no Python)
├── proto/                 protobuf sources (buf module buf.build/helixgitpx/core)
├── gen/go/                generated Go code (committed; see .gitattributes)
└── api/openapi/           generated OpenAPI (committed)
```

## Daily commands

```sh
make lint     # golangci-lint + buf lint
make test     # go test -race
make gen      # regenerate protobuf bindings
make fmt      # gofumpt + goimports
```

See the shared lib docs under `platform/*/README.md`.
```

- [ ] **Step 5: Write `platform/doc.go` and `platform/go.mod`**

Write `impl/helixgitpx/platform/doc.go`:

```go
// Package platform is the shared-library umbrella for HelixGitpx services.
//
// Sub-packages provide logging, telemetry, typed errors, configuration,
// gRPC/HTTP servers, Kafka/Postgres/Redis clients, Temporal wiring, SPIFFE
// integration, OPA evaluation, health endpoints, and a test toolkit.
//
// See the per-package doc.go for details. All constructors accept a context
// and return (client, error); callers own lifecycle (Close).
package platform
```

`go.mod` is already initialized by Step 1. Append dependencies when they first appear in tests.

- [ ] **Step 6: TDD — write the failing test for `errors.New`**

Write `impl/helixgitpx/platform/errors/errors_test.go`:

```go
package errors_test

import (
	stderrors "errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/helixgitpx/platform/errors"
)

func TestNew_RoundTripFields(t *testing.T) {
	cause := stderrors.New("underlying")
	e := errors.New(codes.NotFound, "repo", "branch %q missing", "main").
		Wrap(cause).
		With("ref", "refs/heads/main")

	if e.Code != codes.NotFound {
		t.Errorf("Code = %v, want NotFound", e.Code)
	}
	if e.Domain != "repo" {
		t.Errorf("Domain = %q, want repo", e.Domain)
	}
	if e.Message != `branch "main" missing` {
		t.Errorf("Message = %q, want quoted branch", e.Message)
	}
	if !stderrors.Is(e, cause) {
		t.Errorf("Is(cause) = false, want true")
	}
	if e.Details["ref"] != "refs/heads/main" {
		t.Errorf("Details[ref] = %v, want refs/heads/main", e.Details["ref"])
	}
}

func TestError_HTTPStatus(t *testing.T) {
	cases := []struct {
		code codes.Code
		want int
	}{
		{codes.OK, http.StatusOK},
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.NotFound, http.StatusNotFound},
		{codes.PermissionDenied, http.StatusForbidden},
		{codes.Unauthenticated, http.StatusUnauthorized},
		{codes.ResourceExhausted, http.StatusTooManyRequests},
		{codes.FailedPrecondition, http.StatusPreconditionFailed},
		{codes.Aborted, http.StatusConflict},
		{codes.Unavailable, http.StatusServiceUnavailable},
		{codes.DeadlineExceeded, http.StatusGatewayTimeout},
		{codes.Unimplemented, http.StatusNotImplemented},
		{codes.Internal, http.StatusInternalServerError},
	}
	for _, c := range cases {
		e := errors.New(c.code, "x", "msg")
		if got := e.HTTPStatus(); got != c.want {
			t.Errorf("code %v: HTTPStatus() = %d, want %d", c.code, got, c.want)
		}
	}
}

func TestError_IsByCode(t *testing.T) {
	e1 := errors.New(codes.NotFound, "repo", "x")
	e2 := errors.New(codes.NotFound, "repo", "y")
	if !stderrors.Is(e1, e2) {
		t.Errorf("Is(same code+domain) = false, want true")
	}
	e3 := errors.New(codes.InvalidArgument, "repo", "z")
	if stderrors.Is(e1, e3) {
		t.Errorf("Is(diff code) = true, want false")
	}
}
```

- [ ] **Step 7: Run the test, expect failure**

```sh
cd impl/helixgitpx/platform && go test ./errors/...
```

Expected: compile error — `errors.New`, `*Error.Wrap`, `*Error.With`, `*Error.HTTPStatus` are undefined.

- [ ] **Step 8: Implement `errors.go`**

Write `impl/helixgitpx/platform/errors/errors.go`:

```go
// Package errors defines the canonical HelixGitpx error type.
//
// An Error carries a gRPC status code, a domain tag, a human-readable message,
// an optional cause, and a map of structured details. It implements the Go
// error interface, the standard errors.Is/As contract, and supplies an HTTP
// status mapping per RFC 7807.
package errors

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// Error is HelixGitpx's typed error.
type Error struct {
	Code    codes.Code
	Domain  string
	Message string
	Cause   error
	Details map[string]any
}

// New constructs an Error. The message is formatted with fmt.Sprintf if args are supplied.
func New(code codes.Code, domain, format string, args ...any) *Error {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return &Error{Code: code, Domain: domain, Message: msg}
}

// Wrap attaches a cause and returns the receiver.
func (e *Error) Wrap(cause error) *Error {
	e.Cause = cause
	return e
}

// With adds a structured detail and returns the receiver.
func (e *Error) With(key string, value any) *Error {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s/%s] %s: %v", e.Domain, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s/%s] %s", e.Domain, e.Code, e.Message)
}

// Unwrap supports errors.Is/As.
func (e *Error) Unwrap() error { return e.Cause }

// Is returns true when target is an *Error with the same Code and Domain, or
// when the cause matches target.
func (e *Error) Is(target error) bool {
	var other *Error
	if errors.As(target, &other) {
		return e.Code == other.Code && e.Domain == other.Domain
	}
	return false
}

// HTTPStatus maps the gRPC code to the closest HTTP status per RFC 7807.
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.InvalidArgument, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.DataLoss, codes.Internal, codes.Unknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
```

Write `impl/helixgitpx/platform/errors/problem.go`:

```go
package errors

// Problem is the RFC 7807 problem-details JSON representation.
type Problem struct {
	Type     string         `json:"type"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	Domain   string         `json:"domain,omitempty"`
	Code     string         `json:"code,omitempty"`
	Errors   map[string]any `json:"errors,omitempty"`
}

// ToProblem renders e as an RFC 7807 problem document.
func (e *Error) ToProblem(instance string) Problem {
	return Problem{
		Type:     "https://helixgitpx.dev/errors/" + e.Code.String(),
		Title:    e.Code.String(),
		Status:   e.HTTPStatus(),
		Detail:   e.Message,
		Instance: instance,
		Domain:   e.Domain,
		Code:     e.Code.String(),
		Errors:   e.Details,
	}
}
```

Write `impl/helixgitpx/platform/errors/doc.go`:

```go
// Package errors provides the canonical HelixGitpx error type.
//
// Usage:
//
//	err := errors.New(codes.NotFound, "repo", "ref %q missing", name).
//		Wrap(cause).
//		With("ref", ref)
//	return err
//
// HTTP handlers map errors to RFC 7807 problem documents via err.ToProblem.
package errors
```

- [ ] **Step 9: Fetch deps and run the test**

```sh
cd impl/helixgitpx/platform
go mod tidy
go test ./errors/...
```

Expected: `PASS` for all three tests.

- [ ] **Step 10: Commit**

```sh
git add impl/helixgitpx/go.work impl/helixgitpx/.golangci.yml impl/helixgitpx/Makefile \
        impl/helixgitpx/README.md impl/helixgitpx/platform
git commit -s -m "feat(platform/errors): typed error with grpc code + RFC 7807 mapping"
```

---

### Task 6: `platform/log` package (zap/slog wrapper, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/log/doc.go`
- Create: `impl/helixgitpx/platform/log/log.go`
- Create: `impl/helixgitpx/platform/log/log_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/log/log_test.go`:

```go
package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/helixgitpx/platform/log"
)

func TestNew_EmitsJSON(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "info", Output: &buf, Service: "hello", Version: "test"})
	lg.Info("hello", "name", "world")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v\n%s", err, buf.String())
	}
	if got["msg"] != "hello" {
		t.Errorf("msg = %v, want hello", got["msg"])
	}
	if got["name"] != "world" {
		t.Errorf("name = %v, want world", got["name"])
	}
	if got["service"] != "hello" {
		t.Errorf("service = %v, want hello", got["service"])
	}
	if got["version"] != "test" {
		t.Errorf("version = %v, want test", got["version"])
	}
}

func TestFromContext_ReturnsChildLogger(t *testing.T) {
	var buf bytes.Buffer
	root := log.New(log.Options{Level: "info", Output: &buf, Service: "s"})
	ctx := log.WithContext(context.Background(), root.With("request_id", "abc"))

	log.FromContext(ctx).Info("tick")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if got["request_id"] != "abc" {
		t.Errorf("request_id = %v, want abc", got["request_id"])
	}
}
```

- [ ] **Step 2: Run — expect fail**

```sh
cd impl/helixgitpx/platform && go test ./log/...
```

Expected: compile errors.

- [ ] **Step 3: Implement `log.go`**

Write `impl/helixgitpx/platform/log/log.go`:

```go
// Package log provides a structured, JSON logger built on log/slog.
//
// Options controls level, output, and service/version tags. FromContext/WithContext
// pass a child logger through a context.Context for request-scoped fields.
package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
)

// Options configures New.
type Options struct {
	Level   string    // "debug"/"info"/"warn"/"error"
	Output  io.Writer // defaults to os.Stdout
	Service string    // added as "service" field on every record
	Version string    // added as "version" field on every record
}

// Logger wraps *slog.Logger with With/FromContext/WithContext helpers.
type Logger struct {
	sl *slog.Logger
}

type ctxKey struct{}

var global atomic.Pointer[Logger]

// New constructs a Logger.
func New(opts Options) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	handler := slog.NewJSONHandler(opts.Output, &slog.HandlerOptions{
		Level: parseLevel(opts.Level),
	})
	sl := slog.New(handler)
	if opts.Service != "" {
		sl = sl.With("service", opts.Service)
	}
	if opts.Version != "" {
		sl = sl.With("version", opts.Version)
	}
	lg := &Logger{sl: sl}
	global.Store(lg)
	return lg
}

// Default returns the most recently constructed Logger, or a noop logger if none.
func Default() *Logger {
	if lg := global.Load(); lg != nil {
		return lg
	}
	return &Logger{sl: slog.New(slog.NewJSONHandler(io.Discard, nil))}
}

// Info logs at info level with key/value pairs.
func (l *Logger) Info(msg string, kv ...any)  { l.sl.Info(msg, kv...) }
func (l *Logger) Warn(msg string, kv ...any)  { l.sl.Warn(msg, kv...) }
func (l *Logger) Error(msg string, kv ...any) { l.sl.Error(msg, kv...) }
func (l *Logger) Debug(msg string, kv ...any) { l.sl.Debug(msg, kv...) }

// With returns a child logger with the supplied fields attached.
func (l *Logger) With(kv ...any) *Logger { return &Logger{sl: l.sl.With(kv...)} }

// WithContext stores lg on ctx. Callers retrieve it with FromContext.
func WithContext(ctx context.Context, lg *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, lg)
}

// FromContext returns the logger attached via WithContext, or Default().
func FromContext(ctx context.Context) *Logger {
	if lg, ok := ctx.Value(ctxKey{}).(*Logger); ok {
		return lg
	}
	return Default()
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
```

Write `impl/helixgitpx/platform/log/doc.go`:

```go
// Package log provides the shared structured logger. See the package type Logger
// and the WithContext/FromContext helpers for request-scoped child loggers.
package log
```

- [ ] **Step 4: Run — expect pass**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./log/...
```

Expected: PASS.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx/platform/log
git commit -s -m "feat(platform/log): slog-based structured logger with context helpers"
```

---

### Task 7: `platform/config` package (Viper-ish env→struct loader, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/config/doc.go`
- Create: `impl/helixgitpx/platform/config/config.go`
- Create: `impl/helixgitpx/platform/config/config_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/config/config_test.go`:

```go
package config_test

import (
	"testing"
	"time"

	"github.com/helixgitpx/platform/config"
)

type helloConfig struct {
	HTTPAddr     string        `env:"HTTP_ADDR" default:":8001"`
	GRPCAddr     string        `env:"GRPC_ADDR" default:":9001"`
	Timeout      time.Duration `env:"TIMEOUT" default:"30s"`
	KafkaBrokers []string      `env:"KAFKA_BROKERS" default:"localhost:9092" split:","`
	Enabled      bool          `env:"ENABLED" default:"true"`
	MaxConns     int           `env:"MAX_CONNS" default:"10"`
}

func TestLoad_UsesDefaults(t *testing.T) {
	var c helloConfig
	if err := config.Load(&c, config.Options{Prefix: "HELLO"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.HTTPAddr != ":8001" {
		t.Errorf("HTTPAddr = %q, want :8001", c.HTTPAddr)
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", c.Timeout)
	}
	if len(c.KafkaBrokers) != 1 || c.KafkaBrokers[0] != "localhost:9092" {
		t.Errorf("KafkaBrokers = %v", c.KafkaBrokers)
	}
	if !c.Enabled {
		t.Errorf("Enabled = false, want true")
	}
	if c.MaxConns != 10 {
		t.Errorf("MaxConns = %d, want 10", c.MaxConns)
	}
}

func TestLoad_EnvOverridesDefault(t *testing.T) {
	t.Setenv("HELLO_HTTP_ADDR", ":9999")
	t.Setenv("HELLO_KAFKA_BROKERS", "a:1,b:2")
	t.Setenv("HELLO_ENABLED", "false")
	var c helloConfig
	if err := config.Load(&c, config.Options{Prefix: "HELLO"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.HTTPAddr != ":9999" {
		t.Errorf("HTTPAddr = %q", c.HTTPAddr)
	}
	if len(c.KafkaBrokers) != 2 || c.KafkaBrokers[1] != "b:2" {
		t.Errorf("KafkaBrokers = %v", c.KafkaBrokers)
	}
	if c.Enabled {
		t.Errorf("Enabled = true, want false")
	}
}

type required struct {
	DSN string `env:"DSN" required:"true"`
}

func TestLoad_RequiredFieldMissing(t *testing.T) {
	var c required
	err := config.Load(&c, config.Options{Prefix: "X"})
	if err == nil {
		t.Fatalf("expected error for missing required DSN")
	}
}
```

- [ ] **Step 2: Run — expect fail**

```sh
cd impl/helixgitpx/platform && go test ./config/...
```

- [ ] **Step 3: Implement `config.go`**

Write `impl/helixgitpx/platform/config/config.go`:

```go
// Package config loads typed configuration structs from environment variables
// with struct-tag-driven defaults, required markers, and list splitting.
//
// Struct tags supported:
//
//	env:"NAME"        environment variable name (required if field not anonymous)
//	default:"value"   default when env var unset
//	required:"true"   return error when unset and no default
//	split:","         for []string, split on the given separator
//
// Precedence: env var > default. File-based config and CLI flags may be layered
// by callers; this package intentionally keeps the minimal surface.
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Options control Load.
type Options struct {
	// Prefix is prepended to every env tag as "<Prefix>_<tag>" when set.
	Prefix string
}

// Load populates *dst from the environment.
func Load(dst any, opts Options) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config: Load requires a pointer to a struct, got %T", dst)
	}
	return loadStruct(v.Elem(), opts.Prefix)
}

func loadStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		fv := v.Field(i)

		if sf.Type.Kind() == reflect.Struct && sf.Anonymous {
			if err := loadStruct(fv, prefix); err != nil {
				return err
			}
			continue
		}

		envName := sf.Tag.Get("env")
		if envName == "" {
			continue
		}
		if prefix != "" {
			envName = prefix + "_" + envName
		}

		raw, present := os.LookupEnv(envName)
		if !present {
			raw = sf.Tag.Get("default")
		}
		if raw == "" && sf.Tag.Get("required") == "true" {
			return fmt.Errorf("config: required env %s is unset", envName)
		}
		if raw == "" {
			continue
		}
		if err := assign(fv, raw, sf.Tag.Get("split")); err != nil {
			return fmt.Errorf("config: %s: %w", envName, err)
		}
	}
	return nil
}

func assign(fv reflect.Value, raw, split string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fv.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(raw)
			if err != nil {
				return err
			}
			fv.SetInt(int64(d))
			return nil
		}
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	case reflect.Slice:
		if split == "" {
			split = ","
		}
		parts := strings.Split(raw, split)
		out := reflect.MakeSlice(fv.Type(), len(parts), len(parts))
		for i, p := range parts {
			if err := assign(out.Index(i), strings.TrimSpace(p), ""); err != nil {
				return err
			}
		}
		fv.Set(out)
	default:
		return fmt.Errorf("unsupported kind %s", fv.Kind())
	}
	return nil
}
```

Write `impl/helixgitpx/platform/config/doc.go`:

```go
// Package config is the HelixGitpx env-to-struct loader. See the Load function.
package config
```

- [ ] **Step 4: Run — expect pass**

```sh
cd impl/helixgitpx/platform && go test ./config/...
```

Expected: PASS.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx/platform/config
git commit -s -m "feat(platform/config): env→struct loader with defaults/required/split"
```

---

### Task 8: `platform/health` package (probes + handler, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/health/doc.go`
- Create: `impl/helixgitpx/platform/health/health.go`
- Create: `impl/helixgitpx/platform/health/health_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/health/health_test.go`:

```go
package health_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/helixgitpx/platform/health"
)

func TestHandler_LiveAlwaysOK(t *testing.T) {
	h := health.New()
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.Live(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Live Code = %d, want 200", w.Code)
	}
}

func TestHandler_ReadyReflectsProbes(t *testing.T) {
	h := health.New()
	h.Register("db", func(context.Context) error { return nil })
	h.Register("cache", func(context.Context) error { return errors.New("down") })

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	h.Ready(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Ready Code = %d, want 503 (one probe down)", w.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("body not JSON: %v", err)
	}
	checks, _ := body["checks"].(map[string]any)
	if checks["db"] != "ok" {
		t.Errorf("db = %v, want ok", checks["db"])
	}
	if checks["cache"] == "ok" {
		t.Errorf("cache should be failing")
	}
}
```

- [ ] **Step 2: Run — expect fail**

- [ ] **Step 3: Implement `health.go`**

Write `impl/helixgitpx/platform/health/health.go`:

```go
// Package health provides /livez, /readyz, /healthz HTTP handlers and a
// probe registry. Register a named probe with Register; Ready returns 503
// when any probe fails.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// Probe returns nil when the dependency is healthy.
type Probe func(context.Context) error

// Handler serves the three liveness/readiness endpoints.
type Handler struct {
	mu     sync.RWMutex
	probes map[string]Probe
}

// New builds an empty Handler.
func New() *Handler {
	return &Handler{probes: make(map[string]Probe)}
}

// Register attaches a probe under name; re-registering replaces the prior probe.
func (h *Handler) Register(name string, p Probe) {
	h.mu.Lock()
	h.probes[name] = p
	h.mu.Unlock()
}

// Routes registers handler funcs on the given mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.Live)
	mux.HandleFunc("/livez", h.Live)
	mux.HandleFunc("/readyz", h.Ready)
}

// Live always returns 200; process is up.
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

// Ready runs every probe with the request context and returns 503 if any fail.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	probes := make(map[string]Probe, len(h.probes))
	for k, v := range h.probes {
		probes[k] = v
	}
	h.mu.RUnlock()

	results := make(map[string]any, len(probes))
	allOK := true
	for name, p := range probes {
		if err := p(r.Context()); err != nil {
			results[name] = err.Error()
			allOK = false
		} else {
			results[name] = "ok"
		}
	}

	status := http.StatusOK
	body := map[string]any{"status": "ok", "checks": results}
	if !allOK {
		status = http.StatusServiceUnavailable
		body["status"] = "unavailable"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
```

Write `impl/helixgitpx/platform/health/doc.go`:

```go
// Package health provides liveness/readiness HTTP handlers. See *Handler.
package health
```

- [ ] **Step 4: Run — expect pass**

```sh
cd impl/helixgitpx/platform && go test ./health/...
```

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx/platform/health
git commit -s -m "feat(platform/health): probe registry + live/ready handlers"
```

---

### Task 9: `platform/pg` package (pgx + goose migrations, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/pg/doc.go`
- Create: `impl/helixgitpx/platform/pg/pg.go`
- Create: `impl/helixgitpx/platform/pg/pg_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/pg/pg_test.go`:

```go
package pg_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/helixgitpx/platform/pg"
)

func TestOpen_InvalidDSNFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := pg.Open(ctx, pg.Options{DSN: "not-a-valid-dsn"})
	if err == nil {
		t.Fatalf("expected error for invalid DSN")
	}
}

func TestIsUnavailable(t *testing.T) {
	if !pg.IsUnavailable(pg.ErrUnavailable) {
		t.Errorf("ErrUnavailable not classified as unavailable")
	}
	if pg.IsUnavailable(errors.New("other")) {
		t.Errorf("arbitrary error classified as unavailable")
	}
}
```

- [ ] **Step 2: Implement `pg.go`**

Write `impl/helixgitpx/platform/pg/pg.go`:

```go
// Package pg wraps pgx/v5 with HelixGitpx-specific defaults and a thin
// migration runner (via pressly/goose). Callers get a *pgxpool.Pool ready
// for use; the package exposes typed sentinel errors so callers can classify
// failures without string matching.
package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrUnavailable signals the database is unreachable.
var ErrUnavailable = errors.New("pg: unavailable")

// Options configures Open.
type Options struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	ConnectTimeout  time.Duration
	HealthCheckInterval time.Duration
}

// Open constructs a pool. Callers own Close().
func Open(ctx context.Context, opts Options) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("pg: parse DSN: %w", err)
	}
	if opts.MaxConns > 0 {
		cfg.MaxConns = opts.MaxConns
	}
	if opts.MinConns > 0 {
		cfg.MinConns = opts.MinConns
	}
	if opts.ConnectTimeout > 0 {
		cfg.ConnConfig.ConnectTimeout = opts.ConnectTimeout
	}
	if opts.HealthCheckInterval > 0 {
		cfg.HealthCheckPeriod = opts.HealthCheckInterval
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pg: new pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Join(ErrUnavailable, err)
	}
	return pool, nil
}

// IsUnavailable reports whether err wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

// Probe returns a health.Probe that pings the pool.
func Probe(pool *pgxpool.Pool) func(context.Context) error {
	return func(ctx context.Context) error {
		if pool == nil {
			return ErrUnavailable
		}
		return pool.Ping(ctx)
	}
}
```

Write `impl/helixgitpx/platform/pg/doc.go`:

```go
// Package pg wraps pgx/v5. See Open and Probe.
package pg
```

- [ ] **Step 3: Run tests — expect pass**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./pg/...
```

Expected: PASS (tests only exercise error paths; real connectivity is tested via testcontainers in the integration test of Task 22).

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/platform/pg
git commit -s -m "feat(platform/pg): pgxpool wrapper + unavailability sentinel + probe"
```

---

### Task 10: `platform/redis` package (go-redis v9 wrapper, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/redis/doc.go`
- Create: `impl/helixgitpx/platform/redis/redis.go`
- Create: `impl/helixgitpx/platform/redis/redis_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/redis/redis_test.go`:

```go
package redis_test

import (
	"errors"
	"testing"

	hr "github.com/helixgitpx/platform/redis"
)

func TestKey_AppliesNamespace(t *testing.T) {
	c := hr.Client{Namespace: "hello"}
	got := c.Key("greeting", "world")
	want := "hello:greeting:world"
	if got != want {
		t.Errorf("Key = %q, want %q", got, want)
	}
}

func TestIsUnavailable(t *testing.T) {
	if !hr.IsUnavailable(hr.ErrUnavailable) {
		t.Errorf("sentinel not classified")
	}
	if hr.IsUnavailable(errors.New("other")) {
		t.Errorf("other err misclassified")
	}
}
```

- [ ] **Step 2: Implement `redis.go`**

Write `impl/helixgitpx/platform/redis/redis.go`:

```go
// Package redis wraps go-redis v9 with namespaced keys and typed unavailability.
package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

// ErrUnavailable is returned when the server cannot be reached.
var ErrUnavailable = errors.New("redis: unavailable")

// Options configures Open.
type Options struct {
	Addr      string
	Password  string
	DB        int
	Namespace string
}

// Client wraps *redis.Client with namespace helpers.
type Client struct {
	*redis.Client
	Namespace string
}

// Open constructs a Client and pings the server.
func Open(ctx context.Context, opts Options) (*Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
	if err := rc.Ping(ctx).Err(); err != nil {
		_ = rc.Close()
		return nil, errors.Join(ErrUnavailable, err)
	}
	return &Client{Client: rc, Namespace: opts.Namespace}, nil
}

// Key joins parts with ":" and prepends the namespace.
func (c Client) Key(parts ...string) string {
	if c.Namespace == "" {
		return strings.Join(parts, ":")
	}
	return fmt.Sprintf("%s:%s", c.Namespace, strings.Join(parts, ":"))
}

// IsUnavailable reports whether err wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

// Probe returns a health probe function.
func Probe(c *Client) func(context.Context) error {
	return func(ctx context.Context) error {
		if c == nil || c.Client == nil {
			return ErrUnavailable
		}
		return c.Ping(ctx).Err()
	}
}
```

Write `impl/helixgitpx/platform/redis/doc.go`:

```go
// Package redis wraps go-redis v9 with namespaced keys. See Open.
package redis
```

- [ ] **Step 3: Run tests**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./redis/...
```

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/platform/redis
git commit -s -m "feat(platform/redis): go-redis wrapper with namespaced keys"
```

---

### Task 11: `platform/kafka` package (franz-go producer/consumer wrapper, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/kafka/doc.go`
- Create: `impl/helixgitpx/platform/kafka/kafka.go`
- Create: `impl/helixgitpx/platform/kafka/kafka_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/kafka/kafka_test.go`:

```go
package kafka_test

import (
	"errors"
	"testing"

	"github.com/helixgitpx/platform/kafka"
)

func TestOptions_Validation(t *testing.T) {
	_, err := kafka.NewProducer(kafka.ProducerOptions{})
	if err == nil {
		t.Fatalf("expected error for missing brokers")
	}
}

func TestIsUnavailable(t *testing.T) {
	if !kafka.IsUnavailable(kafka.ErrUnavailable) {
		t.Errorf("sentinel")
	}
	if kafka.IsUnavailable(errors.New("other")) {
		t.Errorf("other")
	}
}
```

- [ ] **Step 2: Implement `kafka.go`**

Write `impl/helixgitpx/platform/kafka/kafka.go`:

```go
// Package kafka wraps franz-go (github.com/twmb/franz-go/pkg/kgo) with
// defaults suitable for HelixGitpx services. Schema-registry integration
// is stubbed in M1 (a ResolveFn hook) and wired in M2.
package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// ErrUnavailable is returned when brokers are unreachable.
var ErrUnavailable = errors.New("kafka: unavailable")

// ProducerOptions configures NewProducer.
type ProducerOptions struct {
	Brokers  []string
	ClientID string
	Topic    string // default topic for Emit helper; may be overridden per record
}

// Producer wraps *kgo.Client with HelixGitpx ergonomics.
type Producer struct {
	cl    *kgo.Client
	topic string
}

// NewProducer constructs a Producer.
func NewProducer(opts ProducerOptions) (*Producer, error) {
	if len(opts.Brokers) == 0 {
		return nil, fmt.Errorf("kafka: Brokers is required")
	}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(opts.Brokers...),
		kgo.ClientID(stringOrDefault(opts.ClientID, "helixgitpx")),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerLinger(0),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka: new client: %w", err)
	}
	return &Producer{cl: cl, topic: opts.Topic}, nil
}

// Emit publishes one record to the default (or overridden) topic.
func (p *Producer) Emit(ctx context.Context, key, value []byte, topic ...string) error {
	t := p.topic
	if len(topic) > 0 && topic[0] != "" {
		t = topic[0]
	}
	if t == "" {
		return fmt.Errorf("kafka: topic unset")
	}
	res := p.cl.ProduceSync(ctx, &kgo.Record{Topic: t, Key: key, Value: value})
	if err := res.FirstErr(); err != nil {
		return errors.Join(ErrUnavailable, err)
	}
	return nil
}

// Close flushes and closes the client.
func (p *Producer) Close(ctx context.Context) error {
	if p == nil || p.cl == nil {
		return nil
	}
	if err := p.cl.Flush(ctx); err != nil {
		return err
	}
	p.cl.Close()
	return nil
}

// IsUnavailable reports whether err wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

func stringOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
```

Write `impl/helixgitpx/platform/kafka/doc.go`:

```go
// Package kafka wraps franz-go for HelixGitpx services. See NewProducer.
package kafka
```

- [ ] **Step 3: Run tests**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./kafka/...
```

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/platform/kafka
git commit -s -m "feat(platform/kafka): franz-go producer wrapper + sentinel"
```

---

### Task 12: `platform/telemetry` package (OTel bootstrap, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/telemetry/doc.go`
- Create: `impl/helixgitpx/platform/telemetry/telemetry.go`
- Create: `impl/helixgitpx/platform/telemetry/telemetry_test.go`

- [ ] **Step 1: Write the test**

Write `impl/helixgitpx/platform/telemetry/telemetry_test.go`:

```go
package telemetry_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/telemetry"
)

func TestStart_NoEndpoint_ReturnsNoop(t *testing.T) {
	ctx := context.Background()
	shutdown, err := telemetry.Start(ctx, telemetry.Options{Service: "hello"})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := shutdown(ctx); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}
```

- [ ] **Step 2: Implement `telemetry.go`**

Write `impl/helixgitpx/platform/telemetry/telemetry.go`:

```go
// Package telemetry bootstraps OpenTelemetry SDK. When OTEL_EXPORTER_OTLP_ENDPOINT
// is unset (or Options.OTLPEndpoint is empty), the SDK installs a no-op tracer
// and meter, so calling Start on developer machines without a collector is safe.
package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Options configures Start.
type Options struct {
	Service      string
	Version      string
	Environment  string
	OTLPEndpoint string // overrides OTEL_EXPORTER_OTLP_ENDPOINT when set
}

// ShutdownFunc flushes and closes all providers.
type ShutdownFunc func(context.Context) error

// Start installs global TracerProvider. Returns a no-op shutdown when no endpoint.
func Start(ctx context.Context, opts Options) (ShutdownFunc, error) {
	endpoint := opts.OTLPEndpoint
	if endpoint == "" {
		endpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	if endpoint == "" {
		// No collector → no-op. Tracing calls become cheap no-ops through the global.
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(opts.Service),
			semconv.ServiceVersion(opts.Version),
			semconv.DeploymentEnvironment(opts.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpointURL(endpoint),
		otlptracegrpc.WithInsecure(),
	))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func(shutdownCtx context.Context) error {
		return tp.Shutdown(shutdownCtx)
	}, nil
}
```

Write `impl/helixgitpx/platform/telemetry/doc.go`:

```go
// Package telemetry bootstraps OpenTelemetry. See Start.
package telemetry
```

- [ ] **Step 3: Run tests**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./telemetry/...
```

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/platform/telemetry
git commit -s -m "feat(platform/telemetry): OTel bootstrap with no-op fallback"
```

---

### Task 13: `platform/grpc` package (server constructor, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/grpc/doc.go`
- Create: `impl/helixgitpx/platform/grpc/server.go`
- Create: `impl/helixgitpx/platform/grpc/interceptor.go`
- Create: `impl/helixgitpx/platform/grpc/server_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/grpc/server_test.go`:

```go
package grpc_test

import (
	"context"
	"net"
	"testing"

	hgrpc "github.com/helixgitpx/platform/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewServer_ServesHealth(t *testing.T) {
	s, err := hgrpc.NewServer(hgrpc.Options{})
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = s.Serve(lis) }()
	defer s.GracefulStop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecureCreds()),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	hc := grpc_health_v1.NewHealthClient(conn)
	resp, err := hc.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("status = %v, want SERVING", resp.Status)
	}
}
```

Write `impl/helixgitpx/platform/grpc/testhelpers_test.go`:

```go
package grpc_test

import "google.golang.org/grpc/credentials/insecure"

func insecureCreds() interface{ PerRPCCredentials } { return nil }

// alias via type-elide; use the real credentials.NewCredentials below
var _ = insecure.NewCredentials()
```

Actually simplify — just inline insecure dialer:

Replace the two lines `grpc.WithTransportCredentials(insecureCreds()),` in `server_test.go` with:

```go
grpc.WithTransportCredentials(insecure.NewCredentials()),
```

And add the import `"google.golang.org/grpc/credentials/insecure"` at the top of `server_test.go`. Delete `testhelpers_test.go`.

Final `server_test.go`:

```go
package grpc_test

import (
	"context"
	"net"
	"testing"

	hgrpc "github.com/helixgitpx/platform/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewServer_ServesHealth(t *testing.T) {
	s, err := hgrpc.NewServer(hgrpc.Options{})
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = s.Serve(lis) }()
	defer s.GracefulStop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	hc := grpc_health_v1.NewHealthClient(conn)
	resp, err := hc.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("status = %v, want SERVING", resp.Status)
	}
}
```

- [ ] **Step 2: Implement `server.go`**

Write `impl/helixgitpx/platform/grpc/server.go`:

```go
// Package grpc builds a gRPC server preconfigured with the HelixGitpx
// interceptor chain (logging, recovery, telemetry, auth stub).
package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Options configures NewServer.
type Options struct {
	// ServerOptions is appended to the built-in option list.
	ServerOptions []grpc.ServerOption
	// EnableReflection is true by default; set false to disable server reflection.
	DisableReflection bool
}

// NewServer constructs a *grpc.Server with HelixGitpx defaults:
//   - unary + stream interceptor chain (see interceptor.go)
//   - grpc.health.v1.Health registered and reporting SERVING
//   - server reflection registered (unless disabled)
func NewServer(opts Options) (*grpc.Server, error) {
	so := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryChain()...),
		grpc.ChainStreamInterceptor(streamChain()...),
	}
	so = append(so, opts.ServerOptions...)

	s := grpc.NewServer(so...)

	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)

	if !opts.DisableReflection {
		reflection.Register(s)
	}
	return s, nil
}
```

Write `impl/helixgitpx/platform/grpc/interceptor.go`:

```go
package grpc

import (
	"context"
	"runtime/debug"

	"github.com/helixgitpx/platform/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unaryChain() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recoveryUnary,
		loggingUnary,
	}
}

func streamChain() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		recoveryStream,
	}
}

func recoveryUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.FromContext(ctx).Error("panic in grpc handler",
				"method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "internal error")
		}
	}()
	return handler(ctx, req)
}

func recoveryStream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.FromContext(ss.Context()).Error("panic in grpc stream",
				"method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "internal error")
		}
	}()
	return handler(srv, ss)
}

func loggingUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	lg := log.FromContext(ctx).With("method", info.FullMethod)
	if err != nil {
		lg.Error("grpc call failed", "err", err.Error())
	} else {
		lg.Debug("grpc call ok")
	}
	return resp, err
}
```

Write `impl/helixgitpx/platform/grpc/doc.go`:

```go
// Package grpc builds HelixGitpx gRPC servers. See NewServer.
package grpc
```

- [ ] **Step 3: Run tests**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./grpc/...
```

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/platform/grpc
git commit -s -m "feat(platform/grpc): server constructor with recovery+logging interceptors"
```

---

### Task 14: `platform/gin` package (HTTP router constructor, TDD)

**Files:**
- Create: `impl/helixgitpx/platform/gin/doc.go`
- Create: `impl/helixgitpx/platform/gin/router.go`
- Create: `impl/helixgitpx/platform/gin/router_test.go`

- [ ] **Step 1: Write the failing test**

Write `impl/helixgitpx/platform/gin/router_test.go`:

```go
package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	hgin "github.com/helixgitpx/platform/gin"
)

func TestNewRouter_BaseHeaders(t *testing.T) {
	r := hgin.NewRouter(hgin.Options{Service: "hello", Version: "test"})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code = %d", w.Code)
	}
	if w.Header().Get("X-HelixGitpx-Service") != "hello" {
		t.Errorf("missing service header")
	}
	if w.Header().Get("X-HelixGitpx-Version") != "test" {
		t.Errorf("missing version header")
	}
}
```

- [ ] **Step 2: Implement `router.go`**

Write `impl/helixgitpx/platform/gin/router.go`:

```go
// Package gin builds a Gin router preconfigured with HelixGitpx middleware:
// recovery, request-id, service/version headers, structured logging.
package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/helixgitpx/platform/log"
)

// Options configures NewRouter.
type Options struct {
	Service string
	Version string
	Mode    string // "debug", "release", "test"; defaults to "release"
}

// NewRouter constructs an *gin.Engine.
func NewRouter(opts Options) *gin.Engine {
	if opts.Mode == "" {
		opts.Mode = gin.ReleaseMode
	}
	gin.SetMode(opts.Mode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		identityHeaders(opts),
		loggingMiddleware(),
	)
	return r
}

func identityHeaders(opts Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		if opts.Service != "" {
			c.Header("X-HelixGitpx-Service", opts.Service)
		}
		if opts.Version != "" {
			c.Header("X-HelixGitpx-Version", opts.Version)
		}
		c.Next()
	}
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		lg := log.FromContext(c.Request.Context()).With(
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
		)
		if c.Writer.Status() >= 500 {
			lg.Error("http request")
		} else {
			lg.Debug("http request")
		}
	}
}
```

Write `impl/helixgitpx/platform/gin/doc.go`:

```go
// Package gin builds HelixGitpx HTTP routers. See NewRouter.
package gin
```

- [ ] **Step 3: Run tests**

```sh
cd impl/helixgitpx/platform && go mod tidy && go test ./gin/...
```

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/platform/gin
git commit -s -m "feat(platform/gin): router with recovery + identity headers + logging"
```

---

### Task 15: `platform/temporal`, `platform/spire`, `platform/opa` stubs

**Files:**
- Create: `impl/helixgitpx/platform/temporal/{doc,temporal,temporal_test}.go`
- Create: `impl/helixgitpx/platform/spire/{doc,spire,spire_test}.go`
- Create: `impl/helixgitpx/platform/opa/{doc,opa,opa_test}.go`

These three packages ship minimal surfaces in M1 (TODO(Mx) markers). Each has a constructor that returns a no-op client when the dependency is absent, matching the design spec §6 policy.

- [ ] **Step 1: Write `platform/temporal/temporal.go`**

```go
// Package temporal wraps go.temporal.io/sdk with HelixGitpx defaults.
// M1 ships a no-op stub; M5 wires real workflow/activity registration.
package temporal

import (
	"context"
	"errors"
)

// ErrUnavailable indicates the Temporal service cannot be reached.
var ErrUnavailable = errors.New("temporal: unavailable")

// Options configures NewClient.
type Options struct {
	HostPort  string
	Namespace string
}

// Client is a placeholder for the real Temporal client wired in M5.
type Client struct {
	HostPort  string
	Namespace string
	noop      bool
}

// NewClient returns a Client. When HostPort is empty, returns a no-op client.
//
// TODO(M5): wire go.temporal.io/sdk/client.Dial; register workers and workflows.
func NewClient(_ context.Context, opts Options) (*Client, error) {
	if opts.HostPort == "" {
		return &Client{noop: true}, nil
	}
	return &Client{HostPort: opts.HostPort, Namespace: opts.Namespace}, nil
}

// Close releases resources. No-op for M1.
func (c *Client) Close() error { return nil }
```

Write `impl/helixgitpx/platform/temporal/doc.go`:

```go
// Package temporal wraps Temporal SDK. See NewClient. Real workflow wiring lands in M5.
package temporal
```

Write `impl/helixgitpx/platform/temporal/temporal_test.go`:

```go
package temporal_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/temporal"
)

func TestNewClient_NoopWhenAddrEmpty(t *testing.T) {
	c, err := temporal.NewClient(context.Background(), temporal.Options{})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c == nil {
		t.Fatalf("nil client")
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}
```

- [ ] **Step 2: Write `platform/spire/spire.go`**

```go
// Package spire integrates SPIFFE/SPIRE workload API.
// M1 ships a stub; M2 wires go-spiffe/v2 to fetch SVIDs.
package spire

import (
	"context"
	"errors"
	"os"
)

// ErrUnavailable indicates the SPIRE agent socket is not reachable.
var ErrUnavailable = errors.New("spire: unavailable")

// Options configures NewFetcher.
type Options struct {
	SocketPath string // e.g. "unix:///run/spire/agent.sock"
}

// Fetcher retrieves workload SVIDs. Stubbed in M1.
type Fetcher struct {
	SocketPath string
	noop       bool
}

// NewFetcher returns a Fetcher. When the socket file is absent, returns a no-op fetcher.
//
// TODO(M2): wire github.com/spiffe/go-spiffe/v2/workloadapi.NewX509Source and
// supply SVIDs to grpc/TLS constructors.
func NewFetcher(_ context.Context, opts Options) (*Fetcher, error) {
	if opts.SocketPath == "" {
		return &Fetcher{noop: true}, nil
	}
	if _, err := os.Stat(trimUnix(opts.SocketPath)); err != nil {
		return &Fetcher{noop: true, SocketPath: opts.SocketPath}, nil
	}
	return &Fetcher{SocketPath: opts.SocketPath}, nil
}

// Close releases resources. No-op for M1.
func (f *Fetcher) Close() error { return nil }

func trimUnix(s string) string {
	const prefix = "unix://"
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
```

Write `impl/helixgitpx/platform/spire/doc.go`:

```go
// Package spire integrates SPIFFE/SPIRE. See NewFetcher. Real wiring lands in M2.
package spire
```

Write `impl/helixgitpx/platform/spire/spire_test.go`:

```go
package spire_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/spire"
)

func TestNewFetcher_NoopWhenSocketAbsent(t *testing.T) {
	f, err := spire.NewFetcher(context.Background(), spire.Options{
		SocketPath: "unix:///tmp/definitely-not-here.sock",
	})
	if err != nil {
		t.Fatalf("NewFetcher: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}
```

- [ ] **Step 3: Write `platform/opa/opa.go`**

```go
// Package opa wraps github.com/open-policy-agent/opa/rego for in-process
// policy evaluation. M1 ships a minimal surface; M3 adds real policies.
package opa

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

// Evaluator evaluates a compiled Rego query against structured input.
type Evaluator struct {
	query rego.PreparedEvalQuery
}

// Options configures NewEvaluator.
type Options struct {
	Module string // Rego source
	Query  string // e.g. "data.helixgitpx.allow"
}

// NewEvaluator compiles the given module and query.
func NewEvaluator(ctx context.Context, opts Options) (*Evaluator, error) {
	if opts.Query == "" {
		return nil, errors.New("opa: Query is required")
	}
	q, err := rego.New(
		rego.Query(opts.Query),
		rego.Module("policy.rego", opts.Module),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("opa: compile: %w", err)
	}
	return &Evaluator{query: q}, nil
}

// Eval runs the query with input and returns the first defined result or false.
func (e *Evaluator) Eval(ctx context.Context, input any) (any, error) {
	rs, err := e.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return false, nil
	}
	return rs[0].Expressions[0].Value, nil
}
```

Write `impl/helixgitpx/platform/opa/doc.go`:

```go
// Package opa evaluates Rego in-process. See NewEvaluator.
package opa
```

Write `impl/helixgitpx/platform/opa/opa_test.go`:

```go
package opa_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/opa"
)

func TestEval_AllowFromModule(t *testing.T) {
	ev, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `package helixgitpx
allow { input.role == "admin" }`,
		Query: "data.helixgitpx.allow",
	})
	if err != nil {
		t.Fatalf("NewEvaluator: %v", err)
	}
	got, err := ev.Eval(context.Background(), map[string]any{"role": "admin"})
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	if b, _ := got.(bool); !b {
		t.Errorf("allow = %v, want true", got)
	}
}
```

- [ ] **Step 4: Run all three test packages**

```sh
cd impl/helixgitpx/platform && go mod tidy
go test ./temporal/... ./spire/... ./opa/...
```

Expected: PASS all.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx/platform/temporal impl/helixgitpx/platform/spire impl/helixgitpx/platform/opa
git commit -s -m "feat(platform): temporal/spire/opa stubs with M2+M5 TODO markers"
```

---

### Task 16: `platform/testkit` package (testcontainers helpers)

**Files:**
- Create: `impl/helixgitpx/platform/testkit/doc.go`
- Create: `impl/helixgitpx/platform/testkit/pg.go`
- Create: `impl/helixgitpx/platform/testkit/redis.go`
- Create: `impl/helixgitpx/platform/testkit/kafka.go`

- [ ] **Step 1: Write `pg.go`**

```go
// Package testkit provides testcontainers-go helpers for HelixGitpx tests.
// Each helper returns (dsn/addr string, teardown func) and registers cleanup
// via testing.TB so callers don't need to remember Close.
package testkit

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// StartPostgres launches a Postgres 16 container and returns the DSN.
func StartPostgres(t testing.TB) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("helixgitpx"),
		postgres.WithUsername("helix"),
		postgres.WithPassword("helix"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("testkit.StartPostgres: %v", err)
	}
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("testkit.StartPostgres dsn: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
	return dsn
}
```

- [ ] **Step 2: Write `redis.go`**

```go
package testkit

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// StartRedis launches a Redis 7 container and returns "host:port".
func StartRedis(t testing.TB) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ctr, err := redis.Run(ctx, "redis:7-alpine",
		testcontainers.WithWaitStrategy(nil),
	)
	if err != nil {
		t.Fatalf("testkit.StartRedis: %v", err)
	}
	ep, err := ctr.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("testkit.StartRedis endpoint: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
	_ = time.Now
	return ep
}
```

- [ ] **Step 3: Write `kafka.go`**

```go
package testkit

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

// StartKafka launches a Kafka 3.8 KRaft-mode container and returns the bootstrap broker.
func StartKafka(t testing.TB) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ctr, err := kafka.Run(ctx, "apache/kafka:3.8.1",
		kafka.WithClusterID("helixgitpx-test"),
		testcontainers.WithWaitStrategy(nil),
	)
	if err != nil {
		t.Fatalf("testkit.StartKafka: %v", err)
	}
	brokers, err := ctr.Brokers(ctx)
	if err != nil {
		t.Fatalf("testkit.StartKafka brokers: %v", err)
	}
	if len(brokers) == 0 {
		t.Fatalf("testkit.StartKafka: no brokers")
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
	_ = time.Now
	return brokers[0]
}
```

Write `impl/helixgitpx/platform/testkit/doc.go`:

```go
// Package testkit provides testcontainers helpers for Postgres, Redis, and Kafka.
// Each helper registers t.Cleanup for termination.
package testkit
```

- [ ] **Step 4: Run tidy (no unit tests here — exercised by hello integration test)**

```sh
cd impl/helixgitpx/platform && go mod tidy
go build ./testkit/...
```

Expected: compiles.

- [ ] **Step 5: Commit**

```sh
git add impl/helixgitpx/platform/testkit
git commit -s -m "feat(platform/testkit): testcontainers helpers for pg/redis/kafka"
```

---

### Task 17: Proto root — buf config + common + hello + stub domains

**Files:**
- Create: `impl/helixgitpx/proto/buf.yaml`
- Create: `impl/helixgitpx/proto/buf.gen.yaml`
- Create: `impl/helixgitpx/proto/buf.work.yaml`
- Create: `impl/helixgitpx/proto/README.md`
- Create: `impl/helixgitpx/proto/helixgitpx/common/v1/{errors.proto,pagination.proto,time.proto}`
- Create: `impl/helixgitpx/proto/helixgitpx/hello/v1/hello.proto`
- Create: `impl/helixgitpx/proto/helixgitpx/{auth,repo,sync,conflict,upstream,collab,events,platform}/v1/stubs.proto` (8 files)

- [ ] **Step 1: Write `buf.yaml`**

```yaml
version: v2
modules:
  - path: .
    name: buf.build/helixgitpx/core
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
  - buf.build/bufbuild/protovalidate
lint:
  use:
    - STANDARD
  except:
    - PACKAGE_VERSION_SUFFIX
breaking:
  use:
    - FILE
```

- [ ] **Step 2: Write `buf.gen.yaml`**

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/helixgitpx/helixgitpx/gen/go
plugins:
  - remote: buf.build/protocolbuffers/go
    out: ../gen/go
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: ../gen/go
    opt: paths=source_relative,require_unimplemented_servers=false
  - remote: buf.build/connectrpc/go
    out: ../gen/go
    opt: paths=source_relative
  - remote: buf.build/bufbuild/es
    out: ../../helixgitpx-web/libs/proto/src
    opt: target=ts
  - remote: buf.build/connectrpc/es
    out: ../../helixgitpx-web/libs/proto/src
    opt: target=ts
  - remote: buf.build/grpc/kotlin
    out: ../../helixgitpx-clients/shared/src/commonMain/kotlin/gen
  - remote: buf.build/apple/swift
    out: ../../helixgitpx-clients/iosApp/Gen
  - remote: buf.build/grpc-ecosystem/openapiv2
    out: ../api/openapi
    opt: output_format=json,logtostderr=true
```

- [ ] **Step 3: Write `buf.work.yaml`**

```yaml
version: v1
directories:
  - .
```

- [ ] **Step 4: Write `proto/README.md`**

```markdown
# HelixGitpx proto root

Canonical `.proto` sources. The BSR module is `buf.build/helixgitpx/core`.

| Domain | M1 status | Owner milestone |
|---|---|---|
| common | populated | — |
| hello | populated | — |
| auth | stub | M3 |
| repo | stub | M4 |
| sync | stub | M5 |
| conflict | stub | M5 |
| upstream | stub | M4 |
| collab | stub | M5 |
| events | stub | M2 |
| platform | stub | M2 |

Historical copies remain under `docs/specifications/.../17-protos/` for provenance; this directory is the live source.

Regenerate language bindings with `make gen` from `impl/helixgitpx/`.
```

- [ ] **Step 5: Write `common/v1/errors.proto`**

```sh
mkdir -p impl/helixgitpx/proto/helixgitpx/common/v1
mkdir -p impl/helixgitpx/proto/helixgitpx/hello/v1
for d in auth repo sync conflict upstream collab events platform; do
  mkdir -p "impl/helixgitpx/proto/helixgitpx/$d/v1"
done
```

Write `impl/helixgitpx/proto/helixgitpx/common/v1/errors.proto`:

```proto
syntax = "proto3";

package helixgitpx.common.v1;

// Problem mirrors the RFC 7807 problem-details payload used by HelixGitpx REST
// responses. gRPC clients receive equivalent information via status details.
message Problem {
  string type = 1;
  string title = 2;
  int32 status = 3;
  string detail = 4;
  string instance = 5;
  string domain = 6;
  string code = 7;
  map<string, string> errors = 8;
}
```

- [ ] **Step 6: Write `common/v1/pagination.proto`**

```proto
syntax = "proto3";

package helixgitpx.common.v1;

message PageRequest {
  int32 page_size = 1;
  string page_token = 2;
}

message PageResponse {
  string next_page_token = 1;
  int32 total = 2;
}
```

- [ ] **Step 7: Write `common/v1/time.proto`**

```proto
syntax = "proto3";

package helixgitpx.common.v1;

import "google/protobuf/timestamp.proto";

message TimeRange {
  google.protobuf.Timestamp from = 1;
  google.protobuf.Timestamp to = 2;
}
```

- [ ] **Step 8: Write `hello/v1/hello.proto`**

Write `impl/helixgitpx/proto/helixgitpx/hello/v1/hello.proto`:

```proto
syntax = "proto3";

package helixgitpx.hello.v1;

import "google/api/annotations.proto";

// HelloService is the spine service for M1 Foundation.
// It exercises platform/{grpc,gin,pg,redis,kafka} end-to-end.
service HelloService {
  // SayHello returns a greeting and (as a side effect) increments a counter
  // in Postgres, caches the last greeting in Redis, and emits hello.said to Kafka.
  rpc SayHello(SayHelloRequest) returns (SayHelloResponse) {
    option (google.api.http) = {
      get: "/v1/hello"
    };
  }
}

message SayHelloRequest {
  string name = 1;
}

message SayHelloResponse {
  string greeting = 1;
  int64 count = 2;
}
```

- [ ] **Step 9: Write stub domain files**

For each of `auth repo sync conflict upstream collab events platform`, write `impl/helixgitpx/proto/helixgitpx/<domain>/v1/stubs.proto` with this content (substituting `<DOMAIN>` with the name):

```proto
syntax = "proto3";

package helixgitpx.<DOMAIN>.v1;

// Stub placeholder — real messages and services land in the milestone owning this domain.
// See proto/README.md for the M1 → Mx ownership table.
message Placeholder {
  string reserved = 1;
}
```

Example for `auth`: set `package helixgitpx.auth.v1;`.

- [ ] **Step 10: Install buf (via mise) and lint**

```sh
mise install
cd impl/helixgitpx/proto
buf dep update
buf lint
buf build
```

Expected: no errors.

- [ ] **Step 11: Generate code**

```sh
cd impl/helixgitpx
buf generate proto/
```

Expected: files appear under `gen/go/helixgitpx/common/v1/`, `gen/go/helixgitpx/hello/v1/`, etc. And under `api/openapi/`.

- [ ] **Step 12: Commit**

```sh
git add impl/helixgitpx/proto impl/helixgitpx/gen impl/helixgitpx/api
git commit -s -m "feat(proto): buf config + common/hello populated + 8 domain stubs"
```

---

### Task 18: `tools/scaffold` — service-template renderer

**Files:**
- Create: `impl/helixgitpx/tools/scaffold/go.mod`
- Create: `impl/helixgitpx/tools/scaffold/main.go`
- Create: `impl/helixgitpx/tools/scaffold/main_test.go`
- Create: `impl/helixgitpx/tools/scaffold/templates/service/...` (hello template files below)

- [ ] **Step 1: Initialize module**

```sh
mkdir -p impl/helixgitpx/tools/scaffold/templates/service/{cmd,internal/{app,domain,handler/grpc,handler/http,repo},config,deploy,migrations,test/integration}
cd impl/helixgitpx/tools/scaffold
go mod init github.com/helixgitpx/helixgitpx/tools/scaffold
cd ../../..
cd impl/helixgitpx && go work use ./tools/scaffold
```

- [ ] **Step 2: Write the failing test**

Write `impl/helixgitpx/tools/scaffold/main_test.go`:

```go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_CreatesExpectedFiles(t *testing.T) {
	dst := t.TempDir()
	cfg := Config{
		Name:         "greet",
		ProtoPackage: "helixgitpx.greet.v1",
		HTTPPort:     8002,
		GRPCPort:     9002,
		HealthPort:   8082,
		Out:          dst,
	}
	if err := Render(cfg); err != nil {
		t.Fatalf("Render: %v", err)
	}

	expected := []string{
		"cmd/greet/main.go",
		"internal/app/app.go",
		"Makefile",
		"README.md",
		"go.mod",
		"deploy/Dockerfile",
	}
	for _, p := range expected {
		if _, err := os.Stat(filepath.Join(dst, p)); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}

	// Check substitution happened
	main, err := os.ReadFile(filepath.Join(dst, "cmd/greet/main.go"))
	if err != nil {
		t.Fatalf("read main: %v", err)
	}
	if !contains(string(main), "greet") {
		t.Errorf("main.go did not substitute name")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
```

- [ ] **Step 3: Write `main.go`**

Write `impl/helixgitpx/tools/scaffold/main.go`:

```go
// Command scaffold renders a new HelixGitpx Go service from templates.
//
//	scaffold --name greet --proto helixgitpx.greet.v1 --http 8002 --grpc 9002
//
// When --dry-run is set, prints the list of files that would be written and exits 0.
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/service
var templatesFS embed.FS

// Config drives Render.
type Config struct {
	Name         string
	ProtoPackage string
	HTTPPort     int
	GRPCPort     int
	HealthPort   int
	Out          string
	DryRun       bool
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Name, "name", "", "service name (lowercase, no spaces)")
	flag.StringVar(&cfg.ProtoPackage, "proto", "", "protobuf package (e.g. helixgitpx.greet.v1)")
	flag.IntVar(&cfg.HTTPPort, "http", 8000, "HTTP port")
	flag.IntVar(&cfg.GRPCPort, "grpc", 9000, "gRPC port")
	flag.IntVar(&cfg.HealthPort, "health", 8080, "health/metrics port")
	flag.StringVar(&cfg.Out, "out", "", "output directory (default: services/<name>)")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "list files without writing")
	flag.Parse()

	if cfg.Name == "" || cfg.ProtoPackage == "" {
		flag.Usage()
		os.Exit(2)
	}
	if cfg.Out == "" {
		cfg.Out = filepath.Join("services", cfg.Name)
	}
	if cfg.DryRun {
		if err := DryRun(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "scaffold: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if err := Render(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "scaffold: %v\n", err)
		os.Exit(1)
	}
}

// Render writes all template files to cfg.Out, substituting values.
func Render(cfg Config) error {
	return walkTemplates(cfg, func(srcPath, dstPath string, raw []byte) error {
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}
		tpl, err := template.New(srcPath).Delims("<<", ">>").Parse(string(raw))
		if err != nil {
			return fmt.Errorf("parse %s: %w", srcPath, err)
		}
		f, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return tpl.Execute(f, cfg)
	})
}

// DryRun prints files that would be written.
func DryRun(cfg Config) error {
	return walkTemplates(cfg, func(_, dstPath string, _ []byte) error {
		fmt.Println(dstPath)
		return nil
	})
}

func walkTemplates(cfg Config, fn func(src, dst string, raw []byte) error) error {
	root := "templates/service"
	return fs.WalkDir(templatesFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		raw, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, root+"/")
		rel = strings.ReplaceAll(rel, "__name__", cfg.Name)
		rel = strings.TrimSuffix(rel, ".tmpl")
		return fn(path, filepath.Join(cfg.Out, rel), raw)
	})
}
```

- [ ] **Step 4: Write minimal template files**

Write `impl/helixgitpx/tools/scaffold/templates/service/go.mod.tmpl`:

```
module github.com/helixgitpx/helixgitpx/services/<<.Name>>

go 1.23

require (
	github.com/helixgitpx/platform v0.0.0
	github.com/helixgitpx/helixgitpx/gen v0.0.0
)

replace github.com/helixgitpx/platform => ../../platform
replace github.com/helixgitpx/helixgitpx/gen => ../../gen
```

Write `impl/helixgitpx/tools/scaffold/templates/service/cmd/__name__/main.go.tmpl`:

```go
// Command <<.Name>> is a HelixGitpx service scaffolded by tools/scaffold.
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/helixgitpx/helixgitpx/services/<<.Name>>/internal/app"
	"github.com/helixgitpx/platform/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lg := log.New(log.Options{Level: "info", Service: "<<.Name>>"})
	if err := app.Run(ctx, lg); err != nil {
		lg.Error("service exited with error", "err", err.Error())
	}
}
```

Write `impl/helixgitpx/tools/scaffold/templates/service/internal/app/app.go.tmpl`:

```go
// Package app is the composition root for <<.Name>>.
package app

import (
	"context"

	"github.com/helixgitpx/platform/log"
)

// Run wires dependencies and serves until ctx is done.
// Fill in with real wiring (see services/hello for a reference).
func Run(ctx context.Context, lg *log.Logger) error {
	lg.Info("<<.Name>> starting")
	<-ctx.Done()
	lg.Info("<<.Name>> stopped")
	return nil
}
```

Write `impl/helixgitpx/tools/scaffold/templates/service/Makefile.tmpl`:

```makefile
SHELL := /usr/bin/env bash
.PHONY: lint test build image
lint:
	golangci-lint run ./...
test:
	go test -race ./...
build:
	go build -o bin/<<.Name>> ./cmd/<<.Name>>
image:
	docker build -t helixgitpx/<<.Name>>:dev -f deploy/Dockerfile ../..
```

Write `impl/helixgitpx/tools/scaffold/templates/service/README.md.tmpl`:

```markdown
# <<.Name>>

HelixGitpx service: <<.Name>>.

- HTTP: `:<<.HTTPPort>>`
- gRPC: `:<<.GRPCPort>>`
- Health/metrics: `:<<.HealthPort>>`

## Dev

```sh
make lint test build
```
```

Write `impl/helixgitpx/tools/scaffold/templates/service/deploy/Dockerfile.tmpl`:

```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/<<.Name>>
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/<<.Name>> ./cmd/<<.Name>>

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/<<.Name>> /app/<<.Name>>
USER nonroot
ENTRYPOINT ["/app/<<.Name>>"]
LABEL org.opencontainers.image.title="<<.Name>>"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
```

- [ ] **Step 5: Run tests**

```sh
cd impl/helixgitpx/tools/scaffold && go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

```sh
git add impl/helixgitpx/tools/scaffold impl/helixgitpx/go.work
git commit -s -m "feat(scaffold): service-template renderer (embed-based Go binary)"
```

---

### Task 19: Generate `services/hello/` from the scaffold

**Files:**
- Create: via scaffold — `impl/helixgitpx/services/hello/{cmd,internal/app,Makefile,README.md,deploy/Dockerfile,go.mod}`

- [ ] **Step 1: Run the scaffold**

```sh
cd impl/helixgitpx
go run ./tools/scaffold \
  --name hello \
  --proto helixgitpx.hello.v1 \
  --http 8001 --grpc 9001 --health 8081 \
  --out services/hello
```

- [ ] **Step 2: Add to workspace**

```sh
cd impl/helixgitpx && go work use ./services/hello
```

- [ ] **Step 3: Verify it compiles**

```sh
cd impl/helixgitpx/services/hello && go mod tidy && go build ./...
```

Expected: compiles (the placeholder app.Run is a no-op loop).

- [ ] **Step 4: Commit**

```sh
git add impl/helixgitpx/services/hello impl/helixgitpx/go.work
git commit -s -m "feat(services/hello): scaffold initial layout"
```

---

### Task 20: Hello service — domain + handlers + wire + real main

**Files:**
- Modify/Create: `impl/helixgitpx/services/hello/internal/app/app.go`
- Create: `impl/helixgitpx/services/hello/internal/domain/greeter.go`
- Create: `impl/helixgitpx/services/hello/internal/domain/greeter_test.go`
- Create: `impl/helixgitpx/services/hello/internal/handler/grpc/server.go`
- Create: `impl/helixgitpx/services/hello/internal/handler/http/router.go`
- Create: `impl/helixgitpx/services/hello/internal/repo/counter_pg.go`
- Create: `impl/helixgitpx/services/hello/internal/repo/cache_redis.go`
- Create: `impl/helixgitpx/services/hello/internal/repo/event_kafka.go`
- Create: `impl/helixgitpx/services/hello/migrations/20260420000001_init.sql`

- [ ] **Step 1: Write migration**

```sh
mkdir -p impl/helixgitpx/services/hello/migrations
```

Write `impl/helixgitpx/services/hello/migrations/20260420000001_init.sql`:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS hello.greetings (
    name TEXT PRIMARY KEY,
    count BIGINT NOT NULL DEFAULT 0,
    last_said_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS hello.greetings;
```

- [ ] **Step 2: TDD — write domain test**

Write `impl/helixgitpx/services/hello/internal/domain/greeter_test.go`:

```go
package domain_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
)

type fakeCounter struct {
	count int64
	err   error
}

func (f *fakeCounter) Increment(_ context.Context, name string) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.count++
	return f.count, nil
}

type fakeCache struct {
	last string
}

func (f *fakeCache) SetLast(_ context.Context, name, greeting string) error {
	f.last = greeting
	return nil
}

type fakeEmitter struct {
	events int
}

func (f *fakeEmitter) Emit(_ context.Context, name, greeting string, count int64) error {
	f.events++
	return nil
}

func TestGreeter_Greet_ReturnsFormattedGreeting(t *testing.T) {
	g := domain.NewGreeter(&fakeCounter{}, &fakeCache{}, &fakeEmitter{})
	resp, err := g.Greet(context.Background(), "world")
	if err != nil {
		t.Fatalf("Greet: %v", err)
	}
	if resp.Greeting != "hello, world" {
		t.Errorf("Greeting = %q", resp.Greeting)
	}
	if resp.Count != 1 {
		t.Errorf("Count = %d, want 1", resp.Count)
	}
}

func TestGreeter_Greet_EmptyNameFails(t *testing.T) {
	g := domain.NewGreeter(&fakeCounter{}, &fakeCache{}, &fakeEmitter{})
	_, err := g.Greet(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error for empty name")
	}
}
```

- [ ] **Step 3: Implement `domain/greeter.go`**

Write `impl/helixgitpx/services/hello/internal/domain/greeter.go`:

```go
// Package domain holds the business logic for the hello service.
// Pure Go, no framework imports.
package domain

import (
	"context"
	"fmt"

	"github.com/helixgitpx/platform/errors"
	"google.golang.org/grpc/codes"
)

// Counter increments a per-name counter and returns the new value.
type Counter interface {
	Increment(ctx context.Context, name string) (int64, error)
}

// Cache stores the last greeting for a name.
type Cache interface {
	SetLast(ctx context.Context, name, greeting string) error
}

// Emitter publishes a hello.said event.
type Emitter interface {
	Emit(ctx context.Context, name, greeting string, count int64) error
}

// Response is the result of Greet.
type Response struct {
	Greeting string
	Count    int64
}

// Greeter is the hello business-logic aggregate.
type Greeter struct {
	counter Counter
	cache   Cache
	emitter Emitter
}

// NewGreeter constructs a Greeter.
func NewGreeter(c Counter, ca Cache, e Emitter) *Greeter {
	return &Greeter{counter: c, cache: ca, emitter: e}
}

// Greet increments the counter, caches the last greeting, emits an event, and returns the response.
func (g *Greeter) Greet(ctx context.Context, name string) (*Response, error) {
	if name == "" {
		return nil, errors.New(codes.InvalidArgument, "hello", "name is required")
	}
	count, err := g.counter.Increment(ctx, name)
	if err != nil {
		return nil, errors.New(codes.Internal, "hello", "counter increment").Wrap(err)
	}
	greeting := fmt.Sprintf("hello, %s", name)
	if err := g.cache.SetLast(ctx, name, greeting); err != nil {
		// Non-fatal — log-and-continue is the domain's choice; here we propagate to caller.
		return nil, errors.New(codes.Unavailable, "hello", "cache set last").Wrap(err)
	}
	if err := g.emitter.Emit(ctx, name, greeting, count); err != nil {
		return nil, errors.New(codes.Unavailable, "hello", "emit event").Wrap(err)
	}
	return &Response{Greeting: greeting, Count: count}, nil
}
```

- [ ] **Step 4: Implement repo adapters**

Write `impl/helixgitpx/services/hello/internal/repo/counter_pg.go`:

```go
package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CounterPG implements domain.Counter using Postgres UPSERT.
type CounterPG struct {
	Pool *pgxpool.Pool
}

func (c *CounterPG) Increment(ctx context.Context, name string) (int64, error) {
	var n int64
	err := c.Pool.QueryRow(ctx, `
        INSERT INTO hello.greetings(name, count, last_said_at)
        VALUES ($1, 1, NOW())
        ON CONFLICT (name) DO UPDATE
          SET count = hello.greetings.count + 1,
              last_said_at = NOW()
        RETURNING count`, name).Scan(&n)
	return n, err
}
```

Write `impl/helixgitpx/services/hello/internal/repo/cache_redis.go`:

```go
package repo

import (
	"context"
	"time"

	hr "github.com/helixgitpx/platform/redis"
)

// CacheRedis implements domain.Cache using HelixGitpx's Redis wrapper.
type CacheRedis struct {
	Client *hr.Client
	TTL    time.Duration
}

func (c *CacheRedis) SetLast(ctx context.Context, name, greeting string) error {
	key := c.Client.Key("last", name)
	ttl := c.TTL
	if ttl == 0 {
		ttl = 10 * time.Minute
	}
	return c.Client.Set(ctx, key, greeting, ttl).Err()
}
```

Write `impl/helixgitpx/services/hello/internal/repo/event_kafka.go`:

```go
package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/helixgitpx/platform/kafka"
)

// EventKafka implements domain.Emitter using HelixGitpx's Kafka wrapper.
type EventKafka struct {
	Producer *kafka.Producer
	Topic    string
}

type helloSaid struct {
	Name     string `json:"name"`
	Greeting string `json:"greeting"`
	Count    int64  `json:"count"`
	At       string `json:"at"`
}

func (e *EventKafka) Emit(ctx context.Context, name, greeting string, count int64) error {
	payload, err := json.Marshal(helloSaid{
		Name: name, Greeting: greeting, Count: count,
		At: time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}
	return e.Producer.Emit(ctx, []byte(name), payload, e.Topic)
}
```

- [ ] **Step 5: Implement gRPC + HTTP handlers**

Write `impl/helixgitpx/services/hello/internal/handler/grpc/server.go`:

```go
package grpc

import (
	"context"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
)

// Server implements hellopb.HelloServiceServer.
type Server struct {
	hellopb.UnimplementedHelloServiceServer
	Greeter *domain.Greeter
}

// SayHello satisfies hellopb.HelloServiceServer.
func (s *Server) SayHello(ctx context.Context, req *hellopb.SayHelloRequest) (*hellopb.SayHelloResponse, error) {
	r, err := s.Greeter.Greet(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &hellopb.SayHelloResponse{Greeting: r.Greeting, Count: r.Count}, nil
}
```

Write `impl/helixgitpx/services/hello/internal/handler/http/router.go`:

```go
package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
)

// Register adds /v1/hello to r.
func Register(r *gin.Engine, g *domain.Greeter) {
	r.GET("/v1/hello", func(c *gin.Context) {
		name := c.Query("name")
		resp, err := g.Greet(c.Request.Context(), name)
		if err != nil {
			c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(nethttp.StatusOK, gin.H{"greeting": resp.Greeting, "count": resp.Count})
	})
}
```

- [ ] **Step 6: Replace scaffolded `app.go` with real wiring**

Overwrite `impl/helixgitpx/services/hello/internal/app/app.go`:

```go
// Package app is the composition root for hello.
package app

import (
	"context"
	"errors"
	"net"
	nethttp "net/http"
	"time"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
	grpchandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/grpc"
	httphandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/http"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/repo"

	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
	"github.com/helixgitpx/platform/config"
	hgin "github.com/helixgitpx/platform/gin"
	hgrpc "github.com/helixgitpx/platform/grpc"
	"github.com/helixgitpx/platform/health"
	"github.com/helixgitpx/platform/kafka"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
	hredis "github.com/helixgitpx/platform/redis"
	"github.com/helixgitpx/platform/telemetry"
)

type cfg struct {
	HTTPAddr     string   `env:"HTTP_ADDR" default:":8001"`
	GRPCAddr     string   `env:"GRPC_ADDR" default:":9001"`
	HealthAddr   string   `env:"HEALTH_ADDR" default:":8081"`
	PostgresDSN  string   `env:"POSTGRES_DSN" required:"true"`
	RedisAddr    string   `env:"REDIS_ADDR" default:"localhost:6379"`
	KafkaBrokers []string `env:"KAFKA_BROKERS" default:"localhost:9092" split:","`
	KafkaTopic   string   `env:"KAFKA_TOPIC" default:"hello.said"`
	OTLPEndpoint string   `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Version      string   `env:"VERSION" default:"m1-dev"`
}

// Run wires dependencies and serves until ctx is done.
func Run(ctx context.Context, lg *log.Logger) error {
	var c cfg
	if err := config.Load(&c, config.Options{Prefix: "HELLO"}); err != nil {
		return err
	}

	shutdownTel, err := telemetry.Start(ctx, telemetry.Options{
		Service: "hello", Version: c.Version, OTLPEndpoint: c.OTLPEndpoint,
	})
	if err != nil {
		return err
	}
	defer func() {
		shctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTel(shctx)
	}()

	pool, err := pg.Open(ctx, pg.Options{DSN: c.PostgresDSN})
	if err != nil {
		return err
	}
	defer pool.Close()

	rc, err := hredis.Open(ctx, hredis.Options{Addr: c.RedisAddr, Namespace: "hello"})
	if err != nil {
		return err
	}
	defer rc.Close()

	prod, err := kafka.NewProducer(kafka.ProducerOptions{
		Brokers: c.KafkaBrokers, ClientID: "hello", Topic: c.KafkaTopic,
	})
	if err != nil {
		return err
	}
	defer func() {
		shctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = prod.Close(shctx)
	}()

	greeter := domain.NewGreeter(
		&repo.CounterPG{Pool: pool},
		&repo.CacheRedis{Client: rc},
		&repo.EventKafka{Producer: prod, Topic: c.KafkaTopic},
	)

	// gRPC server
	grpcSrv, err := hgrpc.NewServer(hgrpc.Options{})
	if err != nil {
		return err
	}
	hellopb.RegisterHelloServiceServer(grpcSrv, &grpchandler.Server{Greeter: greeter})

	// HTTP (REST) server
	router := hgin.NewRouter(hgin.Options{Service: "hello", Version: c.Version})
	httphandler.Register(router, greeter)

	// Health server (separate port; probes all three deps)
	hh := health.New()
	hh.Register("pg", pg.Probe(pool))
	hh.Register("redis", hredis.Probe(rc))
	hmux := nethttp.NewServeMux()
	hh.Routes(hmux)

	// Listeners
	grpcL, err := net.Listen("tcp", c.GRPCAddr)
	if err != nil {
		return err
	}
	httpSrv := &nethttp.Server{Addr: c.HTTPAddr, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	healthSrv := &nethttp.Server{Addr: c.HealthAddr, Handler: hmux, ReadHeaderTimeout: 5 * time.Second}

	errCh := make(chan error, 3)
	go func() { errCh <- grpcSrv.Serve(grpcL) }()
	go func() { errCh <- httpSrv.ListenAndServe() }()
	go func() { errCh <- healthSrv.ListenAndServe() }()

	lg.Info("hello serving",
		"grpc", c.GRPCAddr, "http", c.HTTPAddr, "health", c.HealthAddr)

	select {
	case <-ctx.Done():
		lg.Info("shutdown signalled")
	case err := <-errCh:
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			lg.Error("server exited", "err", err.Error())
		}
	}

	shctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcSrv.GracefulStop()
	_ = httpSrv.Shutdown(shctx)
	_ = healthSrv.Shutdown(shctx)
	return nil
}
```

- [ ] **Step 7: Run domain tests**

```sh
cd impl/helixgitpx/services/hello && go mod tidy && go test ./internal/domain/...
```

Expected: PASS.

- [ ] **Step 8: Build**

```sh
go build ./...
```

Expected: compiles.

- [ ] **Step 9: Commit**

```sh
git add impl/helixgitpx/services/hello
git commit -s -m "feat(services/hello): domain + handlers + wire real platform libs"
```

---

### Task 21: Hello integration test (testcontainers-backed)

**Files:**
- Create: `impl/helixgitpx/services/hello/test/integration/hello_e2e_test.go`

- [ ] **Step 1: Write the integration test**

Write `impl/helixgitpx/services/hello/test/integration/hello_e2e_test.go`:

```go
//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
	grpchandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/grpc"
	httphandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/http"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/repo"

	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
	hgin "github.com/helixgitpx/platform/gin"
	hgrpc "github.com/helixgitpx/platform/grpc"
	"github.com/helixgitpx/platform/kafka"
	"github.com/helixgitpx/platform/pg"
	hredis "github.com/helixgitpx/platform/redis"
	"github.com/helixgitpx/platform/testkit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHelloE2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	dsn := testkit.StartPostgres(t)
	redisAddr := testkit.StartRedis(t)
	kafkaBroker := testkit.StartKafka(t)

	// Migrate — minimal manual CREATE SCHEMA + table since we don't wire goose here.
	pool, err := pg.Open(ctx, pg.Options{DSN: dsn})
	if err != nil {
		t.Fatalf("pg.Open: %v", err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS hello`); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if _, err := pool.Exec(ctx, `SET search_path TO hello`); err != nil {
		t.Fatalf("search_path: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS hello.greetings (
			name TEXT PRIMARY KEY,
			count BIGINT NOT NULL DEFAULT 0,
			last_said_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	rc, err := hredis.Open(ctx, hredis.Options{Addr: redisAddr, Namespace: "hello"})
	if err != nil {
		t.Fatalf("redis.Open: %v", err)
	}

	prod, err := kafka.NewProducer(kafka.ProducerOptions{
		Brokers: []string{kafkaBroker}, ClientID: "hello-test", Topic: "hello.said",
	})
	if err != nil {
		t.Fatalf("kafka producer: %v", err)
	}
	defer prod.Close(ctx)

	greeter := domain.NewGreeter(
		&repo.CounterPG{Pool: pool},
		&repo.CacheRedis{Client: rc},
		&repo.EventKafka{Producer: prod, Topic: "hello.said"},
	)

	// Start gRPC
	grpcSrv, _ := hgrpc.NewServer(hgrpc.Options{})
	hellopb.RegisterHelloServiceServer(grpcSrv, &grpchandler.Server{Greeter: greeter})
	grpcL, _ := net.Listen("tcp", "127.0.0.1:0")
	go grpcSrv.Serve(grpcL)
	defer grpcSrv.GracefulStop()

	// Start HTTP
	router := hgin.NewRouter(hgin.Options{Service: "hello"})
	httphandler.Register(router, greeter)
	httpL, _ := net.Listen("tcp", "127.0.0.1:0")
	httpSrv := &nethttp.Server{Handler: router}
	go httpSrv.Serve(httpL)
	defer httpSrv.Close()

	// HTTP
	resp, err := nethttp.Get(fmt.Sprintf("http://%s/v1/hello?name=world", httpL.Addr().String()))
	if err != nil {
		t.Fatalf("http: %v", err)
	}
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["greeting"] != "hello, world" {
		t.Errorf("http greeting = %v", body["greeting"])
	}

	// gRPC
	conn, err := grpc.NewClient(grpcL.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	cl := hellopb.NewHelloServiceClient(conn)
	r, err := cl.SayHello(ctx, &hellopb.SayHelloRequest{Name: "world"})
	if err != nil {
		t.Fatalf("SayHello: %v", err)
	}
	if r.Greeting != "hello, world" {
		t.Errorf("grpc greeting = %q", r.Greeting)
	}
	if r.Count < 2 {
		t.Errorf("count = %d, want ≥ 2 (HTTP + gRPC calls)", r.Count)
	}
}
```

- [ ] **Step 2: Run the integration test**

```sh
cd impl/helixgitpx/services/hello && go test -tags=integration -timeout=10m ./test/integration/...
```

Expected: PASS. (Requires Podman/Docker socket reachable; may take 2–4 minutes first run.)

- [ ] **Step 3: Commit**

```sh
git add impl/helixgitpx/services/hello/test
git commit -s -m "test(services/hello): end-to-end integration with pg+redis+kafka"
```

---

### Task 22: Hello Dockerfile + compose wiring + end-to-end verification

**Files:**
- Create: `impl/helixgitpx/services/hello/deploy/Dockerfile` (replace scaffolded version if needed)
- Create: `impl/helixgitpx/services/hello/deploy/skaffold.yaml`

- [ ] **Step 1: Write the production Dockerfile**

Overwrite `impl/helixgitpx/services/hello/deploy/Dockerfile`:

```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.work go.work.sum* ./
COPY platform/ ./platform/
COPY gen/ ./gen/
COPY services/hello/ ./services/hello/
RUN cd services/hello && go mod tidy && \
    CGO_ENABLED=0 GOFLAGS=-mod=mod \
    go build -trimpath -ldflags="-s -w" -o /out/hello ./cmd/hello

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/hello /app/hello
USER nonroot
EXPOSE 8001 9001 8081
ENTRYPOINT ["/app/hello"]
LABEL org.opencontainers.image.title="hello"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
LABEL org.opencontainers.image.licenses="Apache-2.0"
```

- [ ] **Step 2: Write `skaffold.yaml`**

Write `impl/helixgitpx/services/hello/deploy/skaffold.yaml`:

```yaml
apiVersion: skaffold/v4beta11
kind: Config
metadata:
  name: hello
build:
  artifacts:
    - image: helixgitpx/hello
      context: ../../..
      docker:
        dockerfile: services/hello/deploy/Dockerfile
manifests:
  helm:
    releases:
      - name: hello
        chartPath: ../../../services/hello/deploy/helm
deploy:
  helm: {}
```

- [ ] **Step 3: Build the image via the wrapper**

```sh
impl/helixgitpx-platform/compose/bin/compose --profile all build hello
```

Expected: image builds.

- [ ] **Step 4: Bring up the full stack**

```sh
make dev
```

Expected: compose starts pg, kafka, redis, jaeger, prom, grafana, hello; health checks pass.

- [ ] **Step 5: Apply migration**

```sh
export POSTGRES_URL="postgres://hello_svc:hello_svc@localhost:5432/helixgitpx?sslmode=disable&search_path=hello"
goose -dir impl/helixgitpx/services/hello/migrations postgres "$POSTGRES_URL" up
```

Expected: `OK   20260420000001_init.sql`.

- [ ] **Step 6: Verify HTTP end-to-end**

```sh
curl -s "http://localhost:8001/v1/hello?name=world" | tee /dev/tty
```

Expected output: `{"count":1,"greeting":"hello, world"}`.

- [ ] **Step 7: Verify gRPC end-to-end**

```sh
grpcurl -plaintext -d '{"name":"world"}' localhost:9001 helixgitpx.hello.v1.HelloService/SayHello
```

Expected: `{"greeting": "hello, world", "count": "2"}`.

- [ ] **Step 8: Verify Kafka event**

```sh
impl/helixgitpx-platform/compose/bin/compose exec kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server kafka:9092 --topic hello.said --from-beginning --max-messages 1
```

Expected: a JSON line containing `"name":"world","greeting":"hello, world"`.

- [ ] **Step 9: Verify Redis cache**

```sh
impl/helixgitpx-platform/compose/bin/compose exec redis redis-cli GET "hello:last:world"
```

Expected: `"hello, world"`.

- [ ] **Step 10: Tear down and commit**

```sh
make dev-down
git add impl/helixgitpx/services/hello/deploy
git commit -s -m "feat(services/hello): Dockerfile + skaffold; verify E2E spine"
```

Phase A exit: **Spine complete.** Hello reachable via HTTP and gRPC, persisting to Postgres, caching in Redis, emitting to Kafka.

---

## Phase B — Breadth

### Task 23: Root `.github/workflows/ci-go.yml` + reusable `_setup-go.yml`

**Files:**
- Create: `.github/workflows/ci-go.yml`
- Create: `.github/workflows/_reusable/_setup-go.yml`

- [ ] **Step 1: Write `_setup-go.yml`**

```yaml
name: _setup-go
on:
  workflow_call:
    inputs:
      go-version:
        required: false
        type: string
        default: "1.23.4"

jobs:
  noop:
    runs-on: ubuntu-latest
    steps:
      - run: echo "callable only"

# Reusable composite logic — actions/setup-go + cache — is inlined at call sites
# because composite actions can't be used via workflow_call. See .github/actions/go-toolchain/
# in a later iteration if reuse grows beyond the two current workflows.
```

- [ ] **Step 2: Write `ci-go.yml`**

```yaml
name: ci-go
on:
  workflow_dispatch:
    inputs:
      service:
        description: "Service to build (or 'all')"
        type: string
        default: all
      run-mutation:
        description: "Run gremlins mutation testing"
        type: boolean
        default: false

concurrency:
  group: ci-go-${{ github.ref }}-${{ inputs.service }}
  cancel-in-progress: true

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: impl/helixgitpx
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.23.4"
          cache-dependency-path: impl/helixgitpx/go.work.sum

      - name: Install tools
        run: |
          go install github.com/bufbuild/buf/cmd/buf@v1.47.2
          go install mvdan.cc/gofumpt@v0.7.0
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

      - name: Buf lint
        run: buf lint proto/

      - name: gofumpt
        run: test -z "$(gofumpt -l .)"

      - name: golangci-lint
        run: golangci-lint run ./...

      - name: Test (race + shuffle)
        run: go test -race -shuffle=on -count=1 -coverprofile=coverage.out ./...

      - name: Coverage summary
        run: go tool cover -func=coverage.out | tail -1

      - name: Mutation (optional)
        if: ${{ inputs.run-mutation }}
        run: |
          go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.5.0
          gremlins unleash ./... || true

      - name: Build
        run: go build ./...

      - uses: actions/upload-artifact@v4
        with:
          name: coverage-go
          path: impl/helixgitpx/coverage.out
```

- [ ] **Step 3: Validate**

```sh
# Local static check
npx -y @action-validator/cli@latest -q .github/workflows/ci-go.yml
```

Expected: valid.

- [ ] **Step 4: Commit**

```sh
git add .github/workflows/ci-go.yml .github/workflows/_reusable/_setup-go.yml
git commit -s -m "ci(go): workflow_dispatch-only Go CI + reusable setup stub"
```

---

### Task 24: Remaining CI workflows (web/clients/docs/platform)

**Files:**
- Create: `.github/workflows/ci-web.yml`
- Create: `.github/workflows/ci-clients.yml`
- Create: `.github/workflows/ci-docs.yml`
- Create: `.github/workflows/ci-platform.yml`
- Create: `.github/workflows/_reusable/_setup-node.yml`
- Create: `.github/workflows/_reusable/_setup-jvm.yml`

- [ ] **Step 1: Write `ci-web.yml`**

```yaml
name: ci-web
on:
  workflow_dispatch:
    inputs:
      affected-only:
        type: boolean
        default: true
      run-e2e:
        type: boolean
        default: false
concurrency:
  group: ci-web-${{ github.ref }}
  cancel-in-progress: true
jobs:
  web:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-web } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: "20.18.1", cache: npm, cache-dependency-path: impl/helixgitpx-web/package-lock.json }
      - run: npm ci
      - name: Nx affected lint+test+build
        if: ${{ inputs.affected-only }}
        run: |
          npx nx affected -t lint,test,build --base=origin/main --head=HEAD
      - name: Nx run-many lint+test+build
        if: ${{ !inputs.affected-only }}
        run: npx nx run-many -t lint,test,build
      - name: Lighthouse (bundle size + perf)
        run: npm run lighthouse --if-present
      - name: Playwright e2e
        if: ${{ inputs.run-e2e }}
        run: npx playwright install --with-deps && npm run e2e --if-present
      - uses: actions/upload-artifact@v4
        if: always()
        with: { name: web-dist, path: impl/helixgitpx-web/dist }
```

- [ ] **Step 2: Write `ci-clients.yml`**

```yaml
name: ci-clients
on:
  workflow_dispatch:
    inputs:
      target:
        type: choice
        options: [all, android, ios, desktop, jvm]
        default: all
concurrency:
  group: ci-clients-${{ github.ref }}-${{ inputs.target }}
  cancel-in-progress: true
jobs:
  clients:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-clients } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { distribution: temurin, java-version: "21" }
      - uses: gradle/actions/setup-gradle@v4
      - name: Detekt + ktlint
        run: ./gradlew detekt ktlintCheck
      - name: Test
        run: |
          case "${{ inputs.target }}" in
            all)     ./gradlew check ;;
            android) ./gradlew :shared:testDebugUnitTest :androidApp:assembleDebug ;;
            ios)     ./gradlew :shared:iosX64Test ;;
            desktop) ./gradlew :desktopApp:test ;;
            jvm)     ./gradlew :shared:jvmTest ;;
          esac
      - name: Compose preview screenshots
        run: ./gradlew :shared:recordPaparazziDebug || true
      - uses: actions/upload-artifact@v4
        if: always()
        with: { name: clients-reports, path: impl/helixgitpx-clients/**/build/reports }
```

- [ ] **Step 3: Write `ci-docs.yml`**

```yaml
name: ci-docs
on: { workflow_dispatch: {} }
concurrency:
  group: ci-docs-${{ github.ref }}
  cancel-in-progress: true
jobs:
  docs:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-docs } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: "20.18.1", cache: npm, cache-dependency-path: impl/helixgitpx-docs/package-lock.json }
      - run: npm ci
      - name: sync-docs from spec suite
        run: node sync-docs.mjs
      - name: markdownlint
        run: npx -y markdownlint-cli2 "docs/**/*.md" "!docs/assets/**"
      - name: vale
        uses: errata-ai/vale-action@reviewdog
        with: { files: docs, fail_on_error: true }
      - name: link-check
        run: npx -y lychee --config ../../.lychee.toml docs
      - name: build
        run: npm run build
      - uses: actions/upload-artifact@v4
        with: { name: docs-build, path: impl/helixgitpx-docs/build }
```

- [ ] **Step 4: Write `ci-platform.yml`**

```yaml
name: ci-platform
on:
  workflow_dispatch:
    inputs:
      environment:
        type: choice
        options: [dev, staging, prod-eu]
        default: dev
concurrency:
  group: ci-platform-${{ github.ref }}-${{ inputs.environment }}
  cancel-in-progress: true
jobs:
  platform:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-platform } }
    steps:
      - uses: actions/checkout@v4
      - uses: azure/setup-helm@v4
        with: { version: v3.16.3 }
      - uses: hashicorp/setup-terraform@v3
        with: { terraform_version: 1.10.1 }
      - name: helm lint
        run: find helm -maxdepth 2 -name Chart.yaml -execdir helm lint . \;
      - name: kustomize build
        run: |
          for overlay in kustomize/overlays/*/; do
            kubectl kustomize "$overlay" > /dev/null
          done
        env:
          KUBECTL_VERSION: "1.31.3"
      - name: terraform fmt+validate+plan
        working-directory: impl/helixgitpx-platform/terraform
        run: |
          terraform fmt -check
          terraform init -backend=false
          terraform validate
          terraform plan -refresh=false -lock=false || true
      - uses: kyverno/action-install-cli@v0.2.0
        with: { release: v1.13.2 }
      - name: kyverno test
        run: kyverno test kyverno/
      - name: Install checkov
        run: pipx install checkov==3.2.334
      - name: checkov
        run: checkov -d . --config-file checkov/.checkov.yml
      - name: conftest
        uses: instrumenta/conftest-action@master
        with:
          files: helm
          policy: kyverno/policies
```

- [ ] **Step 5: Commit**

```sh
git add .github/workflows
git commit -s -m "ci(breadth): web/clients/docs/platform workflows (workflow_dispatch only)"
```

---

### Task 25: Security-scan, supply-chain, release, deploy workflows

**Files:**
- Create: `.github/workflows/security-scan.yml`
- Create: `.github/workflows/supply-chain.yml`
- Create: `.github/workflows/release.yml`
- Create: `.github/workflows/deploy.yml`
- Create: `.github/workflows/_reusable/_vault-creds.yml`
- Create: `.github/workflows/_reusable/_cosign-sign.yml`
- Create: `.github/workflows/_reusable/_sbom.yml`

- [ ] **Step 1: Write `security-scan.yml`**

```yaml
name: security-scan
on:
  workflow_dispatch:
    inputs:
      deep:
        type: boolean
        default: false
        description: "Enable slow scanners (SonarQube, CodeQL)"
permissions:
  contents: read
  security-events: write
jobs:
  semgrep:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: returntocorp/semgrep-action@v1
        env:
          SEMGREP_APP_TOKEN: ${{ secrets.SEMGREP_APP_TOKEN }}
        with:
          config: >-
            p/default
            p/golang
            p/typescript
            p/kotlin
            p/security-audit
  gitleaks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITLEAKS_ENABLE_COMMENTS: "false"
  gosec:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.23.4" }
      - run: |
          go install github.com/securego/gosec/v2/cmd/gosec@v2.21.4
          gosec -fmt=sarif -out=gosec.sarif ./... || true
      - uses: github/codeql-action/upload-sarif@v3
        with: { sarif_file: impl/helixgitpx/gosec.sarif }
  trivy-fs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@0.28.0
        with: { scan-type: fs, scan-ref: ., format: sarif, output: trivy-fs.sarif }
      - uses: github/codeql-action/upload-sarif@v3
        with: { sarif_file: trivy-fs.sarif }
  grype:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: anchore/scan-action@v5
        with: { path: ".", fail-build: false, severity-cutoff: high }
  snyk:
    runs-on: ubuntu-latest
    if: ${{ env.SNYK_TOKEN != '' }}
    env: { SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }} }
    steps:
      - uses: actions/checkout@v4
      - uses: snyk/actions/golang@master
        with: { args: --severity-threshold=high }
  sonarqube:
    runs-on: ubuntu-latest
    if: ${{ inputs.deep && env.SONAR_TOKEN != '' }}
    env:
      SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
      SONAR_HOST_URL: ${{ secrets.SONAR_HOST_URL }}
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: SonarSource/sonarqube-scan-action@v4
  codeql:
    runs-on: ubuntu-latest
    if: ${{ inputs.deep }}
    strategy:
      matrix: { language: [go, javascript, java] }
    steps:
      - uses: actions/checkout@v4
      - uses: github/codeql-action/init@v3
        with: { languages: ${{ matrix.language }} }
      - uses: github/codeql-action/analyze@v3
```

- [ ] **Step 2: Write `supply-chain.yml`**

```yaml
name: supply-chain
on:
  workflow_dispatch:
    inputs:
      ref:
        type: string
        default: main
permissions:
  contents: read
  packages: write
  id-token: write
  attestations: write
jobs:
  sbom:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { ref: ${{ inputs.ref }} }
      - uses: anchore/sbom-action@v0
        with: { path: ., format: cyclonedx-json, output-file: sbom.cdx.json }
      - uses: anchore/sbom-action@v0
        with: { path: ., format: spdx-json, output-file: sbom.spdx.json }
      - uses: actions/upload-artifact@v4
        with: { name: sboms, path: "sbom.*.json" }
  cosign:
    runs-on: ubuntu-latest
    if: ${{ env.COSIGN_PASSWORD != '' }}
    env: { COSIGN_PASSWORD: ${{ secrets.COSIGN_PASSWORD }} }
    steps:
      - uses: actions/checkout@v4
      - uses: sigstore/cosign-installer@v3
      - run: echo "cosign-ready"
  slsa-provenance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/attest-build-provenance@v2
        with: { subject-path: 'impl/helixgitpx/services/*/deploy/Dockerfile' }
```

- [ ] **Step 3: Write `release.yml`**

```yaml
name: release
on:
  workflow_dispatch:
    inputs:
      version: { type: string, required: true }
      dry-run: { type: boolean, default: true }
permissions:
  contents: write
  packages: write
  id-token: write
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Validate version
        run: |
          case "${{ inputs.version }}" in
            v*.*.*) echo ok ;;
            *) echo "expected semver vMAJOR.MINOR.PATCH"; exit 1 ;;
          esac
      - uses: actions/setup-go@v5
        with: { go-version: "1.23.4" }
      - name: Build images (dry-run if requested)
        run: echo "dry-run=${{ inputs.dry-run }} version=${{ inputs.version }}"
      - name: GitHub Release
        if: ${{ !inputs.dry-run }}
        uses: softprops/action-gh-release@v2
        with:
          tag_name: ${{ inputs.version }}
          generate_release_notes: true
```

- [ ] **Step 4: Write `deploy.yml`**

```yaml
name: deploy
on:
  workflow_dispatch:
    inputs:
      environment: { type: choice, options: [dev, staging, prod-eu], default: dev }
      service:     { type: string, default: all }
permissions:
  id-token: write
  contents: read
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: ${{ inputs.environment }}
    steps:
      - uses: actions/checkout@v4
      - name: Vault OIDC login (if configured)
        if: ${{ env.VAULT_ADDR != '' }}
        uses: hashicorp/vault-action@v3
        env: { VAULT_ADDR: ${{ secrets.VAULT_ADDR }} }
        with:
          url: ${{ secrets.VAULT_ADDR }}
          method: jwt
          role: gha-deploy
          secrets: |
            kv/data/deploy KUBECONFIG | KUBECONFIG ;
      - name: Skip if vault unset
        if: ${{ env.VAULT_ADDR == '' }}
        run: echo "VAULT_ADDR unset — skipping kubeconfig fetch (activation planned for M2)"
      - name: Argo CD sync (stub in M1)
        run: |
          echo "Would invoke: argocd app sync helixgitpx-${{ inputs.environment }} --service=${{ inputs.service }}"
```

- [ ] **Step 5: Write `_vault-creds.yml`**

```yaml
name: _vault-creds
on:
  workflow_call:
    inputs:
      role: { type: string, required: true }
      secrets-path: { type: string, required: true }
jobs:
  fetch:
    runs-on: ubuntu-latest
    permissions: { id-token: write, contents: read }
    steps:
      - name: Vault OIDC exchange
        if: ${{ env.VAULT_ADDR != '' }}
        uses: hashicorp/vault-action@v3
        with:
          url: ${{ secrets.VAULT_ADDR }}
          method: jwt
          role: ${{ inputs.role }}
          secrets: ${{ inputs.secrets-path }}
      - name: Skip
        if: ${{ env.VAULT_ADDR == '' }}
        run: echo "VAULT_ADDR unset — skipping"
```

- [ ] **Step 6: Commit**

```sh
git add .github/workflows
git commit -s -m "ci(security,supply-chain,release,deploy): full catalog, secret-gated"
```

---

### Task 26: Kyverno policies + Checkov config

**Files:**
- Create: `impl/helixgitpx-platform/kyverno/policies/disallow-privileged.yaml`
- Create: `impl/helixgitpx-platform/kyverno/policies/require-labels.yaml`
- Create: `impl/helixgitpx-platform/kyverno/policies/require-signed-images.yaml`
- Create: `impl/helixgitpx-platform/kyverno/policies/enforce-resource-limits.yaml`
- Create: `impl/helixgitpx-platform/kyverno/tests/kyverno-test.yaml`
- Create: `impl/helixgitpx-platform/checkov/.checkov.yml`

- [ ] **Step 1: Write `disallow-privileged.yaml`**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: disallow-privileged
  annotations:
    policies.kyverno.io/title: Disallow Privileged Containers
    policies.kyverno.io/category: Pod Security
    policies.kyverno.io/severity: high
spec:
  validationFailureAction: Enforce
  background: true
  rules:
    - name: check-privileged
      match:
        any:
          - resources: { kinds: [Pod] }
      validate:
        message: "Privileged containers are forbidden."
        pattern:
          spec:
            =(securityContext):
              =(privileged): "false"
            containers:
              - =(securityContext):
                  =(privileged): "false"
```

- [ ] **Step 2: Write `require-labels.yaml`**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-labels
spec:
  validationFailureAction: Audit
  background: true
  rules:
    - name: require-app-and-env
      match:
        any:
          - resources: { kinds: [Deployment, StatefulSet, DaemonSet] }
      validate:
        message: "Workloads must declare app.kubernetes.io/name and helixgitpx.dev/env labels."
        pattern:
          metadata:
            labels:
              app.kubernetes.io/name: "?*"
              helixgitpx.dev/env: "?*"
```

- [ ] **Step 3: Write `require-signed-images.yaml`**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-signed-images
spec:
  validationFailureAction: Audit
  webhookTimeoutSeconds: 30
  rules:
    - name: cosign-signed
      match:
        any:
          - resources: { kinds: [Pod] }
      verifyImages:
        - imageReferences: ["ghcr.io/helixgitpx/*"]
          attestors:
            - entries:
                - keyless:
                    subject: "https://github.com/helixgitpx/*"
                    issuer: "https://token.actions.githubusercontent.com"
```

- [ ] **Step 4: Write `enforce-resource-limits.yaml`**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: enforce-resource-limits
spec:
  validationFailureAction: Enforce
  background: true
  rules:
    - name: containers-must-have-limits
      match:
        any:
          - resources: { kinds: [Pod] }
      validate:
        message: "Containers must declare CPU and memory limits."
        pattern:
          spec:
            containers:
              - resources:
                  limits:
                    cpu: "?*"
                    memory: "?*"
```

- [ ] **Step 5: Write a kyverno test case**

Write `impl/helixgitpx-platform/kyverno/tests/disallow-privileged.yaml`:

```yaml
name: disallow-privileged
policies:
  - ../policies/disallow-privileged.yaml
resources:
  - resource.yaml
results:
  - policy: disallow-privileged
    rule: check-privileged
    resource: evil
    kind: Pod
    result: fail
  - policy: disallow-privileged
    rule: check-privileged
    resource: nice
    kind: Pod
    result: pass
```

Write `impl/helixgitpx-platform/kyverno/tests/resource.yaml`:

```yaml
apiVersion: v1
kind: Pod
metadata: { name: evil }
spec:
  containers:
    - name: c
      image: busybox
      securityContext: { privileged: true }
---
apiVersion: v1
kind: Pod
metadata: { name: nice }
spec:
  containers:
    - name: c
      image: busybox
      securityContext: { privileged: false }
```

- [ ] **Step 6: Write `.checkov.yml`**

Write `impl/helixgitpx-platform/checkov/.checkov.yml`:

```yaml
framework:
  - terraform
  - kubernetes
  - dockerfile
  - github_actions
  - helm
skip-check:
  # Noisy checks for a docs-first repo; revisit in M8 hardening.
  - CKV_K8S_35   # credentials in env vars (dev compose uses plaintext)
output: cli
soft-fail: false
```

- [ ] **Step 7: Validate**

```sh
mise x kyverno-cli -- kyverno test impl/helixgitpx-platform/kyverno/tests
checkov -d impl/helixgitpx-platform --config-file impl/helixgitpx-platform/checkov/.checkov.yml
```

Expected: Kyverno tests pass; Checkov reports findings but exits zero if all are in the allow list.

- [ ] **Step 8: Commit**

```sh
git add impl/helixgitpx-platform/kyverno impl/helixgitpx-platform/checkov
git commit -s -m "feat(platform): kyverno policy bundle + checkov config"
```

---

### Task 27: GitHub Actions runner controller + Kata RuntimeClass

**Files:**
- Create: `impl/helixgitpx-platform/github-actions-runner-controller/values.yaml`
- Create: `impl/helixgitpx-platform/github-actions-runner-controller/runner-scale-set.yaml`
- Create: `impl/helixgitpx-platform/github-actions-runner-controller/kata-runtimeclass.yaml`
- Create: `impl/helixgitpx-platform/github-actions-runner-controller/README.md`

- [ ] **Step 1: Write `values.yaml`**

```yaml
# Helm values for actions-runner-controller (ARC)
# Chart: oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller
githubConfigUrl: "https://github.com/helixgitpx"
githubConfigSecret: arc-github-token  # externally supplied (Vault injects in M2)
replicaCount: 1
```

- [ ] **Step 2: Write `runner-scale-set.yaml`**

```yaml
apiVersion: actions.github.com/v1alpha1
kind: AutoscalingRunnerSet
metadata:
  name: helix-kata
  namespace: arc-runners
spec:
  githubConfigUrl: "https://github.com/helixgitpx"
  githubConfigSecret: arc-github-token
  template:
    spec:
      runtimeClassName: kata-qemu
      containers:
        - name: runner
          image: ghcr.io/actions/actions-runner:latest
          resources:
            limits:   { cpu: "4", memory: "8Gi" }
            requests: { cpu: "2", memory: "4Gi" }
      serviceAccountName: arc-runner
  minRunners: 0
  maxRunners: 10
```

- [ ] **Step 3: Write `kata-runtimeclass.yaml`**

```yaml
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: kata-qemu
handler: kata-qemu
scheduling:
  nodeSelector:
    katacontainers.io/kata-runtime: "true"
```

- [ ] **Step 4: Write `README.md`**

```markdown
# actions-runner-controller — Kata

Activation is M2 (requires a running Kubernetes cluster). Artifacts here are
config-only until then. Workflows keep `runs-on: ubuntu-latest` for M1.

## Activation checklist (M2)

1. Install ARC: `helm install arc -n arc-system oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller -f values.yaml`
2. Install Kata runtimeclass: `kubectl apply -f kata-runtimeclass.yaml`
3. Label nodes with `katacontainers.io/kata-runtime=true`.
4. Create GitHub App + secret `arc-github-token` (via Vault).
5. Apply `runner-scale-set.yaml`.
6. Flip workflow `runs-on: ubuntu-latest` → `runs-on: helix-kata`.
```

- [ ] **Step 5: Validate**

```sh
kubectl-kustomize impl/helixgitpx-platform/github-actions-runner-controller 2>/dev/null || \
  kubectl apply --dry-run=client -f impl/helixgitpx-platform/github-actions-runner-controller/
```

Expected: manifests parse (the dry-run may flag missing CRDs, which is expected without ARC installed).

- [ ] **Step 6: Commit**

```sh
git add impl/helixgitpx-platform/github-actions-runner-controller
git commit -s -m "feat(platform): ARC + Kata runner config (M2 activation)"
```

---

### Task 28: Vault — Terraform module + OIDC role + policies

**Files:**
- Create: `impl/helixgitpx-platform/vault/terraform/main.tf`
- Create: `impl/helixgitpx-platform/vault/terraform/variables.tf`
- Create: `impl/helixgitpx-platform/vault/policies/deploy.hcl`
- Create: `impl/helixgitpx-platform/vault/policies/ci.hcl`
- Create: `impl/helixgitpx-platform/vault/oidc-role.json`
- Create: `impl/helixgitpx-platform/vault/README.md`

- [ ] **Step 1: Write `terraform/main.tf`**

```hcl
terraform {
  required_version = ">= 1.10"
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 4.5"
    }
  }
}

provider "vault" {}

resource "vault_jwt_auth_backend" "github_actions" {
  description = "GitHub Actions OIDC"
  path        = "github-actions"
  type        = "jwt"
  oidc_discovery_url = "https://token.actions.githubusercontent.com"
  bound_issuer       = "https://token.actions.githubusercontent.com"
}

resource "vault_jwt_auth_backend_role" "gha_deploy" {
  backend        = vault_jwt_auth_backend.github_actions.path
  role_name      = "gha-deploy"
  token_policies = ["deploy"]
  bound_audiences = ["https://github.com/helixgitpx"]
  bound_claims_type = "string"
  bound_claims = {
    repository = "helixgitpx/*"
    workflow   = "deploy"
  }
  user_claim    = "actor"
  token_ttl     = 900
  token_max_ttl = 1800
  role_type     = "jwt"
}

resource "vault_policy" "deploy" {
  name   = "deploy"
  policy = file("${path.module}/../policies/deploy.hcl")
}

resource "vault_policy" "ci" {
  name   = "ci"
  policy = file("${path.module}/../policies/ci.hcl")
}
```

- [ ] **Step 2: Write `terraform/variables.tf`**

```hcl
variable "vault_addr" {
  description = "Vault server address"
  type        = string
}
```

- [ ] **Step 3: Write `policies/deploy.hcl`**

```hcl
path "kv/data/deploy/*"     { capabilities = ["read"] }
path "kv/metadata/deploy/*" { capabilities = ["list"] }
```

- [ ] **Step 4: Write `policies/ci.hcl`**

```hcl
path "kv/data/ci/*"     { capabilities = ["read"] }
path "kv/metadata/ci/*" { capabilities = ["list"] }
```

- [ ] **Step 5: Write `oidc-role.json`**

```json
{
  "role_type": "jwt",
  "bound_audiences": ["https://github.com/helixgitpx"],
  "bound_claims_type": "string",
  "bound_claims": {
    "repository": "helixgitpx/*",
    "workflow": "deploy"
  },
  "user_claim": "actor",
  "token_policies": ["deploy"],
  "token_ttl": 900,
  "token_max_ttl": 1800
}
```

- [ ] **Step 6: Write `README.md`**

```markdown
# Vault + OIDC configuration

Activation is M2 (requires a running Vault cluster). Artifacts ship here for
GitOps from day one.

## Activation checklist (M2)

1. `export VAULT_ADDR=https://vault.internal`
2. `terraform init && terraform apply` in `terraform/`.
3. Verify: `vault read auth/github-actions/role/gha-deploy`.
4. Flip `.github/workflows/deploy.yml` from the VAULT_ADDR-gated skip into the real fetch step.

## Safety

Tokens issued by `gha-deploy` are TTL-bound (15 min default). Never widen scope
without an ADR and a security review.
```

- [ ] **Step 7: Validate Terraform**

```sh
cd impl/helixgitpx-platform/vault/terraform
terraform fmt -check
terraform init -backend=false
terraform validate
```

Expected: valid.

- [ ] **Step 8: Commit**

```sh
git add impl/helixgitpx-platform/vault
git commit -s -m "feat(platform/vault): OIDC JWT role + deploy/ci policies (M2 activation)"
```

---

### Task 29: k8s-local scripts + Tiltfile + root skaffold.yaml

**Files:**
- Create: `impl/helixgitpx-platform/k8s-local/{up.sh,down.sh,load-images.sh,k3d-config.yaml,kind-config.yaml}`
- Create: `impl/helixgitpx-platform/Tiltfile`

- [ ] **Step 1: Write `k3d-config.yaml`**

```yaml
apiVersion: k3d.io/v1alpha5
kind: Simple
metadata: { name: helix }
servers: 1
agents: 2
ports:
  - { port: "8001:80", nodeFilters: [loadbalancer] }
  - { port: "9001:443", nodeFilters: [loadbalancer] }
options:
  k3s:
    extraArgs:
      - { arg: "--disable=traefik", nodeFilters: [server:*] }
```

- [ ] **Step 2: Write `kind-config.yaml`**

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: helix
nodes:
  - role: control-plane
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraPortMappings:
      - { containerPort: 80, hostPort: 8001, protocol: TCP }
      - { containerPort: 443, hostPort: 9001, protocol: TCP }
  - role: worker
```

- [ ] **Step 3: Write `up.sh`**

```sh
#!/usr/bin/env bash
# Bring up a local Kubernetes cluster. Defaults to k3d; use KIND=1 to pick kind.
# --dry-run prints planned actions without executing.
set -euo pipefail

DRY_RUN=0
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
  esac
done

engine="${KIND:-0}"
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

run() {
  if [ "$DRY_RUN" = "1" ]; then
    printf '[dry-run] %s\n' "$*"
  else
    "$@"
  fi
}

if [ "$engine" = "1" ]; then
  run kind create cluster --config "$here/kind-config.yaml"
else
  run k3d cluster create --config "$here/k3d-config.yaml"
fi

run kubectl cluster-info
```

- [ ] **Step 4: Write `down.sh`**

```sh
#!/usr/bin/env bash
set -euo pipefail
engine="${KIND:-0}"
if [ "$engine" = "1" ]; then
  kind delete cluster --name helix
else
  k3d cluster delete helix
fi
```

- [ ] **Step 5: Write `load-images.sh`**

```sh
#!/usr/bin/env bash
# Load locally-built images into the cluster (avoids pushing to a registry).
set -euo pipefail
engine="${KIND:-0}"
images=("$@")
if [ "${#images[@]}" = 0 ]; then
  images=("helixgitpx/hello:dev")
fi
for img in "${images[@]}"; do
  if [ "$engine" = "1" ]; then
    kind load docker-image "$img" --name helix
  else
    k3d image import -c helix "$img"
  fi
done
```

- [ ] **Step 6: Write `Tiltfile`**

```starlark
# HelixGitpx Tiltfile — alternative to compose for the on-cluster path.
load('ext://restart_process', 'docker_build_with_restart')

docker_build(
    'helixgitpx/hello',
    context='../helixgitpx',
    dockerfile='../helixgitpx/services/hello/deploy/Dockerfile',
)

k8s_yaml(helm(
    '../helixgitpx/services/hello/deploy/helm',
    name='hello',
    namespace='helix',
))

k8s_resource('hello', port_forwards=['8001:8001', '9001:9001', '8081:8081'])
```

- [ ] **Step 7: chmod + shellcheck**

```sh
chmod +x impl/helixgitpx-platform/k8s-local/{up,down,load-images}.sh
shellcheck impl/helixgitpx-platform/k8s-local/*.sh
```

Expected: clean.

- [ ] **Step 8: Dry-run the up script**

```sh
impl/helixgitpx-platform/k8s-local/up.sh --dry-run
```

Expected: prints `[dry-run] k3d cluster create ...` and `[dry-run] kubectl cluster-info`.

- [ ] **Step 9: Commit**

```sh
git add impl/helixgitpx-platform/k8s-local impl/helixgitpx-platform/Tiltfile
git commit -s -m "feat(platform): k8s-local scripts (kind/k3d) + Tiltfile"
```

---

### Task 30: Nx web workspace (Angular 19) scaffold

**Files:**
- Create: `impl/helixgitpx-web/package.json`
- Create: `impl/helixgitpx-web/nx.json`
- Create: `impl/helixgitpx-web/tsconfig.base.json`
- Create: `impl/helixgitpx-web/README.md`
- Create: `impl/helixgitpx-web/.github/workflows/ci.yml`
- Create: `impl/helixgitpx-web/apps/web/` (minimal Angular shell — see steps)
- Create: `impl/helixgitpx-web/libs/proto/` (placeholder — real codegen from Task 17)

- [ ] **Step 1: Write `package.json`**

```json
{
  "name": "helixgitpx-web",
  "private": true,
  "version": "0.0.0",
  "scripts": {
    "postinstall": "nx sync",
    "start": "nx serve web",
    "build": "nx build web",
    "test": "nx test",
    "lint": "nx lint"
  },
  "devDependencies": {
    "@angular/cli": "19.0.6",
    "@angular/core": "19.0.5",
    "@angular/common": "19.0.5",
    "@angular/compiler": "19.0.5",
    "@angular/compiler-cli": "19.0.5",
    "@angular/platform-browser": "19.0.5",
    "@angular/platform-browser-dynamic": "19.0.5",
    "@angular/router": "19.0.5",
    "@bufbuild/protoc-gen-es": "2.2.3",
    "@connectrpc/protoc-gen-connect-es": "1.6.1",
    "@connectrpc/connect-web": "1.6.1",
    "@nx/angular": "20.3.0",
    "@nx/eslint": "20.3.0",
    "@nx/jest": "20.3.0",
    "@nx/workspace": "20.3.0",
    "nx": "20.3.0",
    "typescript": "5.6.3",
    "@playwright/test": "1.49.1",
    "@lhci/cli": "0.14.0"
  }
}
```

- [ ] **Step 2: Write `nx.json`**

```json
{
  "$schema": "./node_modules/nx/schemas/nx-schema.json",
  "namedInputs": {
    "default": ["{projectRoot}/**/*", "sharedGlobals"],
    "production": ["default", "!{projectRoot}/**/*.spec.ts"],
    "sharedGlobals": []
  },
  "targetDefaults": {
    "build":  { "cache": true, "dependsOn": ["^build"] },
    "lint":   { "cache": true },
    "test":   { "cache": true }
  },
  "defaultBase": "origin/main"
}
```

- [ ] **Step 3: Write `tsconfig.base.json`**

```json
{
  "compileOnSave": false,
  "compilerOptions": {
    "baseUrl": ".",
    "target": "ES2022",
    "module": "ES2022",
    "moduleResolution": "node",
    "strict": true,
    "noImplicitOverride": true,
    "noFallthroughCasesInSwitch": true,
    "esModuleInterop": true,
    "paths": {
      "@helixgitpx/proto": ["libs/proto/src/index.ts"]
    }
  }
}
```

- [ ] **Step 4: Minimal Angular app**

```sh
mkdir -p impl/helixgitpx-web/apps/web/src/app
```

Write `impl/helixgitpx-web/apps/web/project.json`:

```json
{
  "name": "web",
  "$schema": "../../node_modules/nx/schemas/project-schema.json",
  "projectType": "application",
  "sourceRoot": "apps/web/src",
  "prefix": "hx",
  "targets": {
    "build": {
      "executor": "@angular/build:application",
      "options": {
        "outputPath": "dist/web",
        "browser": "apps/web/src/main.ts",
        "index": "apps/web/src/index.html",
        "tsConfig": "apps/web/tsconfig.app.json",
        "polyfills": ["zone.js"]
      }
    },
    "serve": {
      "executor": "@angular/build:dev-server",
      "options": { "buildTarget": "web:build" }
    },
    "test": { "executor": "@nx/jest:jest", "options": { "jestConfig": "apps/web/jest.config.ts" } },
    "lint": { "executor": "@nx/eslint:lint" }
  }
}
```

Write minimal `apps/web/src/index.html`, `main.ts`, `app/app.component.ts`, `app/app.config.ts`, `tsconfig.app.json`, `jest.config.ts` — minimal Angular 19 standalone boilerplate. (Engineer: use `npx ng new` idioms and strip to the bare working shell. Keep dependencies aligned to `package.json`.)

- [ ] **Step 5: Write `README.md`**

```markdown
# helixgitpx-web (Nx workspace, Angular 19)

## Quick start

```sh
npm install
npx nx serve web
```

The root `make gen` regenerates `libs/proto/` from `impl/helixgitpx/proto/`.

## Layout

- `apps/web/` — the Angular shell (M6 expands with real screens).
- `libs/proto/` — protobuf + Connect-ES codegen (committed).
- `libs/ui/` — design system (M6).
- `libs/data/` — data-access layers (M6).
```

- [ ] **Step 6: Write `.github/workflows/ci.yml`**

```yaml
name: ci (web subdir)
on: { workflow_dispatch: {} }
jobs:
  ci:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-web } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: "20.18.1", cache: npm, cache-dependency-path: impl/helixgitpx-web/package-lock.json }
      - run: npm ci
      - run: npx nx run-many -t lint,test,build
```

- [ ] **Step 7: Install and smoke-build**

```sh
cd impl/helixgitpx-web && npm install --no-audit
npx nx build web --skip-nx-cache
```

Expected: `dist/web/` appears; build exits 0.

- [ ] **Step 8: Commit**

```sh
git add impl/helixgitpx-web
git commit -s -m "feat(web): Nx workspace + Angular 19 shell + ci workflow"
```

---

### Task 31: KMP clients scaffold — buildSrc + shared module

**Files:**
- Create: `impl/helixgitpx-clients/settings.gradle.kts`
- Create: `impl/helixgitpx-clients/build.gradle.kts`
- Create: `impl/helixgitpx-clients/gradle.properties`
- Create: `impl/helixgitpx-clients/buildSrc/build.gradle.kts`
- Create: `impl/helixgitpx-clients/buildSrc/src/main/kotlin/helix.convention.gradle.kts`
- Create: `impl/helixgitpx-clients/shared/build.gradle.kts`
- Create: `impl/helixgitpx-clients/shared/src/commonMain/kotlin/dev/helixgitpx/Platform.kt`
- Create: `impl/helixgitpx-clients/shared/src/commonTest/kotlin/dev/helixgitpx/PlatformTest.kt`
- Create: `impl/helixgitpx-clients/README.md`
- Create: `impl/helixgitpx-clients/.github/workflows/ci.yml`

- [ ] **Step 1: Write `settings.gradle.kts`**

```kotlin
pluginManagement {
    repositories { gradlePluginPortal(); google(); mavenCentral() }
}
dependencyResolutionManagement {
    repositories { google(); mavenCentral() }
}
rootProject.name = "helixgitpx-clients"
include(":shared")
```

- [ ] **Step 2: Write `build.gradle.kts`**

```kotlin
plugins {
    kotlin("multiplatform") version "2.1.0" apply false
    id("io.gitlab.arturbosch.detekt") version "1.23.7" apply false
    id("org.jlleitschuh.gradle.ktlint") version "12.1.2" apply false
}
allprojects {
    group = "dev.helixgitpx"
    version = "0.1.0"
}
```

- [ ] **Step 3: Write `gradle.properties`**

```
org.gradle.jvmargs=-Xmx4g -XX:MaxMetaspaceSize=1g
org.gradle.parallel=true
org.gradle.caching=true
kotlin.code.style=official
kotlin.mpp.stability.nowarn=true
```

- [ ] **Step 4: Write `buildSrc/build.gradle.kts`**

```kotlin
plugins { `kotlin-dsl` }
repositories { gradlePluginPortal(); mavenCentral() }
```

- [ ] **Step 5: Write `buildSrc/src/main/kotlin/helix.convention.gradle.kts`**

```kotlin
plugins {
    id("org.jetbrains.kotlin.multiplatform")
    id("io.gitlab.arturbosch.detekt")
    id("org.jlleitschuh.gradle.ktlint")
}

kotlin {
    jvmToolchain(21)
    jvm()
    iosX64()
    iosArm64()
    iosSimulatorArm64()
    linuxX64()
    androidTarget {
        // Real android config lands in M6; this keeps commonMain buildable.
    }
}

detekt {
    buildUponDefaultConfig = true
    allRules = false
}

tasks.withType<org.jlleitschuh.gradle.ktlint.tasks.KtLintCheckTask>().configureEach {
    exclude { it.file.path.contains("/gen/") }
}
```

- [ ] **Step 6: Write `shared/build.gradle.kts`**

```kotlin
plugins { id("helix.convention") }

kotlin {
    sourceSets {
        val commonMain by getting {
            dependencies {
                implementation("io.ktor:ktor-client-core:3.0.2")
                implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.9.0")
                implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.7.3")
            }
        }
        val commonTest by getting {
            dependencies { implementation(kotlin("test")) }
        }
    }
}
```

- [ ] **Step 7: Write `shared/src/commonMain/kotlin/dev/helixgitpx/Platform.kt`**

```kotlin
package dev.helixgitpx

expect val platformName: String

object HelixGitpx {
    fun greeting(name: String): String = "hello, $name"
}
```

- [ ] **Step 8: Write `shared/src/commonTest/kotlin/dev/helixgitpx/PlatformTest.kt`**

```kotlin
package dev.helixgitpx

import kotlin.test.Test
import kotlin.test.assertEquals

class GreetingTest {
    @Test
    fun greetingWorld() {
        assertEquals("hello, world", HelixGitpx.greeting("world"))
    }
}
```

Add `expect val` actuals per target under `shared/src/<target>Main/kotlin/dev/helixgitpx/Platform.<target>.kt`:

```kotlin
package dev.helixgitpx
actual val platformName: String = "jvm"  // repeat with "androidUnit", "iosX64", etc.
```

- [ ] **Step 9: Write `README.md`**

```markdown
# helixgitpx-clients (KMP + Compose Multiplatform)

## Quick start

```sh
./gradlew check
```

## Layout

- `shared/` — KMP shared module (domain, network, store). M1 ships a stub.
- `buildSrc/` — convention plugins so per-target modules stay DRY.
- `androidApp/`, `iosApp/`, `desktopApp/` — platform shells added in M6.

iOS Kotlin targets compile on Linux (iOSX64 + iOSArm64 klibs); final iOS linking requires macOS (M6 CI).
```

- [ ] **Step 10: Write `.github/workflows/ci.yml`**

```yaml
name: ci (clients subdir)
on: { workflow_dispatch: {} }
jobs:
  ci:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-clients } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with: { distribution: temurin, java-version: "21" }
      - uses: gradle/actions/setup-gradle@v4
      - run: ./gradlew detekt ktlintCheck check
```

- [ ] **Step 11: Smoke-build**

```sh
cd impl/helixgitpx-clients
./gradlew :shared:jvmTest
```

Expected: PASS.

- [ ] **Step 12: Commit**

```sh
git add impl/helixgitpx-clients
git commit -s -m "feat(clients): KMP scaffold + buildSrc convention plugin + shared stub"
```

---

### Task 32: Docusaurus site + sync-docs script

**Files:**
- Create: `impl/helixgitpx-docs/package.json`
- Create: `impl/helixgitpx-docs/tsconfig.json`
- Create: `impl/helixgitpx-docs/docusaurus.config.ts`
- Create: `impl/helixgitpx-docs/sidebars.ts`
- Create: `impl/helixgitpx-docs/sync-docs.mjs`
- Create: `impl/helixgitpx-docs/src/pages/index.tsx`
- Create: `impl/helixgitpx-docs/src/css/custom.css`
- Create: `impl/helixgitpx-docs/.github/workflows/ci.yml`
- Create: `impl/helixgitpx-docs/README.md`

- [ ] **Step 1: Write `package.json`**

```json
{
  "name": "helixgitpx-docs",
  "version": "0.0.0",
  "private": true,
  "scripts": {
    "sync": "node sync-docs.mjs",
    "start": "npm run sync && docusaurus start",
    "build": "npm run sync && docusaurus build",
    "serve": "docusaurus serve"
  },
  "dependencies": {
    "@docusaurus/core": "3.7.0",
    "@docusaurus/preset-classic": "3.7.0",
    "@mdx-js/react": "3.1.0",
    "clsx": "2.1.1",
    "prism-react-renderer": "2.4.1",
    "react": "18.3.1",
    "react-dom": "18.3.1",
    "redocusaurus": "2.2.0"
  },
  "devDependencies": {
    "@docusaurus/module-type-aliases": "3.7.0",
    "@docusaurus/tsconfig": "3.7.0",
    "@docusaurus/types": "3.7.0",
    "@docusaurus/plugin-content-docs": "3.7.0",
    "typescript": "5.6.3"
  }
}
```

- [ ] **Step 2: Write `sync-docs.mjs`**

```js
#!/usr/bin/env node
// Copy (or symlink under POSIX) the spec tree into Docusaurus-friendly docs/.
// Rewrites intra-spec links to Docusaurus routes.

import fs from "node:fs/promises";
import path from "node:path";
import { existsSync } from "node:fs";

const SRC = path.resolve("../../docs/specifications/main/main_implementation_material/HelixGitpx");
const DST = path.resolve("./docs");

async function copyTree(src, dst) {
  const entries = await fs.readdir(src, { withFileTypes: true });
  await fs.mkdir(dst, { recursive: true });
  for (const entry of entries) {
    const s = path.join(src, entry.name);
    const d = path.join(dst, entry.name);
    if (entry.isDirectory()) {
      await copyTree(s, d);
    } else if (entry.name.endsWith(".md")) {
      const raw = await fs.readFile(s, "utf8");
      const rewritten = raw
        // rewrite relative links to other spec markdowns → Docusaurus-relative
        .replace(/]\(\.\.\/(\d\d-[a-z-]+)\//g, "](../../$1/")
        // strip absolute paths back into the spec
        .replace(/]\(docs\/specifications\/[^)]+\/HelixGitpx\//g, "](/");
      await fs.writeFile(d, rewritten);
    } else {
      await fs.copyFile(s, d);
    }
  }
}

if (!existsSync(SRC)) {
  console.error(`sync-docs: source not found at ${SRC}`);
  process.exit(1);
}
await fs.rm(DST, { recursive: true, force: true });
await copyTree(SRC, DST);
console.log(`sync-docs: copied spec tree → ${DST}`);
```

- [ ] **Step 3: Write `docusaurus.config.ts`**

```ts
import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";

const config: Config = {
  title: "HelixGitpx",
  tagline: "Helix Git Proxy eXtended — federated, privacy-preserving Git proxy",
  favicon: "img/favicon.ico",
  url: "https://docs.helixgitpx.dev",
  baseUrl: "/",
  organizationName: "helixgitpx",
  projectName: "helixgitpx",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  presets: [
    ["classic", {
      docs: {
        sidebarPath: "./sidebars.ts",
        editUrl: "https://github.com/helixgitpx/helixgitpx/edit/main/docs/specifications/main/main_implementation_material/HelixGitpx/",
      },
      theme: { customCss: "./src/css/custom.css" },
    }],
  ],
  themeConfig: {
    navbar: {
      title: "HelixGitpx",
      items: [
        { to: "/docs/intro", label: "Docs", position: "left" },
        { to: "/docs/roadmap/17-milestones", label: "Roadmap", position: "left" },
        { href: "https://github.com/helixgitpx/helixgitpx", label: "GitHub", position: "right" },
      ],
    },
    prism: { theme: prismThemes.github, darkTheme: prismThemes.dracula },
  },
};
export default config;
```

- [ ] **Step 4: Write `sidebars.ts`**

```ts
import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";
const sidebars: SidebarsConfig = {
  default: [
    "intro",
    { type: "category", label: "Core",            items: [{ type: "autogenerated", dirName: "00-core" }] },
    { type: "category", label: "Architecture",    items: [{ type: "autogenerated", dirName: "01-architecture" }] },
    { type: "category", label: "Services",        items: [{ type: "autogenerated", dirName: "02-services" }] },
    { type: "category", label: "Data",            items: [{ type: "autogenerated", dirName: "03-data" }] },
    { type: "category", label: "APIs",            items: [{ type: "autogenerated", dirName: "04-apis" }] },
    { type: "category", label: "Frontend",        items: [{ type: "autogenerated", dirName: "05-frontend" }] },
    { type: "category", label: "Mobile",          items: [{ type: "autogenerated", dirName: "06-mobile" }] },
    { type: "category", label: "AI",              items: [{ type: "autogenerated", dirName: "07-ai" }] },
    { type: "category", label: "Security",        items: [{ type: "autogenerated", dirName: "08-security" }] },
    { type: "category", label: "Observability",   items: [{ type: "autogenerated", dirName: "09-observability" }] },
    { type: "category", label: "Testing",         items: [{ type: "autogenerated", dirName: "10-testing" }] },
    { type: "category", label: "DevOps",          items: [{ type: "autogenerated", dirName: "11-devops" }] },
    { type: "category", label: "Operations",      items: [{ type: "autogenerated", dirName: "12-operations" }] },
    { type: "category", label: "Roadmap",         items: [{ type: "autogenerated", dirName: "13-roadmap" }] },
    { type: "category", label: "Diagrams",        items: [{ type: "autogenerated", dirName: "14-diagrams" }] },
    { type: "category", label: "Reference",       items: [{ type: "autogenerated", dirName: "15-reference" }] },
  ],
};
export default sidebars;
```

- [ ] **Step 5: Write `src/pages/index.tsx` and `src/css/custom.css` (minimal)**

```tsx
// src/pages/index.tsx
import Layout from "@theme/Layout";
export default function Home() {
  return (
    <Layout title="HelixGitpx" description="Federated Git proxy">
      <main style={{ padding: "4rem 2rem" }}>
        <h1>HelixGitpx</h1>
        <p>Helix Git Proxy eXtended — federated, privacy-preserving Git proxy.</p>
        <p><a href="/docs/intro">Read the specs →</a></p>
      </main>
    </Layout>
  );
}
```

```css
/* src/css/custom.css */
:root { --ifm-color-primary: #0b7a75; }
```

- [ ] **Step 6: Write `README.md` + `ci.yml` (same shape as prior tasks)**

Write `impl/helixgitpx-docs/README.md`:

```markdown
# helixgitpx-docs (Docusaurus)

```sh
npm install
npm run start   # sync then serve on :3001 (via make docs)
npm run build
```

The `docs/` folder is regenerated by `sync-docs.mjs` from
`../../docs/specifications/.../HelixGitpx/`. Never edit `docs/` directly — edit the spec suite.
```

Write `impl/helixgitpx-docs/.github/workflows/ci.yml`:

```yaml
name: ci (docs subdir)
on: { workflow_dispatch: {} }
jobs:
  ci:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: impl/helixgitpx-docs } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: "20.18.1", cache: npm, cache-dependency-path: impl/helixgitpx-docs/package-lock.json }
      - run: npm ci
      - run: npm run build
```

- [ ] **Step 7: Build**

```sh
cd impl/helixgitpx-docs && npm install --no-audit && npm run build
```

Expected: `build/` directory created.

- [ ] **Step 8: Commit**

```sh
git add impl/helixgitpx-docs
git commit -s -m "feat(docs): docusaurus scaffold + sync-docs from spec tree"
```

---

### Task 33: Seed ADRs 0001–0005

**Files:**
- Create: `docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0001-github-actions-workflow-dispatch-only.md`
- Create: `.../adr/0002-portable-container-runtime.md`
- Create: `.../adr/0003-single-git-history-impl-subdirs.md`
- Create: `.../adr/0004-mise-toolchain.md`
- Create: `.../adr/0005-spine-first-sequencing.md`

- [ ] **Step 1: Read the template**

```sh
cat docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr-template.md
```

(Use its headings verbatim. The five ADRs below use the same headings.)

- [ ] **Step 2: Write ADR-0001**

`docs/.../adr/0001-github-actions-workflow-dispatch-only.md`:

```markdown
# ADR-0001 — GitHub Actions with `workflow_dispatch`-only triggers

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић
- **Supersedes:** —

## Context
Project is federated across GitHub, GitLab, GitFlic, GitVerse. Every CI host has its own config format; automatic push/PR triggers on multiple hosts would multiply maintenance and produce noisy cross-host reruns.

## Decision
GitHub Actions is the sole CI host. Every workflow (`.github/workflows/*.yml`) declares `on: workflow_dispatch` as its sole trigger. No `push:`, `pull_request:`, or `schedule:` triggers are permitted. `workflow_call` for reusable workflows is allowed because it does not fire automatically.

## Consequences
- Engineers invoke CI manually from the Actions tab (or `gh workflow run`).
- Contributors cannot accidentally burn runner minutes with unreviewed pushes.
- Scheduled jobs (e.g., nightly scans) must be scheduled externally (Render, GH App, or Claude `/schedule`) and call `workflow_dispatch` via the API.

## Alternatives considered
- Push+PR triggers — rejected as explicitly out of scope by the project owner.
- Multi-host CI (GitHub + GitLab + GitFlic + GitVerse) — rejected for maintenance cost.

## Links
- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-2
```

- [ ] **Step 3: Write ADR-0002**

```markdown
# ADR-0002 — Portable container runtime (docker OR podman) via wrapper

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context
Primary dev machines have `podman`; many ecosystem templates assume `docker`. Hardcoding either breaks one class of machines.

## Decision
All local-dev tooling (Makefiles, scripts) invokes `impl/helixgitpx-platform/compose/bin/compose`, which detects the available runtime at invocation time (preferring `docker compose` → `podman compose` → `podman-compose`). Scripts never call `docker` or `podman` directly.

## Consequences
- Works on both docker and podman hosts with no per-machine setup.
- Developers still see the actual runtime in logs via `compose` output.
- Features that differ between runtimes (selinux labels, some network modes) are constrained to the common subset.

## Links
- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-3, §8
```

- [ ] **Step 4: Write ADRs 0003–0005**

Use the same headings template. Subject lines:

- **ADR-0003** — *Single git history with 5 logical subdirs under `impl/`*. Captures C-1. Decision: implementation code lives at `impl/helixgitpx/`, `impl/helixgitpx-web/`, `impl/helixgitpx-clients/`, `impl/helixgitpx-platform/`, `impl/helixgitpx-docs/` with a shared git history. Each is splittable via `git subtree`/`git filter-repo` when team size requires it.

- **ADR-0004** — *`mise` for toolchain management*. Captures C-5. Decision: `mise.toml` at repo root pins every language runtime and CLI the project uses. Chosen over `devbox` (Nix-based, hermetic) because `mise` is faster, simpler, and aligns with tooling engineers already have; on-prem hermetic builds will be re-evaluated in M8.

- **ADR-0005** — *Spine-first M1 sequencing with a completion matrix*. Captures C-6. Decision: M1 builds a thin vertical slice (monorepo → platform libs → hello service → compose → one CI workflow) before fanning out. A completion matrix in the design spec ensures nothing is skipped; the matrix maps each of the 18 roadmap items to a concrete artifact with a status gate.

Write each at `docs/.../adr/000N-<slug>.md` following the ADR-0001/0002 shape.

- [ ] **Step 5: Commit**

```sh
git add docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr
git commit -s -m "docs(adr): seed ADRs 0001-0005 from M1 brainstorming constraints"
```

---

### Task 34: Runbook template + conformance lint script

**Files:**
- Create: `docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_TEMPLATE.md`
- Create: `docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_lint.sh`

- [ ] **Step 1: Read one existing runbook to anchor the headings**

```sh
head -60 docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/RB-010-postgres-primary-down.md
```

Note headings used.

- [ ] **Step 2: Write `_TEMPLATE.md`**

```markdown
# RB-NNN-<slug> — <Title>

- **Severity:** <SEV1|SEV2|SEV3>
- **SLO impact:** <which SLO and how>
- **Owner:** <team/role>
- **Last tested:** <YYYY-MM-DD>

## Summary
One-paragraph description of the incident class this runbook addresses.

## Symptoms
- Bullet list of user-visible or metric-visible symptoms.

## Diagnosis
Step-by-step commands, dashboards, and queries that confirm the class.

## Mitigation
1. Immediate actions to stop the bleed.
2. Rollback/failover steps.
3. Stabilisation.

## Post-mortem checklist
- [ ] Incident timeline documented
- [ ] RCA identified
- [ ] Action items filed with owners + due dates
- [ ] Runbook updated with learnings

## Related alerts
- `<Alertmanager rule name>`
- `<Grafana dashboard link>`

## See also
- Related runbooks, ADRs, dashboards.
```

- [ ] **Step 3: Write `_lint.sh`**

```sh
#!/usr/bin/env bash
# Lint runbooks: every RB-*.md must contain the required headings from _TEMPLATE.md.
set -euo pipefail

here="docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks"
required=("## Summary" "## Symptoms" "## Diagnosis" "## Mitigation" "## Post-mortem checklist" "## Related alerts")

rc=0
shopt -s nullglob
for rb in "$here"/RB-*.md; do
    for h in "${required[@]}"; do
        if ! grep -Fq "$h" "$rb"; then
            echo "$rb: missing heading '$h'" >&2
            rc=1
        fi
    done
done

if [ "$rc" -eq 0 ]; then
    echo "runbook-lint: all runbooks conform"
fi
exit "$rc"
```

- [ ] **Step 4: chmod + run**

```sh
chmod +x docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_lint.sh
make runbook-lint
```

Expected: all nine existing runbooks pass (they already follow this shape). Fix any gaps if the lint surfaces them.

- [ ] **Step 5: Commit**

```sh
git add docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks
git commit -s -m "docs(runbooks): formalize template + conformance lint"
```

---

### Task 35: CODEOWNERS + branch protection config + hello helm chart

**Files:**
- Create: `.github/CODEOWNERS`
- Create: `.github/branch-protection.json`
- Create: `impl/helixgitpx/services/hello/deploy/helm/Chart.yaml`
- Create: `impl/helixgitpx/services/hello/deploy/helm/values.yaml`
- Create: `impl/helixgitpx/services/hello/deploy/helm/templates/{deployment,service,configmap}.yaml`

- [ ] **Step 1: Write `CODEOWNERS`**

```
# Solo operation — every path requires the maintainer. Expand as team grows.
*               @milosvasic

# Sensitive areas (explicit for clarity)
/docs/specifications/main/                 @milosvasic
/.github/workflows/                        @milosvasic
/impl/helixgitpx-platform/vault/           @milosvasic
/impl/helixgitpx-platform/kyverno/         @milosvasic
```

- [ ] **Step 2: Write `branch-protection.json` (config artifact — not enforced under solo operation, see SOLO-NOTES.md)**

```json
{
  "branch": "main",
  "required_status_checks": {
    "strict": true,
    "contexts": ["ci-go", "ci-web", "ci-clients", "ci-docs", "ci-platform"]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 2,
    "require_code_owner_reviews": true,
    "dismiss_stale_reviews": true
  },
  "restrictions": null,
  "required_linear_history": true,
  "required_signatures": true,
  "allow_force_pushes": false,
  "allow_deletions": false
}
```

- [ ] **Step 3: Write the helm chart**

`Chart.yaml`:

```yaml
apiVersion: v2
name: hello
description: HelixGitpx hello service
type: application
version: 0.1.0
appVersion: "0.1.0"
```

`values.yaml`:

```yaml
image:
  repository: helixgitpx/hello
  tag: dev
  pullPolicy: IfNotPresent
env:
  HELLO_HTTP_ADDR: ":8001"
  HELLO_GRPC_ADDR: ":9001"
  HELLO_HEALTH_ADDR: ":8081"
resources:
  limits:   { cpu: "1", memory: "256Mi" }
  requests: { cpu: "100m", memory: "128Mi" }
```

`templates/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}
  labels:
    app.kubernetes.io/name: hello
    helixgitpx.dev/env: {{ .Values.global.env | default "dev" }}
spec:
  replicas: 1
  selector: { matchLabels: { app.kubernetes.io/name: hello } }
  template:
    metadata:
      labels:
        app.kubernetes.io/name: hello
        helixgitpx.dev/env: {{ .Values.global.env | default "dev" }}
    spec:
      containers:
        - name: hello
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - { name: http,    containerPort: 8001 }
            - { name: grpc,    containerPort: 9001 }
            - { name: health,  containerPort: 8081 }
          envFrom: [{ configMapRef: { name: "{{ .Release.Name }}-env" } }]
          readinessProbe: { httpGet: { path: /readyz, port: health } }
          livenessProbe:  { httpGet: { path: /livez,  port: health } }
          resources: {{- toYaml .Values.resources | nindent 12 }}
```

`templates/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata: { name: {{ .Release.Name }} }
spec:
  selector: { app.kubernetes.io/name: hello }
  ports:
    - { name: http,   port: 8001, targetPort: http }
    - { name: grpc,   port: 9001, targetPort: grpc }
    - { name: health, port: 8081, targetPort: health }
```

`templates/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata: { name: "{{ .Release.Name }}-env" }
data:
  {{- range $k, $v := .Values.env }}
  {{ $k }}: {{ $v | quote }}
  {{- end }}
```

- [ ] **Step 4: helm lint**

```sh
helm lint impl/helixgitpx/services/hello/deploy/helm
```

Expected: `0 chart(s) failed`.

- [ ] **Step 5: Commit**

```sh
git add .github/CODEOWNERS .github/branch-protection.json \
        impl/helixgitpx/services/hello/deploy/helm
git commit -s -m "feat(hello): helm chart + repo CODEOWNERS + branch-protection config"
```

---

### Task 36: Completion matrix verification script

**Files:**
- Create: `scripts/verify-m1.sh`

- [ ] **Step 1: Write the verifier**

```sh
mkdir -p scripts
```

Write `scripts/verify-m1.sh`:

```sh
#!/usr/bin/env bash
# Walk the M1 completion matrix and check every row's status gate.
# Exit 0 iff every row passes; print a per-row summary.
set -u

pass=0
fail=0
report() {
    local name="$1" result="$2" msg="${3-}"
    if [ "$result" = "ok" ]; then
        printf '  [ ok ] %s\n' "$name"
        pass=$((pass + 1))
    else
        printf '  [FAIL] %s%s\n' "$name" "${msg:+ — $msg}"
        fail=$((fail + 1))
    fi
}

check() {
    local name="$1"; shift
    if "$@" >/dev/null 2>&1; then
        report "$name" ok
    else
        report "$name" fail
    fi
}

echo "== M1 Completion Matrix =="

check "1  Go monorepo builds"                    bash -c 'cd impl/helixgitpx && go build ./...'
check "2  Web nx build"                           bash -c 'cd impl/helixgitpx-web && npx nx build web --skip-nx-cache'
check "2  Clients gradle check"                   bash -c 'cd impl/helixgitpx-clients && ./gradlew :shared:jvmTest'
check "3  Platform libs tested"                   bash -c 'cd impl/helixgitpx/platform && go test ./...'
check "4  Nx workspace config"                    test -f impl/helixgitpx-web/nx.json
check "4  Gradle convention plugin"               test -f impl/helixgitpx-clients/buildSrc/src/main/kotlin/helix.convention.gradle.kts
check "5  Scaffold tool runs"                     bash -c 'cd impl/helixgitpx && go run ./tools/scaffold --dry-run --name x --proto x.v1'
check "6  Buf lint/build"                         bash -c 'cd impl/helixgitpx && buf lint proto/ && buf build proto/'
check "6  BSR module name in buf.yaml"            grep -q 'buf.build/helixgitpx/core' impl/helixgitpx/proto/buf.yaml
check "7  All CI workflows are workflow_dispatch" bash -c 'for f in .github/workflows/*.yml; do grep -qE "^on:\\s*\\{?\\s*workflow_dispatch" "$f" || grep -qE "^on:\\s*$" "$f" || exit 1; done'
check "8  Kyverno policies parse"                 kyverno test impl/helixgitpx-platform/kyverno/tests
check "8  Checkov config exists"                  test -f impl/helixgitpx-platform/checkov/.checkov.yml
check "9  ARC + Kata manifests present"           test -f impl/helixgitpx-platform/github-actions-runner-controller/runner-scale-set.yaml
check "10 Vault terraform valid"                  bash -c 'cd impl/helixgitpx-platform/vault/terraform && terraform fmt -check && terraform init -backend=false && terraform validate'
check "11 mise.toml exists"                       test -f mise.toml
check "12 Tiltfile + skaffold.yaml"               test -f impl/helixgitpx-platform/Tiltfile
check "13 k8s-local scripts pass shellcheck"      shellcheck impl/helixgitpx-platform/k8s-local/*.sh
check "14 Compose config valid"                   impl/helixgitpx-platform/compose/bin/compose --profile core config --quiet
check "15 Docusaurus builds"                      bash -c 'cd impl/helixgitpx-docs && npm run build'
check "16 ADRs 0001-0005 present"                 bash -c 'for n in 0001 0002 0003 0004 0005; do ls docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/${n}-*.md; done'
check "17 Runbook template + lint"                bash -c 'test -f docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_TEMPLATE.md && bash docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_lint.sh'
check "18 Docs ci workflow present"               test -f .github/workflows/ci-docs.yml

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
```

- [ ] **Step 2: chmod + run**

```sh
chmod +x scripts/verify-m1.sh
bash scripts/verify-m1.sh
```

Expected: every row reports `[ ok ]`, final line `PASS: 21   FAIL: 0`.

- [ ] **Step 3: Commit**

```sh
git add scripts/verify-m1.sh
git commit -s -m "chore(m1): completion-matrix verification script"
```

---

## M1 Exit

All tasks complete ⇒ every row of the completion matrix passes ⇒ M1 Foundation is done. Tag the milestone:

```sh
git tag -a m1-foundation -m "M1 Foundation complete: spine + breadth per 2026-04-20 plan"
```

Proceed to brainstorming M2 Core Data Plane using the same skill loop (brainstorming → spec → plan → execute).

— End of M1 Foundation plan —
