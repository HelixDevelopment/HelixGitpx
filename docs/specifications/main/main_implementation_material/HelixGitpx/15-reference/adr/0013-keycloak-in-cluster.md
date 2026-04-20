# ADR-0013 — Keycloak v26 in-cluster with auto-imported realm

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M3 requires an OIDC identity provider. Options ranged from Keycloak (spec-faithful, heavy) to Dex (lighter, spec deviation) to external SaaS (Okta/Auth0, costs money and adds account setup).

## Decision

Deploy Keycloak v26 in-cluster on its own CNPG Postgres (`helix-kc-pg`). The `helixgitpx` realm is auto-imported from `impl/helixgitpx-platform/helm/keycloak/realm/helixgitpx.json` on first start. Two test users (`user@helixgitpx.local`, `admin@helixgitpx.local`) are pre-provisioned for local use; production overlays strip them.

## Consequences

- Self-contained identity stack; no external SSO account required for M3.
- Keycloak + its Postgres costs ~500 MiB resident locally (within k3d budget).
- Realm updates are GitOps-driven: edit the JSON, Argo CD re-applies the ConfigMap, Keycloak picks up the change on next restart.
- Corporate-SSO brokering (the real M3 value-add for enterprises) is a realm configuration change, not a code change — lands when a customer asks.

## Alternatives considered

- Dex: lighter but the spec explicitly says Keycloak. Deferred.
- External SaaS: rejected for M3 because cloud accounts are out of scope.

## Links

- `docs/superpowers/specs/2026-04-20-m3-identity-orgs-design.md` §4 C-1, §6
