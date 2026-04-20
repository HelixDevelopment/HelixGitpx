#!/usr/bin/env bash
# Load locally-built images into the cluster (avoids pushing to a registry).
set -euo pipefail
engine="${KIND:-0}"
images=("$@")
if [ "${#images[@]}" = 0 ]; then
  images=("helixgitpx/hello:dev")
fi
for img in "${images[@]}"; do
  if [ "$engine" = "1" ]; then
    kind load docker-image "$img" --name helix
  else
    k3d image import -c helix "$img"
  fi
done
