#!/usr/bin/env bash
# Lightweight Rego syntax check (no `opa` CLI required).
# Catches the most common mistakes: missing `package`, mismatched braces,
# unterminated string literals.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

fail=0
total=0

check() {
    local file="$1"
    total=$((total+1))
    python3 - "$file" <<'PY'
import sys, re
p = sys.argv[1]
text = open(p).read()

errs = []

# Rule 1: first non-comment, non-blank line must be "package X".
for line in text.splitlines():
    s = line.strip()
    if not s or s.startswith('#'):
        continue
    if not s.startswith('package '):
        errs.append('first directive must be `package X` (got: ' + s[:40] + ')')
    break

# Rule 2: balanced braces / brackets / parens.
for ch, close in (('{','}'), ('[',']'), ('(',')')):
    if text.count(ch) != text.count(close):
        errs.append(f'unbalanced {ch}{close}: {text.count(ch)} vs {text.count(close)}')

# Rule 3: quoted strings are even-numbered (simple heuristic; ignores escapes).
# Strip strings first and then confirm no stray quotes remain.
stripped = re.sub(r'"(?:[^"\\]|\\.)*"', '', text)
if stripped.count('"') > 0:
    errs.append('unterminated string literal')

if errs:
    for e in errs:
        print(f'  {p}: {e}')
    sys.exit(1)
PY
    rc=$?
    name=$(basename "$file")
    if [ $rc -eq 0 ]; then
        printf '  [ok]   %s\n' "$name"
    else
        printf '  [FAIL] %s\n' "$name"
        fail=$((fail+1))
    fi
}

echo "Rego syntax check:"
for rego in $(find impl/helixgitpx-platform/opa -name '*.rego' -type f 2>/dev/null); do
    check "$rego"
done

echo ""
echo "Checked $total Rego files. $fail failing."
[ "$fail" -eq 0 ]
