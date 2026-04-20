# ADR-0017 — go-git for parsing/policy + git-http-backend CGI for wire protocol

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

`git-ingress` accepts `git push`/`fetch` over the smart-HTTP protocol. Options: (a) pure go-git server; (b) shell out entirely to `git-http-backend`; (c) split — Go for policy, CGI for the wire.

## Decision

Option (c): go-git parses ref updates and signed-push signatures; `git-http-backend` (shipped in the git-ingress Alpine image) handles the actual smart-HTTP pack handshake via CGI over stdin/stdout.

## Consequences

- Battle-tested wire protocol; ergonomic Go surface for business logic.
- git-ingress Dockerfile switches from distroless to `alpine:3.20` (needs `/usr/libexec/git-core/`).
- Two code paths per push (Go pre-filter + CGI shell-out) — acceptable complexity.

## Links

- `docs/superpowers/specs/2026-04-20-m4-git-ingress-adapter-pool-design.md` §4 C-1
- https://git-scm.com/docs/git-http-backend
