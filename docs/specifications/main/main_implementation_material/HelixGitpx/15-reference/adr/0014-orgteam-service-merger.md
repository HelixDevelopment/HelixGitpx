# ADR-0014 — `orgteam-service` merges `org-service` and `team-service`

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

The roadmap §4 lists `org-service` and `team-service` as separate components. Orgs and teams are tightly coupled: every team has an `org_id` parent, cascading deletes on org removal must delete descendant teams and memberships, role resolution walks across both tables, and every mutating RPC lands in the same transactional boundary.

## Decision

Merge org + team into a single `orgteam-service` binary that exposes two gRPC services (`OrgService` and `TeamService`) against schemas `org` and `team`. A new `orgteam_svc` role (cross-schema in `sql/schemas.sql`) owns both.

## Consequences

- One deployment, one migration chain, one OPA bundle consumer, one audit outbox — half the platform boilerplate.
- Future split is a `git subtree` away if a specific team-service feature needs independent scaling (none foreseen for M3-M5).
- Completion-matrix rows 44 (Org CRUD) and 45 (Nested teams) both point at `services/orgteam` — the artifact location differs from the spec, the behaviour is identical.

## Alternatives considered

- Separate services with shared schema: rejected — cross-schema transactions + shared audit outbox become painful.
- Separate services with separate schemas: rejected — cascading deletes across services require distributed transactions.

## Links

- `docs/superpowers/specs/2026-04-20-m3-identity-orgs-design.md` §4 C-2, §8
