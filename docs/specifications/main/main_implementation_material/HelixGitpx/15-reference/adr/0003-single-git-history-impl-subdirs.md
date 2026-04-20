# ADR-0003 — Single git history with 5 logical subdirs under impl/

**Status**: Accepted
**Date**: 2026-04-20
**Deciders**: Милош Васић

---

## Context and Problem Statement

The HelixGitpx spec describes 5 distinct sub-projects (Go monorepo, Angular web, KMP clients, GitOps/infra, Docusaurus). Splitting them into separate git repositories adds coordination overhead while the project is maintained by a single engineer. Keeping them in one repo simplifies cross-cutting changes but risks conflating unrelated concerns.

## Decision Drivers

- Single engineer maintaining the project (no distributed team coordination burden yet)
- Ability to land cross-cutting changes (e.g., proto schema affecting backend and clients) atomically
- Preserve option to split repos later via `git subtree` or `git filter-repo` when team size warrants

## Considered Options

- **A** Five separate git repos — clean separation, high coordination cost
- **B** Single monorepo with 5 subdirs, shared git history — simpler now, easy to split later
- **C** Monorepo with git submodules — unclear ownership, merge conflict pain

## Decision

**We will adopt Option B: a single monorepo with 5 logical subdirectories under `impl/`.** Implementation code lives in:
- `impl/helixgitpx/` — Go monorepo (backend services, shared platform packages)
- `impl/helixgitpx-web/` — Nx/Angular workspace (frontend)
- `impl/helixgitpx-clients/` — KMP + Compose Multiplatform (mobile/desktop)
- `impl/helixgitpx-platform/` — GitOps, infrastructure-as-code, local compose
- `impl/helixgitpx-docs/` — Docusaurus site

All sub-trees share one git history.

## Consequences

### Positive

- Cross-cutting changes (proto updates affecting Go + TypeScript + Kotlin) land in a single, holistic PR.
- One CI config tree (`.github/workflows/` at repo root), plus per-subdir optionality in `.github/workflows/ci-*.yml` for future split-readiness.
- Easy to extract with `git subtree split --prefix=impl/helixgitpx` or `git filter-repo` when team grows.

### Negative / Trade-offs

- Large diffs possible when multiple sub-projects change together (mitigated by disciplined PR boundaries).
- Every PR reviewer must understand multiple language ecosystems; future teams may need stricter CODEOWNERS rules.

### Risks & Mitigations

- Risk: one sub-project's CI failure blocks all PRs → Mitigation: each sub-project owns its CI config; failures are namespace-local
- Risk: git history becomes unwieldy (100k+ commits) → Mitigation: evaluate at M8 (year-2); if needed, split then

## Implementation Notes

- Each sub-project has its own language toolchain and lockfiles (`go.mod`, `package.json`, `gradle.toml`, `Cargo.toml`).
- Root `mise.toml` pins all runtimes (Go, Node, JVM, Kotlin, Rust).
- Root `.github/workflows/` coordinates overall CI; sub-projects have optional per-subdir `ci.yml` for faster feedback.

## Validation & Exit Criteria

- All five `impl/<subdir>` folders have their own build/test targets (e.g., `make -C impl/helixgitpx test`).
- One cross-cutting PR (e.g., proto schema change) lands successfully with reviews from at least one engineer familiar with each sub-project.

## References

- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-1, §5
- `impl/` directory structure
