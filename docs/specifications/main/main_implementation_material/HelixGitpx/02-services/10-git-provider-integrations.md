# 10 — Git Provider Integrations

> **Document purpose**: Specify the **Universal Git Adapter** contract, the per-provider implementations shipped at GA, the **WASM plugin** extensibility point for custom providers, and the operational matrix of what each provider supports.

---

## 1. The Universal Git Adapter

Every Git provider is represented by an implementation of the `UniversalGitAdapter` Go interface. Adapters are deployed inside the `adapter-pool` service; each provider gets its own goroutine pool, rate limiter, and circuit breaker.

```go
// pkg/adapter/interface.go
type UniversalGitAdapter interface {
    // Authentication & limits
    Ping(ctx context.Context) error
    GetRateLimit(ctx context.Context) (RateLimitInfo, error)
    GetIdentity(ctx context.Context) (ProviderIdentity, error)

    // Repo lifecycle
    CreateRepo(ctx context.Context, spec RepoSpec) (RepoRef, error)
    GetRepo(ctx context.Context, ref RepoRef) (Repo, error)
    UpdateRepo(ctx context.Context, update RepoUpdate) error
    DeleteRepo(ctx context.Context, ref RepoRef) error
    ListRepos(ctx context.Context, owner string, filter RepoFilter) ([]Repo, error)

    // Refs
    ListRefs(ctx context.Context, ref RepoRef) ([]Ref, error)
    PushRefs(ctx context.Context, ref RepoRef, updates []RefUpdate, opts PushOptions) (PushResult, error)
    FetchRefs(ctx context.Context, ref RepoRef, opts FetchOptions) (FetchResult, error)

    // Branch protection
    SetBranchProtection(ctx context.Context, ref RepoRef, rule BranchProtectionRule) error
    GetBranchProtection(ctx context.Context, ref RepoRef, pattern string) (BranchProtectionRule, error)

    // Collaboration
    CreatePR(ctx context.Context, spec PRSpec) (PRRef, error)
    UpdatePR(ctx context.Context, pr PRRef, update PRUpdate) error
    GetPR(ctx context.Context, pr PRRef) (PR, error)
    ListPRs(ctx context.Context, ref RepoRef, filter PRFilter) ([]PR, error)
    MergePR(ctx context.Context, pr PRRef, opts MergeOpts) error
    AddPRReview(ctx context.Context, pr PRRef, review ReviewSpec) error
    AddPRComment(ctx context.Context, pr PRRef, comment CommentSpec) error

    CreateIssue(ctx context.Context, spec IssueSpec) (IssueRef, error)
    UpdateIssue(ctx context.Context, issue IssueRef, update IssueUpdate) error
    ListIssues(ctx context.Context, ref RepoRef, filter IssueFilter) ([]Issue, error)
    AddIssueComment(ctx context.Context, issue IssueRef, comment CommentSpec) error

    CreateRelease(ctx context.Context, spec ReleaseSpec) (ReleaseRef, error)
    UploadReleaseAsset(ctx context.Context, rel ReleaseRef, asset AssetSpec, r io.Reader) error

    // Webhooks (registration)
    RegisterWebhook(ctx context.Context, ref RepoRef, spec WebhookSpec) (WebhookRef, error)
    UnregisterWebhook(ctx context.Context, ref WebhookRef) error

    // LFS
    LFSUpload(ctx context.Context, ref RepoRef, spec LFSObjectSpec, r io.Reader) error
    LFSDownload(ctx context.Context, ref RepoRef, oid string) (io.ReadCloser, error)

    // Capabilities introspection
    Capabilities() AdapterCapabilities
}
```

### 1.1 Capabilities

```go
type AdapterCapabilities struct {
    LFS                bool
    WebhooksRegistrable bool
    GraphQL            bool
    PullRequests       bool
    Issues             bool
    Releases           bool
    BranchProtection   bool
    SignedPushes       bool
    RequiredStatusChecks bool
    ForcePushWithLease bool
    SSHAuth            bool
    AppTokenAuth       bool
    OAuthAuth          bool
    SCIMSupported      bool
    MirroringNative    bool
}
```

The rest of HelixGitpx consults `Capabilities()` before dispatching operations so it never asks a provider to do something it can't.

---

## 2. Provider Matrix at GA

| Provider | Kind | Auth | LFS | Webhooks | GraphQL | PR | Issues | Releases | Branch Protection | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| **GitHub** | `github` | OAuth2, App token, PAT, SSH | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | Reference |
| **GitLab** | `gitlab` | OAuth2, PAT, Deploy token, SSH | ✓ | ✓ | ✓ | MR | ✓ | ✓ | ✓ | Saas + self-hosted |
| **Gitee** | `gitee` | OAuth2, PAT | ✓ | ✓ | — | ✓ | ✓ | ✓ | Partial | China-focused |
| **GitFlic** | `gitflic` | PAT | — | ✓ | — | ✓ | ✓ | partial | partial | Russia |
| **GitVerse** | `gitverse` | PAT | ✓ | ✓ | — | ✓ | ✓ | ✓ | partial | Sber (Russia) |
| **Bitbucket Cloud** | `bitbucket` | OAuth2, App password | ✓ | ✓ | — | ✓ | disabled | ✓ | ✓ | Atlassian |
| **Bitbucket DC** | `bitbucket_dc` | PAT | ✓ | ✓ | — | ✓ | ✓ | ✓ | ✓ | Self-hosted |
| **Codeberg** | `codeberg` | OAuth2, PAT | ✓ | ✓ | — | ✓ | ✓ | ✓ | ✓ | Forgejo-based |
| **Gitea** | `gitea` | OAuth2, PAT | ✓ | ✓ | — | ✓ | ✓ | ✓ | ✓ | Self-hosted |
| **Forgejo** | `forgejo` | OAuth2, PAT | ✓ | ✓ | — | ✓ | ✓ | ✓ | ✓ | Gitea fork |
| **SourceHut** | `sourcehut` | OAuth2, PAT | — | ✓ | ✓ | patchsets | tickets | — | — | Hyperkit-style |
| **Azure DevOps** | `azure` | OAuth2, PAT, AAD | ✓ | ✓ | — | ✓ | work items | ✓ | ✓ | Microsoft |
| **AWS CodeCommit** | `aws_codecommit` | IAM SigV4 | — | ✓ | — | pull requests | — | — | approval rules | |
| **Generic Git** | `generic` | SSH / HTTPS basic | via server-side LFS | ❌ (polling) | — | — | — | — | — | Fallback |
| **WASM plugin** | `wasm_plugin` | defined by plugin | defined by plugin | defined by plugin | | | | | | Extensibility |

Any provider missing a capability degrades gracefully: operations that depend on it are no-ops with a logged warning; HelixGitpx still syncs what it can.

---

## 3. Per-Provider Integration Recipes

### 3.1 GitHub

- **Auth**: Prefer **GitHub App** (installation tokens, 1 h TTL, scoped, rate-limited separately). Fallback to OAuth2 + PAT.
- **API**: REST (primary) + GraphQL (batched reads, e.g. PR list with reviewers).
- **Webhooks**: Register repo webhook on `RegisterWebhook` with full event set + HMAC SHA-256 secret.
- **LFS**: Uses GitHub's LFS server via git-lfs protocol; we store our canonical LFS on R2/MinIO and `lfs-transfer` in adapter.
- **Rate limits**: App: 15 k/h per installation; PAT: 5 k/h. Adapter consumes `X-RateLimit-*` and pauses globally when ≤ 100 remaining.
- **Edge cases**:
  - Force-push with branch protection → fall back to `with-lease` or drop.
  - Draft PR state supported (`draft: true`).
  - GitHub Actions workflow files (`.github/workflows/*`) trigger on push — handled by mirroring unchanged.

### 3.2 GitLab

- **Auth**: OAuth2 + Deploy tokens + SSH.
- **Self-hosted variants**: Base URL configurable; adapter introspects version on connect.
- **Rename**: GitLab renames everything (PR → MR, etc.). Adapter translates.
- **CI**: `.gitlab-ci.yml` is mirrored as-is.
- **Webhooks**: system hooks (admin) or project hooks (per repo); we use project hooks.
- **Rate limit**: `RateLimit-*` headers; typical 600/min authenticated.

### 3.3 Gitee

- **Auth**: OAuth2 + PAT.
- **API**: REST v5 only; no GraphQL.
- **Rate limits**: 5 000/h authenticated; stricter on file endpoints.
- **Locale**: Chinese timestamps handled via RFC3339.
- **Known gotchas**: Some fields (e.g. milestones) require separate calls; batch where possible.
- **Webhooks**: password-based HMAC; we store the shared secret in Vault.

### 3.4 GitFlic

- **Auth**: PAT only (as of 2026).
- **API**: REST; coverage limited (no branch protection; we store rules locally and enforce via webhook-side guard).
- **Rate limit**: 1 000/h.
- **Locale**: Russian; no translation needed (we store whatever the upstream returns).

### 3.5 GitVerse (Sber)

- **Auth**: PAT.
- **API**: REST; closer to Gitea model.
- **LFS**: Supported.
- **Unique**: Tight integration with Sber's CI — we mirror build statuses if configured.

### 3.6 Bitbucket Cloud

- **Auth**: OAuth2 consumer + App password (deprecated 2024) → moving to "workspace tokens".
- **API**: v2.0 REST.
- **PRs**: "Pull requests" but issues are disabled by default.
- **Webhooks**: Per-repo.
- **Rate limit**: 1 000/h per user; new plans 2025+ have per-IP bucket too.

### 3.7 Bitbucket Data Center (Server)

- Separate adapter (different API shape: `/rest/api/1.0/...`).
- PAT-based auth only.

### 3.8 Codeberg / Forgejo / Gitea

- All three share the Gitea REST API; one adapter parameterised by base URL is sufficient.
- Webhooks: HMAC.
- LFS: native.
- Federation (ActivityPub) support will be evaluated post-GA for Forgejo's federated PR mechanism.

### 3.9 SourceHut

- Distinct model: patches over email, ticket trackers, build services, git.sr.ht.
- Adapter translates PR → patchset emails, issues → tickets.
- Webhooks configured via `hut` CLI equivalents.

### 3.10 Azure DevOps

- Auth: PAT or AAD.
- API: REST v7.1.
- Uses "Pull Requests" (same concept).
- Work items are translated to issues.
- Pipelines (`azure-pipelines.yml`) mirrored unchanged.

### 3.11 AWS CodeCommit

- SigV4 signed requests.
- Limited features (no issues/PRs in some regions; PR support since 2018).
- Notification via EventBridge → we subscribe via SNS → our webhook gateway.

### 3.12 Generic Git

- Plain `git://`, `https://`, or `ssh://` endpoint.
- No native PR/issues; we simulate with our own shadow PR/issue that only lives in HelixGitpx.
- Change detection: periodic polling (configurable, default 5 min).

---

## 4. WASM Plugin Extensibility

Custom providers (enterprise-internal, future-provider) implement the adapter interface as a **WebAssembly module** compiled from Rust, Go, or AssemblyScript.

### 4.1 ABI

We ship a `wit/helixgitpx-adapter.wit` WIT file defining the adapter interface as component-model exports:

```wit
package helixgitpx:adapter@1.0.0;

interface adapter {
    record repo-ref { owner: string, name: string }
    record ref-update { name: string, old-sha: string, new-sha: string, force: bool }
    record push-result { refs: list<ref-update>, errors: list<string> }

    ping: func() -> result<_, string>;
    list-refs: func(r: repo-ref) -> result<list<ref>, string>;
    push-refs: func(r: repo-ref, updates: list<ref-update>) -> result<push-result, string>;
    create-pr: func(spec: pr-spec) -> result<pr-ref, string>;
    // ... etc.
}
```

### 4.2 Runtime

- Host: `wasmtime-go` v20+ with the component model enabled.
- Sandboxing: CPU + memory + WASI time limits; no network except through explicit host imports (`outbound_http_request`, rate-limited).
- Plugins deployed as OCI artefacts (signed with Cosign); hot-reloaded on new version.
- Per-plugin configuration JSON (validated via JSON Schema shipped with the plugin).

### 4.3 SDK

We publish plugin SDKs for:
- Rust (`helixgitpx-plugin-sdk` crate).
- TinyGo.
- AssemblyScript.

Example (Rust):

```rust
use helixgitpx_plugin_sdk::{export_adapter, AdapterImpl};

struct MyAdapter;

impl AdapterImpl for MyAdapter {
    fn ping(&self) -> Result<(), String> { Ok(()) }
    // ... implement other methods
}

export_adapter!(MyAdapter);
```

---

## 5. Shadow Mode

Every newly connected upstream starts in **shadow mode**:

- Read-only: fetches refs, mirrors state internally, **never pushes** anywhere.
- Runs for a configurable soak (default 72 h).
- Reports:
  - Observed divergence events.
  - Would-have-been actions (counterfactual log).
  - Auth / rate-limit issues.
- Admin reviews the shadow report and graduates the upstream to "active".

This prevents a misconfigured upstream from damaging real data.

---

## 6. Webhook Ingestion

All adapters register webhooks that fire into `webhook-gateway`. The gateway:

1. Verifies HMAC / signature per provider.
2. Deduplicates by `(provider, delivery_id)` with a 24 h Redis window.
3. Translates raw payload into canonical event (`upstream.{kind}.{action}`) via per-provider transform.
4. Produces to Kafka.

For providers without webhooks (generic, some self-hosted), a **polling worker** runs periodically and compares ref lists.

---

## 7. Credential Management

- All secrets stored in **HashiCorp Vault** under `helixgitpx/data/upstreams/{org_id}/{upstream_id}`.
- Short-lived tokens (GitHub App installation tokens, Azure AAD tokens) regenerated on demand; long-lived PATs stored encrypted.
- Rotation events recorded in `upstream.credential_rotations`.
- SSH keypairs generated on registration and never leave Vault.

---

## 8. Observability

| Metric | |
|---|---|
| `helixgitpx_adapter_requests_total{provider, operation, status}` | Per-call counts |
| `helixgitpx_adapter_request_duration_seconds{provider, operation}` | Histogram |
| `helixgitpx_adapter_rate_limit_remaining{provider}` | Gauge |
| `helixgitpx_adapter_circuit_breaker_state{provider}` | 0 closed, 1 half-open, 2 open |
| `helixgitpx_adapter_auth_failures_total{provider}` | Token issues |
| `helixgitpx_adapter_shadow_divergence_total{provider}` | Shadow mode alarms |

Alerts:
- `UpstreamRateLimitLow` — remaining < 5 % for > 5 min.
- `AdapterCircuitOpen` — circuit open > 1 min.
- `AdapterAuthFailureSpike`.

---

## 9. Testing Strategy

- **Provider test harness**: every adapter has a test suite against a recorded cassette (go-vcr) for offline CI + optional live tests against a sandbox account in nightly.
- **Shared contract test** (Go generic): `adapter_contract_test.go` runs the same scenarios against every adapter implementation.
- **Chaos**: rate-limit saturation tests (force `429`s) — adapter must back off, never lose work.
- **WASM plugin**: test harness runs the same contract against a WASM-hosted adapter.

---

## 10. Governance

- Provider SDK versions tracked in `adapter-matrix.yaml`.
- Upstream API changes monitored via `github.com/cli/cli`-style changelog tailing + nightly contract tests.
- Breakages in upstream APIs trigger an `AdapterContractBroken` alert and pin us to the last-known-good adapter version until a fix ships.

---

*— End of Git Provider Integrations —*
