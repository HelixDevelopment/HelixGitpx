#!/usr/bin/env bash
# Walk the M3 completion matrix (15 roadmap items 39-53). Exit 0 iff all pass.
set -u
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

pass=0; fail=0
report() {
    local n="$1" r="$2"
    if [ "$r" = ok ]; then
        printf '  [ ok ] %s\n' "$n"
        pass=$((pass+1))
    else
        printf '  [FAIL] %s\n' "$n"
        fail=$((fail+1))
    fi
}
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M3 Identity & Orgs — Completion Matrix =="

check "39 Keycloak chart + realm"  bash -c 'test -f impl/helixgitpx-platform/helm/keycloak/Chart.yaml && test -f impl/helixgitpx-platform/helm/keycloak/realm/helixgitpx.json'
check "40 RS256 JWT signer"        bash -c 'grep -q "SigningMethodRS256" impl/helixgitpx/platform/auth/jwt.go'
check "41 hpxat_ PAT prefix"       bash -c 'grep -q "patPrefix.*hpxat_" impl/helixgitpx/services/auth/internal/domain/pat.go'
check "42 TOTP enroll"             bash -c 'grep -q "EnrollTOTP" impl/helixgitpx/services/auth/internal/domain/mfa.go'
check "43 auth.sessions table"     bash -c 'grep -q "auth.sessions" impl/helixgitpx/services/auth/migrations/20260420000003_auth.sql'
check "44 OrgService proto"        bash -c 'grep -q "service OrgService" impl/helixgitpx/proto/helixgitpx/org/v1/org.proto'
check "45 Nested teams + cycle"    bash -c 'grep -q "DetectCycle" impl/helixgitpx/services/orgteam/internal/domain/team.go'
check "46 memberships + role_enum" bash -c 'grep -q "role_enum" impl/helixgitpx/services/orgteam/migrations/20260420000004_orgteam.sql'
check "47 OPA bundle v1"           test -f impl/helixgitpx-platform/opa/bundles/v1/authz.rego
check "48 audit.events Kafka consumer" bash -c 'grep -q "audit.events" impl/helixgitpx/services/audit/internal/app/app.go'
check "49 append-only rules"       bash -c 'grep -q "no_update\\|DO INSTEAD NOTHING" impl/helixgitpx/services/audit/migrations/20260420000005_audit.sql'
check "50 Merkle anchoring"        test -f impl/helixgitpx/services/audit/cmd/audit-merkle/main.go
check "51 Angular auth flow"       test -f impl/helixgitpx-web/apps/web/src/app/login/login.component.ts
check "52 Connect-Web clients"     test -f impl/helixgitpx-web/libs/proto/src/helixgitpx/auth/v1/auth_connect.ts
check "53 OTel-web wired"          bash -c 'grep -q "@opentelemetry/sdk-trace-web" impl/helixgitpx-web/apps/web/src/app/app.config.ts'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
