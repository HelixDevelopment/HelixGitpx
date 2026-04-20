# ADR-0011 — Local-friendly external defaults (self-signed, noop DNS, MinIO, placeholder PagerDuty)

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Several M2 components normally depend on external services: cert-manager → Let's Encrypt + DNS-01 provider; external-dns → real DNS provider API; object storage → S3/GCS/Azure Blob; Alertmanager → PagerDuty. Requiring all four on a dev host makes M2 infeasible without cloud accounts.

## Decision

Each has a local-friendly default:
- cert-manager uses `selfsigned-ca` ClusterIssuer locally; Let's Encrypt in staging/prod.
- external-dns uses the `noop` webhook provider locally; real provider in staging.
- Object storage uses MinIO in-cluster; real S3 in staging.
- Alertmanager routes to a null receiver when `PAGERDUTY_INTEGRATION_KEY` Secret is absent.

All four are shipped as both configs (HA/production) in Git and patched down per environment in `kustomize/overlays/`.

## Consequences

- M2 is deployable offline after the initial image pull.
- Staging/prod activation is a single overlay swap, not a config rewrite.
- Local TLS uses self-signed certs; verifiers pass `-k` / `--insecure` (documented).
- MinIO is a single-instance local; no data-durability guarantees locally.

## Links

- `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md` §4 C-6, §7
