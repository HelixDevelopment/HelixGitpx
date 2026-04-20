#!/usr/bin/env bash
# Verify every Argo CD Application's source.path resolves to an existing dir.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

fail=0
total=0

for app in impl/helixgitpx-platform/argocd/applications/*.yaml; do
    total=$((total+1))
    # Pure Python YAML reader — we do NOT depend on helm/kustomize.
    path=$(python3 -c "
import sys, yaml
with open('$app') as f:
    data = yaml.safe_load(f)
src = data.get('spec', {}).get('source', {})
print(src.get('path', ''))
")
    if [ -z "$path" ]; then
        printf '  [skip] %s (no source.path; possibly multi-source)\n' "$(basename "$app")"
        continue
    fi
    if [ -d "$path" ]; then
        printf '  [ok]   %-50s  %s\n' "$(basename "$app" .yaml)" "$path"
    else
        printf '  [FAIL] %-50s  path missing: %s\n' "$(basename "$app" .yaml)" "$path"
        fail=$((fail+1))
    fi
done

echo ""
echo "Checked $total Argo Applications. $fail broken paths."
[ "$fail" -eq 0 ]
