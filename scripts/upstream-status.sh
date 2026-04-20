#!/usr/bin/env bash
# upstream-status.sh — show how far ahead/behind each upstream is from local HEAD.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

BRANCH=${BRANCH:-main}

for script in Upstreams/*.sh; do
    [ -f "$script" ] || continue
    # shellcheck disable=SC1090
    source "$script"
    url="${UPSTREAMABLE_REPOSITORY:-}"
    unset UPSTREAMABLE_REPOSITORY
    [ -n "$url" ] || continue

    name=$(basename "$script" .sh)
    alias_name="upstream-status-$name"
    git remote remove "$alias_name" 2>/dev/null || true
    git remote add "$alias_name" "$url"

    if git fetch --quiet "$alias_name" "$BRANCH" 2>/dev/null; then
        ahead=$(git rev-list --count "$alias_name/$BRANCH..HEAD" 2>/dev/null || echo "?")
        behind=$(git rev-list --count "HEAD..$alias_name/$BRANCH" 2>/dev/null || echo "?")
        printf '  %-12s  ahead %3s  behind %3s  %s\n' "$name" "$ahead" "$behind" "$url"
    else
        printf '  %-12s  UNREACHABLE                 %s\n' "$name" "$url"
    fi
    git remote remove "$alias_name" 2>/dev/null || true
done
