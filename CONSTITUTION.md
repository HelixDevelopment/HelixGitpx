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
