#!/usr/bin/env bash
# test/chaos/run.sh — apply every Litmus experiment and watch for pass/fail.
# NO mocks — runs against a real cluster with ChaosEngine CRDs installed.
set -euo pipefail

OUT=${OUT_DIR:-/tmp/helixgitpx-chaos}
mkdir -p "$OUT"

if ! command -v kubectl >/dev/null 2>&1; then
    echo "kubectl not installed — skipping chaos suite."
    exit 0
fi
if ! kubectl cluster-info >/dev/null 2>&1; then
    echo "No cluster reachable — skipping."
    exit 0
fi

fail=0
for expr in tools/chaos/*.yaml; do
    name=$(basename "$expr" .yaml)
    echo "--- chaos: $name ---"
    kubectl apply -f "$expr"
    # Wait for experiment to reach completed state (Litmus reports via CRD).
    ns=$(grep -E '^\s*namespace:' "$expr" | head -1 | awk '{print $2}')
    kubectl -n "$ns" wait --for=condition=Completed chaosengine/"$name" --timeout=10m || fail=$((fail+1))
    kubectl -n "$ns" describe chaosengine "$name" > "$OUT/$name-describe.txt"
done

echo ""
echo "Chaos suite done. $fail experiment(s) did not complete. See $OUT/."
[ "$fail" -eq 0 ]
