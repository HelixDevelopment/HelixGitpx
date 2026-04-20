# ADR-0005 — Spine-first M1 sequencing with completion matrix

**Status**: Accepted
**Date**: 2026-04-20
**Deciders**: Милош Васић

---

## Context and Problem Statement

M1 Foundation has 18 numbered items across 4 phases. Strict phase-ordered execution leaves the project with no provably-working end-to-end flow until all 18 items land. A thin vertical slice — "hello service responding via gRPC + REST through the compose stack" — provides continuous validation if built first.

## Decision Drivers

- Early detection of cross-cutting risks (toolchain, protobuf codegen, podman vs docker, version drift)
- Continuous end-to-end validation (hello service deployable by mid-M1)
- Parallelizable Phase B work for when team scales
- Explicit completion audit trail (matrix as SSOT)

## Considered Options

- **A** Strict phase order (1.1 → 1.2 → 1.3 → 1.4) — sequential, no early validation
- **B** Spine-first (Phase A: thin vertical slice 1–22, Phase B: breadth 23–36) — early risk surfacing, parallelizable
- **C** Feature branches per item, merge when complete — loses audit trail

## Decision

**We will adopt Option B: spine-first sequencing with two phases.**

1. **Phase A — Spine (Tasks 1–22):** monorepo skeleton, 14 shared Go platform packages, proto root, scaffold tool, hello service with full pg+redis+kafka wiring, integration test, Dockerfile + skaffold.
   - End state: hello service reachable via HTTP and gRPC (subject to local compose runtime availability).

2. **Phase B — Breadth (Tasks 23–36):** 9 CI workflows + 4 reusables, Kyverno policies, Checkov config, ARC/Kata runner config, Vault OIDC terraform, k8s-local scripts + Tiltfile, Nx web scaffold, KMP clients scaffold, Docusaurus site, ADRs, runbook template, CODEOWNERS + branch protection + hello helm chart, verification script.

Each of the 18 roadmap items maps to an artifact in `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §14's completion matrix, which serves as the audit trail for "nothing skipped."

## Consequences

### Positive

- Early-risk items (toolchain pinning, protobuf codegen, podman vs docker) surface in Phase A where they can be diagnosed against a minimal surface.
- Phase B is parallelizable — its tasks touch disjoint file trees.
- Completion matrix's per-row status gate discourages "looks done" commits; each row is either passing or blocking M1 close.
- One engineer can parallelize using git branches; larger team can use multiple agents (Superpowers).

### Negative / Trade-offs

- Phase A implementation is tightly coupled (all-or-nothing); one blocker delays the whole spine.
- Two-phase approach adds mental overhead; must maintain alignment between phases.

### Risks & Mitigations

- Risk: Phase A blocker (e.g., protobuf codegen fails) → Mitigation: daily checkpoints during Phase A; unblock immediately
- Risk: Phase B parallelization creates merge conflicts → Mitigation: each Phase B task owns a subtree; CODEOWNERS prevents overlap

## Implementation Notes

- Phase A: commits tagged `phase-a`. One engineer, focused cadence (target: 7–10 days).
- Phase B: commits tagged `phase-b`. Parallelizable: 4 sub-agents per task group if team grows.
- Completion matrix (`docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §14): status per row updated after each commit.
- Verification script (`scripts/verify-m1.sh`): audit trail; runs after every commit to confirm no regressions.

## Validation & Exit Criteria

- Phase A: `make -C impl/helixgitpx test-hello` succeeds; `docker-compose up postgres` and hello service reachable.
- Phase B: every row in the completion matrix reports `ok` or a documented reason for `deferred`.
- `scripts/verify-m1.sh` returns exit code 0.

## References

- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-6, §14
- `docs/superpowers/plans/2026-04-20-m1-foundation.md`
- `scripts/verify-m1.sh` (implementation)
