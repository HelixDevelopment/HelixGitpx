# Changelog

All notable changes to this repository are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning: per-artifact semver. This file tracks repo-level milestones.

## [1.0.0] — 2026-04-21 — General Availability

### Added

- 8 milestone tags landed (`m1-foundation` → `m8-ga`) and `v1.0.0` cut.
- Go monorepo: 18 services under `impl/helixgitpx/services/` plus
  `platform/` shared libraries and `gen/` proto outputs.
- Web app: Angular 19 + Nx under `impl/helixgitpx-web/`.
- Docs site: Docusaurus 3.10 under `impl/helixgitpx-docs-site/`.
- Marketing site: Astro + Tailwind under `impl/helixgitpx-website/`
  (19 pages shipped).
- KMP clients: shared module + androidApp + desktopApp + iosApp
  under `impl/helixgitpx-clients/`.
- Platform: 30+ Helm charts + 53 Argo CD Applications under
  `impl/helixgitpx-platform/`.
- Protobuf contracts for 20 packages (hello/auth/org/team/repo/upstream
  + collab/conflict/sync/events/ai/search/webhook/adapter/gitingress/
  policy/billing/audit + common/platform).

### Governance

- `CONSTITUTION.md` — supreme policy, 7-type test matrix mandate.
- `AGENTS.md` — rules for AI contributors.
- `CLAUDE.md` — Claude-specific orientation.
- 39 ADRs under `docs/specifications/.../15-reference/adr/`.

### Documentation

- 10 manuals scaffolded under `docs/manuals/src/` with 5 formats
  (HTML/PDF/ePub/DOCX/TXT) exported via `tools/docs-export/build-all.sh`.
- 21 video scripts under `docs/media/video-scripts/`.
- 8 operational runbooks under `docs/operations/runbooks/`.
- 4 security planning docs under `docs/security/`.
- 3 integration-plan docs under `docs/integrations/` including a
  policy review dropping 7 offensive modules.
- `docs/UNFINISHED.md` — living inventory of what's not yet done.

### Testing infrastructure

- Seven test-type directories (`test/{unit,integration,e2e,security,
  stress,ddos,benchmark,chaos}`) + runners that degrade gracefully
  when toolchains are absent.
- `scripts/verify-everything.sh` — 13-phase green-suite.
- Coverage audit, chart lint, Argo path validator, Rego syntax check.
- 4 Go fuzz targets (webhook HMAC, GitHub/GitLab/Gitea canonicalizer,
  git ref validation).
- 10 `Benchmark*` functions across platform/webhook, audit/merkle,
  search RRF fusion, sync retry, conflict classification.
- Playwright config + first e2e spec.
- k6 scenarios for stress/ddos/benchmark/perf-budgets.

### CI / CD

- 13 GitHub Actions workflows, all `workflow_dispatch`-only per
  Constitution §CI, all passing on `main`: `ci-verifiers`, `ci-go`,
  `ci-web`, `ci-docs`, `ci-platform`, `ci-clients`, `security-scan`,
  `supply-chain`, `mutation-testing`, `deploy`, `release`,
  `perf-budgets`, `upstream-sync`.
- 4 reusable callables under `.github/workflows/_reusable/`.
- GitLab pipeline suppressed (identity verification pending on
  gitlab.com); real config preserved in `.gitlab-ci.yml`.
- `scripts/push-to-all-upstreams.sh` — mirror to 4 upstreams.

### Known gaps (see `docs/UNFINISHED.md` for the full inventory)

- 7 of 17 services still at 17-line `app.Run` scaffolds.
- Web unit tests are smoke-only (Angular TestBed not wired under Jest).
- No hosted Service; no customers; no active certifications or
  bug-bounty program — the website pages say this honestly.
- OPA rego tests pass a Python-only syntax check (no `opa test` CLI
  available in CI).

## [0.0.0] — 2026-04-19

### Added

- Initial repository seeded with Git_Proxy_Master_Specification v4.0.0
  (now superseded by the authoritative spec suite).
