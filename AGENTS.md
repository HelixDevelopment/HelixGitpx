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
