#!/usr/bin/env bash
# Bring up a local Kubernetes cluster. Defaults to k3d; use KIND=1 to pick kind.
# --dry-run prints planned actions without executing.
set -euo pipefail

DRY_RUN=0
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
  esac
done

engine="${KIND:-0}"
here=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)

run() {
  if [ "$DRY_RUN" = "1" ]; then
    printf '[dry-run] %s\n' "$*"
  else
    "$@"
  fi
}

if [ "$engine" = "1" ]; then
  run kind create cluster --config "$here/kind-config.yaml"
else
  run k3d cluster create --config "$here/k3d-config.yaml"
fi

run kubectl cluster-info
