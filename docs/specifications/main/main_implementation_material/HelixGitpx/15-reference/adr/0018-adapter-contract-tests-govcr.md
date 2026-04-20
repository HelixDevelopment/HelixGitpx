# ADR-0018 — Adapter contract tests via go-vcr cassettes

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

adapter-pool's three provider implementations (GitHub/GitLab/Gitea) need tests that verify the provider API contract without hitting real APIs on every CI run (rate limits, flakiness, credential exposure).

## Decision

Each provider package ships a `testdata/*.yaml` directory of go-vcr cassettes (recorded HTTP request/response pairs). `go test` replays the cassettes; a one-time `RECORD=1 go test` recording against real credentials re-records them (documented in each provider's README).

## Consequences

- Deterministic CI with no live API dependency.
- Cassettes go stale when providers change API shape; quarterly re-record cadence planned.
- Sensitive headers (tokens) are scrubbed by go-vcr hooks before checking in.

## Links

- `docs/superpowers/specs/2026-04-20-m4-git-ingress-adapter-pool-design.md` §4 C-2, §5.3
- https://github.com/dnaeon/go-vcr
