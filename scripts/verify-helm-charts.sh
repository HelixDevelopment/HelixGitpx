#!/usr/bin/env bash
# Lightweight Helm chart sanity check (no `helm` CLI required).
# Asserts: Chart.yaml has name+version, values.yaml parses, templates/
# reference only values defined in values.yaml.
set -uo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

fail=0
total=0

lint_chart() {
    local dir="$1"
    total=$((total+1))
    local out
    out=$(python3 - <<PY 2>&1
import sys, os, re, yaml, glob

base = "$dir"
try:
    chart = yaml.safe_load(open(os.path.join(base, "Chart.yaml")))
except Exception as e:
    print(f"Chart.yaml invalid: {e}")
    sys.exit(1)
for k in ("name", "version"):
    if not chart.get(k):
        print(f"Chart.yaml missing field: {k}")
        sys.exit(1)

values_path = os.path.join(base, "values.yaml")
values = {}
if os.path.exists(values_path):
    try:
        values = yaml.safe_load(open(values_path)) or {}
    except Exception as e:
        print(f"values.yaml invalid: {e}")
        sys.exit(1)

# Collect top-level keys in values.yaml (single level is enough for this lint).
top = set(values.keys()) if isinstance(values, dict) else set()
pat = re.compile(r'{{\s*\.Values\.([A-Za-z0-9_]+)')
errors = []
for tmpl in glob.glob(os.path.join(base, "templates", "*.yaml")):
    text = open(tmpl).read()
    for m in pat.finditer(text):
        k = m.group(1)
        if k not in top:
            errors.append(f"  {os.path.relpath(tmpl, base)}: references .Values.{k} not in values.yaml")
if errors:
    print("\n".join(errors))
    sys.exit(1)
sys.exit(0)
PY
)
    rc=$?
    name=$(python3 -c "import yaml; print(yaml.safe_load(open('$dir/Chart.yaml'))['name'])" 2>/dev/null || basename "$dir")
    if [ $rc -eq 0 ]; then
        printf '  [ok]   %s\n' "$name"
    else
        printf '  [FAIL] %s\n' "$name"
        printf '%s\n' "$out" | sed 's/^/         /'
        fail=$((fail+1))
    fi
}

# Platform charts.
for dir in impl/helixgitpx-platform/helm/*/; do
    [ -f "$dir/Chart.yaml" ] && lint_chart "$dir"
done
# Service charts.
for dir in impl/helixgitpx/services/*/deploy/helm/; do
    [ -f "$dir/Chart.yaml" ] && lint_chart "$dir"
done
# Website + docs site.
for dir in impl/helixgitpx-website/deploy/helm/ impl/helixgitpx-docs-site/deploy/helm/; do
    [ -f "$dir/Chart.yaml" ] && lint_chart "$dir"
done

echo ""
echo "Checked $total charts. $fail failing."
[ "$fail" -eq 0 ]
