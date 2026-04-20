# 01 — Vision, Scope & Constraints

> **Document purpose**: Define **what HelixGitpx is**, **what it is not**, **who it is for**, and **what non-negotiable technical and product constraints** govern every subsequent engineering decision in this suite.

---

## 1. Product Vision

> **HelixGitpx is a self-healing, AI-assisted, multi-upstream Git federation platform that lets a single source of truth live across every Git host — simultaneously, bi-directionally, and without conflict surprises.**

A developer or organisation publishes code once; HelixGitpx ensures every enabled remote (GitHub, GitLab, Gitee, GitFlic, GitVerse, Bitbucket, Codeberg, SourceHut, Azure DevOps, AWS CodeCommit, self-hosted Gitea/Forgejo, …) is kept in sync. Every *collaboration artefact* — branches, tags, releases, issues, pull requests, reviews, CI/CD metadata — is mirrored, normalised and conflict-resolved. When upstreams disagree, HelixGitpx resolves the conflict according to configurable policy (and in the ambiguous cases, a fine-tuned LLM proposes a resolution that a human signs off in one click).

### 1.1 North-Star Metrics

| Metric | Target at GA | Target at Year 2 |
|---|---|---|
| **Sync lag** (event on any upstream → visible everywhere) | p95 ≤ 5 s | p95 ≤ 1 s |
| **Conflict auto-resolution rate** | ≥ 75 % | ≥ 92 % |
| **Push success rate across all upstreams** | ≥ 99.5 % | ≥ 99.9 % |
| **Data loss events per year** | 0 | 0 |
| **Supported providers out-of-the-box** | 12 | 20+ |
| **Monthly active organisations** | 100 | 10 000 |
| **Repositories under management** | 10 k | 10 M |
| **Security incidents with exposed secrets** | 0 | 0 |

---

## 2. Who It Is For (Personas)

| Persona | Description | Key Jobs-to-Be-Done |
|---|---|---|
| **Multi-Host Maintainer** | Independent OSS author who mirrors their repo to 3+ hosts for resilience and reach. | Push once, appear everywhere; survive one host going down. |
| **Sovereignty-Sensitive Enterprise** | Regulated company that must keep code in-region on self-hosted Git, but also wants public mirror on GitHub for hiring/visibility. | Automated bi-directional sync with audit trail; data residency guarantee. |
| **Dual-Market SaaS Company** | Company selling into both Western and Chinese/Russian markets; must be present on Gitee / GitVerse / GitFlic alongside GitHub / GitLab. | One workflow, N markets. |
| **OSS Foundation / Distro** | Project (Linux distro, Apache, CNCF) that wants signed, federated mirrors for anti-censorship resilience. | Cryptographic provenance; redundancy across jurisdictions. |
| **DevEx Platform Team** | Internal platform team at a large company; uses HelixGitpx to give devs a unified view across legacy Bitbucket Server and new GitHub Enterprise during migration. | Zero-downtime migration; unified permissions. |
| **Solo Developer** | Just wants a great free tier with mirroring to Codeberg + GitHub. | Set-and-forget. |

---

## 3. Scope

### 3.1 In Scope — GA (v1.0)

1. **Git federation**
   - Bi-directional sync of refs (`refs/heads/*`, `refs/tags/*`), LFS, and submodules.
   - Advanced conflict resolution (see §09).
   - Connect / disconnect / enable / disable per upstream, per organisation, per repository.
   - Configurable synchronisation direction (push-only, pull-only, bidirectional).
2. **Collaboration artefacts**
   - Branches, tags, releases.
   - Pull / merge requests with conversation threads, reviews, and status checks.
   - Issues with labels, milestones, assignees.
   - Wiki pages, project boards (where the upstream supports them).
3. **CI/CD integration**
   - Status check forwarding (mirrored in both directions).
   - Workflow triggering across hosts (e.g. push on GitHub triggers GitLab pipeline).
   - Release artefact replication.
4. **User & organisation management**
   - SSO via OIDC (Keycloak, Dex, Authentik, Azure AD, Okta).
   - SCIM 2.0 provisioning.
   - RBAC with scoped tokens (org / repo / ref granularity).
5. **Clients**
   - Web (Angular).
   - Desktop (Windows, macOS, Linux — Compose Multiplatform).
   - Mobile (Android, iOS — Compose Multiplatform).
   - CLI (`helixctl`, cross-platform Go binary).
6. **APIs**
   - gRPC (primary, for services and rich clients).
   - REST (OpenAPI 3.1, for scripts and legacy integrations).
   - Live event subscription (gRPC streaming + WebSocket fallback).
7. **AI assistance**
   - Conflict-resolution proposals.
   - Intent-of-commit summaries.
   - PR review suggestions.
   - Self-learning: fine-tunes on user-accepted suggestions.
8. **Observability & operations**
   - Real-time monitoring with OpenTelemetry.
   - Full audit log.
   - Incident runbooks.
9. **Security**
   - SonarQube + Snyk + Semgrep in CI.
   - SBOM for every release.
   - Cosign keyless signing.
   - SLSA Level 3 build provenance.
10. **Testing**
    - ≥ 100 % code coverage (meaning every line, every branch, plus mutation testing to prove test sensitivity).
    - Security, DDoS, scalability, chaos, benchmarking suites.

### 3.2 Out of Scope (v1.0)

- Hosting Git repositories as the **primary** (HelixGitpx is a federation layer; repos primary-live in Gitea/Forgejo or a selected upstream).
- Running a package registry (delegated to existing registries; HelixGitpx integrates with them).
- Project management beyond what upstreams provide (issues/boards are synced, not re-invented).
- Marketplace / billing platform for third-party apps (post-v1 consideration).
- On-prem “air-gapped” deployment (supported as v1.2 milestone via published OCI bundle).

### 3.3 Explicit Non-Goals

- We do **not** replace GitHub/GitLab. We **augment** them.
- We do **not** store code we don't need to (federation is on-demand; working copies are ephemeral).
- We do **not** make Git itself proprietary. Every artefact is extractable as vanilla Git + open standards.

---

## 4. Non-Negotiable Mandates

These **MUST** hold for every release. A feature that violates any of them is rejected at design review.

### M-1 — Technology Constraints

The following technologies **must** appear in the final production system in the listed role:

| Layer | Technology | Role |
|---|---|---|
| Backend language | **Go 1.23+** | All services |
| HTTP framework | **Gin Gonic** | REST endpoints |
| RPC | **gRPC** + gRPC-Web | Inter-service, service-client (incl. mobile/desktop) |
| Primary datastore | **PostgreSQL 16+** | Transactional state |
| Async messaging | **Apache Kafka** (with Schema Registry) | Event backbone |
| In-memory store | **Redis 7** (or Dragonfly) | Cache, locks, rate limit, Bloom filters |
| Search | **OpenSearch** (logs) + **Meilisearch** (user-facing) + **Qdrant** (vector) | Chosen as preferred open-source alternative to Elasticsearch due to Apache 2.0 license vs ES SSPL. |
| Code quality | **SonarQube** + **Snyk** | CI/CD mandatory gates |
| Architecture | **Microservices** (clean-architecture inside each) | |
| Web | **Angular** 19+ | User-facing |
| Mobile + Desktop | **Kotlin Multiplatform + Compose Multiplatform** | Share UI + logic across 5 platforms |
| IPC on client side | **gRPC** (+ WebSocket fallback where gRPC is impractical) | |
| Supply chain | **Cosign**, **in-toto**, **SLSA L3**, **CycloneDX SBOM** | |
| Observability | **OpenTelemetry** → **Prometheus / Tempo / Loki / Pyroscope** | |

### M-2 — 100 % Test Coverage

- Every Go package: **100 % line + branch coverage**, verified by `go test -cover` + `go-tool-mutation` mutation testing.
- Every Angular feature: **100 % line + branch** with Jest/Karma, **100 % user-path** with Playwright.
- Every KMP module: **100 % line + branch** with Kotest / turbine.
- Every public API endpoint: **contract test** (Pact) + **fuzz test** + **negative-path test**.
- Every Kafka topic: **consumer contract test**.
- Test types required (§15): unit, integration, contract, E2E, fuzz, property, mutation, security (SAST/DAST/SCA/container/IaC), DDoS, load, soak, chaos, DR, accessibility, performance, compatibility, benchmark.

### M-3 — 100 % Schema-Driven APIs

- All APIs defined in `.proto` first (see [17-protos/](../17-protos/)).
- REST generated from protobufs via grpc-gateway (or hand-written but contract-tested against proto).
- OpenAPI generated from protobufs via `protoc-gen-openapiv2`.
- Client SDKs generated for Go / TypeScript / Kotlin / Swift / Python.

### M-4 — Bi-Directional Sync Safety

- **Never** destroy upstream data without explicit, logged, signed authorisation.
- **Never** auto-merge a conflict that the configured policy cannot deterministically handle — always escalate to human or AI-with-human-signoff.
- **Never** lose a push: every received ref is first persisted to the event log before acknowledgement.

### M-5 — Live Reactivity

- Any change anywhere (a push on GitHub, an issue created on GitLab, a label changed on Gitee) **must** be reflected in connected clients within the sync-lag SLO.
- Clients **must** support resume-from-last-event (LSN/offset) after disconnect.

### M-6 — Zero-Trust & Supply-Chain Integrity

- No long-lived secrets in images or on disk. OIDC-based short-lived credentials for all machine-to-machine calls (SPIFFE/SPIRE).
- Every container image, every Git commit produced by the system, every release artefact is signed (Cosign keyless + Gitsign) and recorded in a transparency log (Rekor).
- Mandatory SonarQube + Snyk quality gates in CI; merge is blocked on High/Critical findings.

### M-7 — Multi-Region & Horizontal Scaling

- Every stateful service is sharded or replicated; single-node capacity is never the system's capacity.
- Services are stateless where possible; where stateful, they use leader-election via Kubernetes leases or consensus protocols.
- The system is provably able to scale from a single-VM dev deployment to a 3-region active-active K8s mesh without code changes — only topology changes.

### M-8 — Documentation as Deliverable

- No feature is "done" until its reference API docs, user guide section, runbook, and diagrams are updated.
- Every service owns a README in its repo that **must** cover: purpose, dependencies, config, metrics, SLOs, runbook pointer.
- Every release produces changelogs (Keep a Changelog format) and an **SBOM**.

### M-9 — Accessibility & Internationalisation

- Web and client apps must meet **WCAG 2.2 AA** at launch; AAA for key flows by year 2.
- i18n: English, German, French, Spanish, Portuguese, Serbian, Russian, Chinese (Simplified), Japanese — at GA. RTL (Arabic, Hebrew) by year 2.
- Localisation via ICU MessageFormat; translations managed in Weblate.

### M-10 — Privacy by Design

- User content (code, issues, reviews) is never used to train upstream closed-model LLMs.
- Any on-device telemetry is opt-in, anonymised, and auditable via `helixctl telemetry show`.
- Data residency: organisations select primary region; data and ops stay there unless replication is explicitly enabled.

---

## 5. Guiding Principles

| Principle | Meaning in Practice |
|---|---|
| **Make the right thing the easy thing** | The default config does the secure, conservative thing. |
| **Small deltas, not big bangs** | Changes flow continuously; no quarterly releases. |
| **Own the data contract, not the UI** | Upstream UIs stay; we own the invariants. |
| **Prefer open protocols to proprietary APIs** | Git, WebSocket, OCI, OAuth, OIDC, SCIM. |
| **Everything is an event** | Kafka is the backbone; read models are projections. |
| **Reversible decisions are cheap** | Deploy behind feature flags; kill fast if wrong. |
| **Proof over promise** | Every SLO is monitored; every claim has a dashboard. |
| **Bug reports beat bug fears** | If a risk is real, we write a test that proves it. |

---

## 6. Cross-Cutting Constraints

### 6.1 Performance Envelopes

| Operation | p50 | p95 | p99 |
|---|---|---|---|
| Web page load (shell) | 500 ms | 1.5 s | 3 s |
| REST API call | 50 ms | 250 ms | 500 ms |
| gRPC unary | 20 ms | 100 ms | 250 ms |
| gRPC streaming event delivery | 10 ms | 50 ms | 150 ms |
| Git push (empty commit) → visible everywhere | 2 s | 5 s | 12 s |
| Git push (10 MB) → visible everywhere | 5 s | 15 s | 30 s |
| LLM conflict suggestion (simple) | 2 s | 6 s | 12 s |

### 6.2 Cost Envelopes (per thousand active repos)

| Category | Bootstrap (self-host) | Managed SaaS |
|---|---|---|
| Compute | $0 (Oracle free tier) | $8 |
| Storage | $0 (B2 / R2 free tier) | $2 |
| Egress | $0 (Cloudflare R2) | $0 |
| LLM | $0 (self-hosted Ollama) | $3 |
| **Total** | **$0** | **$13** |

### 6.3 Availability Targets

| Tier | SLO | Error Budget / mo |
|---|---|---|
| Control plane (API) | 99.9 % | 43 min |
| Data plane (sync) | 99.95 % | 21 min |
| Live events stream | 99.9 % | 43 min |
| AI suggestions | 99.0 % (best-effort) | 432 min |

### 6.4 Security & Compliance

| Standard | Commitment |
|---|---|
| SOC 2 Type I | Year 1 |
| SOC 2 Type II | Year 2 |
| ISO/IEC 27001 | Year 2 |
| GDPR | Day 1 (data subject requests via CLI + UI) |
| CCPA | Day 1 |
| SLSA | Level 3 at GA, Level 4 for core by year 2 |

---

## 7. Assumptions

1. Upstream Git providers expose a reasonable public API (rate-limited) and webhook support — this is true of every provider on our target list.
2. Users have the right to mirror to each upstream; HelixGitpx provides the machinery, not the legal basis.
3. Network egress pricing will continue to be the dominant cloud cost — hence the emphasis on Cloudflare R2 / Backblaze for LFS.
4. LLMs continue to improve rapidly; our architecture must be model-agnostic (LiteLLM router) to absorb new capabilities without rewrites.
5. Kubernetes remains the dominant orchestration platform for the 5-year planning horizon.

---

## 8. Dependencies & External Factors

| Dependency | Mitigation |
|---|---|
| Upstream API stability (GitHub, GitLab, Gitee) | Versioned adapter matrix, shadow tests on staging. |
| Kafka availability | Primary on self-hosted Strimzi; commercial fallback to MSK/Confluent. |
| Postgres majors (every ~1 year) | Automated pg_upgrade runbook; logical replication fallback. |
| LLM provider ToS changes | LiteLLM swap; local Ollama always available. |
| Cloud free-tier changes | Documented in [15-reference/free-tiers.md](../15-reference/free-tiers.md); quarterly re-check. |

---

## 9. Change Control

This document is the **product constitution**. Changes require:

1. An RFC in [13-roadmap/rfcs/](../13-roadmap/rfcs/).
2. Sign-off from: Program Lead, Architect, SRE Lead, Security Lead.
3. Updated document with version bump and changelog entry.

---

*— End of Vision, Scope & Constraints —*
