# 29 — API Versioning & Deprecation Policy

> **Document purpose**: Define how HelixGitpx **evolves its public interfaces** — gRPC, REST, WebSocket, webhooks, SDK artefacts, and Kafka event schemas — so customers can rely on them for years with predictable change windows.

---

## 1. Why This Matters

Customers build CI pipelines, IDE plugins, mobile apps, and business-critical automation against our APIs. Breaking those silently is a trust-destroyer. This document is a commitment.

---

## 2. Scope & Versioning Axes

| Surface | Version unit | Example |
|---|---|---|
| Protobuf services | **package major** | `helixgitpx.v1`, `helixgitpx.v2` |
| REST | **URL prefix** | `/api/v1/*`, `/api/v2/*` |
| WebSocket / live events | **header + url prefix** | `/events/v1`, `sec-websocket-protocol: helixgitpx.v1` |
| Webhooks sent by us | **payload schema version field** | `schema_version: 1` |
| Kafka event schemas | **Avro compatibility mode** + `.vN` topic suffix on breaking | `helixgitpx.repo.ref.updated.v2` |
| SDKs | **semver per package** | `@helixgitpx/sdk@1.7.3` |
| CLI | **semver** | `helixctl 1.2.0` |
| Plugin interfaces (WIT) | **WIT version** | `helixgitpx:adapter/universal-git@1.0.0` |

---

## 3. Support Windows

| Surface | Stable lifetime | Announcement before removal |
|---|---|---|
| Protobuf package major | **≥ 24 months** after next major GA | **12 months** |
| REST URL prefix | **≥ 24 months** after next major GA | **12 months** |
| Webhook payload major | **≥ 12 months** | **6 months** |
| Kafka topic major | **≥ 12 months** | **6 months** |
| SDK major | **≥ 18 months** | **12 months** |
| CLI major | **≥ 18 months** | **12 months** |
| Plugin WIT major | **≥ 24 months** | **12 months** |
| Preview / beta features | Not stable; marked `x-beta` / `@Experimental`; can change any time | Where practical |

Anything marked `x-preview` / `x-beta` / `@Experimental` is explicitly outside these guarantees.

---

## 4. What Counts as a Breaking Change?

Any of the following to a stable interface:

- Removing / renaming a field, method, RPC, event type, endpoint, or plugin import.
- Narrowing a type (e.g. `string` → enum) or making an optional field required.
- Changing semantics of an existing field or operation.
- Changing the URL, query param name, or required header of an endpoint.
- Lowering a limit (e.g. max request size) in a way that previously valid traffic fails.
- Strengthening authz (a request that previously worked now requires additional scope) — unless justified by security and announced.

These are **allowed** and non-breaking:

- Adding new optional fields.
- Adding new RPCs / endpoints / event types.
- Adding new enum values (clients must handle "unknown").
- Loosening validation (accepting what was previously rejected).
- Performance improvements with the same observable contract.

---

## 5. How We Enforce

### 5.1 Protobuf

- **buf breaking** runs in CI with `FILE` compatibility.
- PRs that break a stable package are rejected unless they target an `_unstable.proto` file or a new major.
- Reviewers must approve any deprecation annotation; a bot files a ticket to plan removal.

### 5.2 REST

- OpenAPI spec in Git; `oasdiff` runs in CI comparing against latest tagged release.
- Breaking changes require a `BREAKING:` commit prefix + issue with migration notes.
- Contract tests (Dredd / schemathesis) exercise the deprecated surface to ensure it keeps working until sunset.

### 5.3 Kafka

- Karapace registry enforces `BACKWARD` compatibility by default per subject.
- A new incompatible schema requires a new topic (`.v2`) and a consumer migration plan.

### 5.4 WebSocket / Live Events

- Same envelope design permits additive fields without breaking consumers.
- `schema_version` in the envelope allows consumers to branch.

### 5.5 Webhooks

- Our outgoing webhooks include `helixgitpx-schema-version` header.
- Customers subscribe to specific major versions; default latest.
- Dual-deliver during transitions (old + new) for 90 days on request.

---

## 6. Deprecation Mechanics

1. **Announce**:
   - Add `deprecated: true` to proto/openapi.
   - Emit `Deprecation` + `Sunset` HTTP response headers per RFC 8594.
   - Blog post, mailing list, in-app banner.
2. **Instrument**:
   - Counter `helixgitpx_api_deprecated_usage_total{api, version}` per call.
   - Dashboard monitored by API ownership.
3. **Help**:
   - Migration guide committed in the same PR that introduces the deprecation.
   - SDK codemods where possible.
4. **Enforce**:
   - At sunset date: remove, return `410 Gone` (REST) or `UNIMPLEMENTED` (gRPC).
   - Customers at cutoff still calling receive one more 30-day grace via alert, then block.

---

## 7. Release Train & Channels

- **Stable**: this is what customers rely on. Tagged releases; no breaking changes except between major versions.
- **Beta**: opt-in; subject to change; clearly labeled.
- **Alpha / preview**: internal or small-partner access; may disappear.
- **Nightly**: snapshots; no guarantees.

Each channel is separately accessible via SDK version tags and feature flags.

---

## 8. Coordinated Changes

When changes span multiple surfaces, we coordinate:

- Add new surface first.
- Let customers migrate.
- Remove old surface at sunset.

Example, renaming `RepoService.ListRepos` to `RepoService.SearchRepos`:

1. Add `SearchRepos`; keep `ListRepos` delegating to the same handler.
2. Deprecate `ListRepos` with sunset + migration guide.
3. Update SDKs; docs link to `SearchRepos`.
4. After 12 months, remove `ListRepos` in the next major.

---

## 9. Evolving the OpenAPI Surface

- Treat OpenAPI as generated from proto via `grpc-gateway` + overlays. Don't hand-edit.
- Add-only rules apply identically; removals require a new URL prefix.
- Extensions like `x-helixgitpx-deprecated-since` carry metadata.

---

## 10. Plugin Interfaces (WIT)

- WIT interfaces use semver.
- Host supports the last two majors of every interface.
- Breaking changes require a new major with 12-month deprecation.
- WIT diffs automatically checked in plugin SDK CI.

---

## 11. SDK Versioning

- Each SDK (TypeScript, Go, Kotlin, Swift, Python, Rust) follows semver independently.
- Majors track the underlying API package major but can release more often (non-breaking SDK-internal changes).
- Auto-generated code kept in a separate subfolder / module to keep hand-written layers stable.

---

## 12. Customer Communication

- `api-announcements@helixgitpx.example.com` mailing list (opt-in by default for org owners).
- RSS feed of announcements at `https://docs.helixgitpx.example.com/announcements.xml`.
- Status page shows upcoming deprecations with countdowns.
- Individual usage reports available to org owners: which of your PATs / integrations are using deprecated endpoints.

---

## 13. Exceptions

The only permitted exceptions to these windows:

- **Security**: a vulnerability requires immediate removal / hardening. Minimum notice given; alternatives offered.
- **Legal / regulatory**: forced by law. We document the requirement.
- **Data integrity**: keeping something available would corrupt data.

All exceptions logged publicly after the fact.

---

## 14. Internal APIs

Anything under `internal/` (Go), `@internal` annotation, or the `_internal.proto` import is explicitly **not** a public API; no compatibility promises. Documented as such.

---

## 15. Supply Chain of SDKs

- SDK releases signed with Cosign + keyless OIDC.
- Published to npm, Maven Central, PyPI, GitHub Releases, crates.io.
- SBOMs attached.
- Reproducible builds documented.

---

## 16. Governance

- API Council: cross-cutting team that reviews proposed breaking changes and sunsets.
- Meets biweekly; decisions recorded as ADRs.
- Public RFC process for large proposals (`docs/rfcs/`).

---

## 17. Versioning FAQ

**Q**: What if a bug means the "stable" behaviour was wrong?
**A**: We fix the bug in a patch release; document if the fix is potentially observable. If the bug's behaviour was relied on, we provide a transition option for one release.

**Q**: Do you support `Accept-Version` headers?
**A**: No; URL / package majors only. Headers are too easy to forget and create hidden coupling.

**Q**: How long for SDK support?
**A**: 18 months post-next-major. Past that, archived — issues closed; source remains for reference.

**Q**: Is "stable" preview-tagged feature X going to break?
**A**: Preview-tagged features are not stable until they leave preview. They can change any time.

---

*— End of API Versioning & Deprecation —*
