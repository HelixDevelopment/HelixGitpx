#!/usr/bin/env bash
# Publish desktop update metadata to the self-hosted update feed (tus.io).
# ADR-0029 documents the choice.
set -euo pipefail
VERSION=${1:?version required, e.g. 0.1.0}
echo "[publish-update-feed] would post manifest for version $VERSION to $TUS_ENDPOINT"
