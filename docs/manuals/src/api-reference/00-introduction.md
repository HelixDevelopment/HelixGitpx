# HelixGitpx API Reference

## 1. Introduction

This reference describes every public API HelixGitpx exposes. Two wire
formats, one set of contracts:

- **Connect-RPC** — schema-first via Protobuf. Recommended.
- **REST** — auto-generated OpenAPI. For tools that can't speak Connect.

### 1.1 Base URLs

| Environment | URL |
|-------------|-----|
| Production | `https://api.helixgitpx.io` |
| Staging | `https://api.staging.helixgitpx.io` |
| Your self-hosted cluster | `https://<your-ingress>/api` |

### 1.2 Authentication

Bearer token in the `Authorization` header. Tokens come from the OIDC
flow (Keycloak) or from a Personal Access Token minted in the web app.

```http
GET /api/v1/orgs HTTP/1.1
Host: api.helixgitpx.io
Authorization: Bearer eyJhbGciOi...
```

### 1.3 Services

Each service has its own .proto and its own REST subtree.

| Service | Proto | REST base | Purpose |
|---------|-------|-----------|---------|
| hello | `helixgitpx.hello.v1` | `/api/v1/hello` | Liveness echo |
| auth | `helixgitpx.auth.v1` | `/api/v1/auth` | OIDC exchange, token mint |
| orgteam | `helixgitpx.org.v1`, `helixgitpx.team.v1` | `/api/v1/orgs`, `/api/v1/teams` | Tenancy |
| repo | `helixgitpx.repo.v1` | `/api/v1/repos` | Repository CRUD |
| upstream | `helixgitpx.upstream.v1` | `/api/v1/upstreams` | Binding management |
| sync | `helixgitpx.sync.v1` | `/api/v1/sync` | Sync status + retries |
| conflict | `helixgitpx.conflict.v1` | `/api/v1/conflicts` | Conflict inbox |
| collab | `helixgitpx.collab.v1` | `/api/v1/collab` | Real-time docs |
| events | `helixgitpx.events.v1` | `/api/v1/events` | Live event stream |
| ai | `helixgitpx.ai.v1` | `/api/v1/ai` | AI use cases |
| search | `helixgitpx.search.v1` | `/api/v1/search` | Hybrid search |
| billing | `helixgitpx.billing.v1` | `/api/v1/billing` | Plan mgmt |

### 1.4 Versioning

Every service carries a semver triple. Breaking changes bump major and
ship under a new proto package (`v2`, `v3`); old versions remain
available for one major-version cycle.

### 1.5 Pagination

Cursor-based. Every list RPC accepts `page_size` and `page_token`, and
returns `next_page_token`. Never use offsets.

### 1.6 Errors

Connect/gRPC status codes map 1:1 to HTTP codes. Every error carries a
`helixgitpx.common.v1.ErrorDetail` with `code`, `message`, `reason`, and
optional `metadata`.

### 1.7 Rate limits

See [/api on the marketing site](https://helixgitpx.io/api#rate-limits).

---
