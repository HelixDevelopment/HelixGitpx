#!/usr/bin/env bash
# push-to-all-upstreams.sh — push main and all tags to every configured upstream.
# Constitution Article IV §2: regular cadence (daily minimum + every tagged release).
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

BRANCH=${BRANCH:-main}
REMOTES=()

# Collect upstreams from Upstreams/*.sh.
for script in Upstreams/*.sh; do
    [ -f "$script" ] || continue
    # shellcheck disable=SC1090
    source "$script"
    if [ -n "${UPSTREAMABLE_REPOSITORY:-}" ]; then
        REMOTES+=("$UPSTREAMABLE_REPOSITORY")
    fi
    unset UPSTREAMABLE_REPOSITORY
done

if [ ${#REMOTES[@]} -eq 0 ]; then
    echo "No upstreams configured in Upstreams/. Nothing to push." >&2
    exit 0
fi

echo "Pushing branch $BRANCH and tags to ${#REMOTES[@]} upstreams..."

fail=0
for remote_url in "${REMOTES[@]}"; do
    alias_name="upstream-$(echo "$remote_url" | sed 's#[^a-zA-Z0-9]#-#g' | cut -c1-60)"
    printf '\n--- %s (%s) ---\n' "$alias_name" "$remote_url"
    git remote remove "$alias_name" 2>/dev/null || true
    git remote add "$alias_name" "$remote_url"
    if ! git push "$alias_name" "$BRANCH"; then
        echo "FAIL: $remote_url (branch)"
        fail=$((fail + 1))
    fi
    if ! git push "$alias_name" --tags; then
        echo "FAIL: $remote_url (tags)"
        fail=$((fail + 1))
    fi
    git remote remove "$alias_name" 2>/dev/null || true
done

if [ "$fail" -gt 0 ]; then
    echo ""
    echo "$fail push(es) failed. Check SSH keys, credentials, and network." >&2
    exit 1
fi

echo ""
echo "Pushed $BRANCH + tags to all ${#REMOTES[@]} upstreams."
