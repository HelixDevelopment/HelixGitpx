# 20 — Developer Guide

> **Document purpose**: Get a new contributor productive in **one day**. Covers local setup, workflow, coding standards, testing, and PR etiquette.

---

## 1. Prerequisites

Install via your package manager or the bootstrap script:

```bash
curl -fsSL https://raw.githubusercontent.com/vasic-digital/helixgitpx/main/scripts/bootstrap.sh | sh
```

This installs:

- `mise` (polyglot version manager) and uses the repo's `.mise.toml` for:
  - Go 1.23+
  - Node.js 22 LTS + pnpm 9
  - Java 21 LTS + Gradle 8.8
  - Kotlin 2.0
  - Python 3.12
  - Rust 1.80
  - Buf, kubectl, helm, kustomize, argocd, temporal CLI, k6.
- `devbox` (optional) for a Nix-backed shell if you prefer.
- `docker` / `podman`.
- `kind` + `ctlptl` (local cluster).
- `tilt` (incremental dev loop).
- `golangci-lint`, `gofumpt`, `goimports`.
- `helixctl` (our CLI).

---

## 2. Clone & Bootstrap

```bash
git clone https://github.com/vasic-digital/helixgitpx.git
cd helixgitpx
make bootstrap
```

`make bootstrap` performs:
- pre-commit hook installation (lefthook).
- Generates protobuf code (`buf generate`).
- Downloads Go modules, pnpm packages, Gradle deps.
- Creates local Vault dev secrets.
- Starts a Kind cluster + Tilt to deploy dev manifests.

---

## 3. Daily Workflow

### 3.1 Starting work

```bash
git checkout -b feat/<ticket>-<short-name>
make dev                            # Tilt up; live-reload services
```

You'll see a dashboard on <http://localhost:10350>. Each service is visible with logs, metrics, and reload controls.

### 3.2 Running specific things

```bash
# Run one service locally (outside the cluster)
make run SERVICE=repo-service

# Run Angular dev server
cd web && pnpm nx serve web

# Run desktop app
cd clients && ./gradlew desktopApp:run

# Run Android on emulator
cd clients && ./gradlew androidApp:installDebug
```

### 3.3 Tests

```bash
# Full test suite for Go
make test

# One package with race detector + coverage
go test -race -cover ./services/repo-service/...

# Angular unit
pnpm nx test web

# KMP common tests
cd clients && ./gradlew :shared:core:allTests

# E2E (Playwright)
pnpm nx e2e e2e

# Integration tests (Testcontainers)
make test-integration

# Load tests
k6 run tests/load/api.js
```

### 3.4 Code generation

```bash
buf generate                        # proto → Go / TS / Kotlin / Swift
make sqlgen                         # sqlc → typed Go from .sql
./gradlew :shared:network:generateProtos   # KMP side
pnpm nx run sdk:generate            # TS types
```

All generated code is committed; CI verifies it's up-to-date.

---

## 4. Project Layout

See each repo's top-level README; highlights:

```
helixgitpx/                        # backend monorepo
├── services/
│   ├── repo-service/
│   │   ├── cmd/                   # main.go
│   │   ├── internal/              # domain logic; no external imports
│   │   ├── api/                   # gRPC handlers
│   │   ├── db/                    # sqlc + migrations
│   │   └── test/                  # integration helpers
│   └── …
├── pkg/helix-platform/            # shared libs
├── proto/                         # protobufs
├── deploy/                        # Helm charts
├── tools/                         # generators, linters, scripts
└── Makefile
```

---

## 5. Coding Standards

### 5.1 Go

- Follow "Go Code Review Comments" + effective-go + Uber style guide (amended).
- No global state; dependency injection via constructor-passed interfaces.
- Explicit errors; wrap with `%w`; typed domain errors via `helix-platform/errors`.
- Context in every function that does I/O.
- Concurrency with `errgroup`; no raw goroutines in library code.
- Package names: lower_snake_case single word; no `util` packages.

Example minimum idioms:

```go
type Service struct {
    log   slog.Logger
    db    RepoStore
    kafka ProducerClient
}

func New(log slog.Logger, db RepoStore, kafka ProducerClient) *Service { … }

func (s *Service) CreateRepo(ctx context.Context, in CreateRepoInput) (*Repo, error) {
    if err := in.Validate(); err != nil {
        return nil, errors.InvalidArgument("create_repo", err)
    }
    …
}
```

### 5.2 TypeScript / Angular

- Strict TS (`strict: true`, `noUncheckedIndexedAccess: true`).
- Standalone components; OnPush change detection; signals first.
- No `any` in production code; `unknown` + narrowing.
- No direct DOM; use `DOCUMENT` / `Renderer2` only when necessary.
- Imports grouped: node builtins, third-party, aliased, relative.

### 5.3 Kotlin / KMP

- Explicit visibility (`public` never implicit).
- Immutable by default (`val`, `data class`).
- Coroutines everywhere; avoid `runBlocking` in production.
- Sealed interfaces for finite states.
- Shared modules must not depend on any platform.

### 5.4 SQL

- One statement per migration file is ideal.
- Use UUIDv7 via `uuid_generate_v7()` helper (or canonical library).
- Every table has created_at / updated_at timestamps and a Postgres trigger to maintain `updated_at`.

### 5.5 Protobuf

- Snake_case fields; CamelCase messages; no trailing `_` in generated code.
- One RPC per method; no kitchen-sink methods.
- Use `google.protobuf.Empty` for void; `google.protobuf.Timestamp` for times.

### 5.6 Git

- Conventional commits enforced.
- Imperative, <=72 chars subject; body wrapped at 72.
- Reference ticket in footer (`Refs: HGX-123`).

---

## 6. Writing Tests

### 6.1 Principles

- Tests are production code. Same coding standards apply.
- Prefer **table-driven** tests; one assertion per case.
- Integration tests use **Testcontainers**; never rely on pre-existing services.
- Avoid snapshot tests that are opaque; prefer explicit assertions.
- Deterministic: clock injected via `helix-platform/clock`; random via `helix-platform/rand`.

### 6.2 Structure

```go
func TestRepoService_CreateRepo(t *testing.T) {
    t.Parallel()
    cases := []struct{
        name   string
        in     CreateRepoInput
        seeds  func(*testing.T, *fx.Container)
        want   *Repo
        wantErr error
    }{
        {name: "happy path", …},
        {name: "duplicate slug", …, wantErr: errors.CodeAlreadyExists},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // arrange, act, assert
        })
    }
}
```

### 6.3 Integration

```go
func TestSyncFlow_Integration(t *testing.T) {
    ctx := testhelper.NewContext(t)
    env := testhelper.NewEnv(t)    // boots PG + Kafka + Redis + mock adapters
    …
}
```

---

## 7. Pull Request Checklist

Open a PR when you've:

- [ ] Named the branch per convention.
- [ ] Kept the PR focused (< 400 lines of diff ideal; > 1000 requires justification).
- [ ] Included/updated tests; coverage gate green.
- [ ] Updated docs (API, ADR, this guide) where relevant.
- [ ] Added / adjusted dashboards + alerts if behaviour changed.
- [ ] Verified locally: `make lint test build`.
- [ ] Filled in the PR template.

Review expectations:

- Reviewers respond within one business day.
- Use "Request changes" sparingly; prefer comments + conversation.
- Approvals are signatures of trust, not rubber stamps.

---

## 8. Tips & Tricks

- `helixctl tail <service>` — follow structured logs, filter by `--trace-id`.
- `make dev-seed` — populate local env with fixture data (orgs, repos, events).
- Tilt triggers: hit a service name in the dashboard to rebuild just that one.
- `go test -run TestFoo -v -race -count=10` quickly surfaces flakes.
- Angular: `pnpm nx graph` to visualise module dependencies.
- KMP: `./gradlew :shared:tasks --all` to see all available tasks.
- **Record everything**: when you debug something tricky, drop a note in `docs/troubleshooting/`; future-you will thank you.

---

## 9. Getting Help

- `#helixgitpx-dev` Slack channel.
- Office hours: Thursdays, 15:00 CET.
- Pair programming available — post in `#pair-me`.
- Architecture questions: open an RFC PR in `docs/rfcs/`.

---

## 10. Contribution License

- Apache-2.0.
- CLA enforced on external contributions (Developer Certificate of Origin alternative accepted for individuals).

---

*— End of Developer Guide —*
