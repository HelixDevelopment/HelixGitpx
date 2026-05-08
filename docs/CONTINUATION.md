# CONTINUATION.md — Session Continuation Document

> **MANDATORY.** This document is the single source of truth for resuming work
> across sessions, CLI agents, and LLM models. It MUST be updated before any
> work session ends. It MUST NOT be out of sync with current work.
>
> **Governance:** Required by CONST-034 in `CONSTITUTION.md`. Also enforced in
> `CLAUDE.md` and `AGENTS.md`. Any agent continuing work MUST read this file
> first and update it before stopping.
>
> **Last updated:** 2026-05-08 (session start — first creation).

---

## Quick Start for Any Agent

1. Read this file (`docs/CONTINUATION.md`) — you are here.
2. Read `docs/UNFINISHED.md` for the detailed gap inventory.
3. Read `CONSTITUTION.md`, `CLAUDE.md`, `AGENTS.md` for governance rules.
4. Pick the next task from **Current Priority Queue** below.
5. Before stopping: update this file's sections below to reflect what you did,
   what changed, and what remains.

---

## Current Session State

### Session 2026-05-08 (this session)

| Item | Status |
|------|--------|
| AGENTS.md rewrite | Done — rewritten with expanded §1–§13 + Universal Mandatory Constraints |
| CONTINUATION.md creation | Done — this file |
| CONST-034 added to CONSTITUTION.md | Done |
| Continuation constraint added to CLAUDE.md | Done |
| Continuation constraint added to AGENTS.md | Done |
| Git commit | Pending |
| Push to all upstreams | Pending |

### Previous sessions (summary)

| Date | Key work |
|------|----------|
| 2026-04-19 | Initial repo seeded with spec v4.0.0 |
| 2026-04-20 | M1–M8 milestones tagged, v1.0.0 GA cut, Constitution + governance docs created |
| 2026-04-21 | Major execution pass: 6/17 services wired, all 13 GitHub Actions re-enabled and passing, business claims ground-truthed, spec archives refreshed, Go benchmarks added, E2E/chaos bodies created, manual chapters expanded |
| 2026-04-21 (cont.) | Plinius integration policy-classified (7 KEEP, 6 KEEP-GATED, 7 DROP), UNFINISHED.md updated with session deltas |
| 2026-04-26 | CONST-033 host-power-management guard added after auto-suspend incident |
| Post-04-26 | Universal Mandatory Constraints cascaded to all governance docs |

---

## Overall Project Status

| Metric | Value |
|--------|-------|
| Version | v1.0.0 GA tagged |
| Milestones tagged | M1–M8 (all tagged, all have plan files, **0 tasks checked**) |
| Services wired end-to-end | 6 / 17 |
| Services scaffolded (17-line stubs) | 11 |
| Go packages with any tests | 39 / 98 |
| Go packages at 100% coverage | 8 / 98 |
| Go packages at 0% coverage | 59 / 98 |
| Integration tests | 5 (all need env vars + compose stack) |
| E2E suites wired to cluster | 0 |
| Constitution-mandated test types with real tests | 3 / 7 (unit, integration runner, security runner shell) |
| GitHub Actions workflows enabled | 13 / 13 (all passing on main) |
| GitLab pipeline | Suppressed (identity verification pending) |
| Helm charts (artifact lint) | 53 / 53 green |
| Argo CD apps (path validated) | 53 / 53 green |
| Rego files (syntax only) | 3 / 3 green |
| Runtime verification (cluster/compose) | 0 — never done |
| Manual chapters written (beyond intro) | ~12 across 6 manuals |
| Video recordings produced | 0 / 21 scripts |
| Marketing-site pages scaffolded | 19 / 19 |
| Marketing-site deployed | No (no DNS) |
| Real customers | 0 |
| Active certifications (SOC 2, ISO 27001) | None |
| Upstream remotes receiving pushes | 4 / 4 (GitHub, GitLab, GitFlic, GitVerse) |
| External integrations ready to execute | 0 / 1 (Plinius blocked on W0 spike) |

---

## Current Priority Queue

Ordered by impact and dependencies. Work top-to-bottom.

### Priority 1 — Wire Remaining Services (UNFINISHED §1)

**11 services remain as 17-line scaffolds.** Each has domain packages with real
logic and tests but no HTTP/gRPC boundary.

| Service | Key dependency | Effort estimate |
|---------|---------------|-----------------|
| adapter-pool | Provider registry + health RPC + Vault token rotation | 1-2 days |
| ai-service | LiteLLM client + NeMo Guardrails proxy + Kafka feedback | 1-2 days |
| billing-service | Stripe webhook receiver + Postgres repo + outbox publisher | 1-2 days |
| collab-service | Automerge-go doc store + gRPC stream fan-out | 1-2 days |
| conflict-resolver | Temporal worker + ref-divergence detector + AI bridge | 1-2 days |
| git-ingress | go-git smart-HTTP server + per-org quota client | 1-2 days |
| live-events-service | Kafka consumer + gRPC/WS/SSE fan-out + resume-token store | 1-2 days |
| orgteam | Has `residency` handler but `app.Run` does not route to it | 0.5 days |
| sync-orchestrator | Temporal worker + FanoutPush / InboundReconcile workflows | 1-2 days |
| upstream | Binding persistence + OpenAPI/REST surface + adapter-pool dispatcher | 1-2 days |
| webhook-gateway | Signed-body verification HTTP router + outbox producer | 1-2 days |

**Already wired:** auth, hello, audit, opa-bundle-server, search-service, repo-service.

### Priority 2 — Test Coverage (UNFINISHED §2)

Current: 59/98 Go packages at 0%. Target: 100% across all 7 types.

Pragmatic order:
1. Wire scaffolded services (Priority 1) — adds integration + e2e targets
2. Stand up compose-backed integration CI so `make test-integration` spins up stack
3. Backfill unit tests on Postgres adapters (mocks OK in unit tests)
4. Add one ddos/stress scenario per service for baseline capture
5. Backfill benchmark `Benchmark*` functions (6 exist, need per-service coverage)
6. Fix web unit tests — switch from broken Jest+Angular to Karma or proper ESM preset

### Priority 3 — Runtime Verification (UNFINISHED §3)

All 53 Helm charts, 53 Argo apps, and 3 Rego files pass artifact lint but have
**never been deployed or runtime-tested.**

- `helm template` never run (no helm CLI during GA push)
- No `argocd app sync` ever done (no cluster)
- `opa test` / `opa eval` never run (no opa CLI)
- Kafka topics/schemas never created on live Kafka
- Postgres migrations never applied via CI
- Container images never built end-to-end for all services
- Compose stack never smoke-tested in CI

**Action:** Set up a kind/k3d cluster in CI, run `helm template`, deploy compose
stack, smoke test all endpoints.

### Priority 4 — Documentation Completion (UNFINISHED §6)

| Task | Status |
|------|--------|
| User guide chapters 1-5 | Written (only useful manual) |
| Operator guide chapters 3+ | Not written |
| Administrator guide chapter 2+ | Not written |
| Developer guide chapter 2+ | Not written |
| API reference real content | Not written (lives in proto files, not prose) |
| CLI reference real content | Not written |
| Security handbook content | Intro only |
| Deployment cookbook recipes | Not written |
| Troubleshooting real scenarios | Not written |
| Migration guide content | Not written |
| Docs site (helixgitpx.dev) | Builds clean, stub pages only |
| Video recordings | 0 / 21 scripts recorded |
| Spec archives (.zip/.7z) | May be stale — regenerate on release |

### Priority 5 — Business / Compliance (UNFINISHED §9)

| Item | Status |
|------|--------|
| Customers | 0 — "Private beta in progress" is aspirational |
| Billing integration | Schema/proto/domain/chart done; no Stripe account linked |
| SOC 2 Type I | Evidence index written; no evidence collected, no auditor |
| ISO 27001 | Gap analysis written; no formal audit |
| Pen-test | Scope doc for Q2 2026; no vendor engaged |
| Bug bounty | Program docs exist; not active on HackerOne |
| Legal entity | HelixDevelopment UG reference removed (unverified); terms/privacy/DPA templates pending legal review |
| DPO / legal contact | Not set up |

### Priority 6 — External Integrations (UNFINISHED §8)

**Plinius integration (the only one planned):**
- Go port layer does NOT exist — ~22 repos need creating
- Policy classified: 7 KEEP, 6 KEEP-GATED, 7 DROP
- W0 spike required before Phase 1
- 3 modules have no upstream at all (go-tempest, go-gandalf-solutions, go-gitgpt)
- Blocked until W0 spike lands

### Priority 7 — Client Platforms (UNFINISHED §7)

| Surface | Gap |
|---------|-----|
| KMP shared | 476 files, no Connect-RPC client wired; screens show seeded data |
| Android | Min viable shell; no signing, Play Store metadata, push, biometrics |
| iOS | Skeleton; no SwiftUI screens calling shared code verified |
| Desktop | Compose window opens; no tray, auto-update, multi-window, drag-drop |
| Web E2E | Blocked on cluster (Priority 3) |

---

## Known Issues and Bugs

### Engineering

1. **Web unit tests are smoke-only.** `app.component.spec.ts` does `expect(true).toBe(true)`. Real Angular TestBed path not wired — Jest + Angular ESM modules conflict. Need Karma runner or proper Jest ESM preset.
2. **Integration tests require manual setup.** `make test-integration` does NOT spin up compose stack. Operator must run `make dev` + export env vars manually.
3. **Spec archives may be stale.** Last regeneration was 2026-04-21 session. `.zip`/`.7z` may not reflect current spec tree state.
4. **Cluster-dependent verifiers print failures without cluster.** Scripts `verify-m{3,4,5,6,7,8}-cluster.sh` will print `[FAIL]` rows without a reachable cluster. This is intentional (strict mode) but confusing.
5. **`.github/workflows/README.md` discrepancy.** The README lists workflows as `.yml.disabled` but the actual files are `.yml` (they were re-enabled in the 2026-04-21 session). The README inventory table is stale.

### Business / Legal

6. **No legal entity confirmed.** HelixDevelopment UG was removed from the codebase but no replacement entity is established.
7. **Terms/Privacy/DPA are templates.** Not legally reviewed. Website pages say this honestly.
8. **No revenue, no customers.** All marketing claims about customers are aspirational.

### Security

9. **OPA bundle never loaded by running cluster.** Policy-as-code is committed but never runtime-verified.
10. **Supply-chain workflows not run.** SBOM, Cosign signatures, SLSA provenance — all depend on CI being fully operational with container builds.
11. **Bug bounty program docs exist but not active.** Not listed on HackerOne or any platform.

---

## Milestone Phases (M1–M8) — Status

All 8 milestones have been **tagged** in git. All have plan files at
`docs/superpowers/plans/2026-04-20-m*.md`. **0 of ~151+ tasks are checked.**

The milestones represent a plan for future implementation. The GA tag was cut
on the basis of contracts + scaffolding + policy, not runtime completeness.

| Milestone | Name | Plan file | Tasks | Checked | Core focus |
|-----------|------|-----------|-------|---------|------------|
| M1 | Foundation | `2026-04-20-m1-foundation.md` | ~36 tasks | 0 | Go monorepo, platform libs, proto, health, errors |
| M2 | Core Data Plane | `2026-04-20-m2-core-data-plane.md` | ~26 tasks | 0 | Postgres, Redis, Kafka, migrations, repo service |
| M3 | Identity & Orgs | `2026-04-20-m3-identity-orgs.md` | ~20 tasks | 0 | Keycloak, org/team, RBAC, SPIFFE/SPIRE |
| M4 | Git Ingress & Adapter Pool | `2026-04-20-m4-git-ingress-adapter-pool.md` | ~20 tasks | 0 | Smart HTTP, adapter-pool, provider registry |
| M5 | Federation & Conflict Engine | `2026-04-20-m5-federation-conflict-engine.md` | 14 tasks | 0 | Bidirectional sync, CRDT, Temporal workflows |
| M6 | Frontend & Mobile | `2026-04-20-m6-frontend-mobile.md` | 9 tasks | 0 | Angular app, KMP clients, desktop/mobile |
| M7 | AI, Search & Policy | `2026-04-20-m7-ai-search-policy.md` | 10 tasks | 0 | LiteLLM, Meilisearch, OPA, AI features |
| M8 | Scale, Harden, GA | `2026-04-20-m8-ga.md` | 20 tasks | 0 | Performance, chaos, DR, observability, launch |

### Relationship between milestones and current code

The milestones were designed as a sequential implementation plan. Current code
represents partial completion of M1–M8 (scaffolding and contracts done, runtime
wiring incomplete). The milestone plan files should be used as implementation
checklists when wiring each service and subsystem.

---

## E2E Flow Status

From `tools/e2e-gaps.md`:

| Flow | Status | Ticket |
|------|--------|--------|
| Sign-up → first repo → push | Covered | HGX-101 |
| Create org → invite member → accept | Covered | HGX-102 |
| Bind upstream → initial mirror → PR | Covered | HGX-103 |
| Conflict detected → AI proposal → human accept | Partial (no AI path) | HGX-310 |
| Change residency → data migrates | Missing | HGX-311 |
| DR failover → customer impact | Runbook-tested only | HGX-312 |
| Desktop app auto-update | Happy path only | HGX-313 |
| Mobile push notification → deep link | Missing | HGX-314 |
| Billing: plan upgrade, downgrade, cancel | Missing | HGX-315 |
| Trust center page: load, links resolve | Missing | HGX-316 |
| Full chaos recovery matrix | Manual Litmus runs | HGX-317 |

---

## Files That Track Work

| File | Purpose |
|------|---------|
| `docs/CONTINUATION.md` | **This file.** Session continuation state. Updated every session. |
| `docs/UNFINISHED.md` | Detailed gap inventory (488 lines). Updated per milestone. |
| `docs/superpowers/plans/2026-04-20-m*.md` | 8 milestone plan files with task checklists. |
| `docs/superpowers/specs/2026-04-20-m*.md` | 8 milestone design specifications. |
| `tools/e2e-gaps.md` | Per-flow E2E audit. |
| `docs/marketing/launch-checklist.md` | GA launch operations checklist (all items unchecked). |
| `CHANGELOG.md` | Version history (v1.0.0 GA + v0.0.0 initial seed). |
| `RELEASE.md` | Release notes template (GA-DATE placeholder). |
| `SOLO-NOTES.md` | Solo-maintainer deviations from CONTRIBUTING.md. |

---

## Upstream Federation

| Remote | Script | Status |
|--------|--------|--------|
| GitHub | `bash Upstreams/GitHub.sh` | Active |
| GitLab | `bash Upstreams/GitLab.sh` | Active |
| GitFlic | `bash Upstreams/GitFlic.sh` | Active |
| GitVerse | `bash Upstreams/GitVerse.sh` | Active |

All scripts export `UPSTREAMABLE_REPOSITORY` and push `main` + tags.
After every commit to `main`, run: `make upstream-push` or each script individually.

---

## Agent Handoff Checklist

When stopping work, the agent MUST:

- [ ] Update the **Current Session State** section above
- [ ] Update the **Overall Project Status** metrics if anything changed
- [ ] Update the **Current Priority Queue** if priorities shifted
- [ ] Add any new **Known Issues and Bugs** discovered during the session
- [ ] Update the **Last updated** date at the top of this file
- [ ] Commit `docs/CONTINUATION.md` along with any other changes
- [ ] Push to all upstreams (`make upstream-push`)

When starting work, the agent MUST:

- [ ] Read this entire file
- [ ] Read `docs/UNFINISHED.md` for the detailed gap inventory
- [ ] Check `git log --oneline -20` for recent changes
- [ ] Pick the next task from the Current Priority Queue
- [ ] Begin work

---

*This document is maintained per CONST-034 of the Constitution. It is the
authoritative handoff point between sessions, agents, and LLM models.*
