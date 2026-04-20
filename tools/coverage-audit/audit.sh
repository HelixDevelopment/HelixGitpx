#!/usr/bin/env bash
# Enumerate Go packages under impl/helixgitpx/ and report coverage per package.
# Flag any package below $MIN_COVERAGE (default 80%).
set -euo pipefail

MIN_COVERAGE=${MIN_COVERAGE:-80}
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT/impl/helixgitpx"

failures=0
total_pkgs=0

mods=$(find . -name go.mod -not -path './gen/*' -exec dirname {} \; | sort)

for mod in $mods; do
    [ -d "$mod" ] || continue
    pushd "$mod" >/dev/null
    for pkg in $(GOTOOLCHAIN=go1.23.4 go list ./... 2>/dev/null); do
        total_pkgs=$((total_pkgs+1))
        # Extract first numeric coverage value only; default 0 when test output
        # is empty or the package has no test files.
        set +e
        raw=$(GOTOOLCHAIN=go1.23.4 go test -cover -count=1 "$pkg" 2>&1)
        pct=$(printf '%s' "$raw" | grep -oE 'coverage: [0-9]+\.[0-9]+%' | head -n 1 | grep -oE '[0-9]+\.[0-9]+' | head -n 1)
        set -e
        pct=${pct:-0}
        status="OK "
        # shellcheck disable=SC2016
        below=$(awk -v p="$pct" -v t="$MIN_COVERAGE" 'BEGIN { print (p+0 < t+0) ? 1 : 0 }')
        if [ "$below" = "1" ]; then
            status="LOW"
            failures=$((failures+1))
        fi
        printf '  [%-3s] %6.2f%%  %s\n' "$status" "$pct" "$pkg"
    done
    popd >/dev/null
done

echo ""
echo "Audited $total_pkgs packages. Below ${MIN_COVERAGE}%: $failures."
[ "$failures" -eq 0 ]
