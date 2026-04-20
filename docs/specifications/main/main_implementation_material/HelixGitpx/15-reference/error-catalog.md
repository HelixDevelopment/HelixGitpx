# Error Code Catalog

> Every error in HelixGitpx has a **stable machine-readable code** (`domain.reason`), a human message, HTTP/gRPC mappings, and a doc URL. Adding a new error? Add it here. Changing an existing one? Create a new code; the old one is forever.

---

## Conventions

- Codes are lowercase, dot-separated: `<domain>.<reason>`.
- Codes are **stable forever** — never change meaning.
- New error → new code; old one deprecated with `Deprecated` header pointing to successor.
- Messages are for humans. Clients MUST dispatch on `code`, not message.
- Each code maps to:
  - `http` — HTTP status (for REST)
  - `grpc` — gRPC status
  - `retryable` — idempotency hint for clients
  - `doc` — relative link into this suite
- Errors emitted via `helix-platform/errors` helpers; enforced by linting.

---

## Common / Cross-Service

| Code | HTTP | gRPC | Retryable | Description |
|---|---|---|---|---|
| `common.bad_request` | 400 | INVALID_ARGUMENT | No | Payload malformed |
| `common.validation_failed` | 422 | INVALID_ARGUMENT | No | Field-level validation errors (see `fields`) |
| `common.unauthenticated` | 401 | UNAUTHENTICATED | No | Missing/invalid token |
| `common.permission_denied` | 403 | PERMISSION_DENIED | No | Policy denied |
| `common.not_found` | 404 | NOT_FOUND | No | Resource missing |
| `common.already_exists` | 409 | ALREADY_EXISTS | No | Duplicate |
| `common.precondition_failed` | 412 | FAILED_PRECONDITION | No | ETag/state mismatch |
| `common.conflict` | 409 | ABORTED | Sometimes | Optimistic concurrency |
| `common.rate_limited` | 429 | RESOURCE_EXHAUSTED | Yes (Retry-After) | Too many requests |
| `common.quota_exceeded` | 429 | RESOURCE_EXHAUSTED | No | Plan limit hit |
| `common.unavailable` | 503 | UNAVAILABLE | Yes | Transient |
| `common.timeout` | 504 | DEADLINE_EXCEEDED | Yes | Upstream slow |
| `common.internal` | 500 | INTERNAL | Sometimes | Unhandled |
| `common.not_implemented` | 501 | UNIMPLEMENTED | No | Endpoint missing |
| `common.legal_restriction` | 451 | FAILED_PRECONDITION | No | Region/residency block |

---

## Auth / Identity

| Code | Description |
|---|---|
| `auth.token_expired` | Access token past expiry |
| `auth.token_invalid` | Signature/claim invalid |
| `auth.token_revoked` | Explicitly revoked session |
| `auth.refresh_reuse_detected` | Rotating refresh replayed → entire family revoked |
| `auth.mfa_required` | Step-up needed |
| `auth.mfa_invalid` | Wrong OTP / FIDO assertion |
| `auth.account_locked` | Too many failures or admin lock |
| `auth.session_expired` | Idle timeout |
| `auth.oidc_issuer_unknown` | IdP not in allowlist |
| `auth.pat_revoked` | Personal Access Token revoked |
| `auth.pat_scope_insufficient` | PAT missing required scope |
| `auth.device_binding_mismatch` | Token used from unexpected device |

---

## Org / Team

| Code | Description |
|---|---|
| `org.slug_taken` | Duplicate org slug |
| `org.deleted` | Org tombstoned; resurrect-only by support |
| `org.over_plan_limit` | Action would breach plan limit |
| `team.cycle_detected` | Nested team cycle |
| `membership.last_owner` | Cannot remove last owner |
| `invite.expired` | Invite past expiry |
| `invite.already_accepted` | One-use link reused |

---

## Repo / Ref

| Code | Description |
|---|---|
| `repo.slug_taken` | Duplicate slug in org |
| `repo.archived` | Write blocked on archived repo |
| `repo.size_limit` | Over size quota |
| `ref.not_found` | Ref missing |
| `ref.protected` | Protection rule blocks write |
| `ref.force_push_blocked` | Non-fast-forward + protection |
| `ref.signed_required` | Policy requires signature |
| `ref.non_fast_forward` | Client push rejected |
| `lfs.object_too_large` | Over LFS per-object limit |
| `lfs.mime_rejected` | Rejected by content filter |
| `tag.immutable` | Signed / protected tag cannot be retagged |

---

## Upstream / Adapter

| Code | Description |
|---|---|
| `upstream.provider_unknown` | Not in provider matrix |
| `upstream.disabled` | Operation against disabled upstream |
| `upstream.shadow_mode` | Writes not allowed in shadow mode |
| `upstream.auth_failed` | Provider rejected credential |
| `upstream.rate_limited` | Provider returned 429 |
| `upstream.unreachable` | Network error |
| `upstream.capability_unsupported` | Operation not supported by provider |
| `upstream.credential_rotating` | Transient during rotation |
| `adapter.circuit_open` | Circuit breaker tripped |
| `adapter.plugin_unsigned` | WASM plugin missing signature |
| `adapter.plugin_sandbox_violation` | Plugin tried disallowed action |
| `adapter.push_conflict` | Upstream rejected with conflict |

---

## Sync / Conflict

| Code | Description |
|---|---|
| `sync.job_not_found` | Unknown job id |
| `sync.cancelled` | Operation cancelled before completion |
| `sync.paused` | Upstream paused; drained but not applying |
| `sync.quorum_unmet` | Cannot fan-out to minimum upstreams |
| `conflict.case_not_found` | Unknown conflict id |
| `conflict.already_resolved` | Case terminal |
| `conflict.policy_blocks_auto_apply` | Auto-apply denied; escalated |
| `conflict.ai_confidence_low` | Confidence below threshold |
| `conflict.sandbox_failed` | Proposed patch failed sandbox (compile/lint/test) |
| `conflict.lfs_divergence_quarantined` | Different bytes for same OID; manual review needed |
| `crdt.invalid_op` | Operation malformed |
| `crdt.version_mismatch` | Doc version changed under us |

---

## Pull Request / Issue

| Code | Description |
|---|---|
| `pr.number_taken` | Duplicate PR number (race) |
| `pr.state_invalid` | Action not allowed from current state |
| `pr.required_reviews_missing` | Branch protection |
| `pr.required_checks_failing` | Aggregated checks not passing |
| `pr.upstream_merged_elsewhere` | Already merged on another upstream |
| `issue.locked` | Cannot comment on locked issue |
| `comment.edit_window_passed` | Edit after grace period not allowed |

---

## AI / Search

| Code | Description |
|---|---|
| `ai.model_unavailable` | Inference pool exhausted |
| `ai.budget_exhausted` | Org monthly token budget reached |
| `ai.input_too_large` | Prompt exceeds context window |
| `ai.guardrail_rejected` | Input/output rail rejected |
| `ai.schema_violation` | Structured output invalid |
| `ai.cloud_not_allowed` | Org disallows cloud routing |
| `search.index_unavailable` | Backend store down |
| `search.query_too_complex` | Parser/analyzer error |

---

## Billing

| Code | Description |
|---|---|
| `billing.plan_unknown` | Plan id not found |
| `billing.downgrade_blocked` | Usage above target plan |
| `billing.payment_required` | Past-due account |

---

## Security / Audit

| Code | Description |
|---|---|
| `security.origin_suspicious` | Login from flagged ASN/IP |
| `security.csp_violation_report` | Reporting-only code |
| `audit.retention_prevents_delete` | Retention policy blocks |
| `policy.bundle_missing` | OPA bundle not loaded |
| `policy.decision_failed` | OPA error during eval |

---

## Platform / Infra

| Code | Description |
|---|---|
| `kafka.produce_failed` | Retry with backoff |
| `kafka.dlq` | Moved to DLQ; see runbook RB-021 |
| `db.deadlock` | Transient |
| `db.replication_lag` | Read-from-replica denied due to lag |
| `cache.miss` | Not an error per se; used in metrics |
| `spire.svid_unavailable` | Cannot obtain workload identity |
| `vault.sealed` | Vault sealed; bootstrap required |

---

## Deprecated (example)

| Code | Deprecated | Successor | Removal ETA |
|---|---|---|---|
| `common.request_too_large` | 2026-03-01 | `common.payload_too_large` | 2026-09-01 |

---

## Adding a New Code

1. Append row in the right section.
2. Implement via `helix-platform/errors.New("<code>", ...)`.
3. Add test case in the service that can produce it.
4. Link the error documentation page (`doc_url`) to this catalog.
5. Update SDKs' error-mapping layer if needed.
