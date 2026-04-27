# AGENTS.md

Instructions for AI agents (Claude Code, Cursor, Copilot, Aider, or any
future model-driven contributor) working in this repository.

Read this file **fully** before taking any action. If a rule here conflicts
with the [Constitution](./CONSTITUTION.md), the Constitution wins.

---

## 0. The Big Five

These five rules trump everything else an agent might want to do:

1. **No mocks outside unit tests.** Constitution Article II §2. Mocks,
   stubs, placeholder classes, and hardcoded data are allowed **only** in
   unit tests. Integration / e2e / security / stress / ddos / benchmark
   tests must exercise real dependencies.
2. **No skipped or disabled tests.** If you cannot make a test pass today,
   either fix the code or remove the feature. Silent debt is forbidden.
3. **Workflow-dispatch only CI.** Every GitHub Actions workflow must use
   `on: workflow_dispatch:` only. No `push`, no `pull_request`, no `schedule`.
   This is mandatory.
4. **Runtime-portable tooling.** Local dev tooling auto-detects between
   Docker and Podman. Never hardcode either one.
5. **Cite, don't fabricate.** If you claim a file, function, repo, or fact
   exists, verify it first. When uncertain, say so.

---

## 1. Repository layout quick reference

- `impl/helixgitpx/` — Go monorepo (18 services + platform + gen + tools/scaffold).
- `impl/helixgitpx-web/` — Angular 19 web app.
- `impl/helixgitpx-clients/` — KMP + Compose Multiplatform shells.
- `impl/helixgitpx-platform/` — Helm charts, Argo CD apps, Kustomize overlays, SQL, OPA.
- `impl/helixgitpx-docs-site/` — Docusaurus public docs.
- `docs/specifications/` — authoritative spec suite. Do not edit the
  superseded master spec `Git_Proxy_Master_Specification.md`.
- `docs/integrations/` — planning docs for future integrations.
- `docs/operations/runbooks/` — production runbooks.
- `docs/security/` — pen-test scope, bug bounty, SOC 2, ISO 27001.
- `docs/superpowers/` — spec + plan docs driven by the superpowers skill.
- `tools/` — coverage-audit, perf, fuzz, chaos, dr, e2e-gaps.
- `scripts/` — verifiers (`verify-m{1..8}-*.sh`).
- `Upstreams/` — multi-upstream mirroring scripts.

---

## 2. Before you edit

- `git log --oneline -20` to orient.
- Read the most recent spec for the area you're touching, indexed at
  `docs/specifications/main/main_implementation_material/HelixGitpx/README.md`.
- Check for an open plan at `docs/superpowers/plans/` — if one exists for
  your task, follow it.
- Check for an ADR under
  `docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/`
  that covers the subject.

---

## 3. Commit and push discipline

- Conventional Commits (`feat(m7): …`, `fix(auth): …`, `docs(integrations): …`).
- Every commit signed off (`Signed-off-by: Name <email>`).
- Every commit ending with `Co-Authored-By: <model> <noreply@anthropic.com>`
  when an AI agent wrote or materially shaped the patch.
- Never force-push to `main`. Never amend a published commit.
- After each feature PR merge, run `bash Upstreams/GitHub.sh` and
  corresponding scripts to sync to all upstreams.

---

## 4. Test types — reference

Per [Constitution Article II §1](./CONSTITUTION.md#2-article-ii--testing),
every module needs:

| Type         | Location convention                           | Runner                   |
|--------------|-----------------------------------------------|--------------------------|
| unit         | `*_test.go`, `*.spec.ts`, `commonTest/`       | `go test`, jest, kotest  |
| integration  | `test/integration/` + `*_integration_test.go` | compose up, go test -tags=integration |
| e2e          | `test/e2e/` (per app)                         | Playwright, Appium/XCUITest, k6 scenario |
| security     | `test/security/`                              | OWASP ZAP, Nuclei, custom |
| stress       | `tools/perf/scenarios/*_stress.js`            | k6 |
| ddos         | `tools/chaos/` + dedicated ddos run           | Litmus + k6 arrival bursts |
| benchmark    | `*_bench_test.go`, `tools/perf/benchmark/`    | `go test -bench`, k6 |

Coverage targets 100 % per type per module touched. See
`tools/coverage-audit/audit.sh` and the `tests/` subdirs under each
service for templates.

---

## 5. Documentation requirements

- Every new feature ships with prose documentation. No exceptions.
- When you modify a proto, SQL schema, or manifest, update the sibling
  prose doc that references it (see the root README in
  `docs/specifications/…/HelixGitpx/`).
- User-facing features need entries in `impl/helixgitpx-docs-site/docs/`.
  The docs build must stay green.

---

## 6. When you finish work

1. Run the local verifier suite: `bash scripts/verify-m1-artifacts.sh`
   through `bash scripts/verify-m8-cluster.sh`.
2. Run `go test ./...` in every Go module touched.
3. Run `npx nx test web` and `npx nx build web` if you touched the web app.
4. Run `gradle :shared:check` if you touched KMP code.
5. Run `npx docusaurus build` if you touched docs-site.
6. Commit, push, and sync upstreams.

---

## 7. Escalations

- Anything that changes the Constitution → stop, ask the human.
- Anything that changes `mandatory` policies (CI dispatch-only, container
  runtime portability, testing rules) → stop, ask the human.
- Deleting tests or skipping them → stop, ask the human.
- Force-pushing, destructive git, rewriting published history → stop, ask.



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

## Host Power Management — Hard Ban (CONST-033)

**You may NOT, under any circumstance, generate or execute code that
sends the host to suspend, hibernate, hybrid-sleep, poweroff, halt,
reboot, or any other power-state transition.** This rule applies to:

- Every shell command you run via the Bash tool.
- Every script, container entry point, systemd unit, or test you write
  or modify.
- Every CLI suggestion, snippet, or example you emit.

**Forbidden invocations** (non-exhaustive — see CONST-033 in
`CONSTITUTION.md` for the full list):

- `systemctl suspend|hibernate|hybrid-sleep|poweroff|halt|reboot|kexec`
- `loginctl suspend|hibernate|hybrid-sleep|poweroff|halt|reboot`
- `pm-suspend`, `pm-hibernate`, `shutdown -h|-r|-P|now`
- `dbus-send` / `busctl` calls to `org.freedesktop.login1.Manager.Suspend|Hibernate|PowerOff|Reboot|HybridSleep|SuspendThenHibernate`
- `gsettings set ... sleep-inactive-{ac,battery}-type` to anything but `'nothing'` or `'blank'`

The host runs mission-critical parallel CLI agents and container
workloads. Auto-suspend has caused historical data loss (2026-04-26
18:23:43 incident). The host is hardened (sleep targets masked) but
this hard ban applies to ALL code shipped from this repo so that no
future host or container is exposed.

**Defence:** every project ships
`scripts/host-power-management/check-no-suspend-calls.sh` (static
scanner) and
`challenges/scripts/no_suspend_calls_challenge.sh` (challenge wrapper).
Both MUST be wired into the project's CI / `run_all_challenges.sh`.

**Full background:** `docs/HOST_POWER_MANAGEMENT.md` and `CONSTITUTION.md` (CONST-033).

<!-- END host-power-management addendum (CONST-033) -->

