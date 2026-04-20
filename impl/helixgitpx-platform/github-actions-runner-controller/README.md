# actions-runner-controller — Kata

Activation is M2 (requires a running Kubernetes cluster). Artifacts here are
config-only until then. Workflows keep `runs-on: ubuntu-latest` for M1.

## Activation checklist (M2)

1. Install ARC: `helm install arc -n arc-system oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller -f values.yaml`
2. Install Kata runtimeclass: `kubectl apply -f kata-runtimeclass.yaml`
3. Label nodes with `katacontainers.io/kata-runtime=true`.
4. Create GitHub App + secret `arc-github-token` (via Vault).
5. Apply `runner-scale-set.yaml`.
6. Flip workflow `runs-on: ubuntu-latest` → `runs-on: helix-kata`.
