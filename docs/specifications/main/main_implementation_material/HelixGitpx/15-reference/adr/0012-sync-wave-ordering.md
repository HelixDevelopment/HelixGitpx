# ADR-0012 — Argo CD sync-wave ordering for M2 bring-up

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Argo CD applies Applications in parallel by default. M2 has hard dependencies: operators must exist before their CRs, CNI before any workload, and the SPIRE trust bundle before Istio can pull it.

## Decision

Every Application carries an `argocd.argoproj.io/sync-wave` annotation:

| Wave | Components |
|---|---|
| -10 | cilium |
| -5 | ingress-nginx, cert-manager, external-dns, minio |
| -3 | prometheus-stack, mimir, loki, tempo, pyroscope |
| 0 | spire, istio-base, istio-ambient |
| 5 | cnpg-operator, strimzi-operator |
| 7 | cnpg-cluster, kafka-cluster, dragonfly, meilisearch, opensearch, qdrant, vault |
| 9 | karapace, debezium |
| 10 | hello |

Argo CD's sync-wave contract: Applications in lower waves are Synced + Healthy before higher-wave Applications begin reconciling.

## Consequences

- Bring-up is deterministic; no CRD-before-operator flapping.
- Sync time is linear in wave depth — can't parallelise across waves.
- Adding a new chart requires picking a wave; operator-to-CR dependency must be respected.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §5
