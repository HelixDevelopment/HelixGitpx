# ADR-0007 — HA manifests authoritative in Git; local overlay patches to single-replica

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

CNPG, Kafka, OpenSearch, Vault, Mimir, Loki, Tempo, ingress-nginx, and cert-manager are all HA-capable. On a local k3d host running the full stack simultaneously, true HA (3-replica everywhere) exceeds reasonable single-machine capacity. Abandoning HA in Git would make staging/prod bring-up require bespoke re-engineering per component.

## Decision

Every Helm chart's `values.yaml` declares HA-by-default replica counts (3 for most data services, 2 for control-plane components). A per-environment `values-local.yaml` patches replicas down to 1 for local. `kustomize/overlays/{local,staging,prod-eu}/` selects the right file via each Application's `spec.source.helm.valueFiles`.

## Consequences

- GitOps in staging/prod deploys real HA without code changes.
- Local dev fits in one machine without OOMs.
- PodDisruptionBudgets, replicaCount, and related configs appear in two values files per chart — a small duplication cost.
- Staging is where HA is first validated end-to-end; local is not a faithful HA environment.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-2
