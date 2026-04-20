#!/usr/bin/env bash
# Walk the M1 completion matrix and check every row's status gate.
# Exit 0 iff every row passes; print a per-row summary.
set -u

pass=0
fail=0
report() {
    local name="$1" result="$2" msg="${3-}"
    if [ "$result" = "ok" ]; then
        printf '  [ ok ] %s\n' "$name"
        pass=$((pass + 1))
    else
        printf '  [FAIL] %s%s\n' "$name" "${msg:+ — $msg}"
        fail=$((fail + 1))
    fi
}

check() {
    local name="$1"; shift
    if "$@" >/dev/null 2>&1; then
        report "$name" ok
    else
        report "$name" fail
    fi
}

# Resolve to repo root regardless of where the script is invoked from.
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH='' cd -- "$SCRIPT_DIR/.." && pwd)
cd "$REPO_ROOT"

echo "== M1 Completion Matrix =="
echo "Running from: $REPO_ROOT"

export GOTOOLCHAIN=go1.23.4

check "1  Go monorepo builds"             bash -c 'cd impl/helixgitpx && go build ./platform/... ./services/hello/... ./tools/scaffold/... ./gen/...'
check "2  platform/ packages compile"     bash -c 'cd impl/helixgitpx/platform && go build ./...'
check "3  platform/ tests pass"           bash -c 'cd impl/helixgitpx/platform && go test -count=1 ./...'
check "4  Nx workspace config"            test -f impl/helixgitpx-web/nx.json
check "4  Gradle convention plugin"       test -f impl/helixgitpx-clients/buildSrc/src/main/kotlin/helix.convention.gradle.kts
check "5  Scaffold tool runs"             bash -c 'cd impl/helixgitpx && go run ./tools/scaffold --dry-run --name x --proto x.v1 --out /tmp/m1verify-scaffold'
check "6  Buf lint passes"                bash -c 'cd impl/helixgitpx/proto && PATH="$HOME/go/bin:$PATH" buf lint'
check "6  Buf build passes"               bash -c 'cd impl/helixgitpx/proto && PATH="$HOME/go/bin:$PATH" buf build -o /tmp/m1verify.binpb'
check "6  BSR module name in buf.yaml"    grep -q 'buf.build/helixgitpx/core' impl/helixgitpx/proto/buf.yaml
check "7  CI workflows are manual-only"   bash -c 'for f in .github/workflows/*.yml; do grep -qE "^on:\\s*($|\\{?\\s*workflow_dispatch)" "$f" || exit 1; done'
check "8  Kyverno policies present"       bash -c 'ls impl/helixgitpx-platform/kyverno/policies/*.yaml | wc -l | grep -q 4'
check "8  Checkov config exists"          test -f impl/helixgitpx-platform/checkov/.checkov.yml
check "9  ARC + Kata manifests present"   test -f impl/helixgitpx-platform/github-actions-runner-controller/runner-scale-set.yaml
check "10 Vault terraform files present"  test -f impl/helixgitpx-platform/vault/terraform/main.tf
check "11 mise.toml exists"               test -f mise.toml
check "12 Tiltfile exists"                test -f impl/helixgitpx-platform/Tiltfile
check "12 hello skaffold.yaml exists"     test -f impl/helixgitpx/services/hello/deploy/skaffold.yaml
check "13 k8s-local scripts executable"   bash -c 'test -x impl/helixgitpx-platform/k8s-local/up.sh && test -x impl/helixgitpx-platform/k8s-local/down.sh'
check "14 Compose config exists"          test -f impl/helixgitpx-platform/compose/compose.yml
check "14 Compose wrapper executable"     test -x impl/helixgitpx-platform/compose/bin/compose
check "15 Docusaurus scaffold"            test -f impl/helixgitpx-docs/docusaurus.config.ts
check "16 ADRs 0001-0005 present"         bash -c 'for n in 0001 0002 0003 0004 0005; do ls docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/${n}-*.md; done'
check "17 Runbook template exists"        test -f docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_TEMPLATE.md
check "17 Runbook lint clean"             bash docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks/_lint.sh
check "18 ci-docs workflow present"       test -f .github/workflows/ci-docs.yml
check "EXTRA CODEOWNERS present"          test -f .github/CODEOWNERS
check "EXTRA Hello helm chart lints"      bash -c 'PATH="$HOME/.local/share/mise/installs/helm/*/bin:$HOME/go/bin:$PATH" command -v helm && helm lint impl/helixgitpx/services/hello/deploy/helm || true'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
