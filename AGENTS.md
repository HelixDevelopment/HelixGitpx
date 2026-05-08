# AGENTS.md

Instructions for AI agents (Claude Code, Cursor, Copilot, Aider, or any
future model-driven contributor) working in this repository.

Read this file **fully** before taking any action. If a rule here conflicts
with the [Constitution](./CONSTITUTION.md), the Constitution wins.

---

## 0. The Big Five

These five rules trump everything else an agent might want to do:

1. **No mocks outside unit tests.** Constitution Article II §2. Mocks,
   stubs, placeholder classes, and hardcoded data are allowed **only** in
   unit tests. Integration / e2e / security / stress / ddos / benchmark
   tests must exercise real dependencies.
2. **No skipped or disabled tests.** If you cannot make a test pass today,
   either fix the code or remove the feature. Silent debt is forbidden.
3. **Workflow-dispatch only CI.** Every GitHub Actions workflow must use
   `on: workflow_dispatch:` only. No `push`, no `pull_request`, no `schedule`.
   This is mandatory.
4. **Runtime-portable tooling.** Local dev tooling auto-detects between
   Docker and Podman. Never hardcode either one.
5. **Cite, don't fabricate.** If you claim a file, function, repo, or fact
   exists, verify it first. When uncertain, say so.

---

## 1. What This Project Is

HelixGitpx — Helix Git Proxy eXtended. A federated, multi-upstream Git
proxy platform. Polyglot monorepo: Go backend (18 microservices), Angular
19 web app, KMP + Compose Multiplatform clients (Android/iOS/Desktop),
Kubernetes platform manifests, and a Docusaurus docs site.

- **Repository:** solo-maintained by @milosvasic
- **License:** Apache-2.0 (code) / CC-BY-SA-4.0 (docs)
- **Version:** GA-tagged `v1.0.0`, milestones `m1-foundation` through `m8-ga`

---

## 2. Repository Layout

```
├── impl/
│   ├── helixgitpx/              # Go monorepo (platform + 18 services + gen + tools/scaffold)
│   │   ├── go.work              # Go workspace: gen, platform, 18 services, tools/scaffold
│   │   ├── proto/               # Protobuf definitions (buf.build/helixgitpx/core)
│   │   ├── gen/                 # Generated code (Go, TS, Kotlin, Swift, OpenAPI)
│   │   ├── api/openapi/         # Generated OpenAPI JSON specs
│   │   ├── platform/            # Shared libraries (auth, pg, redis, kafka, gin, grpc, health, log, etc.)
│   │   ├── services/            # 18 microservices, each with internal/{domain,handler,repo}
│   │   └── tools/scaffold/      # Service scaffolding tool
│   ├── helixgitpx-web/          # Angular 19 + Nx web app
│   ├── helixgitpx-clients/      # Kotlin Multiplatform + Compose Multiplatform
│   ├── helixgitpx-platform/     # Helm charts, Argo CD apps, Kustomize, SQL, OPA, compose
│   ├── helixgitpx-docs/         # Docusaurus docs site (helixgitpx.dev)
│   ├── helixgitpx-docs-site/    # Public-facing Docusaurus docs site
│   └── helixgitpx-website/      # Astro 5 marketing site (helixgitpx.io)
├── test/                        # Cross-service test suites (integration, e2e, security, stress, ddos, benchmark, chaos)
│   └── go.work                  # Separate workspace for cross-service tests
├── tools/                       # Coverage audit, perf (k6), fuzz, chaos, DR, docs-export
├── scripts/                     # Verification scripts (verify-m{1..8}-*, verify-everything.sh)
├── challenges/scripts/          # Challenge scripts (CONST-032 reproduction-before-fix, CONST-033)
├── Upstreams/                   # Multi-upstream federation scripts (GitHub.sh, GitLab.sh, GitFlic.sh, GitVerse.sh)
├── docs/                        # Specifications, runbooks, security, manuals, integrations
├── .github/workflows/           # GitHub Actions — ALL currently DISABLED (.yml.disabled)
├── CONSTITUTION.md              # Supreme policy document
├── CLAUDE.md                    # Claude-specific orientation (does NOT override Constitution or AGENTS.md)
├── AGENTS.md                    # This file — agent instructions
└── mise.toml                    # Tool version pinning (Go, Node, Java, Gradle, kubectl, helm, etc.)
```

### Go Monorepo (`impl/helixgitpx/`) — Service Structure

Every service follows the same internal layout (exemplified by `hello`):

```
services/<name>/
├── cmd/<name>/main.go           # Entrypoint (signal handling, migrate subcommand)
├── internal/
│   ├── app/app.go               # Composition root (wires all deps, starts servers)
│   ├── domain/                  # Pure business logic, no framework imports
│   │   ├── <aggregate>.go       # Interfaces (Counter, Cache, Emitter) + aggregate struct
│   │   └── <aggregate>_test.go  # Unit tests with in-process fakes (mocks allowed here)
│   ├── handler/
│   │   ├── http/router.go       # Gin routes
│   │   └── grpc/server.go       # gRPC service implementation
│   └── repo/                    # Concrete implementations (Postgres, Redis, Kafka)
└── test/integration/            # Integration tests (build tag: `integration`)
```

### Platform Libraries (`impl/helixgitpx/platform/`)

Shared packages imported by all services:

| Package | Purpose |
|---------|---------|
| `errors` | Canonical error type with gRPC codes + RFC 7807 problem mapping |
| `health` | `/livez`, `/readyz`, `/healthz` probe registry |
| `log` | Structured logger (slog-based) |
| `pg` | Postgres pool + migrations (goose) |
| `redis` | Redis client wrapper |
| `kafka` | Kafka producer/consumer (franz-go) |
| `gin` | Gin router factory with middleware |
| `grpc` | gRPC server factory with SPIFFE mTLS |
| `config` | Environment + Vault config loading |
| `telemetry` | OpenTelemetry setup |
| `spire` | SPIFFE/SPIRE fetcher |
| `auth` | JWT + session validation |
| `opa` | OPA policy evaluation |
| `testkit` | Testcontainers helpers (Postgres, Redis, Kafka) |

---

## 3. Toolchain and Versions

Pinned in `mise.toml`. Key tools:

| Tool | Version |
|------|---------|
| Go | 1.23.4 (`GOTOOLCHAIN=go1.23.4` required) |
| Node.js | 20.18.1 |
| Java | temurin-21.0.5+11 |
| Gradle | 8.10.2 |
| buf | 1.47.2 |
| sqlc | 1.27.0 |
| kubectl | 1.31.3 |
| helm | 3.16.3 |
| kind | 0.25.0 |
| k3d | 5.7.5 |
| golangci-lint | 1.62.2 |
| gofumpt | 0.7.0 |

Always set `GOTOOLCHAIN=go1.23.4` when running Go commands (the scripts and
workflows do this already).

---

## 4. Essential Commands

### Build

```bash
make build          # Build all sub-projects (Go, web, clients, docs)
make gen            # Regenerate protobuf/OpenAPI code (buf generate)
make bootstrap      # Fetch deps for every sub-project
```

### Go-specific (run from `impl/helixgitpx/`)

```bash
cd impl/helixgitpx
make build          # go build ./...
make test           # go test -race -shuffle=on -count=1 ./...
make lint           # golangci-lint run ./... && buf lint proto/
make fmt            # gofumpt -w . && goimports -w .
make cover          # Coverage report per package
make gen            # buf generate proto/
make tidy           # go work sync + go mod tidy in each module
```

### Web (Angular + Nx, run from `impl/helixgitpx-web/`)

```bash
cd impl/helixgitpx-web
npm install
npx nx build web    # Production build
npx nx test web     # Unit tests (Jest)
npx nx lint web     # ESLint
npx playwright test # E2E tests
```

### Clients (KMP, run from `impl/helixgitpx-clients/`)

```bash
cd impl/helixgitpx-clients
./gradlew :shared:compileKotlinJvm
./gradlew check
```

### Docs Site (Docusaurus, run from `impl/helixgitpx-docs/`)

```bash
cd impl/helixgitpx-docs
npm install
npm run build       # Sync + build
npm run start       # Dev server on :3001
```

### Marketing Site (Astro, run from `impl/helixgitpx-website/`)

```bash
cd impl/helixgitpx-website
npm install
npx astro build
npx astro dev
```

### Local Dev Stack (Compose)

```bash
make dev            # Bring up compose stack (postgres, kafka, redis, jaeger, prometheus, grafana, hello)
make dev-down       # Tear down
# hello REST:  http://localhost:8001/v1/hello?name=world
# hello gRPC:  localhost:9001
# Grafana:     http://localhost:3000 (admin/admin)
# Jaeger:      http://localhost:16686
```

### Test Matrix (all seven mandatory types)

```bash
make test-unit          # Unit tests (mocks allowed)
make test-integration   # Integration tests (compose must be up; build tag: integration)
make test-e2e           # End-to-end (Playwright + k3d)
make test-security      # Security scans
make test-stress        # k6 stress scenarios
make test-ddos          # Rate-limit / exhaustion tests
make test-benchmark     # Go micro-benchmarks + k6 budget checks
make test-all           # All seven types
```

### Verification and CI-Local

```bash
make ci-local                   # Full green-suite (verify-everything.sh)
bash scripts/verify-m1-artifacts.sh   # M1 artifact presence checks
bash scripts/verify-m2-artifacts.sh   # M2 artifact checks
bash scripts/verify-argo-paths.sh     # Argo CD Application path validation
bash scripts/verify-helm-charts.sh    # Helm chart lint
bash scripts/verify-rego.sh           # Rego syntax check
bash scripts/verify-proto-gen.sh      # Proto gen drift detection
bash scripts/verify-secrets.sh        # Secret scanning (gitleaks)
make coverage-audit                   # Per-package coverage audit (threshold: 80%)
```

### Upstream Federation

```bash
make upstream-push     # Push main + tags to ALL upstreams
make upstream-status   # Show divergence per upstream
bash Upstreams/GitHub.sh    # Sets UPSTREAMABLE_REPOSITORY for GitHub
```

### Docs Export

```bash
bash tools/docs-export/build-all.sh         # Build all manuals (PDF, ePub, DOCX, TXT)
bash tools/docs-export/build-all.sh --check  # Check toolchain availability
```

---

## 5. Code Style and Conventions

### Go

- **Formatting:** `gofumpt` (extra-rules: true) + `goimports`. Go files use
  **tabs** for indentation (see `.editorconfig`).
- **Linting:** `golangci-lint` with a strict linter set (see
  `impl/helixgitpx/.golangci.yml`). Linters relaxed for `*_test.go` files
  (gosec, errcheck, unparam excluded).
- **Module paths:**
  - Platform: `github.com/helixgitpx/platform`
  - Services: `github.com/helixgitpx/helixgitpx/services/<name>`
  - Generated: `github.com/helixgitpx/helixgitpx/gen/go`
  - Integration tests: `github.com/helixgitpx/helixgitpx/test/integration`
- **Error handling:** Use `github.com/helixgitpx/platform/errors` —
  `errors.New(codes.XXX, "domain", "message").Wrap(cause)`. Maps to
  RFC 7807 problems via `.ToProblem()`.
- **Domain layer:** Pure Go, no framework imports. Define interfaces
  (`Counter`, `Cache`, `Emitter`) and implement in `repo/`.
- **Test file naming:**
  - Unit: `<subject>_test.go` (co-located)
  - Integration: `<subject>_integration_test.go` (build tag `//go:build integration`)
  - Benchmark: `*_bench_test.go`
- **Unit test fakes:** Hand-written structs implementing domain interfaces.
  No mocking frameworks. Fakes are per-test-file, not shared.
- **Config:** Via `platform/config` — struct tags `env:"NAME"`,
  `vault:"kv/path#key"`, `default:"value"`. Prefix with service name
  (e.g., `HELLO_POSTGRES_DSN`).
- **Graceful shutdown:** Use `signal.NotifyContext` + `GracefulStop`/`Shutdown`
  with 10s timeout.
- **Concurrent maps/slices:** Must use `safe.Store[K,V]` / `safe.Slice[T]`
  from the project's concurrency primitives. Bare `sync.Mutex + map/slice`
  is prohibited for new code.

### TypeScript / Angular (web app)

- **Formatting:** 2-space indent (`.editorconfig`).
- **Build:** Nx monorepo (`nx.json` with cached targets for build, lint, test).
- **Proto:** Generated TS clients land in `libs/proto/src/` via buf.
- **Testing:** Jest for unit, Playwright for E2E.
- **Scripts:** Use `npx nx <target> web` to run targets.

### Kotlin (KMP clients)

- **Formatting:** 4-space indent.
- **Proto:** Generated Kotlin to `shared/src/commonMain/kotlin/gen/`.
- **Gradle:** Convention plugin at `buildSrc/src/main/kotlin/helix.convention.gradle.kts`.

### Protobuf

- **Location:** `impl/helixgitpx/proto/helixgitpx/<domain>/v1/<domain>.proto`
- **Module:** `buf.build/helixgitpx/core`
- **Code generation:** Outputs to `gen/go/`, `helixgitpx-web/libs/proto/src/`,
  `helixgitpx-clients/shared/src/commonMain/kotlin/gen/`,
  `helixgitpx-clients/iosApp/Gen/`, `api/openapi/`.
- **Generated code** is committed and marked `linguist-generated` in
  `.gitattributes`. Run `make gen` after proto changes, then verify with
  `bash scripts/verify-proto-gen.sh`.

### SQL

- **Schemas:** `impl/helixgitpx-platform/sql/schemas.sql` — one Postgres
  schema per service domain (`hello`, `auth`, `org`, `team`, `repo`, etc.).
- **Migrations:** `impl/helixgitpx-platform/sql/migrations/`, applied via
  `goose` (invoked through `platform/pg`).
- **Cross-schema roles:** e.g., `orgteam_svc` has access to both `org` and
  `team` schemas (ADR-0014).

### YAML / Kubernetes

- **Helm charts:** `impl/helixgitpx-platform/helm/<chart>/`
- **Argo CD apps:** `impl/helixgitpx-platform/argocd/applications/`
- **Kustomize:** `impl/helixgitpx-platform/kustomize/base/` + `overlays/`
- **Compose:** `impl/helixgitpx-platform/compose/compose.yml` — uses profiles
  (`core`, `observability`, `all`)

---

## 6. Testing Approach

### Seven Mandatory Types (Constitution Article II)

| Type | Location | Runner | Mocks? |
|------|----------|--------|--------|
| unit | `*_test.go` co-located, `*.spec.ts` | `go test`, jest | **Yes** |
| integration | `test/integration/`, `services/*/test/integration/` | `go test -tags=integration` | No |
| e2e | `test/e2e/`, `impl/helixgitpx-web/e2e/` | k6, Playwright | No |
| security | `test/security/` | OWASP ZAP, Nuclei, gosec | No |
| stress | `tools/perf/scenarios/*_stress.js` | k6 | No |
| ddos | `test/ddos/` | k6 + Litmus | No |
| benchmark | `*_bench_test.go`, `tools/perf/` | `go test -bench`, k6 | No |

### Key Testing Patterns

- **Unit tests** use hand-written fakes implementing domain interfaces.
  No mock frameworks.
- **Integration tests** use `platform/testkit` (testcontainers) to spin up
  real Postgres, Redis, Kafka. Build-tagged with `//go:build integration`.
- **Coverage target:** 100% per type per module touched.
- **Resource limits:** All test/challenge runs must be limited to 30-40% of
  host resources: `GOMAXPROCS=2`, `nice -n 19`, `ionice -c 3`, `-p 1` for
  `go test`.
- **No skips:** `t.Skip` / `xit` / `@Ignore` without `SKIP-OK: #<ticket>`
  are forbidden.

### Reproduction-Before-Fix (CONST-032)

Every bug fix must follow this sequence:

1. Write a Challenge script in `challenges/scripts/`
2. Run it — confirm it reproduces the bug (fails)
3. Write the fix
4. Re-run — confirm pass
5. Commit Challenge + fix together

---

## 7. Protobuf and Code Generation

```
impl/helixgitpx/proto/           # Source of truth
  buf.yaml                       # Module: buf.build/helixgitpx/core
  buf.gen.yaml                   # Plugins: go, grpc-go, connect-go, es, connect-es, grpc-kotlin, swift, openapiv2
  helixgitpx/<domain>/v1/*.proto # 21 proto files across 12 domains

make gen  →  buf generate
Outputs:
  impl/helixgitpx/gen/go/                          # Go
  impl/helixgitpx-web/libs/proto/src/              # TypeScript
  impl/helixgitpx-clients/shared/src/.../kotlin/gen/  # Kotlin
  impl/helixgitpx-clients/iosApp/Gen/              # Swift
  impl/helixgitpx/api/openapi/                     # OpenAPI JSON
```

After any proto change: `make gen` then `bash scripts/verify-proto-gen.sh`.

---

## 8. Commit and Branch Discipline

- **Branches:** `feat/…`, `fix/…`, `docs/…`, `chore/…`
- **Commits:** Conventional Commits (`feat(m7): …`, `fix(auth): …`)
- **Sign-off:** Every commit must have `Signed-off-by: Name <email>` (DCO)
- **AI attribution:** `Co-Authored-By: <model> <noreply@anthropic.com>` when
  an AI agent wrote or materially shaped the patch
- **No force-push** to `main`. Never amend a published commit.
- **After merge:** Run `bash Upstreams/GitHub.sh` + other upstream scripts to
  sync to all federated remotes.

---

## 9. GitHub Actions Workflows

All workflows are **currently DISABLED** (files renamed to `.yml.disabled`).

When re-enabled, every workflow uses `on: workflow_dispatch:` only.
No `push`, `pull_request`, or `schedule` triggers are permitted.

| Workflow | Purpose |
|----------|---------|
| ci-go | `go vet` + `go test -race` + `go build` across all modules |
| ci-web | Nx build + Jest + Playwright for Angular app |
| ci-clients | Gradle `compileKotlinJvm` for KMP shared module |
| ci-docs | Docusaurus build + manual export + marketing site build |
| ci-platform | Argo CD path check + Helm chart lint + Rego syntax |
| ci-verifiers | Full `scripts/verify-everything.sh` green-suite |
| security-scan | gosec, Trivy FS, Grype, CodeQL (optional deep) |
| supply-chain | Per-service container build + SBOM + Cosign sign + Trivy scan |
| perf-budgets | k6 scenarios + budget gate |
| mutation-testing | go-mutesting per service |
| release | Tag-triggered release (dispatch-only) |
| deploy | GitOps image promotion |
| upstream-sync | Push main + tags to all federated upstreams |

---

## 10. Before You Edit

1. `git log --oneline -20` to orient.
2. Read the most recent spec for the area you're touching, indexed at
   `docs/specifications/main/main_implementation_material/HelixGitpx/README.md`.
3. Check for an open plan at `docs/superpowers/plans/` — if one exists for
   your task, follow it.
4. Check for an ADR under
   `docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/`
   that covers the subject.

---

## 11. When You Finish Work

1. Run the relevant verifier: `bash scripts/verify-m1-artifacts.sh` through
   `bash scripts/verify-m8-cluster.sh` depending on the milestone.
2. Run `make lint` (or `cd impl/helixgitpx && make lint` for Go).
3. Run `make test` (or the specific test type you need).
4. Run `bash scripts/verify-everything.sh` for the full green-suite.
5. Run `bash challenges/scripts/no_suspend_calls_challenge.sh` (CONST-033).
6. If you touched protos: `make gen` then `bash scripts/verify-proto-gen.sh`.
7. If you touched the web app: `cd impl/helixgitpx-web && npx nx test web && npx nx build web`.
8. If you touched KMP code: `cd impl/helixgitpx-clients && ./gradlew :shared:check`.
9. If you touched the docs site: `cd impl/helixgitpx-docs && npm run build`.
10. Commit, and sync upstreams if on `main`.

---

## 12. Escalations

- Anything that changes the Constitution → stop, ask the human.
- Anything that changes `mandatory` policies (CI dispatch-only, container
  runtime portability, testing rules) → stop, ask the human.
- Deleting tests or skipping them → stop, ask the human.
- Force-pushing, destructive git, rewriting published history → stop, ask.

---

## 13. Gotchas and Non-Obvious Patterns

- **`go.work` is gitignored at root.** The root `.gitignore` ignores
  `go.work` / `go.work.sum` but has exceptions for `impl/**/go.work` and
  `test/go.work`. Go workspaces live only under `impl/` and `test/`.
- **Generated code is committed.** The `gen/` directory is tracked in git
  and marked `linguist-generated` in `.gitattributes`. Don't edit it
  manually; always use `make gen`.
- **Integration tests need compose.** `make test-integration` requires the
  compose stack running (`make dev`). Tests that can't connect to real
  services must skip (not fail).
- **Compose wrapper.** The compose binary is at
  `impl/helixgitpx-platform/compose/bin/compose` (auto-detects docker vs
  podman). Don't call `docker-compose` directly.
- **Environment variables are prefixed.** Each service's config is prefixed
  with its name (e.g., `HELLO_HTTP_ADDR`, `HELLO_POSTGRES_DSN`). The
  `platform/config` package handles this via `config.Options{Prefix: "HELLO"}`.
- **`package-lock.json` is gitignored.** Web and docs sites don't commit
  lockfiles. CI runs `npm install --legacy-peer-deps`.
- **Spec sections marked `[VERIFY-AT-INTEGRATION]`** contain facts that
  drift (pricing, API versions). Treat as TODOs, not errors.
- **Old master spec is frozen.** Don't edit
  `docs/specifications/main/Git_Proxy_Master_Specification.md` — it's
  superseded by the suite under `main_implementation_material/HelixGitpx/`.
- **SQL schemas use per-service roles.** Each service gets its own Postgres
  role (`<domain>_svc`) with access only to its schema. The `orgteam` service
  is special: `orgteam_svc` accesses both `org` and `team` schemas.
- **Kyverno policies enforce cluster standards.** Four policies at
  `impl/helixgitpx-platform/kyverno/policies/`: require-labels,
  require-signed-images, enforce-resource-limits, disallow-privileged.
- **Health endpoints are mandatory.** Every service must expose
  `/healthz`, `/livez`, `/readyz` via `platform/health`.
- **Challenge scripts live at `challenges/scripts/`.** Every bug must have a
  reproduction challenge before the fix is written.

---

## Universal Mandatory Constraints

> Cascaded from the HelixAgent root `CLAUDE.md` via `/tmp/UNIVERSAL_MANDATORY_RULES.md`.
> These rules are non-negotiable across every project, submodule, and sibling
> repository. Project-specific addenda are welcome but cannot weaken or
> override these.

### Hard Stops (permanent, non-negotiable)

1. **NO CI/CD pipelines.** No `.github/workflows/`, `.gitlab-ci.yml`,
   `Jenkinsfile`, `.travis.yml`, `.circleci/`, or any automated pipeline.
   No Git hooks either. All builds and tests run manually or via
   Makefile/script targets.
2. **NO HTTPS for Git.** SSH URLs only (`git@github.com:…`,
   `git@gitlab.com:…`, etc.) for clones, fetches, pushes, and submodule
   updates. Including for public repos. SSH keys are configured on every
   service.
3. **NO manual container commands.** Container orchestration is owned by
   the project's binary/orchestrator (e.g. `make build` → `./bin/<app>`).
   Direct `docker`/`podman start|stop|rm` and `docker-compose up|down`
   are prohibited as workflows. The orchestrator reads its configured
   `.env` and brings up everything.

### Mandatory Development Standards

1. **100% Test Coverage.** Every component MUST have unit, integration,
   E2E, automation, security/penetration, and benchmark tests. No false
   positives. Mocks/stubs ONLY in unit tests; all other test types use
   real data and live services.
2. **Challenge Coverage.** Every component MUST have Challenge scripts
   (`./challenges/scripts/`) validating real-life use cases. No false
   success — validate actual behavior, not return codes.
3. **Real Data.** Beyond unit tests, all components MUST use actual API
   calls, real databases, live services. No simulated success. Fallback
   chains tested with actual failures.
4. **Health & Observability.** Every service MUST expose health
   endpoints. Circuit breakers for all external dependencies.
   Prometheus / OpenTelemetry integration where applicable.
5. **Documentation & Quality.** Update `CLAUDE.md`, `AGENTS.md`, and
   relevant docs alongside code changes. Pass language-appropriate
   format/lint/security gates. Conventional Commits:
   `<type>(<scope>): <description>`.
6. **Validation Before Release.** Pass the project's full validation
   suite (`make ci-local`) plus all challenges
   (`./challenges/scripts/run_all_challenges.sh`).
7. **No Mocks or Stubs in Production.** Mocks, stubs, fakes,
   placeholder classes, TODO implementations are STRICTLY FORBIDDEN in
   production code. All production code is fully functional with real
   integrations. Only unit tests may use mocks/stubs.
8. **Comprehensive Verification.** Every fix MUST be verified from all
   angles: runtime testing (actual HTTP requests / real CLI
   invocations), compile verification, code structure checks,
   dependency existence checks, backward compatibility, and no false
   positives in tests or challenges. Grep-only validation is NEVER
   sufficient.
9. **Resource Limits for Tests & Challenges (CRITICAL).** ALL test and
   challenge execution MUST be strictly limited to 30-40% of host
   system resources. Use `GOMAXPROCS=2`, `nice -n 19`, `ionice -c 3`,
   `-p 1` for `go test`. Container limits required. The host runs
   mission-critical processes — exceeding limits causes system crashes.
10. **Bugfix Documentation.** All bug fixes MUST be documented in
    `docs/issues/fixed/BUGFIXES.md` (or the project's equivalent) with
    root cause analysis, affected files, fix description, and a link to
    the verification test/challenge.
11. **Real Infrastructure for All Non-Unit Tests.** Mocks/fakes/stubs/
    placeholders MAY be used ONLY in unit tests (files ending
    `_test.go` run under `go test -short`, equivalent for other
    languages). ALL other test types — integration, E2E, functional,
    security, stress, chaos, challenge, benchmark, runtime
    verification — MUST execute against the REAL running system with
    REAL containers, REAL databases, REAL services, and REAL HTTP
    calls. Non-unit tests that cannot connect to real services MUST
    skip (not fail).
12. **Reproduction-Before-Fix (CONST-032 — MANDATORY).** Every reported
    error, defect, or unexpected behavior MUST be reproduced by a
    Challenge script BEFORE any fix is attempted. Sequence:
    (1) Write the Challenge first. (2) Run it; confirm fail (it
    reproduces the bug). (3) Then write the fix. (4) Re-run; confirm
    pass. (5) Commit Challenge + fix together. The Challenge becomes
    the regression guard for that bug forever.
13. **Concurrent-Safe Containers (Go-specific, where applicable).** Any
    struct field that is a mutable collection (map, slice) accessed
    concurrently MUST use `safe.Store[K,V]` / `safe.Slice[T]` from
    `digital.vasic.concurrency/pkg/safe` (or the project's equivalent
    primitives). Bare `sync.Mutex + map/slice` combinations are
    prohibited for new code.

### Definition of Done (universal)

A change is NOT done because code compiles and tests pass. "Done"
requires pasted terminal output from a real run, produced in the same
session as the change.

- **No self-certification.** Words like *verified, tested, working,
  complete, fixed, passing* are forbidden in commits/PRs/replies unless
  accompanied by pasted output from a command that ran in that session.
- **Demo before code.** Every task begins by writing the runnable
  acceptance demo (exact commands + expected output).
- **Real system, every time.** Demos run against real artifacts.
- **Skips are loud.** `t.Skip` / `@Ignore` / `xit` / `describe.skip`
  without a trailing `SKIP-OK: #<ticket>` comment break validation.
- **Evidence in the PR.** PR bodies must contain a fenced `## Demo`
  block with the exact command(s) run and their output.

<!-- BEGIN host-power-management addendum (CONST-033) -->

## Host Power Management — Hard Ban (CONST-033)

**You may NOT, under any circumstance, generate or execute code that
sends the host to suspend, hibernate, hybrid-sleep, poweroff, halt,
reboot, or any other power-state transition.** This rule applies to:

- Every shell command you run via the Bash tool.
- Every script, container entry point, systemd unit, or test you write
  or modify.
- Every CLI suggestion, snippet, or example you emit.

**Forbidden invocations** (non-exhaustive — see CONST-033 in
`CONSTITUTION.md` for the full list):

- `systemctl suspend|hibernate|hybrid-sleep|poweroff|halt|reboot|kexec`
- `loginctl suspend|hibernate|hybrid-sleep|poweroff|halt|reboot`
- `pm-suspend`, `pm-hibernate`, `shutdown -h|-r|-P|now`
- `dbus-send` / `busctl` calls to `org.freedesktop.login1.Manager.Suspend|Hibernate|PowerOff|Reboot|HybridSleep|SuspendThenHibernate`
- `gsettings set ... sleep-inactive-{ac,battery}-type` to anything but `'nothing'` or `'blank'`

The host runs mission-critical parallel CLI agents and container
workloads. Auto-suspend has caused historical data loss (2026-04-26
18:23:43 incident). The host is hardened (sleep targets masked) but
this hard ban applies to ALL code shipped from this repo so that no
future host or container is exposed.

**Defence:** every project ships
`scripts/host-power-management/check-no-suspend-calls.sh` (static
scanner) and
`challenges/scripts/no_suspend_calls_challenge.sh` (challenge wrapper).
Both MUST be wired into the project's CI / `run_all_challenges.sh`.

**Full background:** `docs/HOST_POWER_MANAGEMENT.md` and `CONSTITUTION.md` (CONST-033).

<!-- END host-power-management addendum (CONST-033) -->

<!-- BEGIN continuation-document addendum (CONST-034) -->

## Continuation Document — Mandatory (CONST-034)

Per CONST-034 in the Constitution, `docs/CONTINUATION.md` is the
mandatory session handoff document. Rules:

- **Read it first** when starting any work session.
- **Update it before stopping** — what was done, what changed, next steps.
- **Never let it go out of sync** with current work.
- **It must be self-contained** — enough context for any agent or LLM
  to resume without any other files or conversation history.
- **Commit it with your work** — same commit or the very next one.

If `docs/CONTINUATION.md` does not exist, create it as the first action.

<!-- END continuation-document addendum (CONST-034) -->

<!-- BEGIN anti-bluff addendum (CONST-035) -->

## Anti-Bluff — Tests Must Prove Real Functionality (CONST-035)

Per CONST-035 in the Constitution, **a passing test suite MUST guarantee
that the tested features actually work for real end users.** Bluff tests
are forbidden. A bluff test is one that passes regardless of whether the
feature works.

**Forbidden patterns:**

- Tautological assertions (`expect(true).toBe(true)`, `assert 1 == 1`).
- Healthz-only integration tests (must exercise at least one real
  business endpoint and verify response semantics).
- Grep-only challenges (must execute runtime behavior).
- Stub-method tests on hardcoded returns (must verify behavior that
  breaks when the implementation is wrong).
- Always-skip conditionals (if the prerequisite is never met, the
  test is a bluff).
- Status-code-only assertions without body/state verification.

**Every test MUST:**

1. Assert at least one behavioral property that would fail if the
   implementation were buggy.
2. Verify something a real user would notice.
3. Be capable of failing — if no code change could make it fail,
   it is a bluff.

This rule applies to all submodules, sub-projects, and sibling
repositories. See CONST-035 in `CONSTITUTION.md` for the full rule.

<!-- END anti-bluff addendum (CONST-035) -->
