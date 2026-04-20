# ADR-0020 — LFS via MinIO presigned URLs (not streaming through git-ingress)

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Git LFS objects (large binaries) can exceed several GiB. Streaming them through git-ingress doubles bandwidth cost and adds memory pressure; direct upload/download to object storage is the mature industry pattern.

## Decision

On an LFS upload/download request, git-ingress returns a MinIO (S3-compatible) presigned URL with a 5-minute TTL. The client uploads/downloads directly to MinIO; git-ingress only persists object metadata (oid, size) in `repo.lfs_objects`.

## Consequences

- Bandwidth decoupled from git-ingress replica count.
- Quota enforced at presign time (not at byte-stream time); user can exhaust their budget with a single 100 GiB upload — acceptable with daily budget caps and M5 conflict-resolver policy enforcement.
- MinIO per-repo IAM policies ensure clients can only access their own objects; policies provisioned by repo-service on repo Create.

## Links

- `docs/superpowers/specs/2026-04-20-m4-git-ingress-adapter-pool-design.md` §4 C-4, §5.2
- https://github.com/git-lfs/git-lfs/blob/main/docs/api/batch.md
