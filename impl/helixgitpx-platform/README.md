# helixgitpx-platform

GitOps and infrastructure artifacts for HelixGitpx.

## Contents

| Path | Purpose |
|---|---|
| `compose/` | Local development stack (Postgres, Kafka, Redis, observability, hello) — see [ADR-0002] for the runtime wrapper |
| `k8s-local/` | `kind`/`k3d` cluster scripts + Tiltfile (full-platform rehearsal path) |
| `kyverno/policies/` | Cluster admission policies (enforced in M2+) |
| `checkov/` | IaC static analysis config |
| `github-actions-runner-controller/` | Self-hosted Kata runner configuration (M2 activation) |
| `vault/` | Vault policies + Terraform module + OIDC role (M2 activation) |
| `argocd/` | Argo CD Application definitions (M2 activation) |
| `helm/` | Umbrella Helm chart (M2+) |
| `kustomize/` | Kustomize overlays per environment |
| `terraform/` | Root Terraform modules |
| `Tiltfile` | Alternative local dev path on Kubernetes (vs. compose) |

## Local stack

```sh
make dev          # bring up everything
make dev-down     # tear down (removes volumes)
```

The compose file lives at `compose/compose.yml`. The `compose/bin/compose` wrapper auto-detects `docker`, `podman`, or `podman-compose` — never invoke them directly.

[ADR-0002]: ../../docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0002-portable-container-runtime.md
