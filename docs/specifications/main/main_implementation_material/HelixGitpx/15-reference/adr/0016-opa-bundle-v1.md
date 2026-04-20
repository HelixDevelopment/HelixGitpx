# ADR-0016 — OPA bundle v1 is in-process + ConfigMap-delivered

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Access control must be enforced on every mutating RPC in orgteam-service (and M4+ services). Options ranged from hard-coded Go predicates (rigid), to OPA sidecar (latency + complexity), to in-process OPA evaluation via `platform/opa` (flexible, in-binary).

## Decision

Policy lives in `impl/helixgitpx-platform/opa/bundles/v1/authz.rego`. The `helm/opa-bundles/` chart ships the bundle as a ConfigMap at `helix-system/opa-bundle-v1`. Every service that needs authorization mounts the ConfigMap, loads the bundle at startup via `platform/opa.NewEvaluator` (M1-shipped), and re-loads on SIGHUP when the ConfigMap is updated by Argo CD.

## Consequences

- Zero network hop for every RBAC check.
- Policy changes are GitOps-driven: edit `authz.rego`, Argo CD updates the ConfigMap, services rolling-restart to pick it up.
- Versioning: bump `values.yaml:bundleVersion` and add a new `opa/bundles/vN/` folder; running services don't change their bundle until their pod restarts.
- OPA sidecar path remains an option for M8 scale-out: just swap the evaluator in `platform/opa` with a gRPC client.

## Alternatives considered

- Sidecar: adds 10-20 ms per call + a coord-deployment dependency. Not worth it at M3 scale.
- Hard-coded Go: rejected — every policy change becomes a service release.

## Links

- `docs/superpowers/specs/2026-04-20-m3-identity-orgs-design.md` §4 C-4, §8
- `impl/helixgitpx-platform/opa/bundles/v1/authz.rego`
- `impl/helixgitpx-platform/helm/opa-bundles/templates/configmap.yaml`
