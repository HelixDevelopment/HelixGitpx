#!/usr/bin/env bash
# Publish an iOS IPA to App Store Connect / TestFlight. Requires APPLE_API_KEY.
set -euo pipefail
IPA=${1:-impl/helixgitpx-clients/iosApp/build/HelixGitpx.ipa}
[ -f "$IPA" ] || { echo "IPA not found: $IPA"; exit 2; }
echo "[publish-appstore] would upload $IPA via xcrun altool"
# xcrun altool --upload-app -f "$IPA" --apiKey "$APPLE_API_KEY_ID" --apiIssuer "$APPLE_API_ISSUER"
