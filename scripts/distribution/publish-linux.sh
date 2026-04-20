#!/usr/bin/env bash
# Publish AppImage + .deb + .rpm artefacts. GPG-sign each.
set -euo pipefail
BUILD_DIR=${1:-impl/helixgitpx-clients/desktopApp/build/compose/binaries/main}
for pkg in "$BUILD_DIR"/deb/*.deb "$BUILD_DIR"/rpm/*.rpm "$BUILD_DIR"/appimage/*.AppImage; do
    [ -f "$pkg" ] || continue
    echo "[publish-linux] would GPG-sign + upload $pkg"
done
