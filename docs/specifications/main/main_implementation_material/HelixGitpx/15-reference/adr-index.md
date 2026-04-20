# Architecture Decision Records (ADR) Index

> Every significant architectural decision is captured as an immutable ADR. ADRs live in `docs/adr/`. Once **Accepted**, an ADR is never edited — only **Superseded** by a newer one that links back.

---

## Status Definitions

- **Proposed** — under discussion.
- **Accepted** — adopted; implementation may proceed.
- **Superseded by ADR-NNNN** — no longer the current answer.
- **Deprecated** — scheduled for removal; migration in progress.
- **Rejected** — considered and declined; kept for historical context.

---

## ADR-0001 — Microservices over Modular Monolith

**Status**: Accepted · 2026-02-01

**Context**: HelixGitpx must scale to 100 k orgs, 10 M repos, 1 B events/day, across multiple regions, with independently scalable hot-spots (AI inference, adapter pool, live events).

**Decision**: Adopt a microservices architecture with ~18 core services and ~7 platform services (see [02-services/03-microservices-catalog.md](../02-services/03-microservices-catalog.md)). Shared libraries (`helix-platform`) absorb cross-cutting concerns to avoid polyglot re-implementation.

**Consequences**:
- ✓ Independent deploy/scale per service.
- ✓ Fault isolation.
- ✗ Operational complexity (mitigated by shared libs + Kubernetes + Argo CD).
- ✗ Distributed-system debugging (mitigated by OTel end-to-end).

**Alternatives considered**: Modular monolith (rejected — scale + team topology demands independent deploy); serverless (rejected — stateful adapters + streaming workloads map poorly).

---

## ADR-0002 — Apache Kafka over NATS JetStream for Event Backbone

**Status**: Accepted · 2026-02-03

**Context**: Need durable, replayable, compacted event log; CDC from Postgres; cross-region mirror; rich stream-processing ecosystem.

**Decision**: Kafka 3.8+ in KRaft mode managed by Strimzi, with **Karapace** (Apache-2.0) as schema registry and **Debezium** for CDC.

**Consequences**:
- ✓ Proven at scale; mature tooling (kafka-streams, ksqlDB, MM2).
- ✓ Tiered storage to object store.
- ✗ Higher operational overhead than NATS (accepted).

**Rejected alternatives**: NATS JetStream (less mature for event sourcing + CDC); Redpanda (smaller ecosystem at time of decision).

---

## ADR-0003 — OpenSearch over Elasticsearch

**Status**: Accepted · 2026-02-04

**Context**: We need a searchable log/audit store and need license clarity.

**Decision**: OpenSearch (Apache-2.0 fork) for logs, audit, and k-NN hybrid search.

**Consequences**:
- ✓ Free to self-host without commercial license concerns.
- ✓ Active community.
- ✗ Plugin ecosystem slightly smaller than ES; acceptable.

**Rejected**: Elasticsearch SSPL (license ambiguity for multi-tenant SaaS).

---

## ADR-0004 — Meilisearch for User-Facing Search

**Status**: Accepted · 2026-02-04

**Context**: User-facing search needs sub-50 ms response, typo-tolerance, per-tenant indexes, operational simplicity.

**Decision**: Meilisearch for `repos`, `prs`, `issues`, `releases`, `users` indexes.

**Consequences**:
- ✓ Excellent ranking OOTB.
- ✓ Cheap to operate.
- ✗ Less flexible than OpenSearch for analytic queries; OK because we use OpenSearch for those.

---

## ADR-0005 — Qdrant for Vector Database

**Status**: Accepted · 2026-02-04

**Context**: Need vector DB for semantic search, code embeddings, chat memory.

**Decision**: Qdrant — Rust, HNSW + scalar quantisation, payload filtering, Raft clustering.

**Rejected**: Pgvector (rejected for scale at 10+ M vectors); Milvus (heavier ops); Weaviate (GraphQL-centric).

---

## ADR-0006 — Kotlin Multiplatform + Compose Multiplatform for Mobile & Desktop

**Status**: Accepted · 2026-02-07

**Context**: Need one codebase spanning Android, iOS, Windows, macOS, Linux, with shared UI and business logic at enterprise quality.

**Decision**: KMP (Kotlin 2.0 K2) + Compose Multiplatform 1.7+.

**Rejected alternatives**: Flutter (Dart ecosystem; weaker JVM interop); React Native (JS runtime overhead; weak desktop story); native per-platform (5× engineering cost).

---

## ADR-0007 — Temporal.io for Durable Workflows

**Status**: Accepted · 2026-02-10

**Context**: Multi-step operations (sync, conflict apply, repo migration) must survive restarts with deterministic replay.

**Decision**: Temporal workflows for all durable, long-running, multi-step operations. Activities wrap side-effecting work (adapter calls, DB writes).

**Consequences**:
- ✓ Survives pod crashes; automatic retry; human-friendly workflow UI.
- ✓ History provides audit and replay.
- ✗ New concept for devs; training needed.

---

## ADR-0008 — Event Sourcing for Repository Aggregate

**Status**: Accepted · 2026-02-11

**Context**: Repository state must be reconstructible from an append-only log for audit, time-travel, and cross-region replay.

**Decision**: Use Kafka compacted topic `helixgitpx.repo.events` as event store for the Repo aggregate. Postgres `repo.repositories` is a projection.

**Consequences**:
- ✓ Full history; replay; strong audit.
- ✗ Snapshotting needed for fast reads (implemented in 7-day cadence).

---

## ADR-0009 — CQRS with Per-Read-Model Projectors

**Status**: Accepted · 2026-02-11

**Context**: Different read models need different storage (PG, Meilisearch, Qdrant, OpenSearch).

**Decision**: Command side writes to PG + Kafka via outbox. Read models populated by idempotent Kafka-consuming projectors.

**Consequences**:
- ✓ Each read model optimised independently.
- ✗ Eventual consistency (acceptable; sub-second typical).

---

## ADR-0010 — SPIFFE/SPIRE for Workload Identity

**Status**: Accepted · 2026-02-12

**Context**: Every pod needs short-lived, rotated, cryptographically verifiable identity for zero-trust.

**Decision**: SPIRE server per cluster, SPIRE agents per node, X.509 SVIDs rotated hourly; used for gRPC mTLS, Kafka SASL/OAUTHBEARER, Vault auth.

---

## ADR-0011 — Istio Ambient over Sidecar Mesh

**Status**: Accepted · 2026-02-13

**Context**: Sidecars double pod count and complicate upgrades; ambient decouples L4 (ztunnel) from L7 (waypoint).

**Decision**: Adopt Istio Ambient for mTLS and traffic policy.

**Consequences**:
- ✓ Lower per-pod resource overhead.
- ✓ Upgrades decoupled from app lifecycle.
- ✗ Newer tech (tracked via periodic reviews).

---

## ADR-0012 — OPA for AuthZ

**Status**: Accepted · 2026-02-13

**Context**: Authorization rules cross-cut every service; policies must be diff-reviewed, versioned, and deployed independently of code.

**Decision**: Open Policy Agent (Rego) with bundle server pulling from Git. Each service invokes in-process OPA library for low-latency decisions.

---

## ADR-0013 — Self-Hosted Ollama + vLLM Default

**Status**: Accepted · 2026-02-14

**Context**: Customer code must never leave the cluster by default; we need robust local inference.

**Decision**: Default stack is **Ollama** (quick-start) and **vLLM** (high-throughput) self-hosted. Cloud providers are opt-in per org.

---

## ADR-0014 — Argo CD GitOps

**Status**: Accepted · 2026-02-14

**Context**: Production state must be reconstructible from Git; changes must be auditable and reversible.

**Decision**: Argo CD reconciling from `helixgitpx-platform` repo. ApplicationSets generate per-environment Applications.

---

## ADR-0015 — WASM Plugins for Custom Adapters

**Status**: Accepted · 2026-02-15

**Context**: Some enterprises need proprietary provider adapters or need to extend behaviour without forking.

**Decision**: Plugin model via Wasmtime + WIT (component model). Plugins signed with Cosign; loaded dynamically; sandboxed with explicit host imports.

**Rejected**: Native shared libraries (ABI instability); external sidecars (deployment overhead); JVM plugins (heavy runtime).

---

## ADR-0016 — Gin Gonic for HTTP/REST Layer

**Status**: Accepted · 2026-02-15

**Context**: User mandate; REST is generated from protos via grpc-gateway, but we still want a native Go framework for custom middleware and classic REST where gateway isn't ideal.

**Decision**: Gin Gonic is the standard for any service with hand-written HTTP handlers. gRPC-gateway for transpiled REST. **Chi** fallback permitted for single-purpose utilities (internal-only).

---

## ADR-0017 — Debezium CDC Outbox Pattern

**Status**: Accepted · 2026-02-16

**Context**: We must avoid dual writes (DB + Kafka) to prevent inconsistencies.

**Decision**: Each service has an `event_outbox` table written in the same transaction as the domain row. Debezium tails logical replication and emits to Kafka.

---

## ADR-0018 — Karapace as Schema Registry

**Status**: Accepted · 2026-02-16

**Context**: We need an Avro / JSON Schema registry with strict compatibility rules and clear license.

**Decision**: Karapace (Apache-2.0) over Confluent Schema Registry (Confluent Community License, ambiguous for SaaS).

---

## ADR-0019 — UUIDv7 as Primary Key

**Status**: Accepted · 2026-02-18

**Context**: We need monotonic ordering for index locality, but also global uniqueness without coordination.

**Decision**: UUIDv7 (timestamp-prefixed) for all primary keys. Enables time-sortable indexes without a separate `created_at` sort.

---

## ADR-0020 — Connect Protocol over gRPC-Web

**Status**: Accepted · 2026-02-19

**Context**: Browsers can't speak raw gRPC; gRPC-Web is one option; Connect is a newer, HTTP/1.1-compatible protocol that unifies gRPC, gRPC-Web, and Connect under one server-side handler.

**Decision**: Use **Connect-Go** on the server; clients negotiate automatically. gRPC for native; gRPC-Web as a shim; Connect JSON for simple scripts.

---

## ADR-0021 — Automerge for CRDT Metadata

**Status**: Accepted · 2026-02-21

**Context**: Concurrent metadata edits (labels, milestones, issue body) from multiple upstreams must converge without loss.

**Decision**: Automerge 2.x (Rust core via WASM bindings, also available in Go port) as the CRDT library for metadata sync.

---

## ADR-0022 — Cosign Keyless + Rekor

**Status**: Accepted · 2026-02-22

**Context**: Need to verify image authenticity without managing long-lived signing keys.

**Decision**: Cosign keyless (Sigstore Fulcio) signs at build time; Rekor transparency log provides public proof; Connaisseur / Kyverno verifies at admission.

---

## ADR-0023 — Angular 19 with NgRx Signal Store

**Status**: Accepted · 2026-02-23

**Context**: Frontend must be maintainable by a team of 4–8 engineers, highly performant, and scale to dozens of complex feature areas.

**Decision**: Angular 19 standalone components + signals, NgRx Signal Store per feature, Tailwind + shadcn patterns, Nx monorepo.

**Rejected**: React (team expertise in Angular; enterprise tooling); Vue (smaller ecosystem for tooling we need).

---

## ADR-0024 — Monorepo for Backend, Polyrepo at Boundaries

**Status**: Accepted · 2026-02-25

**Context**: Atomic refactors across services are frequent; plugins are third-party-friendly; clients (web/mobile) have different release cadences.

**Decision**: `helixgitpx` backend monorepo; `helixgitpx-web`, `helixgitpx-clients`, `helixgitpx-platform` separate; each plugin its own repo.

---

## ADR-0025 — Buf Schema Registry (Self-Hosted)

**Status**: Accepted · 2026-02-26

**Context**: Proto contracts are public API; we need versioning, linting, breaking-change detection, and generated SDK publishing.

**Decision**: Self-hosted Buf Schema Registry; CI runs `buf lint` and `buf breaking`; SDK artefacts published on every `main` merge.

---

## ADR-0026 — Vector for Log Shipping

**Status**: Accepted · 2026-02-27

**Context**: Fluentbit/Fluentd performance acceptable but configuration verbose; we want typed pipelines.

**Decision**: **Vector** (Rust) as the log shipper, fed by Kubernetes log directories, forwarding to Loki with OpenSearch mirror.

---

## Superseded / Historical

*(placeholder)* None yet.

---

## ADR Template

```markdown
# ADR-NNNN: Title

**Status**: Proposed | Accepted | Superseded by ADR-MMMM | Deprecated | Rejected
**Date**: YYYY-MM-DD
**Deciders**: @...
**Technical story**: Link to the originating issue/PR

## Context
What is the forcing function? What constraints matter?

## Decision
What are we doing?

## Consequences
### Positive
### Negative
### Risks & Mitigations

## Alternatives Considered
- **Alt A** — why not
- **Alt B** — why not

## References
- Related ADRs
- External docs
```
