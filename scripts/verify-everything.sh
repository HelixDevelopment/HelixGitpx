#!/usr/bin/env bash
# verify-everything.sh — one-shot green-suite for HelixGitpx.
# Runs every artifact-level verifier + `go vet` + `go test` + chart/rego lints.
# Does NOT run the live-cluster probes (those short-circuit cleanly on their own).
set -uo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

step() {
    local name="$1"; shift
    printf '── %s ──\n' "$name"
    if "$@"; then
        printf '   [pass]\n'
        return 0
    else
        printf '   [fail]\n'
        return 1
    fi
}

overall=0

step "M1 artifacts"        bash scripts/verify-m1-artifacts.sh >/dev/null 2>&1 || overall=$((overall+1))
step "M2 artifacts"        bash scripts/verify-m2-artifacts.sh >/dev/null 2>&1 || overall=$((overall+1))
for m in 3 4 5 6 7 8; do
    step "M$m cluster probe" bash "scripts/verify-m$m-cluster.sh" >/dev/null 2>&1 || overall=$((overall+1))
done
step "Argo paths"          bash scripts/verify-argo-paths.sh  >/dev/null 2>&1 || overall=$((overall+1))
step "Helm charts"         bash scripts/verify-helm-charts.sh >/dev/null 2>&1 || overall=$((overall+1))
step "Rego syntax"         bash scripts/verify-rego.sh        >/dev/null 2>&1 || overall=$((overall+1))

step "Go vet (platform)" bash -c '
for mod in $(find impl/helixgitpx -name go.mod | xargs -n1 dirname); do
    (cd "$mod" && GOTOOLCHAIN=go1.23.4 go vet ./... 2>&1 | grep -v "^$") | sed "s|^|  $mod: |"
done | (! grep -q .)
' || overall=$((overall+1))

step "Go test (platform + services)" bash -c '
fail=0
for mod in $(find impl/helixgitpx -name go.mod | xargs -n1 dirname); do
    if ! (cd "$mod" && GOTOOLCHAIN=go1.23.4 go test ./... >/dev/null 2>&1); then
        fail=$((fail+1))
    fi
done
exit $fail
' || overall=$((overall+1))

echo ""
if [ "$overall" -eq 0 ]; then
    echo "All phases passed."
else
    echo "$overall phase(s) failed."
fi
exit $overall
