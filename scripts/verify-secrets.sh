#!/usr/bin/env bash
# verify-secrets.sh — scan the repo for committed secrets with gitleaks.
# Part of the defence-in-depth against accidental key / token commits.
# Degrades gracefully when gitleaks is absent.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if ! command -v gitleaks >/dev/null 2>&1; then
    # Try Go's install-path.
    if [ -x "$HOME/go/bin/gitleaks" ]; then
        GITLEAKS="$HOME/go/bin/gitleaks"
    else
        echo "gitleaks not installed — skipping."
        echo "  install: go install github.com/gitleaks/gitleaks/v8@latest"
        exit 0
    fi
else
    GITLEAKS=gitleaks
fi

echo "Scanning tracked files for committed secrets..."
"$GITLEAKS" detect --no-banner --exit-code 1 --source . --redact \
    --report-format json --report-path /tmp/gitleaks-report.json \
    || {
        rc=$?
        echo ""
        echo "gitleaks found potential secrets (exit=$rc). Review /tmp/gitleaks-report.json."
        exit $rc
    }

echo "No secrets found."
