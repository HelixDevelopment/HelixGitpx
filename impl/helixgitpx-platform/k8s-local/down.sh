#!/usr/bin/env bash
set -euo pipefail
engine="${KIND:-0}"
if [ "$engine" = "1" ]; then
  kind delete cluster --name helix
else
  k3d cluster delete helix
fi
