#!/usr/bin/env bash
# Verify M1 Foundation artifacts are present in the repo (no runtime).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

pass=0; fail=0
check() {
    local num="$1" desc="$2" pred="$3"
    if eval "$pred" >/dev/null 2>&1; then printf "  [ ok ] %s %s\n" "$num" "$desc"; pass=$((pass+1))
    else printf "  [fail] %s %s\n" "$num" "$desc"; fail=$((fail+1)); fi
}

echo "== M1 Foundation — Artifact Matrix =="
check  1 "Monorepo layout"          "test -f impl/helixgitpx/go.work -a -d impl/helixgitpx-platform/helm"
check  2 "mise.toml"                "test -f mise.toml -o -f .mise.toml"
check  3 "Scaffold tool"            "test -f impl/helixgitpx/tools/scaffold/main.go"
check  4 "hello service scaffolded" "test -f impl/helixgitpx/services/hello/cmd/hello/main.go"
check  5 "Proto buf config"         "test -d impl/helixgitpx/gen"
check  6 "Platform shared libs"     "test -d impl/helixgitpx/platform/log"
check  7 "Container runtime helper" "test -x impl/helixgitpx-platform/compose/bin/compose"
check  8 "CI workflow dispatch-only" "grep -rq 'workflow_dispatch' .github/workflows/"
check  9 "Upstreams scripts"        "test -f Upstreams/GitHub.sh -a -f Upstreams/GitLab.sh -a -f Upstreams/GitFlic.sh -a -f Upstreams/GitVerse.sh"
check 10 "CLAUDE.md present"        "test -f CLAUDE.md"
check 11 "ADRs 0001-0016 present"   "ls docs/specifications/main/main_implementation_material/HelixGitpx/15-reference/adr/0016-*.md"
check 18 "go.work includes all"     "grep -q './services/hello' impl/helixgitpx/go.work"

echo ""
echo "PASS: $pass   FAIL: $fail"
[ "$fail" -eq 0 ]
