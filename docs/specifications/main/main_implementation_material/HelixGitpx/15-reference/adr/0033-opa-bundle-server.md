# ADR-0033 — Dedicated OPA bundle server

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M3 shipped OPA bundle v1 as a ConfigMap (ADR-0016). As policy surface grows (v2 with enforcement rules from M7), ConfigMap delivery becomes unwieldy for diff-review and non-atomic updates.

## Decision

Dedicated `opa-bundle-server` service that hosts bundles with ETags. Service OPA agents (in-process in orgteam/repo/etc.) pull via the `httpbundle` plugin. Rollback = switch the active bundle pointer.

## Consequences

- ConfigMap approach retired for production (kept only for M3 compatibility).
- Bundle diff-review lands as a CI workflow step (Task T5 in M7 plan).
- Per-tenant bundles possible (future work; out of M7 scope).

## Links

- Spec §LOCKED C-4
- Supersedes parts of ADR-0016
