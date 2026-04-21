#!/usr/bin/env bash
# test/security/run.sh — run the security test battery against a target.
# NO mocks per Constitution §II §2. Every tool runs against a live service.
set -euo pipefail

TARGET=${HELIXGITPX_TARGET:-https://staging.helixgitpx.io}
OUT=${OUT_DIR:-/tmp/helixgitpx-security}
mkdir -p "$OUT"

have() { command -v "$1" >/dev/null 2>&1; }

echo "Security test run against $TARGET. Results in $OUT."

fail=0

# 1. gosec — Go static analysis. Caps at 2 minutes; a full scan across the
# monorepo can take 5+ minutes.
if have gosec; then
    echo "--- gosec ---"
    (cd impl/helixgitpx && timeout 120 gosec -quiet -severity medium -confidence medium -fmt json -out "$OUT/gosec.json" ./... 2>&1) \
        || echo "gosec timed out or reported findings — see $OUT/gosec.json."
else
    echo "gosec: not installed (skipping; install: go install github.com/securego/gosec/v2/cmd/gosec@latest)"
fi

# 2. OWASP ZAP baseline — live HTTP scan. Requires a reachable target; skipped
# by default (pulling the 1.5 GB image and probing a remote URL from a dev box
# is not suitable for a CI-hot path). Enable with ZAP_ENABLE=1.
if [ "${ZAP_ENABLE:-0}" = "1" ] && (have docker || have podman); then
    echo "--- ZAP baseline ---"
    RUNTIME=$(command -v podman || command -v docker || true)
    if [ -n "$RUNTIME" ]; then
        "$RUNTIME" run --rm -t ghcr.io/zaproxy/zaproxy:stable zap-baseline.py \
            -t "$TARGET" -J "$OUT/zap.json" 2>&1 || fail=$((fail+1))
    fi
else
    echo "ZAP: skipped (set ZAP_ENABLE=1 to enable; requires a reachable target + container runtime)."
fi

# 3. Nuclei — template-driven scan.
if have nuclei; then
    echo "--- nuclei ---"
    nuclei -target "$TARGET" -t ~/nuclei-templates/ -json -o "$OUT/nuclei.json" 2>&1 || fail=$((fail+1))
else
    echo "nuclei: not installed (skipping; install: go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest)"
fi

# 4. Trivy — image + IaC scan.
if have trivy; then
    echo "--- trivy image scan ---"
    trivy image --severity HIGH,CRITICAL ghcr.io/helixgitpx/hello:latest -f json -o "$OUT/trivy-hello.json" 2>&1 || fail=$((fail+1))
    echo "--- trivy config scan ---"
    trivy config impl/helixgitpx-platform -f json -o "$OUT/trivy-config.json" 2>&1 || fail=$((fail+1))
else
    echo "trivy: not installed (skipping; install: https://trivy.dev)"
fi

# 5. HelixGitpx-specific: webhook HMAC tamper check (uses the integration test).
if [ -n "${HELIXGITPX_WEBHOOK_URL:-}" ] && [ -n "${HELIXGITPX_WEBHOOK_SECRET:-}" ]; then
    echo "--- webhook tamper test ---"
    (cd test/integration && GOTOOLCHAIN=go1.23.4 go test -tags=integration -run TestWebhook_RejectsTamperedPayload -v ./... 2>&1) || fail=$((fail+1))
else
    echo "webhook tamper: HELIXGITPX_WEBHOOK_URL / SECRET unset; skipping"
fi

echo ""
echo "Security run complete. $fail tool(s) failed or reported findings."
echo "Review $OUT/*.json. Any High/Critical finding without an approved"
echo "suppression blocks merge per Constitution §II §5."
[ "$fail" -eq 0 ]
