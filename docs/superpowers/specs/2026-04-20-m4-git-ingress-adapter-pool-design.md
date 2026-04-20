# M4 Git Ingress & Adapter Pool — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED (pending user review) |
| Author | Милош Васић + Claude (brainstorming session 2026-04-20) |
| Milestone | M4 — Git Ingress & Adapter Pool (Weeks 15–20) |
| Scope | Full 16-item roadmap §5 (items 54–69); nothing skipped |
| Sequencing | repo-service → git-ingress → adapter-pool → webhook-gateway → upstream-service |

---

## 1. Context

M1 Foundation (`m1-foundation` tag) produced monorepo + platform libs + hello. M2 Core Data Plane (`m2-core-data-plane` tag) produced the cluster manifests + hello outbox upgrade. M3 Identity & Orgs (`m3-identity-orgs` tag) produced auth + orgteam + audit + web shell with verify-m3-cluster passing 15/15.

M4 is where HelixGitpx starts **doing Git**: accepting pushes and fans out to upstream forges (GitHub/GitLab/Gitea) as the spec's central value proposition. After M4, a user can configure `helix.local` as their git remote, push once, and see the commit mirrored to every configured upstream.

## 2. Goals

G1. **repo-service** owns repo metadata: CRUD, event-sourced aggregate, refs, branch protection, LFS metadata.

G2. **git-ingress** accepts smart-HTTP Git clients (`git push`/`fetch`) with RBAC + quota + signed-push verification, exec'ing `git-http-backend` for the wire protocol.

G3. **adapter-pool** exposes an `Adapter` interface with three implementations: GitHub (REST+GraphQL), GitLab, Gitea — the last also serving Codeberg/Forgejo.

G4. **webhook-gateway** accepts inbound webhooks from all three providers, verifies HMAC, deduplicates, canonicalises to a shared `WebhookEvent` proto, publishes to `upstream.webhooks` Kafka.

G5. **upstream-service** owns upstream CRUD and repo↔upstream bindings; credentials live only in Vault KV.

G6. End-to-end: a `git push` to HelixGitpx replicates to configured upstreams (through adapter-pool) and a push to GitHub fires a webhook that lands in our Kafka.

## 3. Non-goals

- Conflict resolution between simultaneous upstream writes (M5).
- The remaining 7 providers: Bitbucket, Gitee, GitFlic, GitVerse, SourceHut, AzureDevOps, AWS CodeCommit, Generic Git (M5).
- WASM plugin host (M5).
- CRDT metadata (labels/milestones/assignees) (M5).
- Live events (M5).

## 4. Locked constraints

| ID | Constraint | Source |
|---|---|---|
| C-1 | Git protocol = go-git for parsing/policy + `git-http-backend` CGI for wire protocol | Q1 |
| C-2 | Adapter contract tests use go-vcr cassettes committed under `testdata/` | Design §5 |
| C-3 | Upstream credentials stored only at `kv/upstream/<id>` in Vault; the Postgres row holds only the Vault path | Design §5 |
| C-4 | LFS uses MinIO presigned URLs rewritten by git-ingress; no streaming through the service | Design §5 |
| Inherited | M1-M3 locked constraints (workflow_dispatch CI, portable compose, mise, HA manifests + local overlay, observability-first, outbox pattern for audit) | Prior specs |

## 5. Services

### 5.1 repo-service

**gRPC:** `RepoService.{Create,Get,List,Update,Delete}`, `RefService.{List,Get,Update,Protect,Unprotect}`.
**Schema (`repo.*`):**
- `repos(id uuid pk, org_id fk→org.orgs, slug text, default_branch text, lfs_enabled bool, created_at)` — unique(org_id, slug)
- `refs(id, repo_id fk, name text, sha text, updated_at)`
- `branch_protections(repo_id, pattern text, require_signed bool, required_reviewers int, updated_at)` — pk(repo_id, pattern)
- `lfs_objects(repo_id, oid char(64), size bigint, uploaded_at)` — pk(repo_id, oid)
- `outbox_events` — outbox for `repo.events` topic

**Event-sourced Repo aggregate:** every mutation emits an `RepoEvent` to the outbox; repo-service itself also reads its own events on restart to reconstruct derived state (ref counts, last-activity timestamps) — but the authoritative state is the rows. Event-sourcing is additive here; the aggregate is in-memory cache backed by PG.

### 5.2 git-ingress

**Deployment:** a single pod per replica with both the Go binary and the `git` + `git-http-backend` binaries (Alpine + `apk add git`). Repos on a PVC mounted at `/var/helix/repos/`.

**Flow:**
1. Client: `GET /<org>/<repo>.git/info/refs?service=git-upload-pack` (smart-HTTP protocol v2).
2. `git-ingress` Gin handler: JWT validation via `platform/auth`, OPA check (`data.helixgitpx.authz.allow` with `{user, action: "repo.read", repo}`), quota check (Redis token-bucket).
3. For receive-pack + `require_signed_pushes=true`: parse the `git push --signed` output from the body, `ssh.VerifySignature` against user's Keycloak-stored SSH keys.
4. Approved → `exec git-http-backend` via `os/exec` with stdin/stdout piped from the HTTP request/response.
5. After the process exits: parse stdout for ref updates (captured via a post-receive hook script installed per-repo), emit `RefUpdated` event to the outbox.

**Paths:**
- `/(.+)\.git/info/refs` → smart-HTTP discovery
- `/(.+)\.git/git-upload-pack` → fetch
- `/(.+)\.git/git-receive-pack` → push
- `/(.+)\.git/info/lfs/objects/batch` → LFS API (returns MinIO presigned URLs)

### 5.3 adapter-pool

**gRPC:** `AdapterService.{Push,Fetch,CreatePR,ListRefs,GetRepo,ListWebhooks,RegisterWebhook}` — every RPC takes a `UpstreamRef` (upstream id + provider) as its first field; the pool routes to the provider implementation by the enum.

**Adapter interface (Go):**

```go
package adapter

type Adapter interface {
	Push(ctx context.Context, dst Destination, refs []RefUpdate) error
	Fetch(ctx context.Context, src Source, refs []string) ([]RefValue, error)
	CreatePR(ctx context.Context, src, dst Branch, title, body string) (*PullRequest, error)
	ListRefs(ctx context.Context, src Source) ([]RefValue, error)
	GetRepo(ctx context.Context, src Source) (*RepoInfo, error)
	ListWebhooks(ctx context.Context, src Source) ([]Webhook, error)
	RegisterWebhook(ctx context.Context, src Source, url string, secret string, events []string) (*Webhook, error)
}
```

**Implementations:**
- `internal/providers/github/` — `google/go-github` + GraphQL for richer operations (PR review threads).
- `internal/providers/gitlab/` — `xanzy/go-gitlab`.
- `internal/providers/gitea/` — `code.gitea.io/sdk/gitea`. Same impl powers Codeberg + Forgejo by just changing `base_url`.

**Contract tests:** each provider has `testdata/{list_repos,push,register_webhook}.yaml` recorded with go-vcr. `go test` replays; `RECORD=1 go test` against live creds (documented for re-recording).

### 5.4 webhook-gateway

**HTTP endpoints** at `/webhook/{github,gitlab,gitea}`:
1. Read request body.
2. Verify HMAC-SHA256 using the provider's delivery secret (per-repo, resolved from `upstream.bindings` → Vault).
3. Dedup via Redis SET `webhook:seen:<provider>:<delivery-id>` with 7-day TTL. Duplicate → 200 + no-op.
4. Canonicalise into `helixgitpx.upstream.v1.WebhookEvent` proto (common fields: provider, delivery_id, repo, ref, event_type, body_raw).
5. Publish to `upstream.webhooks` via `platform/kafka.Producer` (synchronous, no outbox — webhook-gateway itself is stateless so exactly-once is bounded by HMAC+dedup at the edge).

### 5.5 upstream-service

**gRPC:** `UpstreamService.{Create,Get,List,Update,Delete}`, `UpstreamService.{Bind,Unbind,ListBindings}`.

**Schema (`upstream.*`):**
- `upstreams(id, slug citext unique, provider enum, base_url text, vault_path text, created_at)` — vault_path e.g. `kv/upstream/acme-github-prod`.
- `bindings(repo_id fk, upstream_id fk, remote_name text, direction enum('push','fetch','mirror'), last_sync_at)` — pk(repo_id, upstream_id, remote_name).

**Vault KV layout at `kv/upstream/<id>`:**
- `token` (GitHub PAT / GitLab access token / Gitea token)
- `ssh_private_key` (if SSH auth)
- `webhook_secret` (the HMAC secret we use when registering webhooks at the upstream)

## 6. Helm charts + Argo CD apps

Five new Helm charts under `impl/helixgitpx-platform/helm/`: `repo-service`, `git-ingress`, `adapter-pool`, `webhook-gateway`, `upstream-service`. All extend the hello pattern. `git-ingress` needs a PVC mount for repo storage (`/var/helix/repos`, 50 GiB local / 500 GiB staging) and installs `git` in its distroless image (switch to `alpine:3.20` base for this one service).

Five Argo CD Applications at wave 9 (peers of orgteam/audit).

## 7. Kafka topics

- `repo.events` (6p, 7d retention) — added to `kafka-cluster/values.yaml`.
- `upstream.webhooks` (6p, 7d retention).

## 8. Error handling

- Git wire errors (stderr from `git-http-backend`) surfaced to the client with the same HTTP status the backend chose, plus a correlation-id header that appears in our logs.
- Adapter provider errors map: 401/403 → `codes.PermissionDenied`, 429 → `codes.ResourceExhausted` with `Retry-After` header, 5xx → `codes.Unavailable` + retry.
- Webhook HMAC failure → 401 + logged audit event `webhook.hmac_failed`.

## 9. Testing

| Layer | Tool | M4 ships |
|---|---|---|
| Unit (Go) | testing + testify | repo domain (refs/branch protection), quota token-bucket, HMAC verify, webhook canonicalisation |
| Contract | go-vcr cassettes | GitHub/GitLab/Gitea adapters; committed `testdata/*.yaml` |
| Integration | testcontainers + real `git` binary | git-ingress end-to-end push via `exec.Command("git", "push", ...)` |
| Helm | `helm unittest` | 5 new charts |
| E2E | `verify-m4-spine.sh` | push → replicate → webhook round-trip |

## 10. Completion matrix (16 items)

| # | Item | Artifact | Gate |
|---|---|---|---|
| 54 | Repo CRUD + event sourcing | `services/repo/` + `repo.outbox_events` table | gRPC round-trip + event on `repo.events` |
| 55 | Refs + branch protection | `RefService` + `repo.branch_protections` | Protect → push rejected on unsigned commit |
| 56 | Presigned LFS | `git-ingress` LFS handler + `repo.lfs_objects` | `git lfs push` → MinIO bucket has object |
| 57 | git-upload/receive-pack proxy | `git-ingress` binary + `git-http-backend` | `git clone` + `git push` succeed |
| 58 | Quota + rate-limits | Redis token-bucket in `git-ingress` | 429 when rate exceeded |
| 59 | Signed-push verify | SSH sig verify against Keycloak keys | Unsigned push rejected when required |
| 60 | Adapter interface + plumbing | `internal/adapter/adapter.go` + dispatch | Interface satisfied by all 3 providers |
| 61 | GitHub adapter | `internal/providers/github/` | Contract tests green (replay) |
| 62 | GitLab adapter | `internal/providers/gitlab/` | Contract tests green |
| 63 | Gitea adapter | `internal/providers/gitea/` | Contract tests green |
| 64 | Adapter contract tests | `testdata/*.yaml` cassettes | `go test` green |
| 65 | Webhook receivers + HMAC + dedup | `services/webhook-gateway/` | `curl -H Signature ...` → 200; replay → no-op |
| 66 | Canonicalisation → Kafka | Common `WebhookEvent` proto + producer | Event on `upstream.webhooks` |
| 67 | Upstream CRUD API | `services/upstream/` | gRPC round-trip |
| 68 | Credentials in Vault | `kv/upstream/<id>` layout | No credentials in Postgres; only `vault_path` |
| 69 | Repo↔upstream bindings | `upstream.bindings` table + RPCs | Bind call → row; Unbind → deleted |

## 11. Exit criteria

- verify-m4-cluster: 16/16 green on artifact presence.
- User journey: create repo via repo-service → bind to a GitHub upstream (creds in Vault) → `git push helix.local/acme/example main` → adapter-pool mirrors to GitHub → GitHub sends webhook → `webhook-gateway` publishes to `upstream.webhooks` → consumer prints the event.

## 12. Risks

| Risk | Mitigation |
|---|---|
| `git-http-backend` CGI process mgmt subtleties | Use `os/exec` with explicit stdin/stdout, strict timeouts, cleanup on ctx cancel |
| go-vcr cassettes go stale when providers change APIs | Quarterly re-record via `RECORD=1 go test`; contract tests gate CI only; prod is validated by staging |
| LFS large uploads bypass quota (presigned URL direct to MinIO) | Presigned URL TTL = 5 min; quota applied at presign-request time |
| Signed-push verify requires Keycloak admin API | auth-service adds a `ListUserKeys(subject)` RPC in M4 Task list |
| MinIO presign for repos not in MinIO ACL | Per-repo IAM policy generated on repo create; documented as M4 Task |

## 13. References

- Roadmap §5: `docs/specifications/.../13-roadmap/17-milestones.md`
- M1/M2/M3 specs under `docs/superpowers/specs/`
- https://git-scm.com/docs/git-http-backend
- https://github.com/go-git/go-git
- https://github.com/dnaeon/go-vcr

— End of M4 design —
