#!/usr/bin/env bash
# Build + sign MSIX. Requires signing cert PFX + password.
set -euo pipefail
SRC=${1:-impl/helixgitpx-clients/desktopApp/build/compose/binaries/main/msi}
echo "[publish-windows] would package and sign MSIX from $SRC"
