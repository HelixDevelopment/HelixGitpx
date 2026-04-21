#!/usr/bin/env bash
# verify-proto-gen.sh — ensures the committed `gen/` tree is in sync with
# the .proto sources. Runs `buf generate` into a temp dir and diffs it
# against impl/helixgitpx/gen/. Exits 1 if any file differs.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if ! command -v buf >/dev/null 2>&1; then
    echo "buf not installed — skipping proto-gen check."
    exit 0
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# Snapshot the current gen/ tree.
cp -r impl/helixgitpx/gen "$tmp/gen-current"

# Regenerate into the current tree, diff, then restore.
(cd impl/helixgitpx/proto && buf generate --template buf.gen.yaml --include-imports)

if ! diff -r -q "$tmp/gen-current" impl/helixgitpx/gen >/dev/null 2>&1; then
    echo "Proto gen drift detected — committed gen/ is out of sync with proto/."
    echo "Run 'cd impl/helixgitpx/proto && buf generate' and commit the result."
    diff -r --brief "$tmp/gen-current" impl/helixgitpx/gen | head -20
    # Restore so a dirty checkout doesn't linger.
    rm -rf impl/helixgitpx/gen
    mv "$tmp/gen-current" impl/helixgitpx/gen
    exit 1
fi

echo "Proto gen in sync."
