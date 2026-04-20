# ADR-0001 — GitHub Actions with workflow_dispatch-only triggers

**Status**: Accepted
**Date**: 2026-04-20
**Deciders**: Милош Васић
**Supersedes**: —

---

## Context and Problem Statement

Project is federated across GitHub, GitLab, GitFlic, GitVerse. Every CI host has its own config format; automatic push/PR triggers on multiple hosts would multiply maintenance and produce noisy cross-host reruns.

## Decision Drivers

- Reduce CI configuration drift across four git-hosting platforms
- Minimize unintended CI executions and runner cost
- Simplify troubleshooting (single source of truth for CI behavior)

## Considered Options

- **A** GitHub Actions with `push:` and `pull_request:` triggers — simple but generates noise across all upstreams
- **B** GitHub Actions with `workflow_dispatch` only — manual invocation, zero automatic executions
- **C** Multi-host CI (GitHub + GitLab + GitFlic + GitVerse) — federated approach, prohibitive maintenance cost

## Decision

**We will adopt GitHub Actions with `workflow_dispatch`-only triggers (Option B).** Every workflow (`.github/workflows/*.yml`) declares `on: workflow_dispatch` as its sole trigger. No `push:`, `pull_request:`, or `schedule:` triggers are permitted. `workflow_call` for reusable workflows is allowed because it does not fire automatically.

## Consequences

### Positive

- Engineers invoke CI manually from the Actions tab (or `gh workflow run`).
- Contributors cannot accidentally burn runner minutes with unreviewed pushes.
- Single, explicit control point: easier to debug why a CI run did or did not trigger.

### Negative / Trade-offs

- Scheduled jobs (e.g., nightly scans) must be scheduled externally and call `workflow_dispatch` via the API.
- Integration test feedback is not automatic on every push; developers must remember to run CI.

### Risks & Mitigations

- Risk: developers forget to run CI before merging → Mitigation: branch-protection rule requires passing CI checks before merge
- Risk: scheduled jobs become stale → Mitigation: external scheduler (e.g., cron or GitHub's scheduled triggers via `schedule: [{cron: '...'}]` calling `workflow_dispatch`) documented in runbooks

## Decision Outcome

Option B was chosen to eliminate cross-platform CI sprawl and keep configuration simple while the team is one person. Option C (multi-host) was deferred pending team growth.

## Implementation Notes

- Enforced by repository lint: `scripts/ci-lint.sh` checks every `.github/workflows/*.yml` for absent `push:` and `pull_request:` keys.
- Developers use `gh workflow run <workflow-id>` or the Actions UI.
- Scheduled tasks (backups, security scans) use external triggers (Temporal, cron) that invoke `workflow_dispatch` via API.

## Validation & Exit Criteria

- No manual `push:` or `pull_request:` triggers in production workflows.
- CI lint passes in every commit (enforced by pre-commit hook).
- Zero unintended CI runs due to automatic triggers (monitored via `gh workflow run --list` audit).

## References

- `docs/superpowers/specs/2026-04-20-m1-foundation-design.md` §4 C-2
- `.github/workflows/` directory
