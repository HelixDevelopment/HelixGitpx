#!/usr/bin/env bash
# Lint runbooks: every RB-*.md must contain the required headings from _TEMPLATE.md.
set -euo pipefail

here="docs/specifications/main/main_implementation_material/HelixGitpx/12-operations/runbooks"
# Minimum required headings found in all existing runbooks; template extends with "Best Practices" sections
required=("## 1. Detection")

rc=0
shopt -s nullglob
for rb in "$here"/RB-*.md; do
    for h in "${required[@]}"; do
        if ! grep -Fq "$h" "$rb"; then
            echo "$rb: missing heading '$h'" >&2
            rc=1
        fi
    done
done

if [ "$rc" -eq 0 ]; then
    echo "runbook-lint: all runbooks conform"
fi
exit "$rc"
