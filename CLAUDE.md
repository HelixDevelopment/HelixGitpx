# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Authority order:** the [Constitution](./CONSTITUTION.md) is the supreme
policy document. Per-agent instructions are in [`AGENTS.md`](./AGENTS.md).
CLAUDE.md (this file) is a Claude-specific orientation guide — it does not
override the Constitution or AGENTS.md.

## Hard rules (from the Constitution)

1. **Testing** — every module carries tests in all seven required types:
   unit, integration, e2e, security, stress, ddos, benchmark. Only unit
   tests may use mocks, stubs, placeholder classes, or hardcoded data.
   No test may be skipped, disabled, broken, or flaky. Target coverage is
   100 % per type per module touched. See Article II of the Constitution.
2. **CI is workflow_dispatch-only.** No push / pull_request / schedule
   triggers are permitted. Mandatory.
3. **Container runtime is portable** — auto-detect docker vs podman; never
   hardcode either one.
4. **Upstream federation is regular** — every change on `main` is pushed to
   all configured upstreams (GitHub, GitLab, GitFlic, GitVerse, …) via
   `Upstreams/<target>.sh` scripts. Daily cadence minimum.
5. **Documentation is source** — every feature ships with HTML + PDF + ePub
   + Markdown + plain-text documentation under `docs/`.

## Repository Status

Implementation code lives under `impl/` (Go monorepo, Angular web app, KMP
clients, Helm/Argo/Kustomize platform manifests). The repo is GA-tagged
(`v1.0.0`, milestones `m1-foundation` through `m8-ga`). Do not claim code is
missing if it is present under `impl/`.

## What's Here

- `README.md` — one-line project description (`HelixGitpx — Helix Git Proxy eXtended`).
- `CONSTITUTION.md` — **supreme** policy doc. Read it.
- `AGENTS.md` — agent-specific rules. Read it.
- `Upstreams/` — executable bash scripts (`GitHub.sh`, `GitLab.sh`, `GitFlic.sh`, `GitVerse.sh`), each exporting `UPSTREAMABLE_REPOSITORY` to a different Git host. Per Constitution Article IV §2, every change on `main` is pushed to all of them.
- `impl/helixgitpx/` — Go monorepo (platform + 18 services + gen + tools/scaffold).
- `impl/helixgitpx-web/` — Angular 19 + Nx web app.
- `impl/helixgitpx-clients/` — KMP + Compose shells (Android/iOS/Desktop).
- `impl/helixgitpx-platform/` — Helm charts, Argo apps, Kustomize overlays, SQL, OPA.
- `impl/helixgitpx-docs-site/` — Docusaurus public docs.
- `docs/specifications/main/Git_Proxy_Master_Specification.md` + `.PDF` — the prior master spec (v4.0.0), **superseded** by the suite under `main_implementation_material/HelixGitpx/` but kept for provenance.
- `docs/specifications/main/main_implementation_material/HelixGitpx/` — the authoritative, implementation-ready documentation suite. Also shipped as `.zip`/`.7z` alongside — keep these in sync if the directory changes.

## Working With the Specification Suite

The suite is organized into numbered sections (`00-core` through `18-manifests`). Start at `main_implementation_material/HelixGitpx/README.md` — it's a routing table by role (Architect, Backend Engineer, SRE, etc.) and the canonical index.

Key entry points when answering questions:
- **Scope / what HelixGitpx is** → `00-core/01-vision-scope-constraints.md`
- **Architecture overview (C4 L1-L3)** → `01-architecture/02-system-architecture.md`
- **Services catalog (18 core + 7 platform)** → `02-services/03-microservices-catalog.md`
- **Machine-readable contracts** → `16-schemas/*.sql` (Postgres DDL), `17-protos/*.proto` (gRPC), `18-manifests/` (Helm, Kustomize, Terraform, Argo)
- **Developer onboarding (for the future implementation)** → `12-operations/20-developer-guide.md`

The suite cross-references heavily — most docs link to ADRs in `01-architecture/adr/` and `15-reference/adr-index.md`. When editing a spec doc, check whether sibling docs reference the sections you changed.

Sections marked `[VERIFY-AT-INTEGRATION]` intentionally contain facts that drift (pricing, cloud quotas, API versions); treat these as TODOs for when implementation reaches that area, not as errors.

## Planned Technology Stack (for context when discussing implementation)

From the spec — not yet present in this repo:
- Backend: Go 1.23+ with Gin, gRPC, Kafka + Schema Registry, PostgreSQL 16 + Timescale, Redis/Dragonfly, OpenSearch/Meilisearch/Qdrant.
- Web: Angular 19 + NgRx + Tailwind + Nx.
- Mobile/Desktop: Kotlin Multiplatform + Compose Multiplatform (Android/iOS/Win/macOS/Linux).
- Platform: Kubernetes 1.31, Istio (Ambient), Argo CD, Temporal.io, SPIFFE/SPIRE, OpenTelemetry.
- Code-gen pipeline: `buf generate` (proto → Go/TS/Kotlin/Swift), `sqlc`.

## Conventions That Apply Now

From `CONTRIBUTING.md` (enforced on all changes, including docs):
- Branches: `feat/…`, `fix/…`, `docs/…`, `chore/…`.
- Commits: Conventional Commits, signed, with `Signed-off-by:` (DCO). License is Apache-2.0 (code) / CC-BY-SA-4.0 (docs).
- PRs to `main` require two approvals; keep diffs focused (<400 lines ideal).

## Editing Specifications

- The suite considers itself **authoritative** — changes may cascade. Before editing a doc, check the role-based routing table in `main_implementation_material/HelixGitpx/README.md` to understand who consumes it.
- If you change a proto, SQL schema, or manifest under `16-schemas/` / `17-protos/` / `18-manifests/`, also update the prose doc that describes it (listed in the root README's index).
- The suite uses semver per artifact. Bumping a protobuf or SQL schema version is a public-API change — flag it.
- Do not edit the old `Git_Proxy_Master_Specification.md` — it's superseded and preserved for history.
