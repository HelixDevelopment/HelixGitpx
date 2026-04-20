# 03 — Microservices Catalog

> **Document purpose**: For every microservice in HelixGitpx, document its **responsibilities, interfaces, data ownership, scaling profile, SLOs, and runbook pointer**. This is the authoritative map of the service graph.

---

## 1. Service Taxonomy

HelixGitpx services fall into four tiers:

| Tier | Purpose | Examples |
|---|---|---|
| **Edge** | Public-facing; TLS termination, routing, AuthN/Z | API Gateway, Webhook Gateway, Git Ingress |
| **Domain Core** | Bounded-context business logic | Org Service, Repo Service, Sync Orchestrator, Conflict Resolver |
| **Platform** | Cross-cutting capability | Auth, Policy, Audit, Notifier, Search, Live Events |
| **Specialist** | Heavy / asynchronous work | AI Service, Adapter Pool, Scheduler, Billing |

All services are implemented in **Go 1.23+** unless noted. Each service uses:
- **Gin Gonic** for REST endpoints
- **gRPC (google.golang.org/grpc)** for internal RPC + streaming
- **sqlx / pgx** for Postgres
- **segmentio/kafka-go** or Sarama for Kafka
- **go-redis/v9** for Redis
- **OpenTelemetry SDK** for telemetry
- **zerolog** or **zap** for structured logging

---

## 2. Service Catalog (Summary Table)

| # | Service | Port (internal) | Owns Data | Key Kafka Topics (produce) | Scaling |
|---|---|---|---|---|---|
| 1 | **api-gateway** | 8080 (HTTP) / 9090 (gRPC) | — | — | HPA: CPU 70 % |
| 2 | **auth-service** | 9091 | `auth.*` PG schema | `auth.events` | HPA |
| 3 | **org-service** | 9092 | `org.*` | `org.events`, `team.events` | Sharded by `org_id` |
| 4 | **repo-service** | 9093 | `repo.*` | `repo.events`, `ref.events` | Sharded by `repo_id` |
| 5 | **upstream-registry** | 9094 | `upstream.*` | `upstream.events` | HPA |
| 6 | **adapter-pool** | 9095 | — (stateless) | `adapter.call.events` | Per-provider pool |
| 7 | **git-ingress** | 443 (HTTPS) / 22 (SSH) | — (writes to repo storage) | `git.push.received`, `git.pack.ingested` | HPA, L4+L7 |
| 8 | **webhook-gateway** | 8443 | `webhook.inbound` (dedup) | `upstream.*.received` | HPA |
| 9 | **sync-orchestrator** | 9096 | Temporal workflow state | `sync.*` | HPA |
| 10 | **conflict-resolver** | 9097 | `conflict.*` | `conflict.*` | Sharded by `repo_id` |
| 11 | **live-events-service** | 9098 (gRPC stream) / 8081 (WSS) | session state in Redis | — | HPA + sticky |
| 12 | **notifier** | 9099 | `notif.*` | `notify.sent`, `notify.failed` | HPA |
| 13 | **search-service** | 9100 | — (façade) | — | HPA |
| 14 | **ai-service** | 9101 | `ai.*` | `ai.*` | GPU, autoscaled 0→N |
| 15 | **policy-service** | 9102 | `policy.*` | `policy.decisions` | HPA |
| 16 | **audit-service** | 9103 | `audit.*` (append-only) | `audit.events` | Horizontal |
| 17 | **billing-service** | 9104 | `billing.*` | `billing.*` | HPA |
| 18 | **scheduler** | 9105 | `scheduler.*` | `scheduler.tick` | Leader-elected |
| P1 | **schema-registry (Karapace)** | 8082 | — | — | 3 replicas |
| P2 | **temporal-server** | 7233 | Temporal PG | — | HA cluster |
| P3 | **otel-collector** | 4317 | — | — | DaemonSet |
| P4 | **opa-policy-agent** | 8181 | — | — | Sidecar |
| P5 | **spire-server** / **spire-agent** | 8081 / — | SPIRE PG | — | HA + DaemonSet |
| P6 | **vault** | 8200 | Vault backend | — | HA |
| P7 | **argo-cd** | 8083 | Argo PG | — | HA |

---

## 3. Service Template

Every service folder follows this template (enforced by `cookiecutter`):

```
service-name/
├── cmd/
│   └── service-name/
│       └── main.go               # wires config, DI, signal handling
├── internal/
│   ├── api/
│   │   ├── grpc/                 # gRPC handlers
│   │   └── rest/                 # Gin handlers
│   ├── domain/                   # pure business rules (entity, value objects, policies)
│   ├── app/                      # use-cases / command handlers
│   ├── infra/
│   │   ├── pg/                   # Postgres repository
│   │   ├── kafka/                # producer, consumer
│   │   ├── redis/                # cache
│   │   ├── clients/              # gRPC clients to other services
│   │   └── telemetry/            # OTel setup
│   └── platform/
│       ├── config/               # viper-based config
│       ├── health/               # /healthz, /readyz, /livez
│       └── errors/               # typed domain errors
├── proto/                        # service's proto (imports from /proto shared)
├── migrations/                   # golang-migrate SQL files
├── test/
│   ├── unit/
│   ├── integration/
│   ├── contract/                 # Pact tests
│   └── fixtures/
├── deploy/
│   ├── helm/                     # Helm chart
│   ├── k8s/                      # raw manifests (for kustomize)
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── Dockerfile.dev
│   └── compose.yaml              # local dev compose
├── docs/
│   ├── README.md                 # index
│   ├── runbook.md                # operational procedures
│   ├── metrics.md                # canonical metrics
│   └── slo.md                    # SLO definitions
├── Taskfile.yml                  # task runner
├── go.mod / go.sum
└── sonar-project.properties      # SonarQube config
```

---

## 4. Service Details

### 4.1 `api-gateway`

**Purpose**: Single public entrypoint for clients. Handles:
- TLS termination
- AuthN (validate OIDC JWT, translate to internal principal)
- AuthZ (OPA quick-deny; full check in downstream service)
- Rate limiting (token + IP + org)
- REST ↔ gRPC translation (via `grpc-gateway`)
- gRPC-Web for browsers
- Request/response observability

**Public interface**:
- `/api/v1/*` — REST (Gin + grpc-gateway)
- `/grpc.*` — gRPC (direct)
- `/events/subscribe` — gRPC server streaming / WebSocket fallback
- `/git/*` — forwarded to git-ingress
- `/webhooks/*` — forwarded to webhook-gateway

**Internal dependencies**: auth-service, policy-service, live-events-service, all domain services.

**Data ownership**: None (stateless). Session state in Redis (TTL).

**Key metrics**:
- `helixgitpx_gateway_requests_total{method, path, status}`
- `helixgitpx_gateway_request_duration_seconds`
- `helixgitpx_gateway_rate_limit_rejected_total{reason}`
- `helixgitpx_gateway_auth_failures_total{reason}`

**SLOs**:
- Availability: 99.95 %
- p95 latency (unary gRPC pass-through): 30 ms
- p95 latency (REST gateway): 50 ms

**Scaling**: HPA (CPU 70 %, min 3, max 100). Sharded by path via Istio routing.

**Runbook pointer**: [12-operations/19-operations-runbook.md#api-gateway](../12-operations/19-operations-runbook.md)

---

### 4.2 `auth-service`

**Purpose**: Identity management and session issuance.

**Responsibilities**:
- OIDC flow (code grant, PKCE) with configured IdPs.
- Issue internal access tokens (JWT, RS256, 15-min TTL).
- Refresh token rotation with reuse detection.
- Personal Access Tokens (PAT) — hashed with Argon2id, scoped.
- SCIM 2.0 server for provisioning.
- Session management (Redis-backed).
- MFA (TOTP, WebAuthn/Passkey).
- Anomaly detection on login (geoip + behavioural).

**gRPC API** (see [06-grpc-api.md](../04-apis/06-grpc-api.md)):
```protobuf
service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Refresh(RefreshRequest) returns (TokenResponse);
  rpc Logout(LogoutRequest) returns (google.protobuf.Empty);
  rpc ValidateToken(ValidateRequest) returns (Principal);
  rpc CreatePAT(CreatePATRequest) returns (PATResponse);
  rpc RevokePAT(RevokePATRequest) returns (google.protobuf.Empty);
  rpc ListSessions(ListSessionsRequest) returns (SessionList);
  // SCIM endpoints exposed via REST (gateway routes)
}
```

**Data**: `auth` schema — `users`, `sessions`, `tokens_pat`, `tokens_refresh`, `mfa_factors`, `login_audit`.

**Scaling**: HPA (CPU 70 %). All reads cacheable (Redis).

**Key metrics**:
- `helixgitpx_auth_logins_total{status}`
- `helixgitpx_auth_token_verifications_total{status}`
- `helixgitpx_auth_pat_created_total`
- `helixgitpx_auth_mfa_challenges_total{method, status}`

---

### 4.3 `org-service`

**Purpose**: Owns organisations, teams, and memberships.

**Responsibilities**:
- CRUD on orgs, teams.
- Member management (add/remove/role).
- Org-level settings (default policies, LLM preferences).
- Org switching (primary → personal ↔ org).

**Aggregates**: `Organisation`, `Team`, `Membership`.

**Events emitted**:
- `org.created`, `org.updated`, `org.deleted`
- `team.created`, `team.updated`, `team.deleted`
- `membership.added`, `membership.role_changed`, `membership.removed`

**gRPC API**:
```protobuf
service OrgService {
  rpc CreateOrg(CreateOrgRequest) returns (Org);
  rpc GetOrg(GetOrgRequest) returns (Org);
  rpc UpdateOrg(UpdateOrgRequest) returns (Org);
  rpc DeleteOrg(DeleteOrgRequest) returns (google.protobuf.Empty);
  rpc ListOrgs(ListOrgsRequest) returns (OrgList);
  rpc CreateTeam(CreateTeamRequest) returns (Team);
  rpc AddMember(AddMemberRequest) returns (Membership);
  rpc UpdateMemberRole(UpdateMemberRoleRequest) returns (Membership);
  rpc RemoveMember(RemoveMemberRequest) returns (google.protobuf.Empty);
  rpc ListMembers(ListMembersRequest) returns (MemberList);
  // Streaming variant for large orgs
  rpc WatchOrg(WatchOrgRequest) returns (stream OrgEvent);
}
```

**Data**: `org` schema — `organisations`, `teams`, `memberships`, `org_settings`. RLS enforced.

**Scaling**: Sharded by `org_id` hash; 16 shards at GA.

---

### 4.4 `repo-service`

**Purpose**: Owns repositories, refs, tags, releases, branch protection rules.

**Responsibilities**:
- CRUD on repos.
- Ref metadata (branch, tag, head SHA, protection).
- Release metadata + artefact references.
- Branch protection policy evaluation (before accepting push).
- Repo-level settings (which upstreams are enabled, policy overrides).

**Aggregates** (event-sourced):
- `Repository` (created, renamed, archived, visibility_changed, …)
- `Ref` (created, updated, deleted, protected)
- `Release` (created, published, yanked)

**gRPC API** includes server-streaming `WatchRepo(repo_id)` that emits every state change in real time (used by Live Events service).

**Data**: `repo` schema:
- `repositories`
- `refs` — composite PK (repo_id, name); snapshot projection
- `ref_events` — append-only event log (snapshot source)
- `releases`
- `release_assets`
- `branch_protection_rules`
- `repo_upstream_bindings` (which upstreams this repo syncs to)

**Scaling**: Sharded by `repo_id`. Hot repos (many pushes/sec) use per-repo mutex in Redis to serialise writes.

---

### 4.5 `upstream-registry`

**Purpose**: Catalogue of providers and their per-org credentials.

**Responsibilities**:
- Register / deregister providers (GitHub, GitLab, …).
- Store credentials encrypted via Vault transit engine.
- Rotate credentials (automated where provider supports it).
- Shadow-mode toggle (new provider enabled read-only first).
- Health-check upstreams; publish `upstream.health_changed`.

**gRPC API**:
```protobuf
service UpstreamRegistry {
  rpc RegisterUpstream(RegisterUpstreamRequest) returns (Upstream);
  rpc UpdateUpstreamCredentials(UpdateCredRequest) returns (google.protobuf.Empty);
  rpc RotateCredentials(RotateRequest) returns (Upstream);
  rpc EnableUpstream(EnableRequest) returns (Upstream);
  rpc DisableUpstream(DisableRequest) returns (Upstream);
  rpc DisconnectUpstream(DisconnectRequest) returns (google.protobuf.Empty);
  rpc ListUpstreams(ListRequest) returns (UpstreamList);
  rpc TestUpstream(TestRequest) returns (TestResult);
  rpc GetHealth(GetHealthRequest) returns (Health);
}
```

**Data**: `upstream` schema — `providers`, `upstream_configs`, `credentials_vault_refs` (no plaintext).

---

### 4.6 `adapter-pool`

**Purpose**: The actual talkers to each Git provider. Implements a common **Universal Git Adapter** interface.

**Key interface** (Go):
```go
type UniversalGitAdapter interface {
    // Ref operations
    ListRefs(ctx context.Context, repo RepoRef) ([]Ref, error)
    PushRefs(ctx context.Context, repo RepoRef, refs []RefUpdate, opts PushOptions) (PushResult, error)
    FetchRefs(ctx context.Context, repo RepoRef) ([]Ref, error)

    // Repo
    CreateRepo(ctx context.Context, spec RepoSpec) (RepoRef, error)
    UpdateRepo(ctx context.Context, update RepoUpdate) error
    DeleteRepo(ctx context.Context, repo RepoRef) error

    // Collaboration
    CreatePR(ctx context.Context, spec PRSpec) (PRRef, error)
    MergePR(ctx context.Context, pr PRRef, opts MergeOpts) error
    ListPRs(ctx context.Context, repo RepoRef, filter PRFilter) ([]PR, error)

    CreateIssue(ctx context.Context, spec IssueSpec) (IssueRef, error)
    UpdateIssue(ctx context.Context, update IssueUpdate) error
    ListIssues(ctx context.Context, repo RepoRef, filter IssueFilter) ([]Issue, error)

    CreateRelease(ctx context.Context, spec ReleaseSpec) (ReleaseRef, error)

    // Webhooks
    RegisterWebhook(ctx context.Context, repo RepoRef, spec WebhookSpec) error

    // Health & limits
    GetRateLimit(ctx context.Context) (RateLimitInfo, error)
    Ping(ctx context.Context) error
}
```

**Provider implementations**: `github/`, `gitlab/`, `gitee/`, `gitflic/`, `gitverse/`, `bitbucket/`, `codeberg/`, `gitea/`, `forgejo/`, `sourcehut/`, `azure/`, `aws/`, `generic/` (for anything that speaks plain Git + optional REST).

**Extensibility**: WebAssembly plugins for custom/enterprise providers. Plugin interface is the adapter interface.

**Scaling**: One goroutine pool per provider, sized by observed RPS. Each pool has its own circuit breaker and rate limiter.

**Key metrics**:
- `helixgitpx_adapter_calls_total{provider, operation, status}`
- `helixgitpx_adapter_duration_seconds{provider, operation}`
- `helixgitpx_adapter_rate_limit_remaining{provider}`
- `helixgitpx_adapter_circuit_open_total{provider}`

---

### 4.7 `git-ingress`

**Purpose**: Terminate Git-protocol connections (git-over-HTTPS and git-over-SSH).

**Responsibilities**:
- HTTPS smart protocol (`/info/refs`, `/git-upload-pack`, `/git-receive-pack`).
- SSH (uses `libsshd` or gliderlabs/ssh; authorizes via dynamic `AuthorizedKeysCommand`).
- On receive-pack: verify policy via repo-service (branch protection, required signatures).
- Stream packfile to **Gitaly** (for actual object storage) or embedded go-git with local storage.
- Persist **received event** to Kafka before acknowledging client (M-4 invariant).
- Run pre-receive hooks (secret scan via Gitleaks, commit signature verification).

**Architecture note**: Uses **Gitaly** as the canonical Git server library (battle-tested in GitLab). Gitaly speaks gRPC natively; git-ingress is a thin policy+event layer on top.

**Scaling**: Horizontal; long-lived connections require sticky routing at ingress (session affinity by repo).

---

### 4.8 `webhook-gateway`

**Purpose**: Receive webhooks from all upstream providers; dedupe, verify, translate to canonical events.

**Per-provider**:
- HMAC verification (each provider has its own scheme; we encode per-provider).
- Idempotency dedupe via Redis SETNX with TTL (keyed by provider + delivery ID).
- Schema translation: e.g. `push` from GitHub → `upstream.ref.updated`.
- Replay protection (timestamp window).

**gRPC API**: Not client-facing; REST-only via Gin with `/webhooks/:provider/:repo_id`.

**Scaling**: HPA by request rate.

**Key metrics**:
- `helixgitpx_webhook_received_total{provider, event_type, status}`
- `helixgitpx_webhook_hmac_failures_total{provider}`
- `helixgitpx_webhook_dedupe_hits_total{provider}`

---

### 4.9 `sync-orchestrator`

**Purpose**: Bi-directional sync choreography. Uses **Temporal** for durable workflows.

**Key workflows**:

- `FanOutPushWorkflow`: one incoming push → push to N upstreams.
- `InboundReconcileWorkflow`: upstream webhook received → update internal state → consider fan-out to other upstreams.
- `RepoMigrationWorkflow`: one-time migration of a repo from provider A to being federated across {A, B, C}.
- `UpstreamOnboardingWorkflow`: configure new upstream, shadow-mode soak, go-live.
- `OrgOnboardingWorkflow`: create org → create teams → connect upstreams → import existing repos.
- `ReRunFailedSyncsWorkflow`: nightly sweep of dead-letter queue.

**Activities** (gRPC to adapter-pool, repo-service, conflict-resolver, notifier).

**Workflow versioning**: Temporal's versioning primitives (patch, getVersion) — mandatory for any workflow change.

**Scaling**: Horizontal workers; Temporal handles state.

---

### 4.10 `conflict-resolver`

**Purpose**: Detect and resolve conflicts between HelixGitpx's view and an upstream's view.

**Conflict classes** (full detail in [09-conflict-resolution.md](09-conflict-resolution.md)):

1. **Ref divergence** — the same branch has different heads on two upstreams.
2. **Concurrent metadata edits** — issue labels/milestones changed on two sides within the window.
3. **Rename collision** — same file renamed differently upstream.
4. **PR state divergence** — PR marked merged on one side, still open elsewhere.
5. **Tag collision** — same tag name, different targets.
6. **LFS object divergence** — same OID, different content (rare but possible in corruption).

**Resolution strategies**:

- **Policy-based** (deterministic): prefer-primary, prefer-newer, prefer-signed, union-merge, three-way, octopus merge.
- **CRDT-based** (for metadata): labels, milestones, issue bodies as G-Sets / LWW-register / Automerge documents.
- **AI-assisted**: for ambiguous cases, ask the fine-tuned LLM for a proposal; emit `conflict.escalated(ai_proposed)`.
- **Human-in-the-loop**: when AI confidence < threshold or org requires signoff.

**Scaling**: Sharded by `repo_id`; per-repo mutex ensures one conflict resolution at a time per repo.

---

### 4.11 `live-events-service`

**Purpose**: Deliver Kafka events to subscribed clients in real time.

**Protocols**:
- **gRPC server streaming** (primary — mobile, desktop, CLI, modern web).
- **WebSocket (Socket.IO-compatible)** fallback for old browsers / restrictive networks.
- **Server-Sent Events (SSE)** as a second fallback.

**Subscription model**:
- Subscribe to scopes: `user:{id}`, `org:{id}`, `repo:{id}`, `global` (admin).
- Filter by event types.
- Resume from `last_event_id` (Kafka offset translated to opaque cursor).

**Backpressure**:
- Per-subscription buffer (bounded).
- Slow consumers: drop-oldest with a "gap detected, re-fetch" notification.
- Cap total subscriptions per token (e.g., 50).

**Session stickiness**: Required (subscription is per-connection). Istio consistent-hash routing.

---

### 4.12 `notifier`

**Purpose**: Multi-channel outbound notifications.

**Channels** (pluggable, driven by **Apprise**-style drivers):
- Email (SMTP; SES/Postmark optional)
- Slack
- Microsoft Teams
- Discord
- Telegram
- Matrix
- Mattermost
- Rocket.Chat
- Webhook (generic)
- SMS (Twilio, Vonage, self-hosted via Plivo/Sinch)
- Push (FCM for Android, APNs for iOS, UnifiedPush alternative)
- PagerDuty / Opsgenie
- ntfy.sh (self-hostable push)
- Gotify
- IRC (legacy but requested)
- XMPP / Jabber
- Native OS notifications (via client apps)

**Features**:
- Templating with Go `text/template` + user-defined.
- Rate-limiting per-user + per-channel.
- Bounce/failure tracking.
- Two-way chatops: commands from channels (`/helix status myrepo`) via webhooks.

---

### 4.13 `search-service`

**Purpose**: Unified query façade over:
- **OpenSearch** — logs, audit, long-retention.
- **Meilisearch** — user-facing text search (repos, issues, PRs, docs).
- **Qdrant** — vector search (semantic, for AI).

**gRPC API** exposes high-level queries; service picks the right backend.

**Synonyms, typo-tolerance, stopwords**: configurable per tenant.

---

### 4.14 `ai-service`

**Purpose**: All LLM interactions.

**Capabilities**:
- Routing to providers via **LiteLLM** (OpenAI, Anthropic, Gemini, self-hosted Ollama, vLLM).
- Prompt templates (DSPy-compiled; versioned in code).
- Guardrails (Guardrails AI / NeMo Guardrails).
- Fine-tuning orchestration (see [10-llm-self-learning.md](../07-ai/10-llm-self-learning.md)).
- Vector embedding generation (via `sentence-transformers` / BGE).
- RAG over Qdrant.
- Cost metering per call.

**Deployment**: Go service + Python sidecars (for DSPy, sentence-transformers). Communicates via gRPC.

**GPU autoscaling**: KEDA + NVIDIA GPU Operator; scale 0→N based on queue depth.

---

### 4.15 `policy-service`

**Purpose**: AuthZ decisions.

**Implementation**: Wraps **OPA** (Open Policy Agent). Bundles of Rego policies are Git-managed; hot-reload on change.

**Policy bundles**:
- `policy/authz` — who can do what to whom.
- `policy/sync` — which repos sync where, force-push rules.
- `policy/admission` — Kubernetes admission (Kyverno alt).
- `policy/llm` — which LLM features available per tier.

**API**:
```protobuf
service PolicyService {
  rpc Decide(DecisionRequest) returns (DecisionResponse);
  rpc DecideBatch(BatchDecisionRequest) returns (BatchDecisionResponse);
  rpc ListPolicies(ListPoliciesRequest) returns (PolicyList);
  rpc GetBundleVersion(Empty) returns (BundleVersion);
}
```

---

### 4.16 `audit-service`

**Purpose**: Immutable append-only audit sink.

**Design**:
- Consumer of **all** `*.events` topics (universal sink).
- Writes to Postgres `audit.log_YYYYMM` partitioned tables.
- Mirrors critical entries to Rekor (transparency log) for tamper-evidence.
- Retention: 1 year hot, 7 years cold (OpenSearch archive → S3 Glacier).

**Query interface**: Search service façade + REST `/api/v1/audit?...`.

---

### 4.17 `billing-service`

**Purpose**: Meter usage and enforce quotas.

**Meters**:
- Repos connected.
- Upstreams per repo.
- Pushes / day.
- LLM tokens consumed.
- Storage (LFS) GB.

**Enforcement**: Soft (warn) → hard (reject) based on plan.

**Integrations**: Stripe / Paddle (managed SaaS offering). Self-host has no billing.

---

### 4.18 `scheduler`

**Purpose**: Cron, delayed jobs, periodic maintenance.

**Jobs**:
- Nightly repo GC (`git gc --aggressive --prune=now`) via adapter-pool.
- Upstream health sweeps.
- PAT expiration warnings.
- LLM fine-tune trigger (weekly).
- Stale branch cleanup.
- Credential rotation reminders.

**Leader election**: Kubernetes `Lease` resource.

---

## 5. Shared Libraries (Go modules)

| Module | Purpose |
|---|---|
| `helix-platform/logging` | zerolog with tenant context auto-injection |
| `helix-platform/telemetry` | OTel setup, trace/span helpers |
| `helix-platform/errors` | typed domain errors with grpc.Status mapping |
| `helix-platform/config` | viper config loading + validation |
| `helix-platform/grpc` | gRPC interceptors (auth, trace, recover) |
| `helix-platform/gin` | Gin middleware suite (auth, trace, recover, ratelimit) |
| `helix-platform/kafka` | producer/consumer with schema registry |
| `helix-platform/pg` | pgx wrapper with RLS setup |
| `helix-platform/redis` | go-redis wrapper with tracing |
| `helix-platform/temporal` | Temporal client with tracing |
| `helix-platform/spire` | SVID retrieval + renewal |
| `helix-platform/opa` | OPA client |
| `helix-platform/health` | /healthz + /readyz + /livez + SIGTERM handler |
| `helix-platform/test` | test fixtures, testcontainers helpers |

---

## 6. Service-to-Service Call Matrix

A quick-reference of **who calls whom** (non-exhaustive, gRPC unless noted):

| Caller → | gateway | auth | org | repo | upstream | adapter | git-ing | webhook | sync | conflict | events | notif | search | ai | policy | audit | billing | sched |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| **api-gateway** | | ✓ | ✓ | ✓ | ✓ | | | | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | |
| **git-ingress** | | ✓ | | ✓ | | | | | ✓(k) | | | | | | ✓ | ✓(k) | | |
| **webhook-gateway** | | | | | ✓ | | | | ✓(k) | | | | | | | ✓(k) | | |
| **sync-orch** | | | | ✓ | ✓ | ✓ | | | | ✓ | | ✓(k) | | | ✓ | ✓(k) | | |
| **conflict-res** | | | | ✓ | | ✓ | | | | | | | | ✓ | ✓ | ✓(k) | | |
| **live-events** | | ✓ | | | | | | | | | | | | | ✓ | | | |
| **notifier** | | ✓ | ✓ | ✓ | | | | | | | | | | | | ✓(k) | | |
| **scheduler** | | | ✓ | ✓ | ✓ | ✓ | | | ✓ | ✓ | | ✓ | | ✓ | | ✓(k) | ✓ | |

(✓ = gRPC, ✓(k) = via Kafka)

---

## 7. Polyglot Sidecars

| Sidecar | Language | Role |
|---|---|---|
| `ai-embedder` (of ai-service) | Python (FastAPI + sentence-transformers + BGE) | Generate embeddings |
| `ai-guardrails` | Python | Guardrails/NeMo wrapper |
| `dspy-compiler` | Python | Compile prompt programs |
| `opa-agent` | Go (upstream) | Policy evaluation |
| `otel-collector` | Go (upstream) | Trace/metric/log forwarding |
| `envoy` | C++ (service mesh) | mTLS, routing |
| `spire-agent` | Go (upstream) | SVID attestation |

All Python sidecars use `uv` (formerly `poetry`) for dep management and `ruff` for linting. Images distroless (chainguard).

---

## 8. Service Lifecycle & Governance

### 8.1 Creating a New Service

1. File RFC in `13-roadmap/rfcs/`.
2. Get ADR signoff from Architect.
3. Bootstrap from `cookiecutter-helix-go` template (enforces structure).
4. Add to `service-catalog.yaml` (source of truth for CI/CD, Argo CD generators).
5. Wire into Helm umbrella chart.
6. Add dashboards, alerts, runbook.
7. First PR must include: proto, unit tests, one integration test, Helm chart, Grafana dashboard JSON, runbook.

### 8.2 Decommissioning a Service

1. RFC + dependency audit (call graph).
2. Migrate functionality to successor; run in shadow for 2 weeks.
3. Redirect callers (feature flag).
4. Soak for 4 weeks.
5. Remove.

### 8.3 Ownership

Each service has a **RACI** owner (primary, secondary, platform). Documented in `service-catalog.yaml`. On-call rotation is service-scoped.

---

*— End of Microservices Catalog —*
