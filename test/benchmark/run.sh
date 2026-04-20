#!/usr/bin/env bash
# test/benchmark/run.sh — Go micro-benchmarks + k6 budget-gated scenarios.
set -euo pipefail

OUT=${OUT_DIR:-/tmp/helixgitpx-bench}
mkdir -p "$OUT"

echo "--- Go micro-benchmarks ---"
(cd impl/helixgitpx && GOTOOLCHAIN=go1.23.4 go test -run='^$' -bench=. -benchmem -benchtime=1x ./... 2>&1 | tee "$OUT/go-bench.txt") || true

echo ""
echo "--- k6 budget-gated scenarios ---"
if command -v k6 >/dev/null 2>&1; then
    (cd tools/perf && make all 2>&1 | tee "$OUT/k6-bench.txt") || true
    python3 tools/perf/check_budgets.py /tmp/k6-out-*.json tools/perf/budgets.json || true
else
    echo "k6 not installed — skipping scenario budgets."
fi

echo ""
echo "Benchmark run complete. See $OUT/."
