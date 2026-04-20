# ADR-0008 — Hello emits Kafka events via transactional outbox + Debezium

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M1 hello wrote to Postgres and Kafka as separate operations. A crash between the two could leave the counter incremented but the event unemitted (lost message), or vice versa (phantom event). At the scale of a demo this is invisible; the spec's M5 sync orchestrator cannot tolerate either.

## Decision

Hello writes events to a `hello.outbox_events` table inside the same pgx transaction as the counter UPSERT. Debezium's PostgreSQL connector streams the WAL via the `pgoutput` plugin and uses the `EventRouter` SMT to publish each row to the topic named in the row's `topic` column. Result: exactly-once from the service's perspective; duplicates only possible on Kafka-side retries (handled by consumer idempotency).

## Consequences

- `platform/kafka.Producer` is no longer needed for hello's happy path. The package remains for services that do need synchronous production.
- The outbox table requires a logical replication slot; `wal_level=logical` is set in CNPG config.
- Debezium tasks.max=1 per connector — sufficient for one service, will need slot-per-service tuning in M4/M5.
- The `hello.said` topic contract is preserved; downstream consumers see the same JSON payload.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-3, §10
