#!/usr/bin/env bash
# Notarize + staple a DMG for macOS. Requires Apple developer credentials.
set -euo pipefail
DMG=${1:-impl/helixgitpx-clients/desktopApp/build/compose/binaries/main/dmg/HelixGitpx.dmg}
[ -f "$DMG" ] || { echo "DMG not found"; exit 2; }
echo "[publish-macos] would submit $DMG to notarytool + stapler"
