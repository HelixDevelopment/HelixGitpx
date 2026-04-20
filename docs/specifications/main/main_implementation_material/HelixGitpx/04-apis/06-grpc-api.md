# 06 — gRPC API Specification

> **Document purpose**: Define the **public gRPC contracts** exposed by HelixGitpx to clients (web, mobile, desktop, CLI, SDKs) and the **internal gRPC contracts** between services. All `.proto` source files live under [17-protos/](../17-protos/) and are the authoritative source.

---

## 1. Protocol Decisions

| Decision | Value |
|---|---|
| Protocol | **gRPC over HTTP/2** with mTLS |
| IDL | **Protocol Buffers v3** |
| Code generation | `buf` (linter + generator) |
| Registry | [Buf Schema Registry (self-hosted)](https://buf.build/product/bsr) |
| REST bridge | `grpc-gateway` generates REST/JSON from the same `.proto` |
| Browser support | `grpc-web` over `Connect-Go` (which speaks gRPC, gRPC-Web, and Connect) |
| Streaming | Server-streaming for events; bidirectional where both ends push |
| Backward compatibility | Buf `breaking` check runs in CI; no breaking changes merged without a new package version |

---

## 2. Versioning Strategy

- Public APIs live in `proto/helixgitpx/v1/`.
- A new major version goes in `v2/` alongside `v1/`; both are served for ≥ 2 quarters.
- Additive changes (new fields, new RPCs) ship within `v1`.
- The Buf Schema Registry enforces: `FILE_SAME_PACKAGE`, `WIRE_JSON`, `RPC_NO_DELETE`, `FIELD_NO_DELETE`, `FIELD_SAME_TYPE`.

---

## 3. Common Types (`helixgitpx/v1/common.proto`)

```proto
syntax = "proto3";
package helixgitpx.v1;
option go_package = "github.com/vasic-digital/helixgitpx/gen/go/helixgitpx/v1;helixgitpxv1";

import "google/protobuf/timestamp.proto";
import "google/protobuf/field_mask.proto";

message UUID { string value = 1; }                     // UUIDv7 string

message Page {
  int32  page_size = 1;                                // 1..200
  string page_token = 2;                               // opaque cursor
  string sort_by = 3;                                  // "created_at desc"
}

message PageResponse {
  string next_page_token = 1;
  int32  total_count = 2;                              // -1 if unknown
}

message Actor {
  oneof kind {
    UUID user_id = 1;
    string service_name = 2;
    bool system = 3;
    bool ai = 4;
  }
  string trace_id = 10;
}

message ErrorDetail {
  string code = 1;                                     // stable machine-readable
  string message = 2;                                  // human readable
  map<string, string> fields = 3;                      // field validation errors
  string doc_url = 4;
}

enum Visibility {
  VISIBILITY_UNSPECIFIED = 0;
  VISIBILITY_PUBLIC      = 1;
  VISIBILITY_PRIVATE     = 2;
  VISIBILITY_INTERNAL    = 3;
}

enum Role {
  ROLE_UNSPECIFIED = 0;
  ROLE_OWNER       = 1;
  ROLE_ADMIN       = 2;
  ROLE_MAINTAINER  = 3;
  ROLE_DEVELOPER   = 4;
  ROLE_VIEWER      = 5;
}
```

---

## 4. Services Overview

| Service | Package | Public? | Purpose |
|---|---|---|---|
| `AuthService` | `helixgitpx.v1` | Yes | Sessions, tokens, PATs |
| `OrgService` | `helixgitpx.v1` | Yes | Orgs, teams, memberships |
| `RepoService` | `helixgitpx.v1` | Yes | Repositories, refs, branch protection |
| `UpstreamService` | `helixgitpx.v1` | Yes | Connect/manage upstream providers |
| `SyncService` | `helixgitpx.v1` | Yes | Trigger and observe sync jobs |
| `ConflictService` | `helixgitpx.v1` | Yes | View and resolve conflicts |
| `PullRequestService` | `helixgitpx.v1` | Yes | PR CRUD + reviews |
| `IssueService` | `helixgitpx.v1` | Yes | Issues + comments |
| `ReleaseService` | `helixgitpx.v1` | Yes | Releases + assets |
| `SearchService` | `helixgitpx.v1` | Yes | Unified search |
| `EventsService` | `helixgitpx.v1` | Yes | Live event subscription |
| `AIService` | `helixgitpx.v1` | Yes | Suggestions, RAG queries, feedback |
| `NotificationService` | `helixgitpx.v1` | Yes | Channels, subscriptions |
| `PolicyService` | `helixgitpx.v1` | Yes | Policy decisions |
| `AuditService` | `helixgitpx.v1` | Yes | Query audit log |
| `BillingService` | `helixgitpx.v1` | Yes | Quotas & usage |
| `AdapterService` | `helixgitpx.internal.v1` | **Internal** | Adapter pool gRPC |

---

## 5. AuthService

```proto
// proto/helixgitpx/v1/auth.proto
service AuthService {
  rpc Login        (LoginRequest)        returns (LoginResponse);
  rpc Refresh      (RefreshRequest)      returns (TokenResponse);
  rpc Logout       (LogoutRequest)       returns (google.protobuf.Empty);
  rpc ValidateToken(ValidateTokenRequest) returns (Principal);
  rpc GetMe        (GetMeRequest)        returns (User);
  rpc EnrollMFA    (EnrollMFARequest)    returns (EnrollMFAResponse);
  rpc VerifyMFA    (VerifyMFARequest)    returns (TokenResponse);

  // Personal Access Tokens
  rpc CreatePAT    (CreatePATRequest)    returns (PATResponse);
  rpc ListPATs     (ListPATsRequest)     returns (PATList);
  rpc RevokePAT    (RevokePATRequest)    returns (google.protobuf.Empty);

  // Sessions
  rpc ListSessions (ListSessionsRequest) returns (SessionList);
  rpc RevokeSession(RevokeSessionRequest) returns (google.protobuf.Empty);
}

message LoginRequest {
  // OIDC code-grant; or "password" for local dev only.
  oneof method {
    OidcCode oidc = 1;
    LocalCredential local = 2;   // dev/testing only
  }
}

message OidcCode {
  string issuer = 1;
  string code = 2;
  string code_verifier = 3;      // PKCE
  string redirect_uri = 4;
}

message LocalCredential {
  string email = 1;
  string password = 2;
}

message LoginResponse {
  TokenResponse tokens = 1;
  User user = 2;
  bool mfa_required = 3;
  string mfa_challenge_id = 4;
}

message TokenResponse {
  string access_token = 1;
  google.protobuf.Timestamp access_token_expires_at = 2;
  string refresh_token = 3;
  google.protobuf.Timestamp refresh_token_expires_at = 4;
  string token_type = 5;           // "Bearer"
}

message User {
  UUID id = 1;
  string email = 2;
  string username = 3;
  string display_name = 4;
  string avatar_url = 5;
  string locale = 6;
  string timezone = 7;
  google.protobuf.Timestamp created_at = 10;
}

message Principal {
  User user = 1;
  repeated string scopes = 2;
  UUID active_org_id = 3;
  string session_id = 4;
  google.protobuf.Timestamp expires_at = 5;
}

message CreatePATRequest {
  string name = 1;
  repeated string scopes = 2;      // e.g. "repo:read", "org:admin"
  google.protobuf.Timestamp expires_at = 3;
}

message PATResponse {
  UUID id = 1;
  string token = 2;                // returned only once
  string name = 3;
  repeated string scopes = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp expires_at = 6;
}
```

### Status Codes

| gRPC status | Meaning | REST map |
|---|---|---|
| `UNAUTHENTICATED` | Missing / invalid token | 401 |
| `PERMISSION_DENIED` | Token lacks scope | 403 |
| `INVALID_ARGUMENT` | Validation failed | 400 |
| `NOT_FOUND` | Resource missing | 404 |
| `ALREADY_EXISTS` | Duplicate | 409 |
| `FAILED_PRECONDITION` | State invalid | 409 |
| `RESOURCE_EXHAUSTED` | Rate limited / quota | 429 |
| `UNAVAILABLE` | Transient | 503 |
| `DEADLINE_EXCEEDED` | Timeout | 504 |

---

## 6. OrgService

```proto
service OrgService {
  rpc CreateOrg   (CreateOrgRequest)    returns (Org);
  rpc GetOrg      (GetOrgRequest)       returns (Org);
  rpc UpdateOrg   (UpdateOrgRequest)    returns (Org);
  rpc DeleteOrg   (DeleteOrgRequest)    returns (google.protobuf.Empty);
  rpc ListOrgs    (ListOrgsRequest)     returns (OrgList);

  rpc CreateTeam  (CreateTeamRequest)   returns (Team);
  rpc GetTeam     (GetTeamRequest)      returns (Team);
  rpc UpdateTeam  (UpdateTeamRequest)   returns (Team);
  rpc DeleteTeam  (DeleteTeamRequest)   returns (google.protobuf.Empty);
  rpc ListTeams   (ListTeamsRequest)    returns (TeamList);

  rpc AddMember   (AddMemberRequest)    returns (Membership);
  rpc UpdateMemberRole(UpdateMemberRoleRequest) returns (Membership);
  rpc RemoveMember(RemoveMemberRequest) returns (google.protobuf.Empty);
  rpc ListMembers (ListMembersRequest)  returns (MemberList);

  // Live updates
  rpc WatchOrg    (WatchOrgRequest)     returns (stream OrgEvent);
}

message Org {
  UUID id = 1;
  string slug = 2;
  string display_name = 3;
  string description = 4;
  Visibility default_visibility = 5;
  string billing_plan = 6;
  string region = 7;
  map<string, string> settings = 8;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
}
```

---

## 7. RepoService

```proto
service RepoService {
  rpc CreateRepo       (CreateRepoRequest)       returns (Repo);
  rpc GetRepo          (GetRepoRequest)          returns (Repo);
  rpc UpdateRepo       (UpdateRepoRequest)       returns (Repo);
  rpc DeleteRepo       (DeleteRepoRequest)       returns (google.protobuf.Empty);
  rpc ArchiveRepo      (ArchiveRepoRequest)      returns (Repo);
  rpc ListRepos        (ListReposRequest)        returns (RepoList);

  rpc ListRefs         (ListRefsRequest)         returns (RefList);
  rpc GetRef           (GetRefRequest)           returns (Ref);
  rpc UpdateRef        (UpdateRefRequest)        returns (Ref);
  rpc DeleteRef        (DeleteRefRequest)        returns (google.protobuf.Empty);

  rpc SetBranchProtection(SetBranchProtectionRequest) returns (BranchProtection);
  rpc GetBranchProtection(GetBranchProtectionRequest) returns (BranchProtection);
  rpc ListBranchProtections(ListBPRequest) returns (BranchProtectionList);
  rpc DeleteBranchProtection(DeleteBPRequest) returns (google.protobuf.Empty);

  rpc WatchRepo        (WatchRepoRequest)        returns (stream RepoEvent);
}

message Repo {
  UUID id = 1;
  UUID org_id = 2;
  string slug = 3;
  string display_name = 4;
  string description = 5;
  Visibility visibility = 6;
  string default_branch = 7;
  string primary_upstream = 8;
  int64 size_bytes = 9;
  bool archived = 10;
  bool template = 11;
  repeated string topics = 12;
  string license = 13;
  google.protobuf.Timestamp created_at = 20;
  google.protobuf.Timestamp updated_at = 21;
}

message Ref {
  UUID repo_id = 1;
  string name = 2;
  string kind = 3;                // "branch", "tag", "other"
  string target_sha = 4;
  bool is_protected = 5;
  google.protobuf.Timestamp last_updated_at = 6;
  string last_origin = 7;         // "local", "upstream:github", …
}
```

### 7.1 WatchRepo Streaming

`WatchRepo` is a server-streaming RPC that emits every event affecting the repo (ref updates, PR changes, issue activity, releases, sync outcomes). Clients send:

```proto
message WatchRepoRequest {
  UUID repo_id = 1;
  string resume_token = 2;              // empty = start from "now"
  repeated string event_types = 3;      // filter (empty = all)
}

message RepoEvent {
  string event_id = 1;
  string event_type = 2;                // "ref.updated", "pr.opened", ...
  google.protobuf.Timestamp occurred_at = 3;
  google.protobuf.Any payload = 4;      // concrete type per event_type
  string resume_token = 5;              // store + send on reconnect
}
```

Server drains Kafka via the live-events-service, applies ACL filtering via policy-service, and flows events to the client. On disconnect, clients reconnect with the last `resume_token`.

---

## 8. UpstreamService

```proto
service UpstreamService {
  rpc RegisterUpstream  (RegisterUpstreamRequest)  returns (Upstream);
  rpc UpdateCredentials (UpdateCredentialsRequest) returns (google.protobuf.Empty);
  rpc RotateCredentials (RotateCredentialsRequest) returns (Upstream);
  rpc EnableUpstream    (EnableUpstreamRequest)    returns (Upstream);
  rpc DisableUpstream   (DisableUpstreamRequest)   returns (Upstream);
  rpc DisconnectUpstream(DisconnectUpstreamRequest) returns (google.protobuf.Empty);
  rpc ListUpstreams     (ListUpstreamsRequest)     returns (UpstreamList);
  rpc TestUpstream      (TestUpstreamRequest)      returns (TestResult);

  // Repo-level binding
  rpc BindRepo          (BindRepoRequest)          returns (Binding);
  rpc UpdateBinding     (UpdateBindingRequest)     returns (Binding);
  rpc UnbindRepo        (UnbindRepoRequest)        returns (google.protobuf.Empty);
  rpc ListBindings      (ListBindingsRequest)      returns (BindingList);

  rpc WatchUpstream     (WatchUpstreamRequest)     returns (stream UpstreamEvent);
}

message Upstream {
  UUID id = 1;
  UUID org_id = 2;
  string provider_kind = 3;     // "github", "gitlab", ...
  string display_name = 4;
  string base_url = 5;
  string auth_method = 6;       // "oauth2", "pat", "app_token", "ssh_key"
  string health_status = 7;
  int32 rate_limit_remaining = 8;
  bool shadow_mode = 9;
  bool enabled = 10;
  google.protobuf.Timestamp created_at = 20;
}
```

---

## 9. SyncService

```proto
service SyncService {
  rpc TriggerSync     (TriggerSyncRequest)     returns (SyncJob);
  rpc GetSyncJob      (GetSyncJobRequest)      returns (SyncJob);
  rpc CancelSyncJob   (CancelSyncJobRequest)   returns (google.protobuf.Empty);
  rpc ListSyncJobs    (ListSyncJobsRequest)    returns (SyncJobList);
  rpc ReplayFailedSync(ReplayFailedSyncRequest) returns (SyncJob);

  rpc WatchSyncJob    (WatchSyncJobRequest)    returns (stream SyncJobEvent);
}

message SyncJob {
  UUID id = 1;
  UUID repo_id = 2;
  string trigger = 3;
  string status = 4;
  repeated SyncStep steps = 5;
  google.protobuf.Timestamp scheduled_at = 10;
  google.protobuf.Timestamp started_at = 11;
  google.protobuf.Timestamp completed_at = 12;
}

message SyncStep {
  UUID upstream_id = 1;
  string operation = 2;
  string status = 3;
  string error_code = 4;
  string error_message = 5;
  int32 refs_updated = 6;
  int64 bytes_in = 7;
  int64 bytes_out = 8;
}
```

---

## 10. ConflictService

```proto
service ConflictService {
  rpc ListConflicts   (ListConflictsRequest) returns (ConflictList);
  rpc GetConflict     (GetConflictRequest)   returns (Conflict);
  rpc ProposeResolution(ProposeResolutionRequest) returns (Resolution);
  rpc ApplyResolution (ApplyResolutionRequest) returns (Resolution);
  rpc AskAI           (AskAIRequest)         returns (AIProposal);
  rpc SubmitFeedback  (SubmitFeedbackRequest) returns (google.protobuf.Empty);
  rpc WatchConflicts  (WatchConflictsRequest) returns (stream Conflict);
}

message Conflict {
  UUID id = 1;
  UUID repo_id = 2;
  string kind = 3;
  string subject = 4;
  google.protobuf.Any left = 5;
  google.protobuf.Any right = 6;
  google.protobuf.Any base = 7;
  string status = 8;
  google.protobuf.Timestamp detected_at = 9;
  repeated AIProposal proposals = 10;
}

message AIProposal {
  string model = 1;
  double confidence = 2;
  google.protobuf.Any proposed = 3;
  string rationale = 4;
  google.protobuf.Timestamp created_at = 5;
}
```

---

## 11. EventsService

This is the **primary live-events gateway** for clients.

```proto
service EventsService {
  rpc Subscribe(SubscribeRequest) returns (stream Event);
  rpc Ack      (AckRequest)       returns (google.protobuf.Empty);
  rpc Resume   (ResumeRequest)    returns (stream Event);
}

message SubscribeRequest {
  repeated Scope scopes = 1;      // what to listen to
  repeated string event_types = 2; // filter (empty = all)
  string resume_token = 3;         // empty = start now
  int32 buffer_hint = 4;           // client desired buffer size
}

message Scope {
  oneof kind {
    UUID user_id = 1;
    UUID org_id = 2;
    UUID repo_id = 3;
    bool global = 4;              // admin only
  }
}

message Event {
  string event_id = 1;
  string event_type = 2;
  string resume_token = 3;
  google.protobuf.Timestamp occurred_at = 4;
  UUID tenant_id = 5;
  google.protobuf.Any payload = 6;
  map<string, string> attributes = 7;
}

message AckRequest {
  repeated string event_ids = 1;
}
```

Connection behaviour, resume semantics, flow control — see [08-live-events.md](08-live-events.md).

---

## 12. AIService

```proto
service AIService {
  rpc SuggestConflictResolution(SuggestConflictRequest) returns (AIProposal);
  rpc SummarisePR   (SummarisePRRequest)   returns (AISummary);
  rpc SuggestLabels (SuggestLabelsRequest) returns (LabelSuggestion);
  rpc SearchSemantic(SearchSemanticRequest) returns (SearchResults);
  rpc Chat          (stream ChatMessage)   returns (stream ChatMessage);  // chatops
  rpc SubmitFeedback(SubmitAIFeedbackRequest) returns (google.protobuf.Empty);
  rpc ListModels    (ListModelsRequest)    returns (ModelList);
}
```

---

## 13. SearchService

```proto
service SearchService {
  rpc Search        (SearchRequest)        returns (SearchResults);
  rpc Suggest       (SuggestRequest)       returns (SuggestResponse);
  rpc CodeSearch    (CodeSearchRequest)    returns (CodeSearchResults);
  rpc SemanticSearch(SemanticSearchRequest) returns (SearchResults);
  rpc HybridSearch  (HybridSearchRequest)  returns (SearchResults);
}
```

---

## 14. Internal (`helixgitpx.internal.v1`) — AdapterService

```proto
service AdapterService {
  rpc PushRefs     (PushRefsRequest)     returns (PushRefsResponse);
  rpc FetchRefs    (FetchRefsRequest)    returns (FetchRefsResponse);
  rpc CreateRemoteRepo(CreateRemoteRepoRequest) returns (RemoteRepo);
  rpc DeleteRemoteRepo(DeleteRemoteRepoRequest) returns (google.protobuf.Empty);
  rpc CreateRemotePR(CreateRemotePRRequest) returns (RemotePR);
  rpc MergeRemotePR (MergeRemotePRRequest) returns (google.protobuf.Empty);
  rpc CreateRemoteIssue(CreateRemoteIssueRequest) returns (RemoteIssue);
  rpc UpdateRemoteIssue(UpdateRemoteIssueRequest) returns (google.protobuf.Empty);
  rpc ListRemoteRepos(ListRemoteReposRequest) returns (stream RemoteRepo);
  rpc Ping(PingRequest) returns (PingResponse);
}
```

---

## 15. Interceptors (Go)

Every service applies this chain:

```go
opts := []grpc.ServerOption{
  grpc.ChainUnaryInterceptor(
    otel.UnaryServerInterceptor(),
    logging.UnaryServerInterceptor(),
    recovery.UnaryServerInterceptor(),
    auth.UnaryServerInterceptor(),      // validate token, populate context
    ratelimit.UnaryServerInterceptor(),
    tenancy.UnaryServerInterceptor(),   // set PG RLS GUC
    policy.UnaryServerInterceptor(),    // OPA pre-check
    metrics.UnaryServerInterceptor(),
    validator.UnaryServerInterceptor(), // protoc-gen-validate
  ),
  grpc.ChainStreamInterceptor(
    otel.StreamServerInterceptor(),
    logging.StreamServerInterceptor(),
    recovery.StreamServerInterceptor(),
    auth.StreamServerInterceptor(),
    tenancy.StreamServerInterceptor(),
    metrics.StreamServerInterceptor(),
  ),
  grpc.Creds(spireTLSCreds),
  grpc.KeepaliveParams(keepalive.ServerParameters{Time: 30 * time.Second}),
}
```

Error mapping via `helix-platform/errors` — typed domain errors translate to gRPC statuses with `status.WithDetails` carrying `ErrorDetail` proto.

---

## 16. Transport & TLS

- **Internal (service-to-service)**: mTLS with SPIFFE SVID; cert rotation hourly.
- **External (client-facing)**: TLS 1.3 with public CA cert; cert-manager.
- **Browser**: Connect-Go supports gRPC-Web transparently.
- **Mobile/desktop**: native gRPC over HTTP/2.
- **Firewall-hostile**: WebSocket tunnel fallback (see [08-live-events.md](08-live-events.md)).

---

## 17. Rate Limiting

- Per-token: default 1000 req/min, higher for enterprise.
- Per-IP: default 100 req/min (anonymous endpoints).
- Streaming RPCs: max 50 concurrent streams per token.
- Enforced at api-gateway with sliding-window Redis counters + response header `X-RateLimit-*`.

---

## 18. Pagination

All `List*` RPCs use `Page` / `PageResponse` with opaque tokens. Cursors are **deterministic** (HMAC-signed to prevent tampering) and encode `(sort_key, last_id)` to survive insertions.

---

## 19. Idempotency

All write RPCs accept an `idempotency_key` in request metadata (gRPC metadata: `x-idempotency-key`). Servers cache the response (keyed by `(token, rpc, key)`) for 24 h in Redis. Replay returns the cached response without re-executing.

---

## 20. Client SDKs

Generated from `.proto` via `buf generate`:

| Language | Output | Publish |
|---|---|---|
| Go | `github.com/vasic-digital/helixgitpx-go` | pkg.go.dev |
| TypeScript (Connect) | `@helixgitpx/client` | npm |
| Kotlin | `io.helixgitpx:client-kotlin` | Maven Central |
| Swift | `HelixGitpxClient` | SwiftPM |
| Python | `helixgitpx` | PyPI |
| Rust | `helixgitpx` | crates.io |

All SDKs include: auth helper, reconnect-with-resume logic, retry with backoff, OTel wiring, typed errors.

---

## 21. Testing

- **Unit**: mock interceptors; table-driven per RPC.
- **Contract**: `buf breaking` in CI; schema diff posted as PR comment.
- **Integration**: run full service with Testcontainers deps; call real gRPC endpoint.
- **Fuzz**: `go-fuzz` against each handler with random valid/invalid protobuf bytes.
- **Property**: `gopter` invariants (idempotency, pagination stability).
- **Load**: `ghz` — 1 k RPS sustained for each endpoint in staging.
- **Chaos**: kill one server mid-stream; client must reconnect and resume without data loss.

---

*— End of gRPC API Specification —*
