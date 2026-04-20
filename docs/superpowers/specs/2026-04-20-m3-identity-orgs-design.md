# M3 Identity & Orgs — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED (pending user review) |
| Author | Милош Васић + Claude (brainstorming session 2026-04-20) |
| Milestone | M3 — Identity & Orgs (Weeks 11–14 in `13-roadmap/17-milestones.md`) |
| Scope | Full 15-item roadmap §4 (items 39–53); nothing skipped |
| Sequencing | Auth first → orgteam → audit → web shell |
| Implements | `docs/specifications/.../13-roadmap/17-milestones.md` §4 |

---

## 1. Context

M1 delivered the monorepo, platform libraries, hello service code, and CI catalog (tag `m1-foundation`). M2 delivered the Core Data Plane manifests + outbox upgrade + 25 Argo CD Applications (tag `m2-core-data-plane`). M3 builds the first real business-logic services: authentication, orgs+teams, audit trail, and a minimal web shell. After M3 the platform has (a) real users, (b) tenant containers (orgs), (c) tamper-evident audit, and (d) a browser entry point — the prerequisites for M4 Git Ingress.

## 2. Goals

G1. Deploy Keycloak v26 in-cluster as the OIDC provider; auto-import the `helixgitpx` realm + two test users on first start.

G2. Build `auth-service` with: OIDC code-flow exchange, RS256 JWT issuance (15-min access + rotating refresh), `hpxat_`-prefixed PATs, TOTP + FIDO2 enrollment, sessions with revocation.

G3. Build `orgteam-service` (auth + orgs + nested teams + memberships + role-based RBAC via OPA bundle v1).

G4. Build `audit-service` that consumes `audit.events` Kafka topic into an append-only PG table and runs an hourly Merkle anchoring CronJob.

G5. Ship a minimal web shell: Angular auth flow via Keycloak, org-list screen, create-org + add-member dialogs, OTel-web traces landing in Tempo.

G6. Verify the exit user journey: a user logs in via OIDC, creates an org, adds a teammate, sees the four resulting audit events in Grafana.

## 3. Non-goals

- Repo, upstream, sync, conflict, collab, live-events services (M4–M5).
- External anchoring of the Merkle root (Bitcoin, Rekor) — M8 work.
- Full accessibility/i18n for the web shell (M6).
- Federation with external SSO providers beyond what Keycloak itself brokers (M3 ships Keycloak; the corporate-SSO brokering is spec-support-only, not configured).

## 4. Locked constraints

| ID | Constraint | Source |
|---|---|---|
| C-1 | Keycloak v26 in-cluster as the OIDC IdP; realm auto-imported from Git | Q1 |
| C-2 | Three services (auth, orgteam, audit) + web-shell updates. `org-service` and `team-service` merged into a single `orgteam-service` binary. | Design §1 rationale |
| C-3 | All mutating RPCs emit audit events via outbox pattern to topic `audit.events` | Design §6 |
| C-4 | OPA bundle v1 enforced in orgteam-service on every mutating RPC | Design §5 |
| C-5 | JWT RS256 with rotating refresh; signing key from Vault | Design §4 |
| C-6 | PATs: `hpxat_` + base62(24B); stored as SHA-256 hash | Design §4 |
| Inherited | GitHub Actions `workflow_dispatch`-only, portable compose wrapper, single git history, `mise.toml`, spine-first sequencing, HA manifests + local overlay, observability-first | M1+M2 |

## 5. Schema additions

Extend `impl/helixgitpx-platform/sql/schemas.sql` with two new schemas: `org`, `team`. Current schemas list already has `auth`, `audit`. CNPG post-sync Job re-applies. New tables owned by per-service roles (`auth_svc`, `orgteam_svc`, `audit_svc`).

New tables (all `ENABLE ROW LEVEL SECURITY`, baseline `USING (true)`):

| Table | Purpose | Columns (summary) |
|---|---|---|
| `auth.users` | Projection of Keycloak identities | id uuid pk, subject text unique, email, display_name, created_at |
| `auth.sessions` | Refresh-token records | id uuid pk, user_id fk, created_at, expires_at, revoked_at nullable, user_agent, ip |
| `auth.pats` | Personal Access Tokens | id uuid pk, user_id fk, name, hashed_secret bytea, scopes jsonb, expires_at, revoked_at |
| `auth.mfa_factors` | TOTP + FIDO2 factors | id, user_id, kind enum('totp','fido2'), secret_or_pubkey bytea, created_at, last_used_at |
| `org.orgs` | Top-level tenants | id uuid pk, slug citext unique, name, created_at |
| `team.teams` | Nested teams in an org | id uuid pk, org_id fk, parent_id fk nullable, slug, name, created_at, unique(org_id, slug) |
| `team.memberships` | User ↔ team with role | id, team_id fk, user_id fk, role enum('viewer','member','admin','owner'), added_at, unique(team_id, user_id) |
| `audit.events` | Append-only audit log | id uuid pk, at timestamptz, actor_user_id, actor_ip inet, action text, target text, details jsonb + DELETE/UPDATE rules do nothing |
| `audit.anchors` | Merkle roots per hour window | id, period tstzrange, merkle_root bytea, external_tx_id nullable, anchored_at |

Append-only enforcement on `audit.events`:

```sql
CREATE RULE audit_events_no_update AS ON UPDATE TO audit.events DO INSTEAD NOTHING;
CREATE RULE audit_events_no_delete AS ON DELETE TO audit.events DO INSTEAD NOTHING;
```

`audit_svc` role bypasses these via a `SECURITY DEFINER` function `audit.append_event(...)` that grants only INSERT. RLS is layered on top.

## 6. Keycloak

New Helm chart `impl/helixgitpx-platform/helm/keycloak/` wrapping `codecentric/keycloakx` v26+. Argo CD wave **4** (before orgteam/audit/auth at wave 8, but after the operators at wave 5 — actually Keycloak belongs in wave 5 as a peer of operators). Final wave assignment:

| Wave | Addition |
|---|---|
| 4 | (new) keycloak-pg (dedicated CNPG Cluster `helix-kc-pg` in `helix-identity`) |
| 5 | (existing M2 operators) + keycloak itself |
| 8 | (new) auth-service |
| 9 | (new) orgteam-service, audit-service |
| 10 | (existing) hello + (new) web-shell updates (same `helixgitpx-web` Application, different values) |

Realm file at `impl/helixgitpx-platform/helm/keycloak/realm/helixgitpx.json`:

- Realm name: `helixgitpx`.
- Clients:
  - `helixgitpx-web` — public, PKCE + authorization code flow. Redirect URIs: `https://web.helix.local/*`, `http://localhost:4200/*`.
  - `auth-service` — confidential, client-credentials flow for Keycloak admin API calls (user lookup).
- Realm roles: `helixgitpx-user` (default), `helixgitpx-admin`.
- Test users (local overlay only; production overlay strips them):
  - `admin@helixgitpx.local` / `admin` / realm role `helixgitpx-admin`
  - `user@helixgitpx.local` / `user`  / realm role `helixgitpx-user`
- Issuer URL: `https://keycloak.helix.local/realms/helixgitpx`.

`up.sh --m3` (an extension of `up.sh --m2`) adds `keycloak.helix.local` and `web.helix.local` + `auth.helix.local` to `/etc/hosts`.

## 7. auth-service

`impl/helixgitpx/services/auth/` scaffolded via `go run ./tools/scaffold --name auth --proto helixgitpx.auth.v1 --http 8002 --grpc 9002 --health 8082`.

### gRPC API (`proto/helixgitpx/auth/v1/auth.proto`)

```proto
service AuthService {
  rpc ExchangeOIDC(ExchangeOIDCRequest) returns (Tokens);
  rpc RefreshToken(RefreshTokenRequest) returns (Tokens);
  rpc IssuePAT(IssuePATRequest) returns (PAT);
  rpc RevokePAT(RevokePATRequest) returns (google.protobuf.Empty);
  rpc ListPATs(google.protobuf.Empty) returns (ListPATsResponse);
  rpc WhoAmI(google.protobuf.Empty) returns (User);
  rpc EnrollTOTP(google.protobuf.Empty) returns (EnrollTOTPResponse);   // returns otpauth URL
  rpc EnrollFIDO2(EnrollFIDO2Request) returns (EnrollFIDO2Response);
  rpc VerifyMFA(VerifyMFARequest) returns (MFAVerification);
  rpc ListSessions(google.protobuf.Empty) returns (ListSessionsResponse);
  rpc RevokeSession(RevokeSessionRequest) returns (google.protobuf.Empty);
}
```

### REST (via `platform/gin`)

- `GET /v1/auth/callback?code&state` — OIDC redirect target. Sets `access_token` cookie (`HttpOnly; Secure; SameSite=Strict`) and `refresh_token` cookie (same attributes, longer expiry).
- `POST /v1/auth/refresh` — rotate tokens.
- `POST /v1/auth/pat` / `DELETE /v1/auth/pat/{id}` / `GET /v1/auth/pats`.
- `POST /v1/auth/mfa/totp/enroll` / `POST /v1/auth/mfa/fido2/enroll` / `POST /v1/auth/mfa/verify`.

### JWT

- RS256, RSA-2048.
- Private key at Vault `kv/auth/jwt#private_pem`. Rotation: new key pair generated quarterly; old key kept for 30 days of validation overlap.
- JWKS published at `GET /.well-known/jwks.json` (derived from public key).
- Access token claims: `sub` (user id), `iat`, `exp`, `iss: "https://auth.helix.local"`, `aud: "helixgitpx"`, `orgs: []` (slug list the user is in).
- Refresh token: opaque 32-byte random, base62-encoded. Maps to `auth.sessions.id`. Rotating: each refresh invalidates the presented token and issues a new one.

### PATs

- Format: `hpxat_` + base62(24 random bytes) — 37 chars total.
- Storage: SHA-256 hash of the secret after the `hpxat_` prefix. Literal token never persisted.
- Comparison: `crypto/subtle.ConstantTimeCompare`.
- Scopes: free-form JSON array; M3 ships `read:repo`, `write:repo`, `admin:org` as validated enums; unknown scopes rejected at issue time.

### MFA

- TOTP: `github.com/pquerna/otp/totp`. Secret 20 bytes. Issuer `HelixGitpx`.
- FIDO2: `github.com/go-webauthn/webauthn` v0.11+. RPID = trust domain (`helixgitpx.local` / `helixgitpx.dev`).
- A user may have multiple factors; login requires any one.
- Mandatory for `helixgitpx-admin` role; optional for `helixgitpx-user` (enforced by auth-service at session-creation time).

### Sessions

- Created at successful OIDC exchange or refresh.
- `auth.sessions` row with `user_agent` + `ip` captured from request.
- `ListSessions` returns active sessions for the caller; `RevokeSession` sets `revoked_at = now()`.
- Refresh-token validation checks `revoked_at IS NULL AND expires_at > now()`.

### New platform package: `platform/auth`

Validates JWTs against a JWKS URL. Used by orgteam-service (and later M4+ services) as a gRPC unary interceptor + Gin middleware.

```go
package auth

type Validator struct { /* caches JWKS, 1h TTL */ }
func New(jwksURL string) *Validator
func (v *Validator) Unary() grpc.UnaryServerInterceptor
func (v *Validator) HTTP() gin.HandlerFunc
```

Claims injected into `context.Context` via the key `auth.UserIDKey`.

## 8. orgteam-service

Single binary exposing two gRPC services.

### gRPC API (proto/helixgitpx/org/v1/org.proto + proto/helixgitpx/team/v1/team.proto)

```proto
service OrgService {
  rpc Create(CreateOrgRequest) returns (Org);
  rpc Get(GetOrgRequest) returns (Org);
  rpc List(ListOrgsRequest) returns (ListOrgsResponse);
  rpc Update(UpdateOrgRequest) returns (Org);
  rpc Delete(DeleteOrgRequest) returns (google.protobuf.Empty);
}

service TeamService {
  rpc Create(CreateTeamRequest) returns (Team);                 // optional parent_id for nesting
  rpc Get(GetTeamRequest) returns (Team);
  rpc List(ListTeamsRequest) returns (ListTeamsResponse);       // org-scoped; optionally filtered to descendants of parent
  rpc Update(UpdateTeamRequest) returns (Team);
  rpc Delete(DeleteTeamRequest) returns (google.protobuf.Empty);
  rpc AddMember(AddMemberRequest) returns (Membership);
  rpc RemoveMember(RemoveMemberRequest) returns (google.protobuf.Empty);
  rpc UpdateMemberRole(UpdateMemberRoleRequest) returns (Membership);
  rpc ListMembers(ListMembersRequest) returns (ListMembersResponse);
}
```

### Nested teams

`parent_id` self-FK with cycle detection: on Create/Update, the service runs:

```sql
WITH RECURSIVE ancestors AS (
    SELECT id, parent_id FROM team.teams WHERE id = $1
  UNION ALL
    SELECT t.id, t.parent_id FROM team.teams t
      JOIN ancestors a ON t.id = a.parent_id
) SELECT 1 FROM ancestors WHERE id = $2;
```

(Rejects if `$2` — the proposed new parent — appears in the ancestor chain of `$1`.)

### Role inheritance

Effective role on a team = max(role across self + all ancestor teams where the user is a member). Computed via a SQL function `team.effective_role(team_id uuid, user_id uuid) RETURNS role_enum` using a recursive CTE. Cached in Redis (`hello:role:<team_id>:<user_id>` — 5-min TTL).

### OPA RBAC

`impl/helixgitpx-platform/opa/bundles/v1/authz.rego`:

```rego
package helixgitpx.authz

default allow := false

# Org owners can do anything in their org
allow if {
    input.action.op == "org.update" or input.action.op == "org.delete"
    some membership in data.memberships
    membership.user_id == input.user.id
    membership.team_id == input.action.org_root_team
    membership.role == "owner"
}

# Team admins can manage members in their team or descendants
allow if {
    startswith(input.action.op, "team.member.")
    some role in input.user.effective_roles
    role.team_id == input.action.team_id
    role.role in {"admin", "owner"}
}

# Viewers can read everything in their team/ancestors
allow if {
    startswith(input.action.op, "read.")
    count(input.user.effective_roles) > 0
}
```

OPA bundle is loaded in-process via `platform/opa.NewEvaluator` (lifted from M1 stub — it's already real). Bundle shipped as ConfigMap at `helm/opa-bundles/`; services mount it and reload on `HUP`.

## 9. audit-service

Single binary. No gRPC surface (read API lands in M4 when there's a consumer). Just:

1. **Kafka consumer** on `audit.events` (NEW topic; added to `helm/kafka-cluster/values.yaml` — partitions: 6, replicas: 3, retention: infinite (deleted only via anchoring policy in M8)).
2. **Insert** into `audit.events` via `audit.append_event(...)` function.
3. **Merkle anchoring CronJob** (separate Deployment, runs hourly): queries `SELECT id, hash(details) FROM audit.events WHERE at >= now() - interval '1 hour' AND at < now() ORDER BY at, id`, builds a binary Merkle tree (SHA-256), writes the root into `audit.anchors(period, merkle_root)`. External anchoring (Rekor/Bitcoin) deferred.

### New topic

Add to `impl/helixgitpx-platform/helm/kafka-cluster/values.yaml`:

```yaml
- name: audit.events
  partitions: 6
  replicas: 3
  retentionMs: -1   # retain forever; cleanup policy "compact" for M8
```

(Local overlay patches to `replicas: 1`.)

### New platform package: `platform/audit`

```go
package audit

type Emitter struct { /* wraps platform/kafka.Producer — or Outbox — to topic audit.events */ }
func NewOutboxEmitter(pool *pgxpool.Pool, aggregate string) *Emitter
func (e *Emitter) Emit(ctx context.Context, action, target string, details map[string]any) error
```

`auth-service` and `orgteam-service` both use the outbox variant so their mutations and audit events land atomically.

## 10. Web shell

`impl/helixgitpx-web/apps/web/src/app/` gains real screens:

- `routes.ts`: `/login`, `/orgs`, `/orgs/:slug` (placeholder), `/auth/callback`.
- `login.component.ts`: renders a "Sign in with HelixGitpx" button; calls `authClient.startOIDC()` which redirects to Keycloak.
- `auth-callback.component.ts`: captures `code`+`state` from query string, posts to auth-service, stores access token in memory + refresh token in HttpOnly cookie (already set by auth-service's REST callback).
- `orgs.component.ts`: lists orgs from `orgteamClient.listOrgs()`; create-org dialog; add-member dialog.
- `app.config.ts`: OTel-web init pointing at Tempo OTLP/HTTP (`https://tempo.helix.local/v1/traces`); trace-id propagation header added to every fetch.
- `nx.json` picks up the new lib targets.

### Connect-Web clients

`buf generate` (from M1) already produces TS bindings for the web workspace. Three new proto domains (auth, org, team) populate under `libs/proto/src/helixgitpx/{auth,org,team}/v1/`.

## 11. Error handling + testing + completion matrix

### Error handling

- All RPC errors go through `platform/errors` (existing). MFA/PAT failures map to `codes.Unauthenticated`. RBAC denials map to `codes.PermissionDenied` with a `policy` detail field.
- Keycloak transient failures → auth-service returns `codes.Unavailable` with retry-after hint.

### Testing

| Layer | Tool | M3 ships |
|---|---|---|
| Unit (Go) | testing + testify | auth JWT sign/verify, PAT hash/verify, OPA policy eval, team nesting cycle-detection |
| Integration (Go) | testcontainers + pgx + kgo | auth token lifecycle against pg/redis/kafka; orgteam CRUD + RBAC end-to-end |
| Contract (proto) | `buf breaking` | green; auth/org/team v1 locked |
| Web | Jest + Angular testing utilities | login + org-list component tests; Playwright placeholder spec for the exit-criteria journey |
| Helm unit | `helm unittest` | keycloak realm import Job; audit-service Deployment |
| E2E | `scripts/verify-m3-spine.sh` | runs the user-journey smoke test against the k3d cluster |

### Completion matrix

| # | Item | Artifact | Gate |
|---|---|---|---|
| 39 | OIDC flow with Keycloak | `helm/keycloak/` + realm file | `curl https://keycloak.helix.local/realms/helixgitpx/.well-known/openid-configuration` returns JSON |
| 40 | JWT issuance RS256 + rotating refresh | `services/auth` + `platform/auth` | `ExchangeOIDC` returns an RS256 JWT validated by the JWKS; `RefreshToken` invalidates the old refresh |
| 41 | PAT endpoints `hpxat_` | `services/auth/internal/handler/*/pat.go` + `auth.pats` table | `IssuePAT` returns a 37-char `hpxat_...` token; `ListPATs` shows it |
| 42 | MFA enrollment (TOTP + FIDO2) | `services/auth/internal/handler/*/mfa.go` | TOTP enroll + verify round-trip succeeds; FIDO2 register endpoint responds with challenge |
| 43 | Sessions + revocation | `auth.sessions` + handler | `ListSessions` returns the current one; `RevokeSession` invalidates refresh |
| 44 | Org CRUD | `services/orgteam` + `org.orgs` | `OrgService.Create/Get/List/Update/Delete` round-trips |
| 45 | Nested teams | `team.teams` + cycle-detection CTE + `TeamService.*` | Creating a cycle returns `FailedPrecondition` |
| 46 | Memberships + roles | `team.memberships` + `AddMember/RemoveMember/UpdateMemberRole` | Role inheritance test: grandparent owner sees child team |
| 47 | OPA bundle v1 | `opa/bundles/v1/authz.rego` + `platform/opa` wired | Viewer blocked from org.update; admin allowed on team.member.add |
| 48 | `audit.events` Kafka consumer | `services/audit` | Consumer group `audit-service` lag < 100 messages |
| 49 | PG append-only table + triggers | `audit.events` + rules | UPDATE/DELETE are no-ops; SELECT works |
| 50 | Merkle anchoring job | `services/audit/internal/merkle.go` + CronJob | `audit.anchors` gains 1 row per hour |
| 51 | Angular auth flow + org list | `apps/web/src/app/{login,auth-callback,orgs}.component.ts` | Playwright test reaches `/orgs` after login |
| 52 | Connect-Go clients generated | `libs/proto/src/helixgitpx/{auth,org,team}/v1/` | `nx build web` succeeds; bindings imported by components |
| 53 | OTel-web wired | `app.config.ts` trace provider | A trace from web→auth→orgteam appears in Tempo |

## 12. Exit criteria

1. `verify-m3-cluster.sh` passes all 15 gates.
2. User journey: visit `https://web.helix.local/` → redirected to Keycloak → log in as `user@helixgitpx.local` → empty `/orgs` → create "Acme" → add `admin@helixgitpx.local` as member → Grafana `audit.events` dashboard shows 4 events (user.login, org.create, team.member.add×2 — the system-created `owners` team + the invite).
3. Four new commits on top of M2's existing ADRs: 0013 (Keycloak in-cluster), 0014 (orgteam merger rationale), 0015 (outbox for audit), 0016 (OPA bundle v1).

## 13. Risks

| Risk | Mitigation |
|---|---|
| Keycloak realm import races with its PG init | Separate `helix-kc-pg` Cluster; readiness-gate + Wait for Kafka via `initContainer`. |
| JWT signing key rotation breaks existing tokens | Dual-key JWKS for 30 days on rotation; access-token TTL is 15 min so only refresh flow is affected, handled by the rotating refresh mechanism. |
| FIDO2 RPID mismatch in local (`helixgitpx.local`) vs cluster DNS | Document that FIDO2 test requires visiting via the `/etc/hosts` alias; TOTP covers the default path. |
| OPA bundle drift between ConfigMap and service | Bundle loaded at startup + re-loaded on SIGHUP; Argo CD triggers a rolling restart on ConfigMap change via `kubectl rollout`. |
| Merkle anchoring with no external ledger means M3 anchoring is trust-us-on-our-word | Explicitly documented in ADR-0015; external anchor lands in M8 with Rekor. |

## 14. Open questions

None — all blocking decisions are locked in §4.

## 15. References

- Roadmap §4: `docs/specifications/.../13-roadmap/17-milestones.md`
- Spec `001_auth.sql`, `004_audit_billing_notify_policy.sql`
- OPA: `docs/specifications/.../08-security/`
- M1 spec: `docs/superpowers/specs/2026-04-20-m1-foundation-design.md`
- M2 spec: `docs/superpowers/specs/2026-04-20-m2-core-data-plane-design.md`

— End of M3 Identity & Orgs design —
