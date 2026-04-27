# HelixGitpx Constitution

> The Constitution is the highest-authority policy document in this repo.
> Every contribution — human or AI — must comply with the articles below.
> Subordinate documents (`CLAUDE.md`, `AGENTS.md`, `CONTRIBUTING.md`, specs)
> may refine but never contradict these rules.
>
> **Version:** 1.0.0 · **Ratified:** 2026-04-20 · **Author:** Милош Васић
> (@milos85vasic).

---

## Article I — Scope and authority

1. The Constitution governs every artifact produced or modified in this
   repository: production code, tests, tooling, documentation, specs,
   manifests, and generated files.
2. Where another document (CLAUDE.md, AGENTS.md, or a skill) conflicts with
   the Constitution, the Constitution wins.
3. Changes to the Constitution require an ADR under
   `docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/`
   and a `governance:` tagged commit.

---

## Article II — Testing

This article is mandatory. No exceptions.

### §1. Test-type coverage requirement

Every shipped module **must** carry tests across the following test types.
Targets are stated as coverage percentages on the modules they apply to.

| Type         | Target | Mocks allowed? | Purpose |
|--------------|--------|---------------|---------|
| unit         | 100 %  | **Yes**        | Function-level correctness |
| integration  | 100 %  | No             | Real collaborators, real deps |
| e2e          | 100 %  | No             | Full user/workflow journeys |
| security     | 100 %  | No             | Authn/z, injection, ASVS L2 |
| stress       | 100 %  | No             | Load up to 3× design target |
| ddos         | 100 %  | No             | Rate-limit, exhaustion, and recovery |
| benchmark    | 100 %  | No             | Latency, throughput, regression |

Additional test types (chaos, mutation, fuzz, smoke, soak, contract,
compatibility, accessibility, localisation, visual-regression,
property-based) are welcome and encouraged. They do **not** reduce the
required seven.

"100 %" means *every public path and every external-facing behaviour is
exercised*. Coverage tools (Go's `-cover`, Jest `--coverage`, Kover, etc.) are
the baseline signal. The bar is behaviour coverage, not just line coverage.

### §2. Mocks, stubs, placeholders, hardcoded data — restricted

Only **unit tests** may use mocks, stubs, placeholder classes, or hardcoded
data. Every other test type **must** exercise real dependencies — real
databases, real Kafka brokers, real Keycloak, real OPA bundles, real
Git servers — run under ephemeral compose / k3d / kind / testcontainers.

Rationale: mocks lie. They mask real integration, ordering, and failure
behaviour. We learned this the hard way; see ADR-0042 (logging and trust).

### §3. Reliability

No test in the repository may be:

- **Skipped** (`t.Skip`, `xit`, `@Ignore`, `pytest.skip`, etc.).
- **Disabled** via annotation, directive, or pragma.
- **Broken** (test that fails when its subject is correct).
- **Faulty / flaky** (non-deterministic within the same build).

If a test cannot pass today, the underlying issue must be fixed or the
feature removed. A skipped test is a silent debt this project will not carry.

### §4. Root-cause discipline

When a test reveals an issue, fix the **root cause**, not the symptom.
Disabling, retrying-until-green, or broadening tolerances are prohibited.
Every fix must include a regression test in the relevant test type(s) so
the same issue cannot reappear undetected.

### §5. CI enforcement

The CI pipeline refuses to merge a PR unless:

- All seven required test types pass with zero skipped.
- Coverage is 100 % on every module touched by the PR (measured per type).
- Mutation score ≥ 60 % on the units that have unit tests.
- Security-test scans (SAST, DAST, SBOM, Secret, Image, IaC) are clean.

The enforcement scripts live at `.github/workflows/ci-*.yml` and
`tools/coverage-audit/`.

---

## Article III — Documentation

### §1. Documentation is source

Every feature, subsystem, and service **must** have documentation under
`docs/` that explains its purpose, inputs, outputs, failure modes, runbook
links, and ADR references. Code without matching documentation does not
ship.

### §2. Multi-format delivery

Public user-facing manuals must be produced in:

- HTML (via Docusaurus, `docs.helixgitpx.io`).
- PDF (via `pandoc`).
- ePub (via `pandoc`).
- Markdown source (authoritative).
- Plain text (for accessibility).

### §3. Media

A parallel video curriculum mirrors every major documentation section.
Scripts live under `docs/media/video-scripts/` and production assets under
`docs/media/video/` (outside git, large files).

---

## Article IV — Versioning and distribution

### §1. Public surfaces are semver

Each public artifact (proto, SQL schema, REST API, CLI flag, Helm chart,
container image) carries its own `X.Y.Z` version. Breaking changes require a
major bump and an ADR-documented migration path.

### §2. Upstream federation

This project practices the federation pattern it specifies. Every change on
`main` is pushed to all configured upstreams (GitHub, GitLab, GitFlic,
GitVerse, and any further targets). See `Upstreams/` scripts and the
`docs/operations/upstream-sync.md` runbook.

### §3. Regular cadence

Pushes to all upstreams happen at **least** daily and on every tagged
release. A scheduled CI job (`workflow_dispatch` + manual trigger) enforces
this.

---

## Article V — Governance of AI contributors

### §1. Equal treatment

Contributions from AI systems (Claude, other agents) are bound by the
Constitution identically to human contributions. AI cannot invoke a "it's
just an agent" exemption for any article.

### §2. Attribution

AI-assisted commits must include a `Co-Authored-By: <model>` trailer
alongside the human `Signed-off-by:` line.

### §3. Honesty

AI agents must not fabricate implementations, simulate passing tests, or
produce "impressive-looking" stubs that are not wired. A scaffold that is
intentionally unwired must be clearly labelled as such with a TODO that
includes the owning milestone.

---

## Article VI — Security and privacy

1. No secrets in git. Ever. Secret-scanning CI must remain green on `main`.
2. mTLS everywhere east-west; TLS 1.3 north-south.
3. Default-deny OPA authorization; every external surface has a policy.
4. Data residency is a per-org choice. The `org.organizations.residency`
   column is the authoritative source.

---

## Article VII — Amendment

1. Propose an amendment by opening a PR that edits this file.
2. Attach an ADR that explains the motivation and alternatives.
3. Require sign-off from two code-owners.
4. On merge, bump the version at the top of this document.



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

### CONST-033 — Host Power Management is Forbidden

**Status:** Mandatory. Non-negotiable. Applies to every project,
submodule, container entry point, build script, test, challenge, and
systemd unit shipped from this repository.

**Rule:** No code in this repository may invoke a host-level power-
state transition (suspend, hibernate, hybrid-sleep, suspend-then-
hibernate, poweroff, halt, reboot, kexec) on the host machine. This
includes — but is not limited to:

- `systemctl {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot,kexec}`
- `loginctl {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot}`
- `pm-{suspend,hibernate,suspend-hybrid}`
- `shutdown {-h,-r,-P,-H,now,--halt,--poweroff,--reboot}`
- DBus calls to `org.freedesktop.login1.Manager.{Suspend,Hibernate,HybridSleep,SuspendThenHibernate,PowerOff,Reboot}`
- DBus calls to `org.freedesktop.UPower.{Suspend,Hibernate,HybridSleep}`
- `gsettings set ... sleep-inactive-{ac,battery}-type` to any value other than `'nothing'` or `'blank'`

**Why:** The host runs mission-critical parallel CLI-agent and
container workloads. On 2026-04-26 18:23:43 the host was auto-
suspended by the GDM greeter's idle policy mid-session, killing
HelixAgent and 41 dependent services. Recurring memory-pressure
SIGKILLs of `user@1000.service` (perceived as "logged out") have the
same outcome. Auto-suspend, hibernate, and any power-state transition
are unsafe for this host.

**Defence in depth (mandatory artifacts in every project):**
1. `scripts/host-power-management/install-host-suspend-guard.sh` —
   privileged installer, manual prereq, run once per host with sudo.
   Masks `sleep.target`, `suspend.target`, `hibernate.target`,
   `hybrid-sleep.target`; writes `AllowSuspend=no` drop-in; sets
   logind `IdleAction=ignore` and `HandleLidSwitch=ignore`.
2. `scripts/host-power-management/user_session_no_suspend_bootstrap.sh` —
   per-user, no-sudo defensive layer. Idempotent. Safe to source from
   `start.sh` / `setup.sh` / `bootstrap.sh`.
3. `scripts/host-power-management/check-no-suspend-calls.sh` —
   static scanner. Exits non-zero on any forbidden invocation.
4. `challenges/scripts/host_no_auto_suspend_challenge.sh` — asserts
   the running host's state matches layer-1 masking.
5. `challenges/scripts/no_suspend_calls_challenge.sh` — wraps the
   scanner as a challenge that runs in CI / `run_all_challenges.sh`.

**Enforcement:** Every project's CI / `run_all_challenges.sh`
equivalent MUST run both challenges (host state + source tree). A
violation in either channel blocks merge. Adding files to the
scanner's `EXCLUDE_PATHS` requires an explicit justification comment
identifying the non-host context.

**See also:** `docs/HOST_POWER_MANAGEMENT.md` for full background and
runbook.

<!-- END host-power-management addendum (CONST-033) -->

