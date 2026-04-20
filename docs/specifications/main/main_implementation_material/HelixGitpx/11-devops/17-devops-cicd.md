# 17 — DevOps & CI/CD

> **Document purpose**: Describe the **build, test, release, and deployment pipelines** that take code from a pull request to a signed artefact running in production — with every quality gate enforced.

---

## 1. Principles

1. **GitOps end-to-end** — production is defined by signed commits in Git. Argo CD reconciles.
2. **Everything signed** — commits (Gitsign), artefacts (Cosign), provenance (in-toto).
3. **Every PR shippable** — `main` is always deployable; releases cut from `main` are frequent.
4. **Small merges, flags, canaries** — feature flags decouple deploy from release.
5. **Reproducible builds** — hermetic, pinned-hash dependencies, stable timestamps.
6. **Fast feedback** — PR CI under 15 min p95; full suite in under 45 min.

---

## 2. SCM & Trunk-Based Development

- Single source of truth: our own HelixGitpx instance (dogfooding), mirrored to `github.com/vasic-digital`.
- Branch model: **trunk-based**; short-lived feature branches; merge queues.
- Commit convention: **Conventional Commits** enforced by `commitlint`.
- Signed commits required on `main` (Gitsign via Sigstore keyless).

### 2.1 Pull Request Flow

1. Developer opens PR from feature branch.
2. PR templates guide description, linked issues, checklists.
3. CI runs (see §3).
4. **Two reviewers** required (or policy-based minimum).
5. Merge queue batches PRs; re-runs CI against HEAD before merge.
6. **Squash merge** to `main` with signed commit.

---

## 3. CI Pipelines

Runners: self-hosted GitHub Actions on **Kata Containers** (hardware-isolated) or **GitLab runners on gVisor**. No third-party code runs outside sandboxes.

### 3.1 PR Pipeline (fast lane)

Targets p95 ≤ 15 minutes.

1. **Lint** — golangci-lint, eslint, ktlint, yamllint, buf lint, proto breaking-check.
2. **Format check** — gofumpt, prettier, ktlint.
3. **Unit tests** — all languages, coverage reported.
4. **Mutation sample** (quick) — Stryker / gremlins / go-mutesting on changed files only.
5. **SAST** — Snyk Code + Semgrep + CodeQL + SonarQube scanner.
6. **SCA** — Snyk Open Source + govulncheck + npm audit.
7. **Secrets scan** — Gitleaks.
8. **Policy check** — OPA conftest on manifests; Checkov; Kubesec.
9. **Build** — containerise each service; produce SBOM and Cosign signature.
10. **Integration (affected)** — only services affected by diff; Testcontainers dependencies.
11. **E2E (smoke)** — small subset tagged `@smoke`.
12. **Docs lint** — markdownlint; link checker.
13. **Frontend checks** — `nx affected`; Storybook + Chromatic visual diff.
14. **Quality Gate aggregation** — SonarQube + coverage + mutation thresholds.

If any gate fails → PR blocked.

### 3.2 Main / Release Pipeline (full lane)

Targets p95 ≤ 45 minutes.

Adds to PR pipeline:

15. **Full test matrix** — unit + integration + contract + E2E + a11y + perf baselines.
16. **Full mutation testing** with thresholds.
17. **DAST** — OWASP ZAP full scan against ephemeral deploy.
18. **Load** — k6 at 1× baseline RPS for 10 min; SLO assertions.
19. **Fuzz** — short-session fuzz (30 min) on target surfaces.
20. **Provenance generation** — SLSA provenance signed to Rekor.
21. **SBOM publish** — CycloneDX 1.5 attached as OCI artefact.
22. **Image promote** — signed image pushed to `ghcr.io/vasic-digital/helixgitpx/<svc>:<semver>` and `:<sha>`.
23. **Helm chart publish** — signed chart to OCI registry.
24. **GitOps PR** — bot opens PR in `helixgitpx-platform` updating image tags for `staging`.

### 3.3 Nightly / Scheduled

- Long-run chaos suite.
- Long-duration fuzz (several hours).
- Supply chain verification audit.
- DR drill in staging (restore backup → run integration).

---

## 4. Artefact Pipeline

```
src ──► build ──► image + SBOM ──► Cosign signed ──► provenance in Rekor
                                      │
                                      ▼
                           OCI registry (immutable tag)
                                      │
                                      ▼
                              Helm chart referencing digest
                                      │
                                      ▼
                        GitOps repo PR → Argo CD reconciles
```

Every artefact is addressable by **digest**, never by mutable tag.

---

## 5. Release Management

- **Semantic versioning** for libraries and services.
- **Calendar versioning** (CalVer) optional for the overall platform (`2026.04.1`).
- **Release cadence**:
  - Services: continuous, multiple times per day.
  - Platform: **weekly** stable releases, monthly LTS-style release summary.
  - Web: daily-to-weekly.
  - Mobile/desktop: biweekly stable, nightly channel for power users.
- **Release notes** auto-generated from commit messages + curated highlights.

---

## 6. Deployment

### 6.1 Argo CD

- Each environment is an **ApplicationSet** generating Applications per service.
- Sync policy: manual approval for `prod`; auto-sync for lower envs.
- Health checks per Application (Pods ready, custom health script).
- **Progressive delivery** via **Argo Rollouts**.

### 6.2 Canary Strategy

Per service, typical rollout:

1. Deploy new version to 1 % of pods.
2. Analysis step:
   - Prometheus queries on success-rate, p99 latency, CPU, error-budget burn.
   - Must stay within thresholds for 10 min.
3. 10 % → 25 % → 50 % → 100 %.
4. Abort on any failing analysis; automatic rollback.

### 6.3 Blue / Green (for DB-migrating services)

- Two copies of the service, old and new.
- Traffic shifted atomically after smoke + migrations pass.
- Rollback is instant (flip the label).

### 6.4 Database Migrations

- Tool: `goose` (Go) with versioned SQL migrations.
- **Expand / contract** pattern: add column → backfill → switch code → drop column (in separate releases). Never destructive in a single migration.
- Migrations run as K8s Job with `preSync` hook (Argo CD) before pod rollout.
- For Postgres online schema change: `pg_repack` + partition swaps.

### 6.5 Message Schema Migrations

- Backward-compatible Avro evolution first choice.
- For breaking: new topic `…v2`, dual-publish during overlap, cut consumers, retire old.

---

## 7. Feature Flags

- **OpenFeature** + **Unleash** self-hosted.
- Every non-trivial change behind a flag with a rollout plan.
- Flags auditable; TTLs encouraged; stale-flag reports in CI.
- Per-org, per-user, per-percentage targeting.
- Kill-switches documented in runbooks.

---

## 8. Infrastructure as Code

- **Terraform / OpenTofu** for cloud resources, Vault bootstrap, DNS.
- **Kustomize** for K8s overlays; **Helm** charts for packaging apps.
- **Tanka/Jsonnet** optional for large templated surfaces.
- All reviewed via PR; `terraform plan` commented on PR; applied by CI after approval.
- State in backend with locking (S3 + DynamoDB / GCS + Cloud Run mutex / Terraform Cloud).

---

## 9. Secrets in CI

- No long-lived secrets in runners.
- **OIDC federation**: GitHub Actions → Vault → dynamic creds (AWS, GCP, container registry).
- Cosign keyless signing via GHA OIDC → Fulcio.
- Deployment creds read from Vault at run-time; short TTL.

---

## 10. Tooling Standards (by language)

### 10.1 Go

- Version managed via `go.toolchain` pin.
- Formatter: gofumpt + goimports.
- Linter: golangci-lint with curated config (.golangci.yml).
- Coverage tool: `go test -cover -coverprofile=cover.out` + `gocovmerge` + `goverreport`.
- Benchmarks: `go test -bench=. -benchmem`; baselines stored as JSON; regression threshold 5 %.
- Security: govulncheck + gosec + gotestsum for test UX.
- Dependency management: `go mod tidy`; vendoring optional.

### 10.2 TypeScript / Angular

- Nx monorepo; `pnpm` as package manager.
- ESLint + Prettier.
- Jest with coverage provider `v8`.
- Storybook + Chromatic.
- Playwright for E2E.

### 10.3 Kotlin / KMP

- Gradle 8.8+ with Convention plugins.
- Detekt + ktlint.
- Kover for coverage.
- Dependency updates via Renovate.

### 10.4 Rust (WASM plugin SDK)

- `rustfmt` + `clippy`.
- `cargo test` + `cargo-mutants`.
- `cargo audit`.

### 10.5 Python (AI sidecars)

- `ruff` + `black`.
- `pytest` + `pytest-cov` (≥ 100 %).
- Poetry for deps.

---

## 11. Monorepo vs. Polyrepo

- Backend: monorepo `helixgitpx` with service-specific subdirs + shared `helix-platform` libs.
- Mobile/Desktop: monorepo `helixgitpx-clients`.
- Frontend: monorepo `helixgitpx-web` (Nx).
- Deployment: `helixgitpx-platform` (GitOps).
- Plugins: per-plugin repo (`helixgitpx-plugin-*`).

Trade-off accepted: larger backend repo but atomic cross-service refactors.

---

## 12. Observability of CI

- CI metrics → Prometheus (GitHub Actions exporter).
- **PR cycle time** (open → merge) tracked; target p50 ≤ 1 day.
- Flaky test dashboard; auto-quarantine after 3 flaky runs.
- Build cache hit ratio (Nx, Gradle build cache, Go build cache).

---

## 13. Cost Controls

- Spot runners for non-critical jobs.
- Nx / Turborepo / Gradle build cache on shared S3.
- Test parallelism via worker autoscaling.

---

## 14. Runbook Links

- Rolling deploy emergency stop.
- Rollback (Argo CD sync to previous revision).
- Hotfix without waiting for next release (cherry-pick + fast-track).
- Secret rotation in CI.
- Break-glass access to prod.

See [19-operations-runbook.md](../12-operations/19-operations-runbook.md).

---

*— End of DevOps & CI/CD —*
