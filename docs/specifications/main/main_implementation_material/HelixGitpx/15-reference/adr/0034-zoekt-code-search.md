# ADR-0034 — Zoekt for code search

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Full-text search exists (Meilisearch), vector search (Qdrant), filter/aggregate (OpenSearch). Code search with regex + positional index is a different problem; options are Zoekt (Google, proven at scale) or Sourcegraph (heavier, adds a full search stack).

## Decision

Zoekt. A dedicated indexer CronJob writes Zoekt shards to a PVC; search-service fans queries to Zoekt for code-shaped queries and to the other engines for non-code.

## Consequences

- Zoekt's regex engine is fast enough for the expected scale (< 10M files).
- Indexing is CPU-heavy but runs as a CronJob — doesn't block the search path.
- Reuses existing infrastructure (MinIO-backed shards for durability).

## Links

- Spec §LOCKED C-5
- https://github.com/sourcegraph/zoekt
