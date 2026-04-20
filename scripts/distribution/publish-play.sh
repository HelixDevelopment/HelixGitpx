#!/usr/bin/env bash
# Publish an Android AAB to Google Play Console. Requires PLAY_SERVICE_ACCOUNT_JSON.
set -euo pipefail
AAB=${1:-impl/helixgitpx-clients/androidApp/build/outputs/bundle/release/androidApp-release.aab}
[ -f "$AAB" ] || { echo "AAB not found: $AAB"; exit 2; }
[ -n "${PLAY_SERVICE_ACCOUNT_JSON:-}" ] || { echo "PLAY_SERVICE_ACCOUNT_JSON env unset"; exit 3; }
echo "[publish-play] would upload $AAB to Play Console (beta track)"
# Real implementation uses `bundletool` + Play Developer API. Deferred to M8
# release automation; this script is the contract.
