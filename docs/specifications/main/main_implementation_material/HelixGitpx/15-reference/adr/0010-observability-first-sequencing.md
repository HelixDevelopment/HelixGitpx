# ADR-0010 — Observability-first sequencing for M2

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Strict spec-phase order (2.1 → 2.5) installs the full data plane before observability lands. Any issue during CNPG or Kafka bring-up is then debugged without metrics/logs/traces — blind.

## Decision

Sync waves install observability (wave -3) before the data plane (wave 5–7) and before hello (wave 10). Prometheus + Mimir + Loki + Tempo + Pyroscope + Grafana + Alertmanager are reconciling and scraping long before the first data service schedules. Every chart ships a `ServiceMonitor` so metrics flow automatically when pods appear.

## Consequences

- Each data service's "done" gate can assert "visible in Grafana" from the first minute.
- Cluster debugging during bring-up is tractable.
- Observability resource cost is paid upfront (~8–12 GiB) — offset by the halved replica counts in the `local` overlay.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-7, §5
