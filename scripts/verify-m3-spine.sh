#!/usr/bin/env bash
# M3 end-to-end spine — exercises the user journey against a running cluster.
# Most gates FAIL without an M2+M3 cluster brought up; that's expected.
set -u
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

if ! command -v kubectl >/dev/null 2>&1 || ! kubectl cluster-info >/dev/null 2>&1; then
    echo "== $(basename "$0" .sh | sed 's/verify-//;s/-spine//') spine — SKIP (no cluster reachable) =="
    exit 0
fi


pass=0; fail=0
check() {
    local n="$1"; shift
    if "$@" >/dev/null 2>&1; then
        printf '  [ ok ] %s\n' "$n"; pass=$((pass+1))
    else
        printf '  [FAIL] %s\n' "$n"; fail=$((fail+1))
    fi
}

echo "== M3 End-to-end Spine =="

check "Keycloak OIDC discovery"   curl -fsSLk https://keycloak.helix.local/realms/helixgitpx/.well-known/openid-configuration
check "auth-service /healthz"     curl -fsSLk https://auth.helix.local/healthz
check "orgteam-service /healthz"  curl -fsSLk https://orgteam.helix.local/healthz
check "audit-service deployment"  kubectl -n helix get deploy audit
check "Grafana /api/health"       curl -fsSLk https://grafana.helix.local/api/health

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
