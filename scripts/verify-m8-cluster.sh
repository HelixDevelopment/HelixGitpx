#!/usr/bin/env bash
# Verify M8 — Scale, Harden, GA (items 138-161).
# Exits 0 if all artifact checks pass.
set -uo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

pass=0
fail=0
check() {
    local num="$1"; local desc="$2"; local pred="$3"
    if eval "$pred" >/dev/null 2>&1; then
        printf "  [ ok ] %s %s\n" "$num" "$desc"
        pass=$((pass+1))
    else
        printf "  [fail] %s %s\n" "$num" "$desc"
        fail=$((fail+1))
    fi
}

echo "== M8 Scale, Harden, GA — Completion Matrix =="

# 8.1 Multi-region
check 138 "Second region overlay"       "test -f impl/helixgitpx-platform/kustomize/overlays/region-b/kustomization.yaml"
check 139 "MirrorMaker2 chart"           "test -f impl/helixgitpx-platform/helm/mirrormaker2/Chart.yaml"
check 139 "Postgres replica chart"       "test -f impl/helixgitpx-platform/helm/postgres-replica/Chart.yaml"
check 140 "GeoDNS chart"                 "test -f impl/helixgitpx-platform/helm/geodns/Chart.yaml"
check 141 "Data residency migration"     "test -f impl/helixgitpx-platform/sql/migrations/20260420_org_residency.sql"

# 8.2 Performance
check 142 "k6 scenarios (4 files)"        "test -f tools/perf/scenarios/api_baseline.js -a -f tools/perf/scenarios/git_push_pull.js -a -f tools/perf/scenarios/websocket_fanout.js -a -f tools/perf/scenarios/ai_chat.js"
check 143 "7d soak workflow"             "test -f tools/perf/soak-7d.yaml"
check 144 "Perf budgets CI"              "test -f .github/workflows/perf-budgets.yml"
check 145 "GPU HPA"                      "test -f tools/perf/gpu-hpa.yaml"

# 8.3 Chaos & DR
check 146 "Chaos experiments (6)"        "test -f tools/chaos/broker-kill.yaml -a -f tools/chaos/pod-evict.yaml -a -f tools/chaos/network-partition.yaml -a -f tools/chaos/disk-full.yaml -a -f tools/chaos/llm-outage.yaml -a -f tools/chaos/upstream-429-storm.yaml"
check 147 "DR drill runbook"             "test -f tools/dr/dr-drill-runbook.md"
check 148 "Operational runbooks"         "test -f docs/operations/runbooks/broker-kill.md -a -f docs/operations/runbooks/dr-failover.md"

# 8.4 Security
check 149 "Pen-test scope"               "test -f docs/security/pentest-scope-2026q2.md"
check 150 "Bug bounty program"           "test -f docs/security/bug-bounty-program.md"
check 151 "SOC 2 Type I index"           "test -f docs/security/soc2-type1-evidence-index.md"
check 152 "ISO 27001 gap analysis"       "test -f docs/security/iso27001-gap-analysis.md"

# 8.5 Coverage / mutation / fuzz / e2e
check 153 "Coverage audit script"        "test -x tools/coverage-audit/audit.sh"
check 154 "Mutation testing CI"          "test -f .github/workflows/mutation-testing.yml"
check 155 "Fuzz corpora"                 "test -f tools/fuzz/corpora/http/seed.txt -a -f tools/fuzz/corpora/webhook/github-push.json"
check 156 "E2E gaps audit"               "test -f tools/e2e-gaps.md"

# 8.6 GA
check 157 "Docs site"                    "test -f impl/helixgitpx-docs-site/docusaurus.config.ts"
check 158 "Trust center route"           "grep -q TrustCenterComponent impl/helixgitpx-web/apps/web/src/app/routes.ts"
check 159 "Billing service + schema"     "test -f impl/helixgitpx/services/billing-service/go.mod -a -f impl/helixgitpx-platform/sql/migrations/20260420_billing.sql"
check 160 "Launch checklist"             "test -f docs/marketing/launch-checklist.md"
check 161 "RELEASE.md present"           "test -f RELEASE.md"

# ADRs
check ADR "ADRs 0035-0039 present"       "test -f docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0035-two-region-active-passive.md -a -f docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0036-mirrormaker2-over-replicator.md -a -f docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0037-coredns-geoip-vs-managed-geodns.md -a -f docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0038-stripe-for-billing.md -a -f docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0039-docusaurus-for-docs.md"

echo ""
echo "PASS: $pass   FAIL: $fail"
[ "$fail" -eq 0 ]
