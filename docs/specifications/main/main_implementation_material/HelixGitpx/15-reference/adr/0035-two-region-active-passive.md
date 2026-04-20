# ADR-0035 — Two-region active-passive for GA

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

GA requires regional HA. Options: single-region + cross-AZ (not enough for region
loss), two-region active-passive (simple, lower cost, well-understood), two-region
active-active (lowest RTO, but requires conflict resolution across writes).

## Decision

Active-passive. eu-central-1 is primary; eu-west-2 is warm standby with MirrorMaker 2
(Kafka) and CNPG logical replication (Postgres). GeoDNS routes customers to the
nearest healthy region; on primary loss the DR runbook triggers failover.

## Consequences

- RTO target 5 min, RPO target 30s — meets SLO.
- Writes only land in active region; no cross-region conflict surface.
- Active-active revisit post-GA (year 2) once we have production data to size it.
- Running costs acceptable at GA (warm standby ~35% of primary).

## Links

- Spec §LOCKED C-6
- tools/dr/dr-drill-runbook.md
