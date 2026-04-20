# M5 Federation & Conflict Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans.

**Goal:** Roadmap §6 items 70-92: Temporal-backed sync workflows, conflict detection+resolution in sandboxed merges, CRDT metadata, 9 additional providers + WASM plugin host, live-events streaming with WS/SSE fallback.

**Tech Stack:** Temporal SDK (go.temporal.io/sdk), tetratelabs/wazero, automerge/automerge-go v2, Goka (Kafka Streams), Connect streaming, TinyGo for example plugin.

**Constraints:** ADRs 0021-0025 per spec §8.

---

## Tasks (compact — covers all 23 items)

- **T1** — `sql/schemas.sql`: add `sync` schema; migration for `sync.sync_runs`, `conflict.resolutions`, `collab.crdt_ops`.
- **T2** — Temporal helm chart + Argo Application (wave 5, dedicated CNPG `helix-temporal-pg`).
- **T3** — 3 new Kafka topics (`sync.dlq`, `conflict.resolved`, `collab.events`).
- **T4** — sync-orchestrator scaffold + `FanOutPush` workflow stub + `InboundReconcile` workflow stub + DLQ producer.
- **T5** — conflict-resolver scaffold + ref-divergence detector (Goka consumer of `upstream.webhooks`) + `conflict.resolutions` repo.
- **T6** — collab-service scaffold + Automerge CRDT store + `collab.crdt_ops` outbox.
- **T7** — live-events-service scaffold + Connect streaming handler + Redis Streams backend.
- **T8** — 9 new adapter-pool providers (gitee/gitflic/gitverse/bitbucket/forgejo/sourcehut/azuredevops/awscodecommit/generic_git).
- **T9** — WASM plugin host package (`adapter-pool/internal/plugin/` with wazero) + example TinyGo plugin under `services/adapter-pool/examples/plugin-hello/`.
- **T10** — 4 new Helm charts (sync-orchestrator, conflict-resolver, collab-service, live-events-service).
- **T11** — 4 new Argo CD Applications at sync-wave 9 + Temporal app at wave 5 (5 apps total).
- **T12** — verify-m5-cluster.sh (23 gates) + verify-m5-spine.sh (multi-upstream conflict scenario).
- **T13** — ADRs 0021-0025.
- **T14** — `m5-federation` tag.

Each task follows the M4 pattern: scaffold via `tools/scaffold`, minimal implementation, TDD where domain logic exists, helm chart copied from hello, Argo CD Application CR at the documented sync-wave.

---

## Exit

verify-m5-cluster.sh PASS 23/23. M5 tag `m5-federation`.
