# ADR-0019 — Upstream credentials live only in Vault

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Every upstream binding needs a Git token or SSH key to authenticate against the remote forge. Options: (a) store encrypted in Postgres; (b) store in Vault KV; (c) hybrid (metadata in PG, secret in Vault).

## Decision

(b) — `upstream.upstreams` holds only `vault_path` (e.g. `kv/upstream/<id>`). The Vault KV value at that path carries `token`, optional `ssh_private_key`, and the `webhook_secret`. adapter-pool reads these via Vault Agent at request time; they never land in Postgres, backups, or log lines.

## Consequences

- Postgres dumps contain no credentials.
- Credential rotation is a Vault KV write; services read freshly on each RPC.
- Requires Vault to be reachable for every adapter call — acceptable because auth-service (M3) already depends on Vault.

## Links

- `docs/superpowers/specs/2026-04-20-m4-git-ingress-adapter-pool-design.md` §4 C-3, §5.5
