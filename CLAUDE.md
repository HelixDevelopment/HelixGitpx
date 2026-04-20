# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Status

**This repo currently contains specifications and documentation only — no implementation code has been checked in yet.** The `.gitignore` is Go-oriented in anticipation of the planned backend, but no `go.mod`, services, or build system exists at the root. Do not fabricate build/test commands; if asked to build or run something, first confirm with the user where the implementation lives (it may be in a separate repo).

## What's Here

- `README.md` — one-line project description (`HelixGitpx — Helix Git Proxy eXtended`).
- `Upstreams/` — four executable bash scripts (`GitHub.sh`, `GitLab.sh`, `GitFlic.sh`, `GitVerse.sh`), each exporting `UPSTREAMABLE_REPOSITORY` to a different Git host. These configure this repo's own multi-upstream mirroring (the project practices the federation pattern it specifies). When adding a new mirror target, create a matching script; don't hardcode remotes elsewhere.
- `docs/specifications/main/Git_Proxy_Master_Specification.md` + `.PDF` — the prior master spec (v4.0.0), **superseded** by the suite under `main_implementation_material/HelixGitpx/` but kept for provenance.
- `docs/specifications/main/main_implementation_material/HelixGitpx/` — the authoritative, implementation-ready documentation suite (v1.0.0, `APPROVED`). Also shipped as `.zip`/`.7z` alongside — keep these in sync if the directory changes.

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
