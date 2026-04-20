# M5 Federation & Conflict Engine — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED (auto-approved per session pattern) |
| Author | Милош Васић + Claude (2026-04-20) |
| Milestone | M5 — Federation & Conflict Engine (Weeks 21–28) |
| Scope | Full 23-item roadmap §6 (items 70–92) |

---

## 1. Context

M1–M4 delivered: foundation + data plane + identity & orgs + git ingress & adapter pool (3 providers). Four tags: `m1-foundation`, `m2-core-data-plane`, `m3-identity-orgs`, `m4-git-ingress`. 35 Argo CD Applications reconciling a local k3d cluster.

M5 is where the project's core promise — federated Git proxying with conflict resolution — comes online. After M5, a repo bound to multiple upstream forges survives divergence correctly.

## 2. Locked constraints

| ID | Constraint |
|---|---|
| C-1 | Temporal in-cluster with dedicated CNPG (`helix-temporal-pg`) |
| C-2 | WASM runtime = `tetratelabs/wazero` (pure Go, no CGO) |
| C-3 | CRDT = `automerge/automerge-go` v2 |
| C-4 | Live events = Connect streaming (gRPC server-streaming) with WS + SSE fallback via Connect-Web |
| C-5 | 3-way merge sandbox = ephemeral K8s Jobs (seccomp, egress=none) |
| C-6 | Remaining 9 providers = new directories under `adapter-pool/internal/providers/<name>/`, same `Adapter` interface |
| Inherited | M1–M4 |

## 3. New services

- **sync-orchestrator** — Temporal worker. Workflows: `FanOutPush(repo_id, ref_update)` + `InboundReconcile(webhook)`. DLQ topic `sync.dlq`. Retries exponential, max 5.
- **conflict-resolver** — Kafka Streams (via Goka) consumer. Detects ref divergence, loads per-repo Rego policy from `kv/conflict/<repo_id>/policy.rego`, launches merge Job, publishes `conflict.resolved`.
- **collab-service** — Automerge CRDT store for labels/milestones/assignees/issue body. `collab.crdt_ops` table outboxed to `collab.events` topic.
- **live-events-service** — Connect streaming `LiveEvents.Subscribe`. Backend = Redis Streams. Resume token = stream-id.

## 4. New / extended components

- **WASM plugin host** — new package `adapter-pool/internal/plugin/` using wazero. Loads plugins from `kv/plugins/<name>`, dispatches `Adapter` through a stable ABI. Example plugin in TinyGo under `services/adapter-pool/examples/plugin-hello/`.
- **9 new providers** under `adapter-pool/internal/providers/`: `gitee`, `gitflic`, `gitverse`, `bitbucket`, `forgejo`, `sourcehut`, `azuredevops`, `awscodecommit`, `generic_git`.

## 5. New Kafka topics

- `sync.dlq` (7d)
- `conflict.resolved` (infinite)
- `collab.events` (7d)

## 6. New schemas

```sql
CREATE SCHEMA IF NOT EXISTS sync;
-- (conflict, collab schemas already exist from M2 schemas.sql)

-- sync.sync_runs
CREATE TABLE sync.sync_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id TEXT NOT NULL,
    repo_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL,
    attempts INT NOT NULL DEFAULT 0
);

-- conflict.resolutions
CREATE TABLE conflict.resolutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID NOT NULL,
    ref TEXT NOT NULL,
    chosen_sha TEXT NOT NULL,
    policy_verdict JSONB NOT NULL,
    resolved_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- collab.crdt_ops
CREATE TABLE collab.crdt_ops (
    repo_id UUID NOT NULL,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    actor TEXT NOT NULL,
    seq BIGINT NOT NULL,
    op_bytes BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (repo_id, aggregate_type, aggregate_id, seq)
);
```

## 7. Exit criteria

1. verify-m5-cluster: 23/23 PASS.
2. Conflict scenario: repo bound to 5 upstream mirrors; diverging concurrent pushes to three of them → conflict-resolver detects → OPA policy picks winner → merge in sandbox → `conflict.resolved` event → streamed to web via live-events → visible in Grafana `audit.events` dashboard.

## 8. ADRs 0021–0025

- 0021 — Temporal in-cluster
- 0022 — wazero for WASM
- 0023 — Automerge-go v2 for CRDT
- 0024 — Merge sandbox via ephemeral Jobs
- 0025 — Connect streaming with WS+SSE fallback

## 9. References

- Roadmap §6, M1-M4 specs, Temporal helm chart
- https://github.com/tetratelabs/wazero
- https://github.com/automerge/automerge-go
- https://connectrpc.com/docs/web/supported-browsers-and-features/

— End of M5 design —
