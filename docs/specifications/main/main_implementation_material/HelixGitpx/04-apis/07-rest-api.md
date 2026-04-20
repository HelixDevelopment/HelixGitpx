# 07 — REST API Specification

> **Document purpose**: Specify the **public REST API** exposed by HelixGitpx. REST is generated from the same `.proto` sources as gRPC (via `grpc-gateway`) and serves: scripts, legacy integrations, web forms, webhook senders, and anyone who needs HTTP/JSON.

---

## 1. Design Principles

1. **OpenAPI 3.1** single source of truth (generated from protobuf via `protoc-gen-openapiv2`, hand-augmented for descriptions and examples).
2. **Resource-oriented URLs**: `/api/v1/{resource}/{id}/{sub-resource}`.
3. **HTTP verbs**: `GET` read, `POST` create, `PATCH` update (partial), `PUT` replace, `DELETE` remove. `HEAD` for existence. `OPTIONS` for CORS.
4. **JSON request/response**: UTF-8, `application/json`. No custom media types except for events (`text/event-stream` for SSE fallback).
5. **Content negotiation** via `Accept` header. Protobuf over HTTP/1.1 also supported (`application/x-protobuf`) for size-conscious callers.
6. **Consistent errors** in [RFC 7807 Problem Details](https://datatracker.ietf.org/doc/html/rfc7807) format.
7. **Versioned path**: `/api/v1/`. New major version = new path; old versions deprecated with `Deprecation` + `Sunset` headers (RFC 8594 / RFC 9745).
8. **Idempotency**: any non-idempotent call accepts `Idempotency-Key` header; server caches response for 24 h.
9. **Pagination** via opaque cursors; never by page number.
10. **Rate limits** declared in response headers.

---

## 2. Base URL & Paths

- Production: `https://api.helixgitpx.example.com/api/v1/`
- Staging: `https://api.staging.helixgitpx.example.com/api/v1/`
- Per-region: `https://api.eu-west.helixgitpx.example.com/api/v1/` (for data-residency customers)

### Resource Map (partial)

```
/api/v1/
├── auth/
│   ├── login                                   POST
│   ├── refresh                                 POST
│   ├── logout                                  POST
│   ├── me                                      GET
│   ├── mfa/enrol                               POST
│   ├── mfa/verify                              POST
│   ├── pats                                    GET POST
│   ├── pats/{id}                               DELETE
│   ├── sessions                                GET
│   └── sessions/{id}                           DELETE
├── orgs                                        GET POST
├── orgs/{id}                                   GET PATCH DELETE
├── orgs/{id}/teams                             GET POST
├── orgs/{id}/teams/{team_id}                   GET PATCH DELETE
├── orgs/{id}/members                           GET
├── orgs/{id}/members                           POST
├── orgs/{id}/members/{user_id}                 GET PATCH DELETE
├── orgs/{id}/upstreams                         GET POST
├── orgs/{id}/upstreams/{up_id}                 GET PATCH DELETE
├── orgs/{id}/upstreams/{up_id}/enable          POST
├── orgs/{id}/upstreams/{up_id}/disable         POST
├── orgs/{id}/upstreams/{up_id}/rotate          POST
├── orgs/{id}/upstreams/{up_id}/test            POST
├── repos                                       GET POST
├── repos/{id}                                  GET PATCH DELETE
├── repos/{id}/archive                          POST
├── repos/{id}/refs                             GET
├── repos/{id}/refs/{name…}                     GET PATCH DELETE
├── repos/{id}/branch-protections               GET POST
├── repos/{id}/branch-protections/{bp_id}       GET PATCH DELETE
├── repos/{id}/bindings                         GET POST
├── repos/{id}/bindings/{binding_id}            PATCH DELETE
├── repos/{id}/sync                             POST    (trigger)
├── repos/{id}/sync-jobs                        GET
├── repos/{id}/sync-jobs/{job_id}               GET
├── repos/{id}/sync-jobs/{job_id}/cancel        POST
├── repos/{id}/conflicts                        GET
├── repos/{id}/conflicts/{case_id}              GET
├── repos/{id}/conflicts/{case_id}/resolutions  POST
├── repos/{id}/conflicts/{case_id}/ai-proposals POST
├── repos/{id}/pull-requests                    GET POST
├── repos/{id}/pull-requests/{num}              GET PATCH
├── repos/{id}/pull-requests/{num}/reviews      GET POST
├── repos/{id}/pull-requests/{num}/comments     GET POST
├── repos/{id}/pull-requests/{num}/merge        POST
├── repos/{id}/issues                           GET POST
├── repos/{id}/issues/{num}                     GET PATCH
├── repos/{id}/issues/{num}/comments            GET POST
├── repos/{id}/releases                         GET POST
├── repos/{id}/releases/{id}                    GET PATCH DELETE
├── repos/{id}/releases/{id}/assets             GET POST
├── search                                      GET
├── search/code                                  GET
├── search/semantic                              POST
├── search/suggest                               GET
├── events/subscribe                              GET  (SSE or WebSocket upgrade)
├── notifications/channels                        GET POST
├── notifications/channels/{id}                   PATCH DELETE
├── notifications/subscriptions                   GET POST
├── audit/events                                  GET
├── audit/events/{id}                             GET
├── policies                                      GET PATCH
├── billing/usage                                 GET
└── healthz / readyz / livez
```

---

## 3. Authentication

Every request (except public endpoints) **MUST** carry one of:

| Scheme | Header | Notes |
|---|---|---|
| Bearer (OIDC access token) | `Authorization: Bearer <jwt>` | Primary for users |
| Bearer (PAT) | `Authorization: Bearer hpxat_…` | Prefix `hpxat_` to distinguish |
| Cookie (web only) | `Cookie: hpx_session=…` | `HttpOnly`, `Secure`, `SameSite=Lax` |
| mTLS (internal) | certificate | SPIFFE SVID |

Scopes (issued with token) limit what each token can do. Per-scope mapping is in [11-security-compliance.md](../08-security/11-security-compliance.md).

### Login (OIDC code-grant with PKCE)

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "oidc": {
    "issuer": "https://auth.example.com",
    "code": "…",
    "code_verifier": "…",
    "redirect_uri": "https://app.helixgitpx.example.com/callback"
  }
}
```

**Response `200 OK`**:

```json
{
  "tokens": {
    "access_token": "eyJ…",
    "access_token_expires_at": "2026-04-19T12:45:00Z",
    "refresh_token": "rrt_…",
    "refresh_token_expires_at": "2026-05-19T12:30:00Z",
    "token_type": "Bearer"
  },
  "user": { "id": "018f1b8a-…", "email": "…", "username": "…" },
  "mfa_required": false
}
```

---

## 4. Content Types & Encoding

- All bodies are UTF-8 JSON by default.
- `Accept: application/x-protobuf` toggles protobuf wire.
- Binary payloads (LFS, release assets) use `application/octet-stream` with presigned URLs rather than streaming through the API.
- `Content-Encoding: gzip, br, zstd` accepted; server picks best.

---

## 5. Error Format (RFC 7807)

```http
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/problem+json

{
  "type": "https://errors.helixgitpx.example.com/repo/slug-taken",
  "title": "Slug already in use",
  "status": 422,
  "detail": "An active repository with slug 'acme' already exists in this organisation.",
  "code": "repo.slug_taken",
  "instance": "/api/v1/orgs/018f…/repos",
  "trace_id": "c0ffee01…",
  "errors": [
    { "field": "slug", "code": "duplicate", "message": "must be unique" }
  ],
  "doc_url": "https://docs.helixgitpx.example.com/errors/repo/slug-taken"
}
```

Common HTTP codes:

| Code | Meaning |
|---|---|
| 200 | OK |
| 201 | Created (with `Location`) |
| 202 | Accepted (async) |
| 204 | No Content |
| 400 | Bad request (malformed) |
| 401 | Unauthenticated |
| 403 | Permission denied |
| 404 | Not found |
| 409 | Conflict (duplicate, version mismatch) |
| 410 | Gone (deprecated endpoint after sunset) |
| 412 | Precondition failed (If-Match, If-Unmodified-Since) |
| 415 | Unsupported media type |
| 422 | Validation error |
| 423 | Locked (resource busy) |
| 425 | Too early (idempotency reuse detected) |
| 428 | Precondition required |
| 429 | Rate limited |
| 451 | Legal (residency/geo restriction) |
| 500 | Internal error |
| 502/503/504 | Upstream / unavailable |

---

## 6. Pagination

```http
GET /api/v1/orgs/{id}/repos?page_size=50&sort_by=updated_at%20desc
```

**Response**:

```json
{
  "items": [ /* ... */ ],
  "next_page_token": "opaque-cursor",
  "total_count": 2345
}
```

Cursor is an HMAC-signed opaque token encoding `(sort_key, last_id)`. Clients must not attempt to parse it.

---

## 7. Filtering, Sorting, Field Masks

- **Filter** (subset of CEL / [AIP-160](https://google.aip.dev/160)): `?filter=status%3D%22open%22%20AND%20labels%3A%22bug%22`
- **Sort**: `?sort_by=updated_at%20desc,created_at%20asc`
- **Field mask** (sparse responses): `?fields=id,name,created_at`
- **Expand** (include sub-resources): `?expand=author,reviewers`

---

## 8. Conditional Requests

- `ETag` on every GET; clients pass `If-None-Match` for 304s.
- `Last-Modified` + `If-Modified-Since` supported.
- `If-Match` on PATCH/PUT → 412 if the resource version has changed (optimistic concurrency).

---

## 9. Idempotency

Non-idempotent methods (POST creating a resource, POST triggering work) accept:

```
Idempotency-Key: 018f1b8a-6b2c-7f3e-a000-0000deadbeef
```

Server stores `(token, method, path, idempotency_key) → response` for 24 h. A replayed request returns the cached result (HTTP 200 with `X-Idempotent-Replay: true`).

---

## 10. Rate Limiting Headers

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 783
X-RateLimit-Reset: 1713557300
X-RateLimit-Resource: api-token
Retry-After: 30      (on 429)
```

---

## 11. Observability Headers

- `X-Request-Id` (echoed / generated) — matches trace id.
- `X-Trace-Id` — OTel trace id if different.
- `Server-Timing` — phase timings (`db;dur=12.3`, `upstream;dur=140.2`).

---

## 12. CORS

Default allowed origins: `*.helixgitpx.example.com`, configurable per-org. Preflight cached 10 min. Credentials allowed (cookies) for the web app domain only.

---

## 13. Endpoint Examples

### 13.1 Create Repository

```http
POST /api/v1/orgs/018f1b8a.../repos
Authorization: Bearer eyJ...
Idempotency-Key: 018f1b8c...
Content-Type: application/json

{
  "slug": "infra-terraform",
  "display_name": "infra-terraform",
  "description": "Platform IaC",
  "visibility": "private",
  "default_branch": "main",
  "topics": ["terraform","aws","platform"],
  "primary_upstream": "github",
  "auto_bind_all_enabled_upstreams": true
}
```

`201 Created` with `Location: /api/v1/repos/{id}` and the full repo resource.

On success: `repo.created` event published to Kafka → fan-out creates the remote repo on each enabled upstream via Temporal workflow.

### 13.2 List Pull Requests

```http
GET /api/v1/repos/018f1b8a.../pull-requests
  ?filter=state%3D%22open%22
  &sort_by=updated_at%20desc
  &page_size=50
Authorization: Bearer eyJ...
```

### 13.3 Subscribe to Live Events (SSE fallback)

```http
GET /api/v1/events/subscribe?scopes=repo:018f1b8a...,user:self
Accept: text/event-stream
Authorization: Bearer eyJ...
Last-Event-ID: 01HPXK9...
```

**Stream**:

```
id: 01HPXKA1...
event: ref.updated
data: {"repo_id":"018f...","ref":"refs/heads/main","old":"abc","new":"def"}

id: 01HPXKA2...
event: pr.opened
data: {"repo_id":"018f...","number":42,"title":"Add caching"}
```

Primary protocol is gRPC streaming or WebSocket; SSE is for simple scripts.

### 13.4 Resolve Conflict

```http
POST /api/v1/repos/018f1b8a.../conflicts/018f1bc0.../resolutions
Content-Type: application/json

{
  "strategy": "take_ai_proposal",
  "proposal_id": "018f1bc1...",
  "signoff": {
    "by_user_id": "018f1a00...",
    "comment": "Looks correct; no rebase needed."
  }
}
```

---

## 14. Webhook Senders (Incoming)

Every upstream provider has a webhook receiver under `/webhooks/{provider}/{repo_id}`:

- `POST /webhooks/github/{repo_id}` — HMAC SHA-256 via `X-Hub-Signature-256`
- `POST /webhooks/gitlab/{repo_id}` — `X-Gitlab-Token` + payload signature
- `POST /webhooks/gitee/{repo_id}`
- `POST /webhooks/gitflic/{repo_id}`
- `POST /webhooks/gitverse/{repo_id}`
- `POST /webhooks/bitbucket/{repo_id}`
- `POST /webhooks/codeberg/{repo_id}`
- `POST /webhooks/gitea/{repo_id}`
- `POST /webhooks/azure/{repo_id}`

All verify HMAC, deduplicate by `X-Delivery-ID`, and enqueue into Kafka. See [02-services/10-git-provider-integrations.md](../02-services/10-git-provider-integrations.md).

---

## 15. OpenAPI Artefact

- Generated at build: `gen/openapi/helixgitpx.v1.yaml` (and `.json`).
- Published to `https://docs.helixgitpx.example.com/openapi/v1`.
- Swagger UI and Redoc both hosted.
- **Stoplight-style lint** (`spectral`) in CI: every operation has description + example + error cases.

---

## 16. Deprecation Policy

- Endpoints marked `deprecated: true` in OpenAPI.
- Response header: `Deprecation: Wed, 01 Jan 2027 00:00:00 GMT`, `Sunset: Wed, 01 Apr 2027 00:00:00 GMT`.
- At sunset: endpoint returns 410 with a `Link: rel="successor-version"` pointing to the replacement.
- Minimum 6-month window between deprecation and sunset.

---

## 17. SDKs

REST SDKs co-exist with gRPC SDKs; generated via **OpenAPI Generator** for:

- Bash / curl (playground)
- PowerShell
- Python (alternative to gRPC client)
- Ruby
- Java (Spring-friendly)
- PHP
- Terraform provider (separate project but uses REST SDK)

---

## 18. Testing

- **Schema conformance**: `schemathesis` + `Dredd` against staging.
- **Contract**: Pact consumer-driven tests (web, mobile).
- **Fuzz**: Restler + custom corpus.
- **Security**: OWASP ZAP full scan in CI.
- **Performance**: k6 suites per endpoint with SLO assertions.

---

*— End of REST API Specification —*
