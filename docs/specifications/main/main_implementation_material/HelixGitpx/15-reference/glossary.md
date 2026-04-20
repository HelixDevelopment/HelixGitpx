# Glossary

> Canonical definitions. If a term is used ambiguously anywhere else in this suite, **this file wins**. Add terms here when introducing them.

---

### Adapter
A Go implementation (or WASM plugin) of the `UniversalGitAdapter` interface for a specific Git provider. Lives inside the `adapter-pool` service.

### Aggregate (DDD)
A consistency boundary containing a root entity and related objects. In HelixGitpx, the **Repository aggregate** is event-sourced.

### Apply Plan
An ordered sequence of atomic operations produced by the conflict resolver to enact a chosen resolution across upstreams.

### Binding
An (org, repo, upstream) triple with per-upstream settings (enabled, shadow, push/pull flags).

### BSR (Buf Schema Registry)
Self-hosted registry for `.proto` files with linting, breaking-change detection, and SDK publishing.

### Canary
A partial rollout (1 %, 10 %, 50 %…) of a new version with automated SLO-driven analysis before full promotion.

### CDC (Change Data Capture)
Postgres logical replication tailed by Debezium, translating row changes into Kafka events. The HelixGitpx outbox pattern relies on CDC.

### CMP / Compose Multiplatform
JetBrains' cross-platform UI toolkit; enables sharing the same `@Composable` code across Android, iOS, Desktop (JVM), and Web.

### Conflict Case
A persisted record of observable divergence between two or more upstream views of the same logical object (ref, label set, PR state, etc.).

### Connect
A protocol by Buf that unifies gRPC, gRPC-Web, and a JSON-over-HTTP style under one server-side handler.

### Consumer Group
Kafka consumer group; HelixGitpx names them `<service>-<purpose>` (e.g. `sync-orchestrator-main`).

### CQRS
Command Query Responsibility Segregation — write path and read path use different models. We use CQRS with event sourcing + projector-populated read models.

### CRDT
Conflict-free Replicated Data Type. HelixGitpx uses Automerge for metadata (labels, issue bodies) to guarantee loss-free concurrent edits.

### DLQ (Dead-Letter Queue)
Kafka topic holding events that exceeded retry limits. Suffix `.dlq`.

### DPO (Direct Preference Optimisation)
LLM fine-tuning technique using pairs of preferred/rejected outputs; more stable and compute-efficient than PPO.

### Envelope
Canonical Avro schema wrapping every Kafka event with correlation, trace, and tenant metadata.

### Event Sourcing
Storing domain state as a sequence of events rather than current-state snapshots. Repo aggregate is event-sourced on Kafka.

### Fan-out Push
Temporal workflow that replicates a local commit to every enabled upstream in parallel.

### Federation
Treating multiple Git hosts as equivalent views of the same repo; HelixGitpx federates bidirectionally.

### Forgejo
Gitea fork maintained by Codeberg's community; uses Gitea API — one adapter covers both.

### GitOps
Declaring infrastructure in Git and letting Argo CD reconcile the cluster. HelixGitpx's `helixgitpx-platform` repo is the single source of truth for prod.

### Goka
Go library for Kafka Streams-style stream processing; used in the conflict detector.

### grpc-gateway
Tool that generates a REST/JSON HTTP reverse proxy from protobuf service definitions.

### HelixGitpx
This product. A multi-upstream Git federation platform with AI-assisted conflict resolution.

### Idempotency Key
Client-provided UUID that lets servers deduplicate replayed requests; stored with 24 h TTL.

### KMP (Kotlin Multiplatform)
Mechanism for sharing Kotlin code (business logic, data, networking) across Android, iOS, JVM, and Native targets.

### Kyverno
K8s policy engine used for admission control (require signed images, resource limits, labels).

### LoRA (Low-Rank Adaptation)
Parameter-efficient fine-tuning technique. HelixGitpx ships per-task (optionally per-org) LoRA adapters over a shared base model.

### Meilisearch
User-facing search engine; indexes `repos`, `prs`, `issues`, etc.

### MirrorMaker 2
Kafka's cross-cluster replication tool; used for cross-region data sync.

### Multi-tenant
A single deployment serving multiple organisations with strict data isolation (row-level security, per-tenant tokens, per-tenant limits).

### OPA (Open Policy Agent)
Authorization engine executing Rego policies; called from every service before applying user-invoked actions.

### Outbox
Per-service Postgres table storing events to be published; Debezium tails it and emits to Kafka. Avoids dual-writes.

### PAT (Personal Access Token)
Long-lived API token. HelixGitpx PATs are prefixed `hpxat_` for easy identification.

### PDB (PodDisruptionBudget)
K8s resource limiting voluntary disruptions (drains, upgrades).

### PKCE
Proof Key for Code Exchange — OAuth 2.0 / OIDC extension required for public clients (SPAs, mobile).

### Policy-as-Code
Authorization + admission rules expressed as versioned, reviewable, deployable code (Rego for OPA, YAML for Kyverno).

### Projector
A stateless service that consumes Kafka events and projects them into a read model (PG, OpenSearch, Meilisearch, Qdrant).

### Qdrant
Vector database used for embeddings (code, issues, chat memory).

### Quorum
Minimum number of replicas that must acknowledge a write for it to be considered durable (Kafka `min.insync.replicas`, Postgres synchronous replication).

### RAG (Retrieval-Augmented Generation)
Pattern where an LLM is given retrieved context (from Qdrant hybrid search) before generating.

### Rego
The language used by OPA for policy rules.

### Rekor
Sigstore's transparency log for software artefacts; HelixGitpx uses it for Cosign signatures and audit-log anchoring.

### Resume Token
Opaque cursor the live-events service issues so clients can resume a stream after disconnect.

### RLAIF
Reinforcement Learning from AI Feedback; similar to RLHF but using AI-generated rewards where appropriate. HelixGitpx's self-learning loop combines human feedback and synthetic preferences.

### RLS (Row-Level Security)
Postgres feature ensuring queries see only rows matching the caller's tenant predicate. HelixGitpx sets `SET helix.tenant_id = '…'` per request.

### SBOM (Software Bill of Materials)
CycloneDX 1.5 JSON; generated by syft; published as OCI artefact for every image.

### Shadow Mode
Read-only soak mode for a newly connected upstream or a new AI model — runs alongside production but applies nothing.

### Signal Store (NgRx)
Reactive state container for Angular 17+ built around signals. Used per-feature.

### SLO / SLI / Error Budget
Service Level Objective / Indicator / remaining tolerance. Defined in [09-observability/18-observability.md](../09-observability/18-observability.md).

### SLSA (Supply-chain Levels for Software Artifacts)
Framework defining build integrity levels. HelixGitpx targets **Level 3** at GA, **Level 4** by Year 2.

### SPIFFE / SPIRE
Secure Production Identity Framework For Everyone / its reference implementation. Issues short-lived X.509 SVIDs to workloads.

### Strimzi
Kubernetes operator for Apache Kafka.

### SVID
SPIFFE Verifiable Identity Document — X.509 certificate or JWT issued by SPIRE.

### Temporal
Durable workflow engine. Workflow functions run inside workers; Temporal Server persists history for deterministic replay.

### Tenant
An organisation (`org_id`). All HelixGitpx data is tenant-scoped via RLS or application checks.

### Three-Way Merge
Classical Git merge with common ancestor + both diverged sides; the conflict resolver can execute this in a sandbox.

### Transparency Log
Append-only log (e.g. Rekor) whose entries produce a Merkle proof; used to detect tampering.

### Upstream
A connected Git host (GitHub, GitLab, Gitee, …) that a HelixGitpx repo synchronises with.

### vLLM
High-throughput LLM inference server with tensor parallelism; used for GPU pods.

### WIT (WebAssembly Interface Types)
Interface description format used to declare WASM plugin APIs in the component model.

### Workload Identity
A cryptographic identity bound to a pod/service rather than a user; provided by SPIFFE/SPIRE.

### Zoekt
Sourcegraph's fast code search engine; integrated for per-repo code search.

### Zonal / Regional / AZ
Deployment topology: a zone/AZ is a failure domain within a region; multiple zones per region, multiple regions in total. HelixGitpx runs multi-AZ per region.
