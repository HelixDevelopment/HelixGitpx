# 17 — Milestones, Phases & Task Breakdown

> **Document purpose**: Turn the specification into a **step-by-step delivery plan** — 8 milestones, 32 phases, and 500+ concrete tasks. Each task is small enough to be completed by one engineer in a week and is traceable to a specific document in this suite.

---

## 1. Summary Timeline

| # | Milestone | Duration | Outcome |
|---|---|---|---|
| **M1** | Foundation | 4 weeks | Empty mono-repos, shared libs, CI skeleton, K8s dev cluster |
| **M2** | Core Data Plane | 6 weeks | Postgres, Kafka, Redis, OpenSearch/Meili/Qdrant, Vault, SPIRE, Karapace running with GitOps |
| **M3** | Identity & Orgs | 4 weeks | Auth, Org, Team services; web app login; PAT; audit |
| **M4** | Git Ingress & Adapter Pool | 6 weeks | Repo service, git-ingress, adapter-pool with GitHub + GitLab + Gitea; webhook gateway; first real push |
| **M5** | Federation & Conflict Engine | 8 weeks | Sync orchestrator, conflict resolver, CRDT metadata, multi-upstream working end-to-end; remaining 9 providers at GA; live events |
| **M6** | Frontend & Mobile | 8 weeks | Angular web full-feature, KMP+Compose clients for Android/iOS/Win/macOS/Linux |
| **M7** | AI, Search, Policy | 6 weeks | AI service with self-learning loop; hybrid search; OPA policy-as-code end-to-end |
| **M8** | Scale, Harden, GA | 8 weeks | Multi-region; DR drills; SOC2 Type I; 100% coverage enforced; public GA |

**Total: ~50 weeks (about 12 months) with a team of ~15 engineers.**

---

## 2. Milestone 1 — Foundation (Weeks 1–4)

### Phase 1.1 — Repository & tooling
1. Create `helixgitpx` monorepo with module layout from [02-services/03-microservices-catalog.md].
2. Create `helixgitpx-web`, `helixgitpx-clients`, `helixgitpx-platform`, `helixgitpx-docs` repos.
3. Set up `helix-platform` shared Go libs (logging, telemetry, errors, config, grpc, gin, kafka, pg, redis, temporal, spire, opa, health, test).
4. Set up Nx workspace for web; Gradle convention plugins for KMP.
5. Scaffold `cookiecutter-helix-go` service template with README, Dockerfile, Helm chart, skaffold.yaml.
6. Configure Buf Schema Registry instance and `proto/` root.

### Phase 1.2 — CI skeleton
7. GitHub Actions pipelines per repo: lint, format, test (empty initially), SBOM, Cosign sign, Snyk, SonarQube, Semgrep, Gitleaks, CodeQL.
8. Kyverno policy bundle; Checkov for IaC.
9. Self-hosted runner pool on Kata.
10. OIDC → Vault wiring for short-lived creds.

### Phase 1.3 — Dev environment
11. `devbox`/`mise` bootstrap file committed.
12. `tilt.yaml` or `skaffold.yaml` for local dev cluster.
13. Kind/k3d cluster scripts.
14. Postgres, Kafka, Redis docker-compose for quickest local start.

### Phase 1.4 — Docs baseline
15. Publish the documentation suite (this repo).
16. ADR registry template.
17. Runbook template.
18. API docs site scaffold (Docusaurus).

**Exit criteria**: new engineer onboards, runs `make dev`, sees a hello-world service responding via gRPC and REST, and their PR runs through green CI.

---

## 3. Milestone 2 — Core Data Plane (Weeks 5–10)

### Phase 2.1 — Kubernetes platform
19. Provision staging cluster (K8s 1.31+, Cilium, Istio Ambient).
20. Argo CD installed; `helixgitpx-platform` repo reconciling.
21. cert-manager + external-dns.
22. SPIRE server + agents deployed; sample workload fetches SVID.

### Phase 2.2 — Postgres
23. Deploy Postgres 16 via CNPG (CloudNativePG) with HA.
24. Per-service schemas; RLS policies.
25. Goose migrations wired.
26. PITR to object store; monthly restore drill.

### Phase 2.3 — Kafka + Karapace + Debezium
27. Strimzi Kafka cluster with KRaft.
28. Karapace schema registry.
29. Debezium Connect cluster.
30. Outbox pattern sample using `dummy-service`.

### Phase 2.4 — Redis / Dragonfly + Meilisearch + OpenSearch + Qdrant
31. Dragonfly cluster.
32. Meilisearch cluster.
33. OpenSearch cluster with ILM + snapshots.
34. Qdrant cluster.

### Phase 2.5 — Vault + Observability
35. Vault deployed with Raft; Shamir keys.
36. Prometheus + Mimir + Loki + Tempo + Pyroscope.
37. Grafana provisioning from Git.
38. Alertmanager wired to PagerDuty sandbox.

**Exit criteria**: a smoke test service persists to Postgres, produces to Kafka, reads Redis, exposes metrics/logs/traces, and all are visible in Grafana.

---

## 4. Milestone 3 — Identity & Orgs (Weeks 11–14)

### Phase 3.1 — auth-service
39. OIDC flow with Keycloak dev IdP.
40. JWT issuance (RS256, 15 min access; rotating refresh).
41. PAT endpoints with prefix `hpxat_`.
42. MFA enrolment (TOTP + FIDO2).
43. Sessions table + revocation.

### Phase 3.2 — org-service + team-service
44. CRUD for orgs.
45. Teams with nested support.
46. Memberships with roles.
47. OPA policy bundle v1 (baseline RBAC).

### Phase 3.3 — audit-service
48. `audit.events` Kafka consumer.
49. Postgres append-only table + triggers.
50. Merkle anchoring job.

### Phase 3.4 — Minimal web shell (login only)
51. Angular app scaffold; auth flow; org list screen.
52. Connect-Go clients generated.
53. OpenTelemetry-web wired.

**Exit criteria**: a user logs in via OIDC, creates an org, adds a teammate, and sees audit entries in Grafana.

---

## 5. Milestone 4 — Git Ingress & Adapter Pool (Weeks 15–20)

### Phase 4.1 — repo-service
54. CRUD for repos; event sourcing for Repo aggregate.
55. Refs + branch protection.
56. Presigned upload/download for LFS.

### Phase 4.2 — git-ingress
57. `git-upload-pack` / `git-receive-pack` proxy.
58. Quota + rate-limits.
59. Signed push verification hook.

### Phase 4.3 — adapter-pool (first three providers)
60. Adapter interface + shared plumbing.
61. GitHub adapter (REST+GraphQL).
62. GitLab adapter.
63. Gitea adapter (also powers Codeberg / Forgejo).
64. Adapter contract tests.

### Phase 4.4 — webhook-gateway
65. Incoming webhook receivers per provider with HMAC verify + dedup.
66. Canonicalisation → Kafka.

### Phase 4.5 — Upstream CRUD & binding
67. upstream-service API.
68. Credentials in Vault.
69. Repo↔upstream bindings.

**Exit criteria**: a push to HelixGitpx replicates to GitHub and GitLab; a push to GitHub triggers a webhook that lands in our Kafka.

---

## 6. Milestone 5 — Federation & Conflict Engine (Weeks 21–28)

### Phase 5.1 — sync-orchestrator (Temporal)
70. `FanOutPush` workflow.
71. `InboundReconcile` workflow.
72. DLQ + retries + replays.

### Phase 5.2 — conflict-resolver (core classes)
73. Ref divergence detector (Kafka Streams / Goka).
74. Policy engine hookup (OPA Rego per repo).
75. Three-way merge executor in sandbox.

### Phase 5.3 — CRDT metadata
76. Automerge-go integration.
77. Labels / milestones / assignees / issue body.
78. Per-upstream replay of CRDT ops.

### Phase 5.4 — Remaining providers
79. Gitee.
80. GitFlic.
81. GitVerse.
82. Bitbucket Cloud + DC.
83. Forgejo (if not covered by Gitea adapter param).
84. SourceHut.
85. Azure DevOps.
86. AWS CodeCommit.
87. Generic Git.
88. WASM plugin host + SDK + example plugin.

### Phase 5.5 — live-events-service
89. gRPC streaming endpoint.
90. WebSocket fallback (Connect).
91. SSE fallback.
92. Resume-token persistence.

**Exit criteria**: a repo bound to 5+ upstreams survives a curated conflict scenario (ref divergence, label race, rename collision) with correct, audited resolution and live events visible in the web app.

---

## 7. Milestone 6 — Frontend & Mobile (Weeks 29–36)

### Phase 6.1 — Web full feature set
93. Dashboard + repo list + repo details + code browser.
94. PR flows: list, detail, diff, review.
95. Issue flows.
96. Conflicts inbox + resolver UI.
97. Upstream config UI.
98. Settings, members, org admin.
99. Search UI (hybrid).
100. i18n for 8 locales.
101. a11y pass + Lighthouse thresholds.
102. PWA + offline shell.

### Phase 6.2 — Shared KMP library
103. Core / network / data / domain / store layers.
104. SQLDelight schemas.
105. Connect + gRPC clients per platform.
106. Offline outbox + replay.

### Phase 6.3 — Compose Multiplatform UI
107. Design tokens + theme.
108. Shared screens for repos/PRs/issues/conflicts/settings.
109. Adaptive layout (phone / tablet / desktop).
110. Mobile-specific: push notifications, widgets, biometrics.
111. Desktop-specific: tray, menubar, multi-window, drag-drop.

### Phase 6.4 — Distribution
112. Play Store submissions + F-Droid.
113. App Store + TestFlight.
114. MSIX + DMG + AppImage + .deb/.rpm builds.
115. Auto-update on desktop; self-hosted update feed.

**Exit criteria**: feature parity between web and mobile/desktop for core flows; enterprise-grade polish; accessibility audits passed.

---

## 8. Milestone 7 — AI, Search, Policy (Weeks 37–42)

### Phase 7.1 — AI platform
116. ai-service with LiteLLM router.
117. Ollama + vLLM inference stacks.
118. Guardrails (NeMo).
119. Structured output enforcement.
120. Prompt library + tests.

### Phase 7.2 — Use cases
121. Conflict resolution proposals with sandbox validation.
122. PR summary + review.
123. Label suggestion.
124. Semantic / code search.
125. Chatops.

### Phase 7.3 — Self-learning
126. Feedback collection (accept/reject/edit).
127. Curator with PII scrubbing.
128. DPO training pipeline (Ray on GPU pool).
129. LoRA adapters per task.
130. Shadow-mode evaluation + promotion.

### Phase 7.4 — Policy-as-code end-to-end
131. OPA bundle server.
132. Rego rules for all major enforcement points.
133. Policy diff-review gate in CI.
134. Kyverno admission policies in cluster.

### Phase 7.5 — Hybrid search
135. Projectors to Meilisearch + Qdrant + OpenSearch live.
136. Reciprocal rank fusion.
137. Code search via Zoekt.

**Exit criteria**: AI acceptance rate ≥ 0.7 on internal eval; search p99 ≤ 150 ms; policy denies unintended actions with friendly error.

---

## 9. Milestone 8 — Scale, Harden, GA (Weeks 43–50)

### Phase 8.1 — Multi-region
138. Second region deployed.
139. MirrorMaker 2 + Postgres logical replication.
140. GeoDNS + failover tests.
141. Data residency toggle per org.

### Phase 8.2 — Performance
142. Load tests (k6) at targets (§ capacity model).
143. Soak for 7 d.
144. Perf budgets enforced in CI.
145. GPU inference autoscaling at scale.

### Phase 8.3 — Chaos & DR
146. Chaos experiments: broker kill, pod evict, partition, disk full, LLM out, upstream 429 storms.
147. Full DR drill: simulate region loss; failover; data integrity checks.
148. Runbooks validated.

### Phase 8.4 — Security hardening
149. Full pen-test (external).
150. Bug bounty launch.
151. SOC 2 Type I evidence collection.
152. ISO 27001 gap analysis.

### Phase 8.5 — Coverage 100 %
153. Every module audited for test coverage.
154. Mutation testing thresholds enforced.
155. Fuzz corpora expanded.
156. E2E gaps closed.

### Phase 8.6 — GA
157. Public documentation site live.
158. Trust center page.
159. Billing + plan management integrated.
160. Marketing + partners primed.
161. Public launch. 🎉

**Exit criteria**: GA-ready; zero P1/P2 open; SLOs green; customers onboarded during beta stable for 30 days.

---

## 10. Post-GA (Year 2)

- SOC 2 Type II.
- ISO 27001 certification.
- SLSA Level 4.
- Additional Git providers (ActivityPub / Forgefed, Phabricator successors).
- On-device AI for desktop.
- Confidential compute option.
- Regional expansion (APAC).
- FedRAMP Moderate consideration.

---

## 11. Cross-Cutting Tracks (Always-On)

- Security posture (continuous).
- Performance regression watch.
- Cost review.
- Documentation refresh (per change).
- Accessibility audit (quarterly).
- Dependency updates (Renovate / Dependabot).

---

## 12. Risk Register

| Risk | Severity | Mitigation |
|---|---|---|
| Upstream API breakage | High | Shadow mode + contract tests; adapter version pinning |
| LLM regression | Medium | Canary + shadow eval + rollback |
| Kafka data loss | High | RF=3, min-ISR=2, idempotent producers, backups |
| Postgres corruption | High | PITR + logical backups; monthly drills |
| Regulatory change | Medium | Compliance watch list; quick pivots via feature flags |
| Engineer burnout | Medium | Scope discipline; ruthless backlog grooming; on-call rotations |

---

## 13. Definition of Done (per task)

- [ ] Code merged to `main` via PR with ≥ 2 approvals.
- [ ] Tests added / updated; coverage + mutation gates green.
- [ ] Docs updated (API, runbook, relevant doc in this suite).
- [ ] Dashboards + alerts updated if behaviour changes.
- [ ] ADR written if decision is architectural.
- [ ] Deployed to staging for ≥ 24 h without regression.
- [ ] Feature flag plan documented.

---

*— End of Milestones & Roadmap —*
