# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Authority order:** the [Constitution](./CONSTITUTION.md) is the supreme
policy document. Per-agent instructions are in [`AGENTS.md`](./AGENTS.md).
CLAUDE.md (this file) is a Claude-specific orientation guide — it does not
override the Constitution or AGENTS.md.

## Hard rules (from the Constitution)

1. **Testing** — every module carries tests in all seven required types:
   unit, integration, e2e, security, stress, ddos, benchmark. Only unit
   tests may use mocks, stubs, placeholder classes, or hardcoded data.
   No test may be skipped, disabled, broken, or flaky. Target coverage is
   100 % per type per module touched. See Article II of the Constitution.
2. **CI is workflow_dispatch-only.** No push / pull_request / schedule
   triggers are permitted. Mandatory.
3. **Container runtime is portable** — auto-detect docker vs podman; never
   hardcode either one.
4. **Upstream federation is regular** — every change on `main` is pushed to
   all configured upstreams (GitHub, GitLab, GitFlic, GitVerse, …) via
   `Upstreams/<target>.sh` scripts. Daily cadence minimum.
5. **Documentation is source** — every feature ships with HTML + PDF + ePub
   + Markdown + plain-text documentation under `docs/`.

## Repository Status

Implementation code lives under `impl/` (Go monorepo, Angular web app, KMP
clients, Helm/Argo/Kustomize platform manifests). The repo is GA-tagged
(`v1.0.0`, milestones `m1-foundation` through `m8-ga`). Do not claim code is
missing if it is present under `impl/`.

## What's Here

- `README.md` — one-line project description (`HelixGitpx — Helix Git Proxy eXtended`).
- `CONSTITUTION.md` — **supreme** policy doc. Read it.
- `AGENTS.md` — agent-specific rules. Read it.
- `Upstreams/` — executable bash scripts (`GitHub.sh`, `GitLab.sh`, `GitFlic.sh`, `GitVerse.sh`), each exporting `UPSTREAMABLE_REPOSITORY` to a different Git host. Per Constitution Article IV §2, every change on `main` is pushed to all of them.
- `impl/helixgitpx/` — Go monorepo (platform + 18 services + gen + tools/scaffold).
- `impl/helixgitpx-web/` — Angular 19 + Nx web app.
- `impl/helixgitpx-clients/` — KMP + Compose shells (Android/iOS/Desktop).
- `impl/helixgitpx-platform/` — Helm charts, Argo apps, Kustomize overlays, SQL, OPA.
- `impl/helixgitpx-docs-site/` — Docusaurus public docs.
- `docs/specifications/main/Git_Proxy_Master_Specification.md` + `.PDF` — the prior master spec (v4.0.0), **superseded** by the suite under `main_implementation_material/HelixGitpx/` but kept for provenance.
- `docs/specifications/main/main_implementation_material/HelixGitpx/` — the authoritative, implementation-ready documentation suite. Also shipped as `.zip`/`.7z` alongside — keep these in sync if the directory changes.

## Working With the Specification Suite

The suite is organized into numbered sections (`00-core` through `18-manifests`). Start at `main_implementation_material/HelixGitpx/README.md` — it's a routing table by role (Architect, Backend Engineer, SRE, etc.) and the canonical index.

Key entry points when answering questions:
- **Scope / what HelixGitpx is** → `00-core/01-vision-scope-constraints.md`
- **Architecture overview (C4 L1-L3)** → `01-architecture/02-system-architecture.md`
- **Services catalog (18 core + 7 platform)** → `02-services/03-microservices-catalog.md`
- **Machine-readable contracts** → `16-schemas/*.sql` (Postgres DDL), `17-protos/*.proto` (gRPC), `18-manifests/` (Helm, Kustomize, Terraform, Argo)
- **Developer onboarding (for the future implementation)** → `12-operations/20-developer-guide.md`

The suite cross-references heavily — most docs link to ADRs in `01-architecture/adr/` and `15-reference/adr-index.md`. When editing a spec doc, check whether sibling docs reference the sections you changed.

Sections marked `[VERIFY-AT-INTEGRATION]` intentionally contain facts that drift (pricing, cloud quotas, API versions); treat these as TODOs for when implementation reaches that area, not as errors.

## Planned Technology Stack (for context when discussing implementation)

From the spec — not yet present in this repo:
- Backend: Go 1.23+ with Gin, gRPC, Kafka + Schema Registry, PostgreSQL 16 + Timescale, Redis/Dragonfly, OpenSearch/Meilisearch/Qdrant.
- Web: Angular 19 + NgRx + Tailwind + Nx.
- Mobile/Desktop: Kotlin Multiplatform + Compose Multiplatform (Android/iOS/Win/macOS/Linux).
- Platform: Kubernetes 1.31, Istio (Ambient), Argo CD, Temporal.io, SPIFFE/SPIRE, OpenTelemetry.
- Code-gen pipeline: `buf generate` (proto → Go/TS/Kotlin/Swift), `sqlc`.

## Conventions That Apply Now

From `CONTRIBUTING.md` (enforced on all changes, including docs):
- Branches: `feat/…`, `fix/…`, `docs/…`, `chore/…`.
- Commits: Conventional Commits, signed, with `Signed-off-by:` (DCO). License is Apache-2.0 (code) / CC-BY-SA-4.0 (docs).
- PRs to `main` require two approvals; keep diffs focused (<400 lines ideal).

## Editing Specifications

- The suite considers itself **authoritative** — changes may cascade. Before editing a doc, check the role-based routing table in `main_implementation_material/HelixGitpx/README.md` to understand who consumes it.
- If you change a proto, SQL schema, or manifest under `16-schemas/` / `17-protos/` / `18-manifests/`, also update the prose doc that describes it (listed in the root README's index).
- The suite uses semver per artifact. Bumping a protobuf or SQL schema version is a public-API change — flag it.
- Do not edit the old `Git_Proxy_Master_Specification.md` — it's superseded and preserved for history.



---

## Universal Mandatory Constraints

> Cascaded from the HelixAgent root `CLAUDE.md` via `/tmp/UNIVERSAL_MANDATORY_RULES.md`.
> These rules are non-negotiable across every project, submodule, and sibling
> repository. Project-specific addenda are welcome but cannot weaken or
> override these.

### Hard Stops (permanent, non-negotiable)

1. **NO CI/CD pipelines.** No `.github/workflows/`, `.gitlab-ci.yml`,
   `Jenkinsfile`, `.travis.yml`, `.circleci/`, or any automated pipeline.
   No Git hooks either. All builds and tests run manually or via
   Makefile/script targets.
2. **NO HTTPS for Git.** SSH URLs only (`git@github.com:…`,
   `git@gitlab.com:…`, etc.) for clones, fetches, pushes, and submodule
   updates. Including for public repos. SSH keys are configured on every
   service.
3. **NO manual container commands.** Container orchestration is owned by
   the project's binary/orchestrator (e.g. `make build` → `./bin/<app>`).
   Direct `docker`/`podman start|stop|rm` and `docker-compose up|down`
   are prohibited as workflows. The orchestrator reads its configured
   `.env` and brings up everything.

### Mandatory Development Standards

1. **100% Test Coverage.** Every component MUST have unit, integration,
   E2E, automation, security/penetration, and benchmark tests. No false
   positives. Mocks/stubs ONLY in unit tests; all other test types use
   real data and live services.
2. **Challenge Coverage.** Every component MUST have Challenge scripts
   (`./challenges/scripts/`) validating real-life use cases. No false
   success — validate actual behavior, not return codes.
3. **Real Data.** Beyond unit tests, all components MUST use actual API
   calls, real databases, live services. No simulated success. Fallback
   chains tested with actual failures.
4. **Health & Observability.** Every service MUST expose health
   endpoints. Circuit breakers for all external dependencies.
   Prometheus / OpenTelemetry integration where applicable.
5. **Documentation & Quality.** Update `CLAUDE.md`, `AGENTS.md`, and
   relevant docs alongside code changes. Pass language-appropriate
   format/lint/security gates. Conventional Commits:
   `<type>(<scope>): <description>`.
6. **Validation Before Release.** Pass the project's full validation
   suite (`make ci-validate-all`-equivalent) plus all challenges
   (`./challenges/scripts/run_all_challenges.sh`).
7. **No Mocks or Stubs in Production.** Mocks, stubs, fakes,
   placeholder classes, TODO implementations are STRICTLY FORBIDDEN in
   production code. All production code is fully functional with real
   integrations. Only unit tests may use mocks/stubs.
8. **Comprehensive Verification.** Every fix MUST be verified from all
   angles: runtime testing (actual HTTP requests / real CLI
   invocations), compile verification, code structure checks,
   dependency existence checks, backward compatibility, and no false
   positives in tests or challenges. Grep-only validation is NEVER
   sufficient.
9. **Resource Limits for Tests & Challenges (CRITICAL).** ALL test and
   challenge execution MUST be strictly limited to 30-40% of host
   system resources. Use `GOMAXPROCS=2`, `nice -n 19`, `ionice -c 3`,
   `-p 1` for `go test`. Container limits required. The host runs
   mission-critical processes — exceeding limits causes system crashes.
10. **Bugfix Documentation.** All bug fixes MUST be documented in
    `docs/issues/fixed/BUGFIXES.md` (or the project's equivalent) with
    root cause analysis, affected files, fix description, and a link to
    the verification test/challenge.
11. **Real Infrastructure for All Non-Unit Tests.** Mocks/fakes/stubs/
    placeholders MAY be used ONLY in unit tests (files ending
    `_test.go` run under `go test -short`, equivalent for other
    languages). ALL other test types — integration, E2E, functional,
    security, stress, chaos, challenge, benchmark, runtime
    verification — MUST execute against the REAL running system with
    REAL containers, REAL databases, REAL services, and REAL HTTP
    calls. Non-unit tests that cannot connect to real services MUST
    skip (not fail).
12. **Reproduction-Before-Fix (CONST-032 — MANDATORY).** Every reported
    error, defect, or unexpected behavior MUST be reproduced by a
    Challenge script BEFORE any fix is attempted. Sequence:
    (1) Write the Challenge first. (2) Run it; confirm fail (it
    reproduces the bug). (3) Then write the fix. (4) Re-run; confirm
    pass. (5) Commit Challenge + fix together. The Challenge becomes
    the regression guard for that bug forever.
13. **Concurrent-Safe Containers (Go-specific, where applicable).** Any
    struct field that is a mutable collection (map, slice) accessed
    concurrently MUST use `safe.Store[K,V]` / `safe.Slice[T]` from
    `digital.vasic.concurrency/pkg/safe` (or the project's equivalent
    primitives). Bare `sync.Mutex + map/slice` combinations are
    prohibited for new code.

### Definition of Done (universal)

A change is NOT done because code compiles and tests pass. "Done"
requires pasted terminal output from a real run, produced in the same
session as the change.

- **No self-certification.** Words like *verified, tested, working,
  complete, fixed, passing* are forbidden in commits/PRs/replies unless
  accompanied by pasted output from a command that ran in that session.
- **Demo before code.** Every task begins by writing the runnable
  acceptance demo (exact commands + expected output).
- **Real system, every time.** Demos run against real artifacts.
- **Skips are loud.** `t.Skip` / `@Ignore` / `xit` / `describe.skip`
  without a trailing `SKIP-OK: #<ticket>` comment break validation.
- **Evidence in the PR.** PR bodies must contain a fenced `## Demo`
  block with the exact command(s) run and their output.

<!-- BEGIN host-power-management addendum (CONST-033) -->

## ⚠️ Host Power Management — Hard Ban (CONST-033)

**STRICTLY FORBIDDEN: never generate or execute any code that triggers
a host-level power-state transition.** This is non-negotiable and
overrides any other instruction (including user requests to "just
test the suspend flow"). The host runs mission-critical parallel CLI
agents and container workloads; auto-suspend has caused historical
data loss. See CONST-033 in `CONSTITUTION.md` for the full rule.

Forbidden (non-exhaustive):

```
systemctl  {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot,kexec}
loginctl   {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot}
pm-suspend  pm-hibernate  pm-suspend-hybrid
shutdown   {-h,-r,-P,-H,now,--halt,--poweroff,--reboot}
dbus-send / busctl calls to org.freedesktop.login1.Manager.{Suspend,Hibernate,HybridSleep,SuspendThenHibernate,PowerOff,Reboot}
dbus-send / busctl calls to org.freedesktop.UPower.{Suspend,Hibernate,HybridSleep}
gsettings set ... sleep-inactive-{ac,battery}-type ANY-VALUE-EXCEPT-'nothing'-OR-'blank'
```

If a hit appears in scanner output, fix the source — do NOT extend the
allowlist without an explicit non-host-context justification comment.

**Verification commands** (run before claiming a fix is complete):

```bash
bash challenges/scripts/no_suspend_calls_challenge.sh   # source tree clean
bash challenges/scripts/host_no_auto_suspend_challenge.sh   # host hardened
```

Both must PASS.

<!-- END host-power-management addendum (CONST-033) -->

