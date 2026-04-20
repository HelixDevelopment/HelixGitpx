#!/usr/bin/env bash
# M2 preflight: refuse to proceed without the resources we need.
set -euo pipefail

MIN_FREE_MEM_GIB=48
FREE_KIB=$(awk '/MemAvailable/ { print $2 }' /proc/meminfo)
FREE_GIB=$((FREE_KIB / 1024 / 1024))

if [ "$FREE_GIB" -lt "$MIN_FREE_MEM_GIB" ]; then
    printf 'preflight: need >= %d GiB free, have %d GiB. Close apps and retry.\n' \
        "$MIN_FREE_MEM_GIB" "$FREE_GIB" >&2
    exit 2
fi

for tool in k3d kubectl helm; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        printf 'preflight: %s not on PATH. Run: mise install\n' "$tool" >&2
        exit 3
    fi
done

# Container runtime — must be docker OR podman (M1 ADR-0002)
if ! (command -v docker >/dev/null 2>&1 || command -v podman >/dev/null 2>&1); then
    printf 'preflight: need docker or podman on PATH\n' >&2
    exit 4
fi

printf 'preflight: ok (%d GiB free, tools present)\n' "$FREE_GIB"
