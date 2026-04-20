#!/usr/bin/env bash
# test/stress/run.sh — drive sustained load at 3× design capacity for 30 min.
set -euo pipefail

TARGET=${HELIXGITPX_TARGET:-https://staging.helixgitpx.io}
OUT=${OUT_DIR:-/tmp/helixgitpx-stress}
mkdir -p "$OUT"

if ! command -v k6 >/dev/null 2>&1; then
    echo "k6 not installed. See https://k6.io/docs/get-started/installation/" >&2
    exit 2
fi

for script in tools/perf/scenarios/*.js; do
    name=$(basename "$script" .js)
    echo "--- stress: $name ---"
    K6_TARGET_ENV="$TARGET" HELIXGITPX_TARGET="$TARGET" \
        k6 run --stage 2m:500,20m:1500,2m:0 --out json="$OUT/${name}-stress.json" "$script" || true
done

echo ""
echo "Stress run complete. See $OUT/*-stress.json."
echo "Pass criteria: p99 ≤ 2× baseline; error rate ≤ SLO; no OOM;"
echo "recovery ≤ 5 min after load removed."
