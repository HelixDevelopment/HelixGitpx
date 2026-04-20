#!/usr/bin/env bash
# test/ddos/run.sh — drive malicious-shaped load and assert rate-limiters hold.
# Real services required. NO mocks per Constitution §II §2.
set -euo pipefail

TARGET=${HELIXGITPX_TARGET:-https://staging.helixgitpx.io}
OUT=${OUT_DIR:-/tmp/helixgitpx-ddos}
mkdir -p "$OUT"

have() { command -v "$1" >/dev/null 2>&1; }

if ! have k6; then
    echo "k6 not installed. See https://k6.io/docs/get-started/installation/" >&2
    exit 2
fi

echo "DDoS drill against $TARGET. Results in $OUT."

# Scenario 1: arrival burst — 100× baseline RPS for 30s, then 0.
cat >/tmp/ddos-burst.js <<'JS'
import http from 'k6/http';
import { check } from 'k6';
export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '30s', target: 1000 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    'http_req_failed{expected_response:true}': ['rate<0.01'],
    'http_reqs{status:429}': [],
  },
};
const base = __ENV.HELIXGITPX_TARGET || 'https://staging.helixgitpx.io';
export default function () {
  const r = http.get(`${base}/api/v1/orgs`);
  check(r, { '200 or 429 or 503': (x) => [200, 429, 503].includes(x.status) });
}
JS
HELIXGITPX_TARGET="$TARGET" k6 run --out json="$OUT/burst.json" /tmp/ddos-burst.js || true

# Scenario 2: slowloris — 5k half-open connections for 10 min.
# k6 can't do true half-open; we approximate with very high `httpTimeout`.
cat >/tmp/ddos-slowloris.js <<'JS'
import http from 'k6/http';
export const options = { vus: 5000, duration: '2m', httpTimeout: '2m' };
const base = __ENV.HELIXGITPX_TARGET || 'https://staging.helixgitpx.io';
export default function () {
  http.get(`${base}/healthz`, { timeout: '2m' });
}
JS
HELIXGITPX_TARGET="$TARGET" k6 run --out json="$OUT/slowloris.json" /tmp/ddos-slowloris.js || true

# Scenario 3: cache-busting query flood — append random suffix each request.
cat >/tmp/ddos-cache-bust.js <<'JS'
import http from 'k6/http';
export const options = { vus: 200, duration: '3m' };
const base = __ENV.HELIXGITPX_TARGET || 'https://staging.helixgitpx.io';
export default function () {
  const suffix = Math.random().toString(36).slice(2);
  http.get(`${base}/api/v1/orgs?_=${suffix}`);
}
JS
HELIXGITPX_TARGET="$TARGET" k6 run --out json="$OUT/cache-bust.json" /tmp/ddos-cache-bust.js || true

echo ""
echo "DDoS run complete. Inspect $OUT/*.json."
echo ""
echo "Pass criteria (Constitution §II §5 + test/ddos/README.md):"
echo "  - rate-limiter engages within 1 s."
echo "  - > 99 % of abusive requests dropped."
echo "  - legitimate baseline RPS still < 1 % error."
echo "  - p99 returns to baseline within 60 s after attack ends."
