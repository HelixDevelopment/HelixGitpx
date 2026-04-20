# Service Template & Scaffold Catalog

> Quick-start templates for common scenarios inside the HelixGitpx codebase. Copying from a good starting point saves hours and keeps services consistent.

---

## 1. New Go Microservice

Run from repo root:

```bash
make scaffold-service NAME=my-new-service BOUNDED_CONTEXT=repo
```

The `cookiecutter-helix-go` template creates:

```
services/my-new-service/
├── cmd/
│   └── my-new-service/
│       └── main.go               # wires deps; calls fx.New + App.Run
├── internal/
│   ├── app/
│   │   └── app.go                # Fx module wiring
│   ├── domain/
│   │   ├── types.go              # pure domain entities
│   │   └── errors.go             # typed domain errors via helix-platform/errors
│   ├── service/
│   │   └── service.go            # use-case orchestration
│   ├── repo/
│   │   └── repo.go               # persistence interface + sqlc impl
│   └── api/
│       ├── grpc.go               # gRPC handlers
│       ├── http.go               # Gin handlers where needed
│       └── middleware.go
├── api/                           # auto-generated from proto/
├── db/
│   ├── migrations/                # goose files
│   └── queries/                   # sqlc .sql files
├── test/
│   ├── fixtures/
│   ├── integration/
│   └── e2e/
├── deploy/
│   ├── Chart.yaml
│   └── values.yaml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

Generated files include:

- Full dependency injection via **fx** from shared `helix-platform/appx`.
- Pre-wired OTel (traces, metrics, logs), health, readiness, liveness endpoints.
- Postgres, Kafka, Redis, SPIRE SVID clients ready to use.
- gRPC + gRPC-gateway boilerplate.
- CI workflow (`.github/workflows/my-new-service.yml`) inheriting from `shared-ci.yml`.
- Helm chart inheriting from `helix-platform/chart-base`.

---

## 2. New Proto Contract

```bash
make scaffold-proto PACKAGE=helixgitpx.v1 NAME=widget
```

Creates `proto/helixgitpx/v1/widget.proto` with:

- Standard options (go_package, java_package, csharp_namespace).
- Imports of common.proto.
- Skeleton `WidgetService` + a CRUD set.
- Registered in `buf.yaml` and `buf.gen.yaml`.

After editing, run `buf generate` to produce Go/TS/Kotlin/Swift bindings.

---

## 3. New Event Type

```bash
make scaffold-event TOPIC=helixgitpx.widget.events TYPE=widget.created
```

- Registers subject in Karapace via idempotent script.
- Adds Avro record to `16-schemas/events-avro.json`.
- Updates producer helper in `helix-platform/events/envelope`.
- Adds consumer group template in affected services.

---

## 4. New Angular Feature

Nx generators are first-class:

```bash
pnpm nx g @helixgitpx/angular:feature \
  --name widget \
  --domain repo \
  --routing \
  --store signal-store
```

Produces:

- `libs/features/widget/`:
  - `feature.routes.ts`
  - `feature.component.ts` (standalone, OnPush)
  - `feature.store.ts` (NgRx Signal Store)
  - `feature.spec.ts`
  - `feature.a11y.spec.ts` (axe)
  - `feature.e2e.spec.ts` (Playwright)

And adds the route into `apps/web/src/app/app.routes.ts`.

---

## 5. New KMP Screen

```bash
cd clients && ./gradlew :tools:scaffold:kmp-screen \
  -Pname=Widget \
  -PparentComponent=RepoDetails
```

Produces:

- `shared/ui/src/commonMain/kotlin/.../widget/WidgetScreen.kt`
- `shared/store/src/commonMain/kotlin/.../widget/WidgetComponent.kt`
- Android, iOS, Desktop test files.
- Integration points in the parent Decompose component.

---

## 6. New WASM Plugin

```bash
helixctl plugin scaffold \
  --name com.example.my-plugin \
  --shape git-adapter \
  --language rust \
  --output ./my-plugin
```

Produces a Cargo workspace with the WIT world, lib.rs skeleton, tests, GitHub Actions release pipeline using Cosign keyless, and a `helixgitpx-plugin.toml` manifest.

---

## 7. New Helm Service Chart

```bash
./tools/scaffold-helm-chart my-new-service
```

Based on `helix-platform/chart-base` with env-specific overlays pre-populated.

---

## 8. Fixtures for Local Dev

```bash
helixctl seed \
  --orgs=5 \
  --repos=50 \
  --users=20 \
  --prs=200 \
  --issues=500 \
  --conflicts=30
```

Idempotent; safe to re-run. Uses fixture data so tests are stable.

---

## 9. Observability Starter

Every scaffolded service includes starter dashboards and alerts committed to `observability/my-new-service/`:

- `dashboard.json` (Grafana).
- `alerts.yaml` (Prometheus rule).
- `runbook.md` (in `12-operations/runbooks/`).

Plumbs the service's RED metrics and links them to the service's page in Grafana.

---

## 10. Third-Party Integration Catalog

Official integrations built on top of HelixGitpx APIs:

| Integration | Purpose | Source |
|---|---|---|
| **GitHub Action `helixgitpx-status`** | Report CI status to HelixGitpx | marketplace |
| **GitLab CI template** | Same, for GitLab | public snippets |
| **Jenkins plugin** | Native plugin | plugin portal |
| **VSCode extension** | In-editor PR review, live events, conflict view | marketplace |
| **IntelliJ / JetBrains plugin** | Same for JetBrains IDEs | marketplace |
| **Slack / Teams / Discord bots** | Notifications + slash commands | respective marketplaces |
| **Terraform provider** | orgs/repos/upstreams/branch-protection as code | registry |
| **Pulumi provider** | Same | registry |
| **Kubernetes Operator** | CRDs for repo-definitions-as-code | GitHub |
| **Backstage plugin** | Service catalog integration | Backstage plugins |

Each integration has its own repo under `vasic-digital/helixgitpx-<name>` with semver lifecycle and an SDK pin it targets.

---

## 11. Contribution: Adding to the Catalog

- Scaffold templates live under `tools/scaffold/` — PRs welcome to improve.
- New integrations: create a repo, open a "Add to catalog" PR in `helixgitpx-docs` updating this file.
- Must declare supported HelixGitpx major versions, licence, maintenance owner.
- Listed only if tests run against the current stable API.

---

*— End of Scaffold / Integration Catalog —*
