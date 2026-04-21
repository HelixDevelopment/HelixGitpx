## 2. Authentication

Every request to the HelixGitpx API requires a token. Two token types:

- **Bearer (OIDC)** — short-lived (1 hour), refreshed via OIDC.
- **Personal Access Token (PAT)** — long-lived, user-scoped.

### 2.1 Obtaining a PAT

Web app → Settings → Personal access tokens → **New token**. Pick a
scope set (`repo:read`, `repo:write`, `admin:org`, `admin:*`, etc.) and
an expiry. Copy the token immediately — it is not stored in plaintext.

CLI equivalent:

```bash
helixgitpx auth pat create \
  --name ci-runner \
  --scopes repo:read,repo:write \
  --expires-in 90d
```

### 2.2 Using a PAT

```http
GET /api/v1/orgs HTTP/1.1
Authorization: Bearer hgx_pat_XXXXXXXXXXXXXXXXXXXX
```

### 2.3 Revoking

`helixgitpx auth pat revoke --id <pat-id>`. Revocation is immediate.

### 2.4 OIDC flow (web / mobile / desktop)

Clients use the standard OIDC authorization-code flow with PKCE. The
Keycloak authorization endpoint is at:

```
https://auth.helixgitpx.io/realms/helixgitpx/protocol/openid-connect/auth
```

Required parameters:

- `client_id` — `helixgitpx-web` / `helixgitpx-cli` / `helixgitpx-mobile` / `helixgitpx-desktop`.
- `response_type=code`.
- `scope=openid email profile offline_access`.
- `redirect_uri` — must be registered on the client.
- `code_challenge` + `code_challenge_method=S256`.

### 2.5 Rate limits

Per [§1.7 in the API reference intro](./00-introduction.md). 429 responses
include `Retry-After` and `X-RateLimit-*` headers.

### 2.6 Service-to-service (SPIFFE)

East-west calls inside the cluster never use PATs. Services carry a
SPIFFE SVID; the auth-service validates the SVID at the ingress
boundary. No human-facing API needed.

---
