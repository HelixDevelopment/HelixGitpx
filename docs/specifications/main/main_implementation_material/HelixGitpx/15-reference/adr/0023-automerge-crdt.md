# ADR-0023 — Automerge-go v2 for CRDT metadata

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Labels, milestones, assignees, and issue bodies can be edited concurrently across multiple upstreams. We need a conflict-free replicated data type (CRDT) that lets each upstream apply local changes and merge deterministically on sync.

## Decision

`automerge/automerge-go` v2. One Automerge document per `(repo_id, aggregate_type, aggregate_id)`. Each change is an op appended to `collab.crdt_ops` (part of the primary key includes `seq`); the service reconstructs state by replaying ops on load.

## Consequences

- Automerge handles all merge conflicts deterministically — no custom resolver for metadata.
- Binary op representation (not JSON) — logs show op counts, not contents.
- Automerge v2 API is stable but less mature than CRDT libraries for other languages; acceptable for M5 where the alternative is hand-rolled OT.

## Links

- `docs/superpowers/specs/2026-04-20-m5-federation-conflict-engine-design.md` §2 C-3
- https://github.com/automerge/automerge-go
