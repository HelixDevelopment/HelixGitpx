# UNFINISHED.md — what is not done, and why

> **Purpose.** One honest inventory of everything that is **not** finished
> in this repository, plus the reason each gap exists. Read this before
> planning the next sprint or answering the question "are we GA-ready?".
>
> **Snapshot date:** 2026-04-21. Verified against git HEAD `ce382ce`.
> **Owner:** refresh on every milestone tag.

Linked companions:

- [`tools/e2e-gaps.md`](../tools/e2e-gaps.md) — per-flow E2E audit.
- [`docs/integrations/helixagent-plinius-verification.md`](integrations/helixagent-plinius-verification.md)
  — gaps specific to that integration plan.
- [`docs/integrations/helixagent-plinius-w0-spike.md`](integrations/helixagent-plinius-w0-spike.md)
  — Week-0 work blocking the plinius integration.
- [`docs/superpowers/plans/2026-04-20-m*.md`](superpowers/plans/) — milestone
  plans (describe what was scoped IN).
- [`docs/marketing/launch-checklist.md`](marketing/launch-checklist.md) —
  post-GA launch operations.

---

## Legend

- **Deferred** — scoped out on purpose, will be done later.
- **Scaffolded** — file / module / type exists, but isn't wired end-to-end.
- **Stub** — placeholder content; needs real prose, real code, real tests.
- **Blocked** — waiting on an external decision, asset, or credential.
- **Broken path** — works in isolation but won't run end-to-end without
  unstated prerequisites.

---

## 1. Services not fully wired

Implementation code lives under `impl/helixgitpx/services/`. Five services
have real `app.Run` wiring (HTTP handlers, ServeMux, graceful shutdown).
**Twelve services are 16–17-line scaffolds** whose `app.Run` only blocks
on `ctx.Done()`. The scaffolded services have domain packages with real
logic and TDD tests; they do not yet expose that logic over the network.

| Service | app.go lines | Status | Why not done |
|---------|-------------|--------|--------------|
| auth | 160 | **wired** | — |
| hello | 129 | **wired** | — |
| audit | 96 | **wired** | — |
| opa-bundle-server | 79 | **wired** | — |
| search-service | 61 | **wired** | — |
| adapter-pool | 17 | **scaffolded** | Needs provider registry wiring + health RPC + token rotation against Vault. |
| ai-service | 17 | **scaffolded** | Needs LiteLLM client + NeMo Guardrails proxy + feedback ingest → Kafka. |
| billing-service | 16 | **scaffolded** | Needs Stripe webhook receiver + Postgres repo + outbox publisher. |
| collab-service | 17 | **scaffolded** | Needs Automerge-go doc store + gRPC stream fan-out. |
| conflict-resolver | 17 | **scaffolded** | Needs Temporal worker + ref-divergence detector + AI bridge. |
| git-ingress | 17 | **scaffolded** | Needs go-git-backed smart-HTTP server + per-org quota client. |
| live-events-service | 17 | **scaffolded** | Needs Kafka consumer + gRPC/WS/SSE fan-out + resume-token store. |
| orgteam | 17 | **scaffolded** | Has `residency` handler, but `app.Run` does not route to it. |
| repo | 17 | **scaffolded** | Needs Postgres repo + Connect-RPC handlers for RepoService. |
| sync-orchestrator | 17 | **scaffolded** | Needs Temporal worker + FanoutPush / InboundReconcile workflow impls. |
| upstream | 17 | **scaffolded** | Needs binding persistence + OpenAPI/REST surface + adapter-pool dispatcher. |
| webhook-gateway | 17 | **scaffolded** | Needs signed-body verification HTTP router + outbox producer. |

**Why this exists.** The v1.0.0 GA tag emphasizes *contracts + scaffolding
+ policy* over end-to-end runtime. The Constitution defines what must be
shipped; the milestones (`m1-foundation` → `m8-ga`) tagged the contracts
and infrastructure. Full service wiring is a post-GA track, currently
un-scheduled. Each service carries domain tests so the business-rule
layer is testable today; the HTTP/gRPC boundary is the gap.

**Impact.** A cluster brought up with today's artifacts will serve the
five wired services plus the scaffold services' `/healthz`. Nothing else.

---

## 2. Test coverage vs. the Constitution

Constitution Article II mandates **seven test types at 100% coverage per
module touched**: unit, integration, e2e, security, stress, ddos, benchmark.

### 2.1 Current state

| Type | Test files present | Gap |
|------|-------------------|-----|
| unit | 26 `*_test.go` across services + platform; 1 `.spec.ts` web | **59/98 Go packages at 0%.** Only 8 Go packages at 100%. |
| integration | 5 files (auth, hello, orgteam, repo, webhook) | All require env vars + live compose stack. None run in CI. |
| e2e | 0 files (Playwright `01-login-and-orgs.spec.ts` exists at `impl/helixgitpx-web/e2e/`, not under `test/e2e/`) | Mobile + desktop have no e2e suites. Requires live cluster. |
| security | `test/security/run.sh` orchestrator | No first-party security tests yet; script shells out to gosec/ZAP/Nuclei/Trivy, all graceful-skip. |
| stress | `test/stress/run.sh` | No calibrated `_stress.js` scenarios exist. |
| ddos | `test/ddos/run.sh` with 3 inline scenarios | Inline only; no repeatable baseline capture. |
| benchmark | `test/benchmark/run.sh` + micro-bench runner | No Go `Benchmark*` funcs exist yet. |
| chaos | 0 files (6 Litmus YAMLs at `tools/chaos/`) | Litmus experiments defined; no automation that runs them + asserts. |

### 2.2 Go-package coverage distribution

Of 98 testable Go packages:

```
0%:      59  (60%)  no tests at all
0-50%:    8
50-80%:   9
80-100%: 14
100%:     8
```

Most at 0% are `cmd/` main packages, `internal/app/` composition roots,
and `internal/repo/` Postgres adapters (which need integration tests,
not unit — per §II §2 mocks are forbidden outside unit tests).

### 2.3 Why the gap

The Constitution was ratified mid-session, after most scaffolding was
already committed. Retrofitting 100% coverage across seven types for
every service is weeks of work per service. A pragmatic order:

1. Wire scaffolded services (§1) — adds integration + e2e targets.
2. Write real integration tests that boot a compose stack in CI.
3. Backfill unit tests on the Postgres adapters (not forbidden — mocks
   OK in unit tests).
4. Promote one of ddos/stress scenarios per service into a recorded
   baseline so CI can detect regression.

### 2.4 Known broken paths

- **Web unit tests** pass as a *pure TS* smoke test (`app.component.spec.ts`
  does `expect(true).toBe(true)`). The real Angular `TestBed` path is
  not wired — jest + jest-preset-angular + Angular ESM modules fight
  each other in this repo. A switch to Karma (`@angular/build:karma`)
  or proper Jest ESM preset is the fix; neither is done.
- **Integration tests** (`test/integration/*.go`) call `t.Fatal` when
  their env vars are unset. They require a compose stack plus a seeded
  Keycloak token. The Makefile target `make test-integration` does
  **not** spin up the stack first — operators must run `make dev` and
  export env vars manually.

---

## 3. Platform artifacts — never runtime-verified

Artifact verifiers are green (53/53 Helm charts, 53/53 Argo paths, 3/3
Rego files, M1/M2 artifacts). **Runtime verification is missing.**

| Artifact class | Artifact lint | Runtime probe | Why gap |
|----------------|---------------|---------------|---------|
| Helm charts (53) | Python-based lint (Chart.yaml + values.yaml + templates) | `helm template` never run | `helm` CLI absent on the dev host used during GA push; CI workflow exists but is disabled (§4). |
| Argo Applications (53) | `verify-argo-paths.sh` confirms `source.path` exists | No `argocd app sync` verification | No Argo CD cluster reachable from this workstation. |
| OPA bundles | `verify-rego.sh` (syntax only) | `opa test`, `opa eval` never run | `opa` CLI absent. Real Rego test runner is a dedicated CI job. |
| Kafka topics / schemas | Referenced in charts + proto | Never created on a live Kafka | No Kafka cluster brought up during GA push. |
| Postgres migrations | SQL files committed | `goose up` never run | CNPG cluster brought up only in operator's local dev. |
| Container images | Dockerfile per service | Never built end-to-end in CI | docker / podman workflow disabled; manual `podman build` never invoked for all 17 services. |
| Compose stack | `impl/helixgitpx-platform/compose/compose.yml` with `all` profile | Never brought up in this repo's CI | `make dev` is an operator-local action; no automated smoke. |

**Why.** GA was declared on the basis of artifact-presence and domain
correctness. Runtime verification requires bringing up a real cluster
or compose stack, which is an operator action, not an author action.
Every verifier script includes an escape hatch that skips cleanly
when tooling is absent — that's deliberate, so CI stays green on dev
laptops. It also hides runtime bugs.

---

## 4. CI / automation — disabled

As of commit `ce382ce` (2026-04-21), **all 13 GitHub Actions workflows
and the 4 reusable callables are renamed to `.yml.disabled`**. GitLab
pipelines are explicitly blocked via `workflow: rules: - when: never`.

| Workflow | Purpose | Disabled because |
|----------|---------|------------------|
| ci-clients | Gradle KMP build | operator decision 2026-04-21 |
| ci-docs | Docusaurus build | operator decision |
| ci-go | `go vet` + `go test` | operator decision |
| ci-platform | Chart + OPA checks | operator decision |
| ci-verifiers | `verify-everything.sh` | operator decision |
| ci-web | Nx build + jest + Playwright | operator decision |
| deploy | GitOps image promotion | operator decision |
| mutation-testing | `go-mutesting` | operator decision |
| perf-budgets | k6 + budget gate | operator decision |
| release | Tagged release workflow | operator decision |
| security-scan | SAST + DAST + IaC | operator decision |
| supply-chain | SBOM + Cosign + SLSA + Trivy | operator decision |
| upstream-sync | Push to all federated hosts | operator decision |
| _reusable/* (4) | Shared callables | operator decision |

**Why.** The operator asked for an explicit disable. Local enforcement
via `bash scripts/verify-everything.sh` covers the green-suite path. Any
merge-time enforcement (mutation thresholds, perf budgets, SBOM
attestation) is temporarily on hold.

**Re-enable recipe** lives in [`.github/workflows/README.md`](../.github/workflows/README.md).

---

## 5. Cluster-dependent verifiers

Six `scripts/verify-m*-cluster.sh` + `-spine.sh` scripts probe a live
cluster. Four of them (M2 cluster, M2 spine, M3 spine, M4 spine) now
short-circuit cleanly when no `kubectl` context is available. **Three
still run raw** (M3/M4/M5/M6/M7/M8 cluster) and will print 0-or-N
`[FAIL]` rows without a reachable cluster.

Why: these three scripts were kept strict intentionally — when an
operator runs them in an environment that *should* have a cluster,
silent skip would hide real regressions. The `-cluster.sh` scripts are
operator-run, not CI-run.

---

## 6. Documentation — stubs and placeholders

### 6.1 Manuals

All ten manuals have an introduction (`00-introduction.md`). Two have
further chapters:

| Manual | Chapters | Status |
|--------|----------|--------|
| user-guide | 5 (00-intro, 02-first-repo, 03-pushing-prs, 04-conflicts, 05-ai) | **Only useful manual so far.** |
| operator-guide | 2 (00-intro, 02-cluster-prereqs) | Chapter 3+ not written. |
| developer-guide | 1 | Intro only. |
| administrator-guide | 1 | Intro only. |
| api-reference | 1 | Intro only; real reference lives in proto files, not prose. |
| cli-reference | 1 | Intro only. |
| security-handbook | 1 | Intro only. |
| deployment-cookbook | 1 | Intro only; actual recipes not written. |
| troubleshooting | 1 | Intro only. No real "what to do when X" prose. |
| migration-guide | 1 | Intro only. |

**Why.** Writing production-quality documentation across ten manuals
is a dedicated documentation project. Intros establish structure so
chapter authors know where to plug in. No shortcuts exist.

### 6.2 Video curriculum

**21 scripts written, zero recorded.** Scripts follow a template with
shot list, narration beats, and companion-doc links. Production (OBS
scenes, editing in DaVinci Resolve, captioning, publishing to
Vimeo/YouTube/MinIO) is a separate workstream described in
[`docs/media/README.md`](media/README.md) §Production-pipeline. No
recording has taken place; no brand assets for the video-specific
intro/outro are rendered.

**Why.** Recording video is a human-in-the-loop task. The scripts are
the deliverable that unblocks the humans.

### 6.3 Docs site (`docs.helixgitpx.io`)

- **Builds cleanly** via `npx docusaurus build`.
- All sidebar entries have a stub page.
- Chapter content beyond intros is absent.
- Chapter content here must mirror the manuals (§6.1); same dependency
  chain.

### 6.4 Marketing website (`helixgitpx.io`)

- All 19 pages referenced in the nav/footer exist.
- Build verified (`npx astro build`).
- **Not deployed anywhere** — Argo Application exists at sync-wave 10,
  but no live `helixgitpx.io` domain resolves to a cluster yet.
- **No imagery** beyond the inline SVG logo / favicon / OG image. No
  screenshots, no founder photos, no customer logos (there are no
  customers yet — §9).

### 6.5 Spec archives out of sync

`docs/specifications/main/main_implementation_material/HelixGitpx.zip`
and `.7z` are dated **2026-04-20 09:35** — before dozens of session
commits that touched the spec tree. They do not reflect current state.

**Why.** CLAUDE.md asks "keep these in sync" — regenerating is a
straightforward `zip -r` / `7z a` job and should happen as part of
every release. It wasn't done this session.

---

## 7. Clients (KMP, Android, iOS, desktop, web)

| Surface | Status | Gap |
|---------|--------|-----|
| KMP shared | 476 files (mostly generated proto) | Compose `ui.App` exists; **no Connect-RPC client wired**. Every screen shows seeded data. |
| androidApp | 5 files (manifest, theme, strings, MainActivity) | Minimum viable Compose shell. No signing config. No Play Store metadata. No push-notification registration. No biometrics. |
| iosApp | 26 files | Skeleton; **no SwiftUI screens that call shared code** verified. TestFlight pipeline not wired. |
| desktopApp | 2 files (build + Main.kt) | Compose window opens with shared UI. **No tray, no auto-update, no multi-window, no drag-drop** (roadmap items 110-111). |
| web | Angular 19 + Nx | Builds cleanly. Components created during M6 exist. **E2E against a live backend** is blocked on §3 (no cluster). |

**Why.** KMP + Compose Multiplatform delivery across four OSes is a
project in itself. M6 tagged the shells; GA does not promise feature
parity with the web app on mobile/desktop. See §111 in the roadmap.

### 7.1 Web test pipeline

- Jest + ts-jest runs 8 tests, all passing.
- Two of those tests are **smoke tests** that don't exercise Angular.
- The `app.component.spec.ts` smoke test explicitly says "not a real
  test" in a comment.
- Real Angular TestBed integration awaits either a Karma runner or an
  ESM-capable Jest preset.

---

## 8. Integrations (external)

### 8.1 HelixAgent × go-elder-plinius

See the [verification report](integrations/helixagent-plinius-verification.md)
and [W0 spike](integrations/helixagent-plinius-w0-spike.md). Summary of
gaps:

- **Go port layer does not exist.** The plan assumes 20+ Go modules
  under `vasic-digital/go-elder-plinius-*`; none are public today.
- **3 modules have no upstream at all** (`go-tempest`,
  `go-gandalf-solutions`, `go-gitgpt`) — drop, rescope, or build.
- **6 HelixAgent submodules referenced** in the integration table are
  not publicly visible (DebateOrchestrator, Agentic, HelixSpecifier,
  MCP-Servers, BackgroundTasks, BootManager). May be private; access
  not confirmed.
- **Phase 1 blocked** until W0 spike lands.

### 8.2 Other planned integrations

None documented under `docs/integrations/`. Any future integration
gets its own file per the README convention there.

---

## 9. Business-layer gaps (not engineering)

### 9.1 No customers

The `customers.astro` page says "Private beta in progress". There are
no case studies, no logos, no named reference customers. The marketing
claim "our customers say" has nothing to back it. Post-GA work.

### 9.2 No billing integration running

`billing-service` has the schema, the proto, the domain layer, the
Helm chart, and the Argo Application. It has **no Stripe webhook
endpoint wired to a live account**. `STRIPE_API_KEY` secret is
declared in the chart; no Stripe account is linked.

### 9.3 Compliance artifacts are plan-shaped

- **SOC 2 Type I** — `docs/security/soc2-type1-evidence-index.md`
  lists the control mapping. **No evidence has been collected**, no
  auditor engaged, no report issued.
- **ISO 27001 gap analysis** — written, not yet a formal audit.
- **Pen-test** — scope document exists for Q2 2026; no vendor engaged.
- **Bug bounty** — program docs exist; **not active on HackerOne**.

### 9.4 Legal / org

- `docs/manuals/src/security-handbook/` references a legal entity
  "HelixDevelopment UG (haftungsbeschränkt)". No evidence this entity
  exists.
- Terms / Privacy / DPA pages on the marketing site look legitimate
  but are **not legally reviewed**.
- No DPO / legal contact inbox actually monitored.

**Why.** GA-shaped content anticipates operational readiness the
company hasn't achieved yet. These pages are ready to be *made* true,
not representations that they *are* true.

---

## 10. Security features — policy vs runtime

The OPA bundle v2 (`enforcement.rego`) exists. It is **not** loaded by
any running cluster (see §3 + §4). Every other security claim depends
on either:

- Runtime artifacts not yet built (Cosign signatures, SBOMs) — §4
  supply-chain workflow is disabled.
- External vendors not yet engaged (pen-test, bug bounty, SOC 2) — §9.3.
- Mandatory Constitution gates not enforced because CI is disabled —
  §4.

---

## 11. TODOs in code

Historic M5 TODO markers (kafka Karapace client, Temporal dialer,
wazero plugin host) were all resolved in commit `b856cf5`.

As of HEAD `ce382ce`: **zero `TODO` or `FIXME` markers** in
`impl/helixgitpx/`, `impl/helixgitpx-web/apps/`, or
`impl/helixgitpx-web/libs/`.

Absence of TODOs is not the same as absence of unfinished work; most
gaps in this document are *missing code*, not commented placeholders.

---

## 12. Summary counters

| Metric | Value |
|--------|-------|
| Services wired end-to-end | 5 / 17 |
| Services with domain tests | 17 / 17 (all have a `domain/` pkg with tests) |
| Go packages with any tests | 39 / 98 |
| Go packages at 100% coverage | 8 / 98 |
| Go fuzz targets | 4 (webhook HMAC + 3 canonicalizers) |
| Integration tests | 5 (all need env vars) |
| E2E suites wired to a cluster | 0 |
| Manual chapters written beyond intro | 6 (5 in user-guide + 1 in operator-guide) |
| Video recordings produced | 0 / 21 scripts |
| Marketing-site pages | 19 / 19 scaffolded (no imagery) |
| CI workflows enabled | 0 / 17 |
| Cluster-probe verifiers cleanly skippable | 4 / 7 |
| External integration plans executable today | 0 / 1 (plinius needs W0) |
| Upstream targets receiving pushes | 4 / 4 |
| Constitution-mandated test types with any real test | 3 / 7 (unit, integration, security runner shell) |
| Real customers on record | 0 |

---

## 13. Suggested prioritization

When this gets scheduled, a pragmatic order:

1. **Re-enable CI** (§4) — disabling was a session-scoped ops call.
2. **Wire one scaffolded service end-to-end** (§1) — use `repo-service`
   as pilot because its schema, proto, and domain are all in place. Add
   real integration tests as part of the wiring.
3. **Stand up compose-backed integration CI** so `make test-integration`
   actually spins up the stack before running tests.
4. **Write the docs-ch-01** (`user-guide/01-signup-orgs.md`) then spread
   to operator-guide 03+, administrator-guide 02+.
5. **Draft one real case study** — blocker for §9.1 and several
   marketing claims.
6. **W0 plinius spike** so the integration plan either progresses or
   gets formally dropped.
7. **Refresh spec archives** (`zip` + `7z` regeneration) — 10 minutes.
8. **Backfill `Benchmark*` Go funcs** on the platform libraries (audit
   Merkle, webhook HMAC, RRF fusion) — gives the benchmark test type
   something real to measure.
9. **Switch web unit tests to Karma** — unblocks real Angular testing.
10. **Runtime smoke** of all 53 Helm charts via `helm template` + a
    kind cluster in CI.

Everything else in this document can wait until these ten are green.

---

*This file is a living audit. Update on every milestone tag, every
ADR merge, and whenever a claim here stops being true.*
