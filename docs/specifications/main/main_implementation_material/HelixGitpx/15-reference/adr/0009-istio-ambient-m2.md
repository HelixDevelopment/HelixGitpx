# ADR-0009 — Istio Ambient mesh is installed in M2

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

The spec lists Istio Ambient in Phase 2.1 alongside cluster provisioning. Deferring it saves resources but delays mesh-dependent features (zero-code mTLS between services, L7 policy for M3 auth).

## Decision

Istio Ambient is installed in M2. Ambient mode (zTunnel DaemonSet + namespace opt-in) is used — not sidecar. Namespaces `helix`, `helix-data`, `helix-cache`, `helix-secrets` are labeled `istio.io/dataplane-mode: ambient`. Istio consumes SPIRE SVIDs via the SPIFFE CSRA integration.

## Consequences

- Hello-to-Postgres and hello-to-Dragonfly connections gain automatic mTLS with zero code change.
- Observability and system namespaces opt out (ingress-nginx and Prometheus need plaintext reach to non-mesh pods).
- zTunnel DaemonSet + Istio CNI add ~3 pods and ~500 MiB resident on each node.
- M3 auth work can rely on L4 mTLS being already in place.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-4, §8
