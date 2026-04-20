# M1 Foundation — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED (pending user review) |
| Author | Милош Васић + Claude (brainstorming session 2026-04-20) |
| Milestone | M1 — Foundation (Weeks 1–4 in `13-roadmap/17-milestones.md`) |
| Scope | Full 18-item roadmap; nothing skipped |
| Sequencing | Approach 2 — spine-first, then fill in breadth |
| Supersedes | — |
| Implements | `docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md` §§1.1–1.4 |

---

## 1. Context & problem statement

The HelixGitpx repository currently contains authoritative specifications only (`docs/specifications/.../HelixGitpx/`, 15,767 lines, 90 files) and four upstream-mirror scripts under `Upstreams/`. No implementation code exists. The user has asked to begin implementing the spec, starting at M1 Foundation, with nothing skipped and with checkpoints between milestones.

This document is the design spec for M1. It does not attempt to cover M2–M8; each of those milestones will get its own brainstorming → spec → plan → execute loop after M1 is validated.

## 2. Goals

G1. Produce a working monorepo skeleton with every artifact called for by roadmap §§1.1–1.4 (18 items).

G2. Prove the end-to-end spine in one working service (`hello`) that exercises Postgres, Kafka, Redis, gRPC, REST, health, telemetry, and container runtime portability.

G3. Make the local developer experience reach the exit criterion: `mise install && make bootstrap && make dev`, then `curl` / `grpcurl` get a greeting back.

G4. Ship every CI workflow the spec calls for, manually triggered only (`workflow_dispatch`), with secret-gated steps that skip gracefully when secrets are absent — so the full catalog is present from day one.

G5. Ship configuration-only artifacts (Kyverno, Kata runners, Vault+OIDC, self-hosted BSR) that activate when their infra arrives in M2+, so that "nothing skipped" is honored without requiring infra that doesn't yet exist.

G6. Seed the ADR registry with the five architectural decisions made during this brainstorming (Q1 through Q5).

## 3. Non-goals

- Implementing any of the 24 non-`hello` services (those are M3–M8).
- Standing up a production or staging Kubernetes cluster (M2).
- Deploying Vault, SPIRE, self-hosted SonarQube, Snyk on-prem, or the self-hosted BSR (they ship as config; activation is M2+).
- Multi-region, DR drills, mutation testing thresholds, SOC2 evidence (M8).
- Federated Git provider adapters (M4–M5).
- Rewriting or editing any existing spec document under `docs/specifications/`. The spec is immutable input to M1.

## 4. Locked constraints from brainstorming

| ID | Constraint | Source |
|---|---|---|
| C-1 | Implementation code lives under `impl/<repo>/` subdirs of this repository, one git history, 5 logical sub-projects | Q1 |
| C-2 | GitHub Actions is the sole CI host; every workflow uses `on: workflow_dispatch` only (no `push:`, `pull_request:`, `schedule:`) — **mandatory** | Q2 |
| C-3 | Local dev stack runs on compose; a wrapper script auto-detects docker vs podman; no Kubernetes required for M1 exit criterion | Q3 |
| C-4 | M1 ships the full 18 roadmap items; what can't be deployed now ships as config + scripts + docs that activate when infra arrives | Q4 |
| C-5 | Toolchain pinning via `mise.toml` at repo root (not devbox) | Q5 |
| C-6 | Sequencing is spine-first (Approach 2): thin end-to-end vertical slice first, then breadth | Approaches |

## 5. Repository layout

```
HelixGitpx/                          (this repo, existing)
├── docs/                            (existing specs — untouched)
├── Upstreams/                       (existing mirror scripts — untouched)
├── impl/
│   ├── helixgitpx/                  Go monorepo (the 25 services + shared libs)
│   │   ├── go.work
│   │   ├── platform/                helix-platform shared libs (module: github.com/helixgitpx/platform)
│   │   │   ├── log/  telemetry/  errors/  config/
│   │   │   ├── grpc/  gin/  kafka/  pg/  redis/
│   │   │   ├── temporal/  spire/  opa/  health/  testkit/
│   │   ├── services/
│   │   │   └── hello/               first service (gRPC + REST), generated from template
│   │   ├── proto/
│   │   │   ├── buf.yaml  buf.gen.yaml  buf.work.yaml  buf.lock
│   │   │   └── helixgitpx/<domain>/v1/*.proto
│   │   ├── gen/go/                  committed codegen output
│   │   ├── tools/scaffold/          Go binary replacing cookiecutter
│   │   ├── api/openapi/             generated OpenAPI
│   │   └── Makefile
│   ├── helixgitpx-web/              Nx workspace (Angular 19)
│   │   ├── nx.json  package.json  tsconfig.base.json
│   │   ├── apps/web/
│   │   └── libs/proto/              TS codegen output
│   ├── helixgitpx-clients/          KMP + Compose Multiplatform
│   │   ├── settings.gradle.kts  build.gradle.kts
│   │   ├── buildSrc/                Gradle convention plugins
│   │   └── shared/ androidApp/ iosApp/ desktopApp/
│   ├── helixgitpx-platform/         GitOps / infra
│   │   ├── argocd/  helm/  kustomize/  terraform/
│   │   ├── kyverno/policies/
│   │   ├── github-actions-runner-controller/
│   │   ├── vault/
│   │   ├── k8s-local/               kind/k3d scripts
│   │   └── compose/
│   │       ├── compose.yml
│   │       ├── observability/
│   │       └── bin/compose          runtime-detect wrapper
│   └── helixgitpx-docs/             Docusaurus site
│       ├── docusaurus.config.ts  sidebars.ts  package.json
│       ├── src/  static/
│       └── docs/                    sync'd from ../../docs/specifications/...
├── .github/workflows/               umbrella workflows (manual-dispatch)
├── mise.toml                        toolchain pins
└── Makefile                         thin orchestrator
```

Each `impl/<subdir>/` also has its own `.github/workflows/ci.yml` so the logical sub-repo is splittable later.

## 6. `helix-platform` shared Go libraries

Module `github.com/helixgitpx/platform` under `impl/helixgitpx/platform/`.

| Package | M1 content | Grows in |
|---|---|---|
| `log` | zap/slog wrapper, JSON output, `WithContext(ctx) *Logger` | — |
| `telemetry` | OTel SDK bootstrap (traces, metrics, logs) via env vars; no-op exporter default | M2 (real exporters) |
| `errors` | typed error codes mapped to gRPC `codes.Code` + HTTP status; RFC 7807 problem details; `errors.Is/As` friendly | — |
| `config` | Viper-based: env > flags > file; typed struct loader; validation via `go-playground/validator` | — |
| `grpc` | server constructor with interceptor chain (logging, recovery, telemetry, auth-stub); client dialer with retries; TLS optional | M2 (SPIFFE/mTLS), M3 (auth) |
| `gin` | router constructor with middleware chain parallel to gRPC | same |
| `kafka` | franz-go producer/consumer wrappers; Karapace-compatible schema-registry hook (stubbed in M1) | M2 (Karapace wired) |
| `pg` | pgx pool + `sqlc`-generated query helper + migration runner (goose) | M2 |
| `redis` | go-redis v9 client; namespaced keys | — |
| `temporal` | client constructor, worker registration helper | M5 (real workflows) |
| `spire` | SVID fetcher scaffold; workload API client creation (no-op if not present) | M2 |
| `opa` | bundle loader + evaluator wrapper (in-process) | M3+ |
| `health` | `/healthz`, `/readyz`, `/livez` handlers + background dependency probes | — |
| `testkit` | testcontainers-go helpers for Postgres/Kafka/Redis; golden-file assert; httptest+grpc harness | — |

Each sub-package: its own `doc.go`, tests, README, and `Example_` test entries. Packages whose real dependency isn't deployed until later milestones still ship in M1 with (a) constructor functions returning no-op clients when the dependency is absent, and (b) `// TODO(Mx): …` markers keyed to the milestone that wires them up.

Coverage ≥80% per package from day one. Mutation testing (`gremlins`) runs but threshold is advisory in M1, enforced in M8.

## 7. Service scaffold + hello service

**Scaffold tool** at `impl/helixgitpx/tools/scaffold/` — a Go binary (no Python cookiecutter runtime) that renders `text/template` files from `templates/` into `services/<name>/`. Inputs: service name, proto package, default ports. Single binary, no external deps.

**Generated service layout:**

```
services/<name>/
├── cmd/<name>/main.go              thin entrypoint
├── internal/
│   ├── app/                        composition root
│   ├── handler/grpc/               gRPC handler impls
│   ├── handler/http/               Gin REST handler impls
│   ├── domain/                     business logic, no framework imports
│   ├── repo/                       persistence adapters
│   └── wire.go                     explicit DI (hand-wired)
├── api/                            OpenAPI generated from proto
├── migrations/                     goose SQL
├── config/{config.yaml, config.schema.json}
├── deploy/
│   ├── Dockerfile                  distroless, multi-stage, non-root, tini, SBOM label
│   ├── Dockerfile.dev              with delve
│   ├── helm/                       minimal chart
│   └── skaffold.yaml
├── test/{integration/, load/k6/}
├── .golangci.yml                   include from root
├── Makefile  go.mod  README.md  CHANGELOG.md
```

**Hello service** (`services/hello/`) is the first generation:
- gRPC: `HelloService.SayHello(name) → greeting` on `:9001`
- REST: `GET /v1/hello?name=X` on `:8001`
- `/healthz`, `/readyz`, `/metrics` on `:8081`
- Persists hit-counter to Postgres, caches last greeting per name in Redis, emits `hello.said` to Kafka.

Hello is the spine: every shared-lib package either is or becomes exercised by it.

## 8. Local compose stack

Compose file at `impl/helixgitpx-platform/compose/compose.yml`. Services:

| Service | Image | Purpose | Ports |
|---|---|---|---|
| `postgres` | `postgres:16-alpine` | per-service schemas; hello uses schema `hello` | 5432 |
| `kafka` | `apache/kafka:3.8` (KRaft) | event bus | 9092 |
| `redis` | `redis:7-alpine` | cache | 6379 |
| `karapace` | `ghcr.io/aiven-open/karapace:latest` | Schema Registry facade | 8081 |
| `jaeger` | `jaegertracing/all-in-one:latest` | OTel receiver + UI | 4317, 16686 |
| `prometheus` | `prom/prometheus:latest` | scrapes `/metrics` | 9090 |
| `grafana` | `grafana/grafana:latest` | dashboards (provisioned) | 3000 |
| `bufstream` | `bufbuild/bufstream:latest` | self-hosted BSR (profile `registry`) | 8082 |
| `hello` | built from `services/hello/Dockerfile` | spine service | 8001, 9001, 8081 |

Profiles: `core`, `observability`, `registry`, `all`. Default `make dev` uses `all`. Named volumes for Postgres/Kafka/Redis/Prometheus; `tmpfs` for Kafka metadata in CI profile.

**Runtime detection.** Every Makefile target calls `impl/helixgitpx-platform/compose/bin/compose`, never `docker` or `podman` directly. The wrapper prefers `docker compose`, falls back to `podman compose`, then `podman-compose`, and errors clearly if none exist.

**kind/k3d scripts** (item #13) ship at `impl/helixgitpx-platform/k8s-local/`: `up.sh`, `down.sh`, `load-images.sh`, `k3d-config.yaml`, `kind-config.yaml`. Documented as the "full-platform rehearsal" path; not required for M1 exit. A `Tiltfile` at repo root supports that path when chosen.

## 9. GitHub Actions workflows (all `workflow_dispatch`)

**Umbrella-repo workflows** (`.github/workflows/`):

| File | Purpose | Inputs |
|---|---|---|
| `ci-go.yml` | golangci-lint, gofumpt, `go test -race` + coverage, gremlins, build | `service`, `run-mutation` |
| `ci-web.yml` | Nx lint + test + build + bundle-size + Lighthouse; Playwright e2e gated | `affected-only`, `run-e2e` |
| `ci-clients.yml` | Gradle lint + detekt + ktlint + test; Compose preview screenshot tests | `target` |
| `ci-docs.yml` | Docusaurus build + markdownlint + vale + link-check | — |
| `ci-platform.yml` | helm-lint, kustomize build, `terraform fmt/validate/plan`, `kyverno test`, checkov, conftest | `environment` |
| `security-scan.yml` | Semgrep + Gitleaks + gosec + Trivy (fs+image) + Snyk + Grype + CodeQL | `deep` |
| `supply-chain.yml` | SBOM (syft → CycloneDX + SPDX), Cosign sign, SLSA provenance, in-toto attestations | `ref` |
| `release.yml` | tag → changelog → multi-arch image build+push+sign → Helm publish → GitHub Release | `version`, `dry-run` |
| `deploy.yml` | Argo CD sync trigger for env; uses Vault OIDC for short-lived kubeconfig | `environment`, `service` |

**Reusable workflows** (`.github/workflows/_reusable/`): `_setup-go.yml`, `_setup-node.yml`, `_setup-jvm.yml`, `_scan.yml`, `_build-image.yml`, `_push-image.yml`, `_cosign-sign.yml`, `_sbom.yml`, `_vault-creds.yml`.

**Per-subdir workflows** under `impl/<subdir>/.github/workflows/ci.yml` call the matching reusable. So if a subdir is ever split into its own repo, its CI works standalone.

**Secret-gated skips.** Workflows reference `SNYK_TOKEN`, `SONAR_TOKEN`, `SONAR_HOST_URL`, `VAULT_ADDR`, `COSIGN_PASSWORD`, `GHCR_TOKEN`. Each step `if: ${{ env.XXX != '' }}`; otherwise logs "skipped: <SECRET> not set" and exits success. Full catalog ships from day one; scanners light up as secrets are added.

**Self-hosted Kata runners** (item #9): config at `impl/helixgitpx-platform/github-actions-runner-controller/`. Workflows use `runs-on: ubuntu-latest` by default with a commented alternative; flipping is trivial when the cluster exists.

**OIDC → Vault** (item #10): `_vault-creds.yml` uses `hashicorp/vault-action@v3` with `id-token: write`. Terraform module at `impl/helixgitpx-platform/vault/terraform/` provisions the OIDC auth method + roles.

**CODEOWNERS + branch protection config** ships at repo root and in each subdir. Solo-operation deviation from the 2-approver rule is documented in `SOLO-NOTES.md` at the repository root.

## 10. Proto root & BSR strategy

Proto root at `impl/helixgitpx/proto/helixgitpx/<domain>/v1/`. Domains: `common`, `hello`, `auth`, `repo`, `sync`, `conflict`, `upstream`, `collab`, `events`, `platform`. M1 populates `common` and `hello` fully; the rest are stubs (empty `service`/`message` declarations) with `// TODO(Mx)` markers.

Existing `.proto` files under `docs/.../17-protos/` are normalized into the new layout. The spec tree keeps those files unchanged as historical source of truth; `proto/README.md` points at the live location.

**Buf config:**
- `buf.yaml` — module name `buf.build/helixgitpx/core`, deps: googleapis, grpc-ecosystem, buf.validate
- `buf.gen.yaml` — plugins: `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-connect-go`, `@bufbuild/protoc-gen-es`, `@connectrpc/protoc-gen-connect-es`, `protoc-gen-kotlin`, `connect-kotlin`, `protoc-gen-swift`, `connect-swift`, `protoc-gen-openapi`, `protoc-gen-doc`
- `buf.lock` — pinned deps
- `buf breaking` runs in CI against an empty baseline in M1; becomes the real baseline for M2 onward.

**Codegen destinations** (all committed; `make gen` is idempotent; CI verifies `git diff --exit-code`):
- Go: `impl/helixgitpx/gen/go/`
- Connect-Go: same tree
- TS (`@bufbuild/protoc-gen-es` + `@connectrpc/protoc-gen-connect-es`): `impl/helixgitpx-web/libs/proto/`
- Kotlin: `impl/helixgitpx-clients/shared/src/commonMain/kotlin/gen/`
- Swift: `impl/helixgitpx-clients/iosApp/Gen/`
- OpenAPI: `impl/helixgitpx/api/openapi/`

**BSR instance (item #6):** hosted BSR at `buf.build/helixgitpx/core` (free tier). Self-hosted `bufbuild/bufstream` lives in the compose file under profile `registry` for on-prem rehearsal.

## 11. Docs baseline

**Docusaurus site** at `impl/helixgitpx-docs/`. Sidebar categories mirror the spec's `00-core` through `18-manifests`. A `sync-docs.mjs` script run by `make docs` copies content from `docs/specifications/.../HelixGitpx/` into `impl/helixgitpx-docs/docs/`, rewriting cross-links to Docusaurus routes. The spec directory stays the single source of truth.

**OpenAPI & protobuf reference pages** are rendered via Redocly (OpenAPI) and `protoc-gen-doc` (protobuf). Both appear under the "Reference" sidebar.

**ADR template** is the existing `docs/.../15-reference/adr-template.md`. M1 seeds five ADRs capturing the locked constraints:

- **ADR-0001** — GitHub Actions with `workflow_dispatch`-only triggers (C-2)
- **ADR-0002** — Portable container runtime (docker or podman) via compose wrapper (C-3)
- **ADR-0003** — Single git history with 5 logical subdirs under `impl/` (C-1)
- **ADR-0004** — `mise` for toolchain management (C-5)
- **ADR-0005** — Spine-first M1 sequencing with completion matrix (C-6)

**Runbook template** at `docs/.../12-operations/runbooks/_TEMPLATE.md` using the headings of the existing 9 runbooks. No new runbooks in M1; the nine existing ones are in place already.

**Spec suite publish** (item #15) — a root `README.md` routes newcomers either to the spec or to the Docusaurus site. Root `CHANGELOG.md` follows keepachangelog.

## 12. Error handling

- `platform/errors` defines `Error{Code, Domain, Message, Cause, Details}` where `Code` is `codes.Code`. REST mapping to RFC 7807 problem details.
- No `panic` in library code. Recovery middleware in both `grpc` and `gin` converts panics to `codes.Internal` with sanitized message + logged stack.
- Every external dependency wrapped with typed sentinel errors (`ErrDBUnavailable`, `ErrKafkaBrokerDown`, `ErrRedisTimeout`) so callers can `errors.Is` without string matching.
- Context cancellation always respected; no library call swallows `ctx.Err()`.

## 13. Testing strategy

| Layer | Tool | M1 ships |
|---|---|---|
| Unit (Go) | stdlib `testing` + testify + gomock | ≥80% coverage per `platform/*`; hello logic 100% |
| Integration (Go) | testcontainers-go via `platform/testkit` | hello E2E: pg+kafka+redis roundtrip |
| Contract (proto) | `buf breaking` | Green; empty baseline locked |
| Fuzz | Go 1.23 native fuzz | One seeded target per `platform/{errors,config}` |
| Mutation | gremlins | Wired; threshold advisory in M1 |
| Web unit | Jest via Nx | Hello page component test + 1 snapshot |
| Web e2e | Playwright | Placeholder gated behind `run-e2e` input |
| KMP | JUnit5 + Kotest | One shared test compiling on JVM + Android + iOSX64 + desktop |
| Helm | `helm unittest` | Hello chart templates |
| Terraform | `terraform validate` + tflint | Platform module |
| Kyverno | `kyverno test` | All policies under `impl/helixgitpx-platform/kyverno/policies/` |
| Security | Semgrep, Gitleaks, Trivy, gosec, CodeQL | `security-scan.yml`; non-blocking in M1 |
| Docs | Vale + markdownlint + link-check | Blocks `ci-docs.yml` |

Every Go package is TDD — tests authored before code within each PR.

## 14. Completion matrix (the "nothing skipped" audit)

| # | Spec item | Artifacts produced in M1 | Status gate |
|---|---|---|---|
| 1 | `helixgitpx` monorepo layout | `impl/helixgitpx/{go.work,platform/,services/hello/,proto/,tools/scaffold/,Makefile}` | `go build ./...` green |
| 2 | `helixgitpx-web`, `-clients`, `-platform`, `-docs` | All 4 subdirs scaffolded with build system + README + CI | Each subdir's `make lint` green |
| 3 | `helix-platform` shared Go libs (14 packages) | `impl/helixgitpx/platform/{log,telemetry,errors,config,grpc,gin,kafka,pg,redis,temporal,spire,opa,health,testkit}` | Each has godoc + tests; coverage ≥80% |
| 4 | Nx workspace + Gradle convention plugins | `impl/helixgitpx-web/nx.json`, `impl/helixgitpx-clients/buildSrc/` | `nx run-many -t lint,test` and `./gradlew check` green |
| 5 | Service template | `impl/helixgitpx/tools/scaffold/` + `templates/` + `services/hello/` generated | `go run ./tools/scaffold --dry-run` produces expected file list |
| 6 | BSR instance + proto root | `impl/helixgitpx/proto/{buf.yaml,buf.gen.yaml,buf.lock,helixgitpx/}` + BSR module pushed + self-hosted bufstream in compose profile `registry` | `buf lint && buf build` green; BSR module resolvable |
| 7 | CI pipelines (lint, format, test, SBOM, Cosign, Snyk, SonarQube, Semgrep, Gitleaks, CodeQL) | 9 workflows + `_reusable/` — all `workflow_dispatch` only | Each runs to completion (skipping steps missing secrets) |
| 8 | Kyverno policy bundle + Checkov for IaC | `impl/helixgitpx-platform/kyverno/policies/*.yaml` (disallow-privileged, require-labels, require-signed-images, enforce-resource-limits) + `.checkov.yml` | `kyverno test` green; `checkov -d impl/helixgitpx-platform/` green |
| 9 | Self-hosted runner pool on Kata | `impl/helixgitpx-platform/github-actions-runner-controller/{values.yaml, runner-scale-set.yaml, kata-runtimeclass.yaml}` | Manifests pass `helm template` and kube-linter |
| 10 | OIDC → Vault | `_reusable/_vault-creds.yml` + `impl/helixgitpx-platform/vault/{terraform/*, policies/*.hcl, oidc-role.json}` + activation guide | `terraform plan` on vault module succeeds |
| 11 | `mise` bootstrap | `mise.toml` pinning: go 1.23, node 20, java 21, gradle 8.10, protoc 28, buf 1.45, sqlc 1.27, kind 0.24, k3d 5.7, kubectl 1.31, helm 3.16, tilt 0.33, skaffold 2.13, cosign 2.4, syft 1.14, grype 0.81, kyverno-cli 1.13, checkov 3 | `mise install` succeeds; `mise doctor` green |
| 12 | Tilt/Skaffold | `Tiltfile` + `services/hello/deploy/skaffold.yaml` + root `skaffold.yaml` | `skaffold diagnose` green; `tilt ci` runs to ready |
| 13 | kind/k3d scripts | `impl/helixgitpx-platform/k8s-local/{up.sh,down.sh,load-images.sh,k3d-config.yaml,kind-config.yaml}` | Scripts pass `shellcheck`; `up.sh --dry-run` prints planned steps without error; `kubectl kustomize` works on bundled configs. Actually spinning up a cluster is optional in M1 (C-3) and validated in M2 when the real platform cluster is provisioned. |
| 14 | Postgres/Kafka/Redis compose | `impl/helixgitpx-platform/compose/compose.yml` + `compose/bin/compose` wrapper | `compose up -d` healthy; hello connects to all three |
| 15 | Publish documentation suite | Repo-root `README.md` routing + `docs/README.md` index + Docusaurus build of the existing spec | `npm run build` green in `helixgitpx-docs` |
| 16 | ADR registry template | Existing `adr-template.md` normalized + five seeded ADRs (0001–0005) | Links valid; registry listable |
| 17 | Runbook template | `docs/.../12-operations/runbooks/_TEMPLATE.md` + conformance check script | `make runbook-lint` green |
| 18 | Docusaurus API docs site | `impl/helixgitpx-docs/` full site + OpenAPI plugin + proto docs plugin + `sync-docs.mjs` + `ci-docs.yml` | `npm run build` green; link-check + vale green |

## 15. Exit criteria

1. Fresh checkout: `mise install && make bootstrap` completes without manual steps.
2. `make dev` brings up the compose stack, detects container runtime, and `curl http://localhost:8001/v1/hello?name=world` returns `{"greeting":"hello, world"}`.
3. gRPC equivalent: `grpcurl -plaintext localhost:9001 helixgitpx.hello.v1.HelloService/SayHello -d '{"name":"world"}'` returns the same greeting.
4. Every workflow in `.github/workflows/*.yml` and per-subdir workflows runs to completion when manually triggered; non-blocking scanners skip gracefully when their secrets are absent.
5. `buf lint`, `buf build`, `go build ./...`, `go test ./...`, `nx run-many -t lint,test,build`, `./gradlew check`, `npm --prefix impl/helixgitpx-docs run build`, `helm lint`, `kubectl kustomize`, `terraform validate`, `kyverno test`, `checkov -d impl/helixgitpx-platform/` — all exit zero.
6. Completion-matrix rows 1–18 are all at "gate passing"; any row not passing blocks M1 close.
7. ADR-0001 through ADR-0005 are merged and cross-linked from the relevant docs.
8. `make docs` builds the Docusaurus site; link-check and vale green; site served locally on `:3001`.

## 16. Risks & mitigations

| Risk | Mitigation |
|---|---|
| KMP + iOS toolchain on Linux host — iOS targets need macOS for final compile. | Compile KMP shared module to JVM + Android + iOSX64 + linuxX64 + desktop on Linux. iOS target placeholder with `expect`/`actual` stubs; full iOS compile deferred to a macOS CI runner in M6. Shipped config is iOS-ready. |
| Podman compose lacks feature parity with docker compose (healthchecks, profiles, volumes). | Use the subset supported by both (profiles yes; `depends_on: condition: service_healthy` yes in recent Podman); detect podman version in the wrapper and warn on <4.7. |
| Buf codegen committed to git balloons the repo. | Codegen is small in M1 (only `common` + `hello`). Use `.gitattributes` to mark `gen/**` as `linguist-generated=true` and exclude from diffs in the UI. Review policy: codegen PRs auto-approved if `make gen` diff matches. |
| Kata runners + GitHub ARC require a real cluster. | M1 ships config artifacts only; they activate in M2 after the cluster exists. M1 workflows use `ubuntu-latest`. |
| Self-hosted BSR / Vault / SonarQube require running servers. | Same pattern — config ships, secret-gated workflow steps skip gracefully until the servers exist. |
| "Nothing skipped" pressure leads to thin/broken artifacts. | Every row in the completion matrix has a concrete status gate (a command that either passes or fails). Nothing counts as "done" until the gate passes. |
| Solo operation vs. 2-approver branch protection in spec. | `SOLO-NOTES.md` documents the deviation; branch protection config ships but is not enforced on this repo until team size permits. |

## 17. Open questions

None. All blocking decisions are locked in §4.

## 18. References

- `docs/specifications/main/main_implementation_material/HelixGitpx/13-roadmap/17-milestones.md` §§1.1–1.4 (scope source)
- `docs/specifications/main/main_implementation_material/HelixGitpx/02-services/03-microservices-catalog.md` (service module layout)
- `docs/specifications/main/main_implementation_material/HelixGitpx/01-architecture/02-system-architecture.md` (architecture overview)
- `docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr-template.md` (ADR template)
- `docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/20-developer-guide.md` (dev onboarding reference)
- `docs/specifications/main/main_implementation_material/HelixGitpx/17-protos/*.proto` (proto source)
- `docs/specifications/main/main_implementation_material/HelixGitpx/16-schemas/*.sql` (SQL schema source)
- `docs/specifications/main/main_implementation_material/HelixGitpx/18-manifests/` (manifest source)

— End of M1 Foundation design —
