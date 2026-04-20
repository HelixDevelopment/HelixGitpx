# HelixGitpx — Documentation Suite

<p align="center">
  <img src="assets/logo.png" alt="HelixGitpx" width="240"/>
</p>

<p align="center">
  <strong>A self-healing, AI-assisted, multi-upstream Git federation platform.</strong><br>
  <em>Part of the <a href="https://github.com/vasic-digital">Helix Development</a> ecosystem.</em>
</p>

---

## Document Classification

| Attribute | Value |
|---|---|
| **Document Class** | Enterprise Engineering Specification Suite |
| **Project Codename** | `helixgitpx` |
| **Parent Program** | Helix Development Platform |
| **Version** | v1.0.0 (Implementation-Ready) |
| **Status** | `APPROVED` — supersedes `Git_Proxy_Master_Specification.md` |
| **Owner** | vasic-digital |
| **License** | Apache-2.0 (code) / CC-BY-SA-4.0 (docs) |

---

## How to Read This Suite

HelixGitpx is a large, multi-service, multi-platform system. The documentation is structured so that **each role can find everything they need without reading everything**. Use this index as your routing table.

| Role | Start Here | Then Read |
|---|---|---|
| **Executive / Product** | [01-vision-scope-constraints.md](00-core/01-vision-scope-constraints.md) | [13-roadmap/17-milestones.md](13-roadmap/17-milestones.md) |
| **Architect** | [02-system-architecture.md](01-architecture/02-system-architecture.md) | [03-microservices-catalog.md](02-services/03-microservices-catalog.md), [14-diagrams/](14-diagrams/) |
| **Backend Engineer** | [03-microservices-catalog.md](02-services/03-microservices-catalog.md) | [04-data-model.md](03-data/04-data-model.md), [05-event-streaming.md](03-data/05-event-streaming.md), [06-grpc-api.md](04-apis/06-grpc-api.md) |
| **Frontend Engineer (Web)** | [12-frontend-angular.md](05-frontend/12-frontend-angular.md) | [07-rest-api.md](04-apis/07-rest-api.md), [08-live-events.md](04-apis/08-live-events.md) |
| **Mobile / Desktop Engineer** | [13-mobile-desktop-kmp.md](06-mobile/13-mobile-desktop-kmp.md) | [06-grpc-api.md](04-apis/06-grpc-api.md) |
| **AI / ML Engineer** | [10-llm-self-learning.md](07-ai/10-llm-self-learning.md) | [09-conflict-resolution.md](02-services/09-conflict-resolution.md) |
| **Security Engineer** | [11-security-compliance.md](08-security/11-security-compliance.md) | [15-testing-strategy.md](10-testing/15-testing-strategy.md) |
| **SRE / DevOps** | [16-infrastructure-scaling.md](11-devops/16-infrastructure-scaling.md) | [17-devops-cicd.md](11-devops/17-devops-cicd.md), [18-observability.md](09-observability/18-observability.md), [19-operations-runbook.md](12-operations/19-operations-runbook.md) |
| **QA Engineer** | [15-testing-strategy.md](10-testing/15-testing-strategy.md) | All API docs |
| **Contributor / OSS Dev** | [20-developer-guide.md](12-operations/20-developer-guide.md) | [02-system-architecture.md](01-architecture/02-system-architecture.md) |

---

## Complete Document Index

### 00 — Core

| # | Document | Purpose |
|---|---|---|
| [01](00-core/01-vision-scope-constraints.md) | Vision, Scope & Constraints | Product mandate, non-negotiable requirements, success metrics, out-of-scope |
| [02](00-core/glossary.md) | Glossary | Terminology reference |

### 01 — Architecture

| # | Document | Purpose |
|---|---|---|
| [02](01-architecture/02-system-architecture.md) | System Architecture | C4 L1-L3, bounded contexts, cross-cutting concerns |
| — | [ADR Log](01-architecture/adr/) | Architecture Decision Records |

### 02 — Services

| # | Document | Purpose |
|---|---|---|
| [03](02-services/03-microservices-catalog.md) | Microservices Catalog | All 18 services with responsibilities, interfaces, scaling profile |
| [09](02-services/09-conflict-resolution.md) | Conflict Resolution Engine | Multi-way merge, CRDT metadata, three-phase reconciliation |
| [10](02-services/10-git-provider-integrations.md) | Git Provider Integrations | Adapter matrix for GitHub / GitLab / Gitee / GitFlic / GitVerse / BitBucket / … |
| [23](02-services/23-wasm-plugin-sdk.md) | WASM Plugin SDK | Third-party-extensible adapters, normalisers, validators |

### 03 — Data & Events

| # | Document | Purpose |
|---|---|---|
| [04](03-data/04-data-model.md) | Data Model (PostgreSQL) | Full schema, indexes, partitions, RLS policies |
| [05](03-data/05-event-streaming.md) | Event Streaming (Kafka) | Topics, schemas, consumer groups, retention, compaction |
| [05b](03-data/05b-search-indexing.md) | Search & Indexing (OpenSearch/Meilisearch/Qdrant) | Index strategy, synonyms, vector search |

### 04 — APIs

| # | Document | Purpose |
|---|---|---|
| [06](04-apis/06-grpc-api.md) | gRPC API Specification | Protobuf schemas, services, streaming |
| [07](04-apis/07-rest-api.md) | REST API Specification | OpenAPI 3.1, auth, pagination, idempotency |
| [08](04-apis/08-live-events.md) | Live Events & WebSocket Fallback | Subscription model, resume tokens, backpressure |
| [29](04-apis/29-api-versioning-deprecation.md) | API Versioning & Deprecation | Support windows, breaking change policy, rollout mechanics |

### 05 — Frontend (Web)

| # | Document | Purpose |
|---|---|---|
| [12](05-frontend/12-frontend-angular.md) | Angular Web Application | Architecture, state, design system, a11y, PWA |
| [28](05-frontend/28-accessibility-standards.md) | Accessibility Standards (WCAG 2.2 AA) | Design, testing, manual/automated audits, feedback intake |
| [32](05-frontend/32-i18n-l10n.md) | Internationalisation & Localisation | 20 locales at GA, ICU MessageFormat, RTL, Crowdin process |
| [37](05-frontend/37-design-system.md) | Design System & Tokens | Colour / typography / spacing tokens, cross-platform pipeline |
| [37](05-frontend/37-design-system.md) | Design System & Tokens | Colour, typography, spacing, motion; cross-platform token pipeline |

### 06 — Mobile & Desktop

| # | Document | Purpose |
|---|---|---|
| [13](06-mobile/13-mobile-desktop-kmp.md) | Mobile & Desktop (Kotlin Multiplatform + Compose) | Shared UI + logic across Android / iOS / Windows / macOS / Linux |

### 07 — AI / ML

| # | Document | Purpose |
|---|---|---|
| [10](07-ai/10-llm-self-learning.md) | Self-Learning LLM Platform | Continuous feedback, RLAIF, model registry, drift detection |

### 08 — Security

| # | Document | Purpose |
|---|---|---|
| [11](08-security/11-security-compliance.md) | Security & Compliance | SonarQube, Snyk, zero-trust, SPIFFE/SPIRE, SLSA L3, SOC 2 readiness |
| [24](08-security/24-threat-model.md) | Threat Model | STRIDE per trust boundary, MITRE ATT&CK mapping, residual risk |
| [27](08-security/27-data-retention-privacy.md) | Data Retention, Privacy & DSR | GDPR / CCPA, retention matrix, data-subject rights workflows |
| [36](08-security/36-trust-center.md) | Trust Center | Public security + privacy posture, certifications, disclosures |

### 09 — Observability

| # | Document | Purpose |
|---|---|---|
| [18](09-observability/18-observability.md) | Observability Platform | OTel end-to-end, SLOs, SLIs, dashboards, alerts, runbooks |

### 10 — Testing

| # | Document | Purpose |
|---|---|---|
| [15](10-testing/15-testing-strategy.md) | Testing Strategy (100% Coverage) | Unit / integration / E2E / fuzz / property / security / DDoS / chaos / perf / DR |
| [33](10-testing/33-performance-baselines.md) | Performance Baselines & Benchmarks | Targets, measured baselines, regression guardrails |

### 11 — DevOps

| # | Document | Purpose |
|---|---|---|
| [16](11-devops/16-infrastructure-scaling.md) | Infrastructure & Scaling | K8s topology, multi-region, auto-scaling, cost optimization |
| [17](11-devops/17-devops-cicd.md) | DevOps & CI/CD | GitOps (Argo CD), pipelines, release trains, feature flags |
| [26](11-devops/26-on-prem-deployment.md) | On-Prem & Air-Gapped Deployment | Customer-managed installations, signed bundles, FIPS mode |
| [31](11-devops/31-release-management.md) | Release Management & Branching | Trunk-based, release trains, canaries, hotfix flow, rollback |
| [34](11-devops/34-cost-optimisation.md) | Cost Optimisation | FinOps principles, tenant attribution, levers with impact |

### 12 — Operations

| # | Document | Purpose |
|---|---|---|
| [19](12-operations/19-operations-runbook.md) | Operations Runbook | Incident response, on-call rotation, maintenance windows |
| [20](12-operations/20-developer-guide.md) | Developer Guide | Local dev setup, contribution workflow, code style |
| [21](12-operations/21-user-guide.md) | User Guide | End-user documentation with screenshots |
| [22](12-operations/22-billing-metering.md) | Billing & Metering | Plans, quotas, usage events, invoicing, dunning |
| [25](12-operations/25-migration-guide.md) | Migration Guide | Onto (and off) HelixGitpx from GitHub / GitLab / Bitbucket / etc. |
| [30](12-operations/30-sla.md) | Service Level Agreements | Per-plan uptime, support response, credits, compliance |
| [35](12-operations/35-support-handbook.md) | Support Handbook | Intake, triage, canned responses, playbooks, escalation |
| — | [runbooks/](12-operations/runbooks/) | Standalone pages: RB-010, RB-011, RB-021, RB-100, RB-120, RB-130, RB-140, RB-150, RB-200, RB-301 |

### 13 — Roadmap

| # | Document | Purpose |
|---|---|---|
| [17](13-roadmap/17-milestones.md) | Milestones, Phases & Tasks | 8 milestones, 32 phases, 500+ tasks, dependency graph |
| — | [RFCs](13-roadmap/rfcs/) | Open RFCs for future capabilities |

### 14 — Diagrams

Mermaid source for every diagram referenced in the docs. See [14-diagrams/README.md](14-diagrams/README.md).

### 15 — Reference

| Content | Purpose |
|---|---|
| [ADR Index](15-reference/adr-index.md) | All 26 ADRs with context, decisions, and alternatives |
| [ADR Template](15-reference/adr-template.md) | Standalone template with writing tips |
| [Glossary](15-reference/glossary.md) | Canonical terms |
| [Error Catalog](15-reference/error-catalog.md) | Stable error codes by domain |
| [Metrics Catalog](15-reference/metrics-catalog.md) | Full metrics inventory with labels |
| [Alerts Catalog](15-reference/alerts-catalog.md) | Every alert with expression and runbook link |
| [Kafka Topic Catalog](15-reference/topic-catalog.md) | Every topic: producers, consumers, keys, retention, ACLs |
| [Scaffold Catalog](15-reference/scaffold-catalog.md) | Service / proto / event / UI scaffolds + official integrations |
| [ADR Template](15-reference/adr-template.md) | Standalone ADR template for future decisions |
| [SDK Examples](15-reference/sdk-examples/) | Go, TypeScript, Kotlin, Swift, Python example apps |
| [Feature Flag Catalog](15-reference/feature-flag-catalog.md) | Every flag with owner, strategy, sunset date |
| [Compliance Controls](15-reference/compliance-controls.md) | SOC 2 / ISO 27001 / NIST CSF / GDPR mapping |
| [Chaos Playbook](15-reference/chaos/playbook.md) | Chaos experiments by tier + Game Days |
| [Policies: Rego](15-reference/policies/) | OPA authz + conflict-resolution policies |
| [Policies: Kyverno](15-reference/policies/kyverno-admission.yaml) | K8s admission policies (signing, PSS, hardening) |
| [Load Tests: k6](15-reference/load-tests/) | API mix + git-push fan-out scripts with SLO thresholds |

### 16-18 — Machine-Readable Artefacts

| Path | Content |
|---|---|
| [16-schemas/001_auth.sql](16-schemas/001_auth.sql) | Auth context DDL (users, sessions, MFA, PATs, login attempts) |
| [16-schemas/002_repo.sql](16-schemas/002_repo.sql) | Org / repo / upstream / binding DDL with RLS |
| [16-schemas/003_conflict_collab.sql](16-schemas/003_conflict_collab.sql) | Conflict, PR, issue, CRDT, AI DDL |
| [16-schemas/004_audit_billing_notify_policy.sql](16-schemas/004_audit_billing_notify_policy.sql) | Audit (append-only), billing, notify, policy, sync DDL |
| [16-schemas/events-avro.json](16-schemas/events-avro.json) | Avro subjects for Kafka topics |
| [17-protos/common.proto](17-protos/common.proto) | Shared types (UUID, Page, Actor, Error, enums, Envelope) |
| [17-protos/auth.proto](17-protos/auth.proto) | AuthService (login, refresh, MFA, PAT, sessions, OIDC) |
| [17-protos/repo.proto](17-protos/repo.proto) | RepoService (CRUD, refs, protection, tree/blob, commits, WatchRepo) |
| [17-protos/upstream.proto](17-protos/upstream.proto) | UpstreamService + BindingService + provider capabilities |
| [17-protos/sync.proto](17-protos/sync.proto) | SyncService (triggers, jobs, DLQ) |
| [17-protos/conflict.proto](17-protos/conflict.proto) | ConflictService (propose, apply, undo, escalate, feedback) |
| [17-protos/events.proto](17-protos/events.proto) | EventsService (bidi stream, replay) |
| [17-protos/platform.proto](17-protos/platform.proto) | PolicyService + AuditService + NotifyService + BillingService |
| [17-protos/collab.proto](17-protos/collab.proto) | PRService, IssueService, AIService |
| [18-manifests/docker-compose.yml](18-manifests/docker-compose.yml) | Local dev stack |
| [18-manifests/Chart.yaml](18-manifests/Chart.yaml) | Helm umbrella chart with dependencies |
| [18-manifests/values-staging.yaml](18-manifests/values-staging.yaml) | Staging environment values |
| [18-manifests/values-prod-eu.yaml](18-manifests/values-prod-eu.yaml) | Production eu-west-1 values with FIPS, residency pinning |
| [18-manifests/deployment-sample.yaml](18-manifests/deployment-sample.yaml) | Production Deployment + PDB + ConfigMap |
| [18-manifests/hpa-keda-samples.yaml](18-manifests/hpa-keda-samples.yaml) | Autoscaling: HPA + KEDA ScaledObject + ScaledJob |
| [18-manifests/network-policy-samples.yaml](18-manifests/network-policy-samples.yaml) | Zero-trust K8s + Cilium policies |
| [18-manifests/argo-application.yaml](18-manifests/argo-application.yaml) | GitOps ApplicationSet + canary analysis |
| [18-manifests/kustomize-overlay.yaml](18-manifests/kustomize-overlay.yaml) | Environment overlay patches on top of Helm-rendered base |
| [18-manifests/terraform/](18-manifests/terraform/) | AWS provisioning module (VPC + EKS + RDS/MSK + S3 + KMS + IRSA) |
| [18-manifests/service-template/](18-manifests/service-template/) | Dockerfile + Makefile + .golangci.yml for new Go services |
| [18-manifests/observability/](18-manifests/observability/) | Prometheus recording + alerting rules + Grafana RED dashboard |
| [18-manifests/examples/](18-manifests/examples/) | Temporal FanOutPush workflow + Debezium connector configs |
| [18-manifests/service-template/](18-manifests/service-template/) | Canonical service Dockerfile, Makefile, golangci-lint config |
| [18-manifests/observability/](18-manifests/observability/) | Prometheus recording + alerting rules, Grafana dashboard JSON |
| [18-manifests/examples/](18-manifests/examples/) | Temporal workflow (FanOutPush) + Debezium connector configs |

### Root

| Document | Purpose |
|---|---|
| [CONTRIBUTING.md](CONTRIBUTING.md) | How to contribute (issues, PRs, CLA / DCO, review process) |
| [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) | Contributor Covenant v2.1 adapted + enforcement ladder |

---

## Quick Facts

| Item | Value |
|---|---|
| **Language (backend)** | Go 1.23+ |
| **HTTP framework** | Gin Gonic |
| **RPC** | gRPC + gRPC-Web |
| **Async messaging** | Apache Kafka + Schema Registry |
| **Primary datastore** | PostgreSQL 16 + Timescale extension |
| **Cache / coordination** | Redis 7 / Dragonfly |
| **Search** | OpenSearch (logs, full-text), Meilisearch (UI), Qdrant (vector) |
| **AI** | Self-learning fine-tuned LLMs + LiteLLM router + RLAIF |
| **Web** | Angular 19 + NgRx + Tailwind + Nx |
| **Mobile / Desktop** | Kotlin Multiplatform + Compose Multiplatform (Android / iOS / Windows / macOS / Linux) |
| **Container runtime** | Kubernetes 1.31 + Istio + cert-manager + Argo CD |
| **Code quality** | SonarQube, Snyk, Semgrep, Trivy, Gitleaks |
| **Observability** | OpenTelemetry → Prometheus + Tempo + Loki + Pyroscope |
| **Supply chain** | SLSA L3, Cosign keyless, in-toto, SBOM (CycloneDX) |
| **Test coverage target** | ≥ 100% (per §15) |
| **Microservices count** | 18 core + 7 platform |
| **Supported Git providers** | 12 at GA, extensible via WASM plugins |
| **Scale target** | 100 k orgs / 10 M repos / 1 B events/day |

---

## Versioning

All documents use **semver**. A **MAJOR** bump means breaking changes to public APIs or on-disk formats. A **MINOR** bump means additive changes. A **PATCH** is editorial.

| Artifact | Current Version |
|---|---|
| Documentation Suite | v1.0.0 |
| Protobuf APIs | v1 |
| REST API | v1 |
| PostgreSQL schema | v1.0.0 |
| Kafka event schemas | v1 |
| Web client | v1.0.0 |
| KMP client | v1.0.0 |

---

## How We Verified the Plan

Every recommendation in this suite has been cross-checked against:

1. The prior **Git Proxy Master Specification v4.0.0** (which this supersedes).
2. Known-good production patterns from large-scale Git platforms.
3. Current (2026) state of the OSS ecosystem — every tool named is actively maintained.
4. The explicit technical constraints set by the project owner.
5. Enterprise requirements (SOC 2, ISO 27001, SLSA, SBOM, a11y WCAG 2.2 AA).

Any section marked **[VERIFY-AT-INTEGRATION]** is a place where the underlying facts change frequently (pricing, API versions, cloud quotas) and must be re-checked when that work is actually scheduled.

---

*— End of README —*
