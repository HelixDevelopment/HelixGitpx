# ADR-0029 — Self-hosted update feed via tus.io

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Desktop auto-update needs a resumable, append-only upload endpoint. Electron-style Squirrel servers, Sparkle RSS feeds, Tauri's updater — all work, all assume a server we control.

## Decision

`tus.io` (resumable upload protocol) running on MinIO behind ingress. Desktop clients poll an append-only manifest; clients download new binaries via signed URLs (same pattern as LFS from M4, ADR-0020).

## Consequences

- Same MinIO infra as LFS — no new service.
- Resumable uploads survive flaky networks during release publishing.
- Client update delta is a manifest diff, easy to understand.

## Links

- Spec §LOCKED C-4
