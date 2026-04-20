#!/usr/bin/env bash
# Bring up a local Kubernetes cluster. Defaults to k3d; use KIND=1 to pick kind.
# --dry-run prints planned actions without executing.
# --m2 runs the full M2 data-plane bootstrap (cluster + Argo CD app-of-apps).
set -euo pipefail

DRY_RUN=0
M2=0
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
    --m2)      M2=1 ;;
  esac
done

engine="${KIND:-0}"
here=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH='' cd -- "$here/../../.." && pwd)

run() {
  if [ "$DRY_RUN" = "1" ]; then
    printf '[dry-run] %s\n' "$*"
  else
    "$@"
  fi
}

if [ "$M2" = "1" ]; then
    "$here/m2/preflight.sh"
    "$here/m2/prepull-images.sh"
fi

# k3d cluster creation — for M2 we disable flannel + k3s network-policy so Cilium can take over
if [ "$M2" = "1" ]; then
    run k3d cluster create \
        --config "$here/k3d-config.yaml" \
        --k3s-arg '--flannel-backend=none@server:*' \
        --k3s-arg '--disable-network-policy@server:*'
else
    if [ "$engine" = "1" ]; then
        run kind create cluster --config "$here/kind-config.yaml"
    else
        run k3d cluster create --config "$here/k3d-config.yaml"
    fi
fi

run kubectl cluster-info

if [ "$M2" = "1" ]; then
    # Namespaces (8 for M2)
    for ns in helix-system helix-identity istio-system helix-data helix-cache helix-secrets helix-observability helix; do
        run kubectl create ns "$ns" --dry-run=client -o yaml | kubectl apply -f -
    done

    # Cilium via Helm
    run helm repo add cilium https://helm.cilium.io
    run helm upgrade --install cilium cilium/cilium \
        -n kube-system \
        -f "$repo_root/impl/helixgitpx-platform/helm/cilium/values-local.yaml" \
        --version 1.16.3

    # Wait for Cilium Ready
    run kubectl -n kube-system rollout status ds/cilium --timeout=5m

    # /etc/hosts entries (non-destructive — skips if already present)
    "$here/m2/etc-hosts.sh"

    # Argo CD bootstrap — delegated to a separate step (Task 17) that kustomize-applies argocd/bootstrap
    printf 'cluster ready. Next: apply impl/helixgitpx-platform/argocd/bootstrap/ (Task 17).\n'
fi
