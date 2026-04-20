# 05b — Search & Indexing (OpenSearch / Meilisearch / Qdrant)

> **Document purpose**: Specify **which search store owns which query pattern**, how indexes are built (projectors), and how we keep them consistent with the event log.

---

## 1. Three-Store Strategy

No single store is best for all needs. We deliberately split:

| Store | Role | Why |
|---|---|---|
| **OpenSearch** | Logs, audit long-retention, power-user search | Mature, Apache-2.0 (vs ES SSPL), tiered storage, k-NN for hybrid |
| **Meilisearch** | Primary user-facing search (repos, issues, PRs, code symbols) | Lightning-fast, typo-tolerant, ranking out of the box, trivial to operate |
| **Qdrant** | Vector search (semantic, code embeddings, RAG) | Pure-Rust, high recall, payload filtering, HNSW + scalar quantisation |

All three are populated by **dedicated projectors** that consume from Kafka. No service writes to two stores directly.

---

## 2. OpenSearch

### 2.1 Cluster

- 3 dedicated masters + 2 hot data nodes + 2 warm + cold S3-searchable (plugin `repository-s3`).
- Index template with ILM: hot (7 d) → warm (30 d) → cold (1 y) → frozen (7 y).
- Snapshot to S3 nightly.

### 2.2 Index Templates

| Index pattern | Source | Rotation | Fields |
|---|---|---|---|
| `logs-app-*` | Vector → Loki → OpenSearch mirror | daily | service, level, msg, trace_id |
| `audit-*` | `helixgitpx.audit.events` | monthly | actor, action, resource, outcome |
| `metrics-*` | OTel gauge/counter (optional; Prom is primary) | daily | |
| `repo-*` (search) | Populated in Meili primarily; OpenSearch only for power-user / admin | — | — |
| `codesearch-*` | Extracted tokens from commits/blobs (language-aware) | by repo | path, content, sha, lang |

### 2.3 Code Search

- **Zoekt** integration (Sourcegraph's open source fast code search) for per-repo code search, federated with OpenSearch for free-text.
- Tokeniser: tree-sitter language grammars produce symbol-aware tokens.
- Incremental indexing on every `ref.updated` event.

---

## 3. Meilisearch

### 3.1 Deployment

- Single-master with 2 read replicas (Meilisearch replication is read-only; writes go to master).
- Per-tenant index namespace: `<tenant>_<index>` with **multi-tenant tokens** (API key scoped to tenant filter).

### 3.2 Indexes

| Index | Source events | Searchable fields | Filterable | Sortable |
|---|---|---|---|---|
| `repos` | `repo.created`, `repo.updated` | name, slug, description, topics | org_id, visibility, language | updated_at, stars |
| `prs` | `pr.*` | title, body, author | repo_id, state, labels | updated_at |
| `issues` | `issue.*`, `comment.*` | title, body | repo_id, state, labels | updated_at |
| `releases` | `release.published` | name, tag, body | repo_id | published_at |
| `users` | `auth.events` (user.created/updated) | username, display_name | — | joined_at |

### 3.3 Ranking Rules

Customised per index. Example for `repos`:

```
[
  "words",
  "typo",
  "proximity",
  "attribute",
  "sort",
  "exactness",
  "stars:desc",
  "updated_at:desc"
]
```

### 3.4 Projector

Go service `search-projector-meili`:
- Consumes relevant events.
- Transforms to Meili document.
- Batches updates (100 docs or 1 s, whichever first).
- Uses `meilisearch-go` client.
- Idempotent (document key = event's aggregate id).

---

## 4. Qdrant

### 4.1 Deployment

- 3-node cluster (Raft-based).
- HNSW index (m=16, ef_construct=100, ef=64 at search).
- Scalar quantisation (int8) for memory efficiency.
- Per-collection TTL where relevant.

### 4.2 Collections

| Collection | Vector size | Distance | Payload | Source |
|---|---|---|---|---|
| `code_embeddings` | 768 | Cosine | repo_id, file_path, sha, language | embedder on `ref.updated` |
| `issue_embeddings` | 768 | Cosine | repo_id, issue_id, state | on `issue.*` |
| `pr_embeddings` | 768 | Cosine | repo_id, pr_id, state | on `pr.*` |
| `conflict_embeddings` | 768 | Cosine | repo_id, case_id, kind | on `conflict.detected` |
| `chat_memory` | 1536 | Cosine | user_id, org_id | chatops interactions |

### 4.3 Embedding Model

- Default: **BGE-M3** (multilingual, 1024d truncated to 768).
- Provided by `ai-embedder` sidecar (FastAPI + sentence-transformers).
- GPU-accelerated where available; CPU fallback acceptable for lower RPS.
- Embedding version (e.g. `bge-m3-v1`) stored as payload metadata; re-embedding is possible with migration.

### 4.4 Hybrid Search

For best recall, we combine Meilisearch (lexical) + Qdrant (semantic):
- Query runs against both; results are merged with **Reciprocal Rank Fusion**.
- Configurable weights per org.
- Useful for "find similar PRs", "find related issues", code-NL search.

---

## 5. Projectors (shared design)

```
┌─────────────┐   ┌──────────────────┐   ┌────────────┐
│  Kafka      │→→│ Projector Service │→→│ Search Store│
│ (event log) │   │ (idempotent)      │   │  (OS/M/Q) │
└─────────────┘   └──────────────────┘   └────────────┘
                          │
                          ↓
                  ┌───────────────┐
                  │ Cursor table  │  (offset checkpoint)
                  └───────────────┘
```

Every projector:
- Owns a Postgres table `proj_cursor(consumer, topic, partition, offset, updated_at)`.
- Commits cursor **after** successful write to the target store (at-least-once).
- Idempotent by key (upsert).
- Exposes `/replay?from=offset` to rebuild the index from scratch.
- Emits `search.projector.lag_seconds{consumer}` metric.

---

## 6. Hot-Restart / Re-Index

To rebuild any index from zero:

1. Drop target index.
2. Reset projector cursor to `earliest`.
3. Scale projector up (≈ partition count).
4. Monitor lag until caught up.
5. Run data-quality asserts (counts match repo-service ground truth within tolerance).

This is practised **monthly** in staging.

---

## 7. Failure Modes

| Failure | Detection | Mitigation |
|---|---|---|
| Meilisearch unavailable | healthcheck + projector backoff | buffer in Redis list, flush on recovery |
| OpenSearch disk full | shard allocation warning | ILM auto-migrate to warm; on-call page |
| Qdrant out-of-memory | telemetry | scalar-quantise + split collection |
| Projector lag spike | `lag_seconds` alert | scale out; inspect DLQ |
| Embedding model change | version mismatch | dual-write period; cutover |

---

## 8. Security

- Meilisearch multi-tenant tokens per user (filter `org_id in (…)`).
- OpenSearch: index-pattern-based role mapping (via OpenSearch Security plugin).
- Qdrant: JWT-based payload filter; every query forced to include `org_id` filter via proxy.
- All stores only reachable via service mesh; no public ingress.

---

*— End of Search & Indexing —*
