# ADR-0021 — Temporal in-cluster with dedicated CNPG Postgres

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

sync-orchestrator (M5) needs a durable workflow engine for `FanOutPush` + `InboundReconcile` orchestration. Options: cloud Temporal (costs money, needs internet), self-hosted via Helm (control, complexity), build-your-own state machine (not M5-feasible).

## Decision

Deploy Temporal via the official Helm chart in `helix-data` namespace. Back it with its own CNPG Postgres cluster `helix-temporal-pg` to isolate workflow state from application state. `platform/temporal.NewClient` lifts the M1 TODO(M5) to a real `go.temporal.io/sdk/client.Dial`.

## Consequences

- Battle-tested workflow primitives (retries, timers, signals, queries).
- ~500 MiB resident overhead for Temporal server + frontend + history + matching + worker pods.
- Temporal Web UI available at `temporal.helix.local` for debugging.
- M8 hardening can migrate to Temporal Cloud for staging/prod while keeping local on self-hosted.

## Links

- `docs/superpowers/specs/2026-04-20-m5-federation-conflict-engine-design.md` §2 C-1
- https://github.com/temporalio/helm-charts
