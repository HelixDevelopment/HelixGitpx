#!/usr/bin/env bash
# Pre-pull heavy images before k3d cluster create to avoid first-boot timeouts.
set -euo pipefail

RUNTIME=docker
if ! command -v docker >/dev/null 2>&1; then
    RUNTIME=podman
fi

IMAGES=(
    ghcr.io/cloudnative-pg/postgresql:16.4-3
    quay.io/strimzi/kafka:0.44.0-kafka-3.8.0
    opensearchproject/opensearch:2.18.0
    hashicorp/vault:1.18.3
    grafana/mimir:2.14.0
    grafana/loki:3.3.0
    grafana/tempo:2.6.1
    grafana/pyroscope:1.10.0
    grafana/grafana:11.4.0
    prom/prometheus:v3.0.1
    quay.io/cilium/cilium:v1.16.3
    registry.k8s.io/ingress-nginx/controller:v1.12.0
    quay.io/jetstack/cert-manager-controller:v1.16.2
    quay.io/minio/minio:RELEASE.2025-01-20T14-49-07Z
    ghcr.io/spiffe/spire-server:1.11.0
    docker.io/istio/pilot:1.24.2
    docker.io/istio/ztunnel:1.24.2
    docker.redpanda.com/redpandadata/connect:4.40.0
    ghcr.io/aiven-open/karapace:latest
    docker.dragonflydb.io/dragonflydb/dragonfly:v1.25.1
    getmeili/meilisearch:v1.12
    qdrant/qdrant:v1.12.4
    quay.io/argoproj/argocd:v3.0.0
    debezium/connect:3.0.0.Final
)

printf 'Pre-pulling %d images using %s...\n' "${#IMAGES[@]}" "$RUNTIME"

for img in "${IMAGES[@]}"; do
    "$RUNTIME" pull "$img" &
done
wait
printf 'prepull: done\n'
