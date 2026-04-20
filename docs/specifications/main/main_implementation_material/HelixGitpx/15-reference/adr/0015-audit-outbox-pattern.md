# ADR-0015 — Audit events use the transactional outbox pattern

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Every mutating RPC in auth-service and orgteam-service must emit an audit event. Naive Kafka production is not exactly-once: a crash between the business write (Postgres) and the Kafka produce would leave audit gaps or phantom events.

## Decision

Use the same outbox pattern M2 introduced for hello. Each service writes audit events to its local `<schema>.outbox_events` table inside the same pgx transaction as the business write. Debezium's PostgreSQL connector streams the WAL (pgoutput plugin) and uses the EventRouter SMT to publish each row to the topic named in its `topic` column — `audit.events` in this case.

## Consequences

- Exactly-once from the service's perspective; duplicates only possible on Kafka-side retries (handled by `audit-service` consumer idempotency via `audit.append_event`).
- Every new service with mutating endpoints adds: (a) an `outbox_events` table in its schema, (b) a KafkaConnector CR pointing at it, (c) calls to `platform/audit.Emitter.EmitInTx` from every mutating RPC.
- `hello.said` topic pattern remains the exemplar; `audit.events` follows the identical route.

## Alternatives considered

- Dedicated `audit-sdk` Kafka producer: rejected — same reliability gap as hello's original M1 naive emitter.
- Change-data-capture on the business table directly: rejected — couples audit format to table schema, brittle.

## Links

- `docs/superpowers/specs/2026-04-20-m3-identity-orgs-design.md` §4 C-3, §9
- https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html
