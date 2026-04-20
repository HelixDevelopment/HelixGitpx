#!/usr/bin/env bash
# Enumerate Go packages under impl/helixgitpx/ and report coverage per package.
# Flag any package below $MIN_COVERAGE (default 80%).
set -euo pipefail

MIN_COVERAGE=${MIN_COVERAGE:-80}
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT/impl/helixgitpx"

failures=0
total_pkgs=0

for mod in $(go list -m -json ./... 2>/dev/null | python3 -c 'import sys,json; [print(m["Dir"]) for m in map(json.loads, sys.stdin.read().split("}\n{")) if "Dir" in m]' 2>/dev/null || find . -name go.mod -exec dirname {} \;); do
    [ -d "$mod" ] || continue
    pushd "$mod" >/dev/null
    for pkg in $(GOTOOLCHAIN=go1.23.4 go list ./... 2>/dev/null); do
        total_pkgs=$((total_pkgs+1))
        out=$(GOTOOLCHAIN=go1.23.4 go test -cover -count=1 "$pkg" 2>&1 | grep -oE 'coverage: [0-9.]+%' | head -1 | grep -oE '[0-9.]+' || echo "0")
        pct=${out:-0}
        status="OK"
        cmp=$(python3 -c "print(1 if float('$pct') < float('$MIN_COVERAGE') else 0)")
        if [ "$cmp" = "1" ]; then
            status="LOW"
            failures=$((failures+1))
        fi
        printf "  [%-3s] %6.2f%%  %s\n" "$status" "$pct" "$pkg"
    done
    popd >/dev/null
done

echo ""
echo "Audited $total_pkgs packages. Below ${MIN_COVERAGE}%: $failures."
[ "$failures" -eq 0 ]
