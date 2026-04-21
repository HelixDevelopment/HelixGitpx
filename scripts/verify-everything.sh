#!/usr/bin/env bash
# verify-everything.sh — one-shot green-suite for HelixGitpx.
# Runs every artifact-level verifier + `go vet` + `go test` + chart/rego lints
# + fuzz/benchmark smokes + proto-gen drift. Cluster probes short-circuit
# cleanly when no cluster is reachable.
set -uo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

overall=0

step() {
    local name="$1"; shift
    printf '── %-45s ' "$name"
    if "$@" >/dev/null 2>&1; then
        printf '[pass]\n'
        return 0
    else
        printf '[fail]\n'
        overall=$((overall+1))
        return 1
    fi
}

step "M1 artifacts"        bash scripts/verify-m1-artifacts.sh
step "M2 artifacts"        bash scripts/verify-m2-artifacts.sh
for m in 3 4 5 6 7 8; do
    step "M$m cluster probe" bash "scripts/verify-m$m-cluster.sh"
done
step "Argo paths"          bash scripts/verify-argo-paths.sh
step "Helm charts"         bash scripts/verify-helm-charts.sh
step "Rego syntax"         bash scripts/verify-rego.sh

step "Go vet (all modules)" bash -c '
set -e
for mod in $(find impl/helixgitpx -name go.mod | xargs -n1 dirname); do
    (cd "$mod" && GOTOOLCHAIN=go1.23.4 go vet ./...) || exit 1
done
'

step "Go test (all modules)" bash -c '
set -e
for mod in $(find impl/helixgitpx -name go.mod | xargs -n1 dirname); do
    (cd "$mod" && GOTOOLCHAIN=go1.23.4 go test ./...) || exit 1
done
'

step "Fuzz smoke (2s × 2)" bash -c '
(cd impl/helixgitpx/platform/webhook && GOTOOLCHAIN=go1.23.4 go test -run=^$ -fuzz=FuzzVerifyHMAC -fuzztime=2s) && \
(cd impl/helixgitpx/services/webhook-gateway/internal/canonical && GOTOOLCHAIN=go1.23.4 go test -run=^$ -fuzz=FuzzCanonicalizeGitHub -fuzztime=2s)
'

step "Benchmark smoke (3x)" bash -c '
GOTOOLCHAIN=go1.23.4
(cd impl/helixgitpx/platform/webhook && go test -run=^$ -bench=. -benchtime=3x) && \
(cd impl/helixgitpx/services/search-service/internal/domain && go test -run=^$ -bench=. -benchtime=3x) && \
(cd impl/helixgitpx/services/audit/internal/merkle && go test -run=^$ -bench=. -benchtime=3x)
'

step "Proto gen in sync"   bash scripts/verify-proto-gen.sh

echo ""
if [ "$overall" -eq 0 ]; then
    echo "All phases passed."
else
    echo "$overall phase(s) failed."
fi
exit $overall
