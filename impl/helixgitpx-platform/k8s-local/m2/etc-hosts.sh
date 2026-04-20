#!/usr/bin/env bash
# Write /etc/hosts entries for M2 ingress hostnames.
# Requires sudo.
set -euo pipefail

HOSTS=(
    hello.helix.local
    grafana.helix.local
    argocd.helix.local
    vault.helix.local
    kafka.helix.local
    prometheus.helix.local
    minio.helix.local
)

TAG='# helixgitpx-m2'
HOSTS_FILE=${HOSTS_FILE:-/etc/hosts}

if ! grep -Fq "$TAG" "$HOSTS_FILE"; then
    {
        printf '\n%s\n' "$TAG"
        for h in "${HOSTS[@]}"; do
            printf '127.0.0.1\t%s\n' "$h"
        done
    } | sudo tee -a "$HOSTS_FILE" >/dev/null
    printf 'etc-hosts: wrote %d entries to %s\n' "${#HOSTS[@]}" "$HOSTS_FILE"
else
    printf 'etc-hosts: already present (tag %s found)\n' "$TAG"
fi
