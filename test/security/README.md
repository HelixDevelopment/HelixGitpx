# Security tests

NO mocks. NO stubs. Every assertion runs against a live, seeded stack.

## Suites

- `owasp-zap/` — baseline + API scans against staging.
- `nuclei/` — curated templates for HelixGitpx surfaces.
- `gosec/` — static Go analysis (per-service).
- `semgrep/` — SAST ruleset (Go + TS + Kotlin).
- `trivy/` — image + IaC scans.
- `custom/` — HelixGitpx-specific: upstream-federation token scope,
  OPA bypass attempts, webhook HMAC tampering, residency leak.

## Running

```bash
make test-security
```

Fails the build on any High/Critical finding without an approved
suppression (suppressions live at `tools/security/suppressions.yaml` and
require an ADR).

## Corpora

See:

- `tools/fuzz/corpora/webhook/` — webhook HMAC fuzz seeds.
- `tools/chaos/upstream-429-storm.yaml` — rate-limit exhaustion scenario.
- `docs/security/pentest-scope-2026q2.md` — external pen-test scope.
- `docs/security/bug-bounty-program.md` — bug bounty.

## Constraint

A security test that passes with mocks proves nothing. Constitution §II §2.
