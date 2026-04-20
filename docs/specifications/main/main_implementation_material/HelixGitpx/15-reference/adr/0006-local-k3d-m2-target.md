# ADR-0006 — Local k3d as the M2 cluster target

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M2 requires a running Kubernetes cluster for real deployment (its exit criterion depends on metrics/logs/traces visible in Grafana). Staging GKE/EKS/AKS costs money and requires account setup; a cloud-less path keeps M2 self-contained on a single workstation.

## Decision

M2 uses `k3d` on the developer's host as the cluster target. `up.sh --m2` provisions a 3-node k3d cluster (1 server + 2 agents), installs Cilium as CNI, ingress-nginx as L7, and then applies the Argo CD bootstrap kustomization.

## Consequences

- Full spine deployable on a 62 GiB host without cloud credentials.
- Real resource pressure surfaces single-host limits early — preflight script refuses to run below 48 GiB free RAM.
- Staging/prod overlays still exist (`kustomize/overlays/{staging,prod-eu}/`) and are validated via `kubectl kustomize`; they activate when real clusters arrive.
- macOS and Windows hosts are explicitly out of scope for M2 (per the on-prem-deployment spec).

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-1, §5
