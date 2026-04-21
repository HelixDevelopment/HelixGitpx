# .github/workflows — disabled

All GitHub Actions workflows in this directory are **currently disabled**.

## How

Files were renamed from `*.yml` → `*.yml.disabled`. GitHub Actions only
picks up `.yml` and `.yaml` files, so nothing runs — not even manual
`workflow_dispatch` — while they remain suffixed.

## Why

Operator decision (2026-04-21). Re-enable individually by renaming back
to `.yml`, or all at once:

```bash
for f in .github/workflows/*.yml.disabled .github/workflows/_reusable/*.yml.disabled; do
    mv "$f" "${f%.disabled}"
done
git add -A && git commit -m "ops: re-enable GitHub Actions workflows"
```

## Inventory of disabled workflows

| File | Purpose |
|------|---------|
| `ci-clients.yml.disabled` | Gradle build + Kotlin tests for KMP clients |
| `ci-docs.yml.disabled` | Docusaurus site build check |
| `ci-go.yml.disabled` | `go vet` + `go test` across the monorepo |
| `ci-platform.yml.disabled` | Helm chart + OPA bundle checks |
| `ci-verifiers.yml.disabled` | `scripts/verify-everything.sh` green-suite |
| `ci-web.yml.disabled` | Nx build + jest + Playwright for the web app |
| `deploy.yml.disabled` | GitOps image promotion |
| `mutation-testing.yml.disabled` | `go-mutesting` per service |
| `perf-budgets.yml.disabled` | k6 scenarios + budget gate |
| `release.yml.disabled` | Tag-triggered release workflow (also dispatch-only) |
| `security-scan.yml.disabled` | SAST + DAST + container + IaC scans |
| `supply-chain.yml.disabled` | Per-service SBOM (Syft) + Cosign sign + SLSA + Trivy |
| `upstream-sync.yml.disabled` | Push `main` + tags to all federated upstreams |
| `_reusable/*.yml.disabled` | Shared callable workflows (Cosign, SBOM, Go setup, Vault creds) |

## Constitutional note

Constitution Article II §5 ("CI enforcement") remains the policy target.
While workflows are disabled, enforcement runs locally via
`bash scripts/verify-everything.sh`. Re-enable the workflow suite before
merging production changes.
