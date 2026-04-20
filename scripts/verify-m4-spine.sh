#!/usr/bin/env bash
# M4 end-to-end spine — exercises the user journey against a running cluster.
# Most gates FAIL without cluster bring-up; that's expected.
set -u
pass=0; fail=0
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then printf '  [ ok ] %s\n' "$n"; pass=$((pass+1)); else printf '  [FAIL] %s\n' "$n"; fail=$((fail+1)); fi; }

echo "== M4 End-to-end Spine =="

check "repo-service /healthz"      curl -fsSLk https://repo.helix.local/healthz
check "git-ingress /healthz"       curl -fsSLk https://git.helix.local/healthz
check "adapter-pool /healthz"      curl -fsSLk https://adapter-pool.helix.local/healthz
check "webhook-gateway /healthz"   curl -fsSLk https://webhook.helix.local/healthz
check "upstream-service /healthz"  curl -fsSLk https://upstream.helix.local/healthz
check "git push → upstream mirror" kubectl -n helix-data get kafkaconnector hello-outbox -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
