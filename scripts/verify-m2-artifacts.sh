#!/usr/bin/env bash
# Verify M2 Core Data Plane artifacts are present in the repo (no runtime).
# Runtime gates (cluster probes) live in verify-m2-cluster.sh.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

pass=0; fail=0
check() {
    local num="$1" desc="$2" pred="$3"
    if eval "$pred" >/dev/null 2>&1; then printf "  [ ok ] %s %s\n" "$num" "$desc"; pass=$((pass+1))
    else printf "  [fail] %s %s\n" "$num" "$desc"; fail=$((fail+1)); fi
}

echo "== M2 Core Data Plane — Artifact Matrix =="
check 19 "Cilium chart"          "test -d impl/helixgitpx-platform/helm/cilium"
check 20 "Argo CD bootstrap"     "test -d impl/helixgitpx-platform/argocd/bootstrap"
check 20 "App-of-apps"           "test -d impl/helixgitpx-platform/argocd/applicationset -o -d impl/helixgitpx-platform/argocd/applications"
check 21 "cert-manager chart"    "test -d impl/helixgitpx-platform/helm/cert-manager"
check 22 "SPIRE chart"           "test -d impl/helixgitpx-platform/helm/spire"
check 23 "CNPG chart"            "test -d impl/helixgitpx-platform/helm/cnpg-operator -a -d impl/helixgitpx-platform/helm/cnpg-cluster"
check 24 "SQL schemas"           "test -f impl/helixgitpx-platform/sql/schemas.sql"
check 25 "Goose migrations dir"  "test -d impl/helixgitpx-platform/sql/migrations"
check 26 "MinIO chart"           "test -d impl/helixgitpx-platform/helm/minio"
check 27 "Strimzi Kafka chart"   "test -d impl/helixgitpx-platform/helm/kafka-cluster"
check 28 "Karapace chart"        "test -d impl/helixgitpx-platform/helm/karapace"
check 29 "Debezium KafkaConnect" "test -d impl/helixgitpx-platform/helm/debezium"
check 30 "Outbox connector cfg"  "find impl/helixgitpx-platform -name '*outbox*' -print -quit | grep -q ."
check 31 "Dragonfly chart"       "test -d impl/helixgitpx-platform/helm/dragonfly"
check 32 "Meilisearch chart"     "test -d impl/helixgitpx-platform/helm/meilisearch"
check 33 "OpenSearch chart"      "test -d impl/helixgitpx-platform/helm/opensearch"
check 34 "Qdrant chart"          "test -d impl/helixgitpx-platform/helm/qdrant"
check 35 "Vault chart"           "test -d impl/helixgitpx-platform/helm/vault"
check 36 "Observability stack"   "test -d impl/helixgitpx-platform/helm/prometheus-stack -a -d impl/helixgitpx-platform/helm/loki -a -d impl/helixgitpx-platform/helm/tempo"
check 37 "Grafana dashboards"    "find impl/helixgitpx-platform/helm/prometheus-stack -name '*dashboard*' -print -quit | grep -q . || find impl/helixgitpx-platform -name 'dashboards' -type d -print -quit | grep -q ."
check 38 "Alertmanager config"   "find impl/helixgitpx-platform -iname 'alert*' -print -quit | grep -q ."

echo ""
echo "PASS: $pass   FAIL: $fail"
[ "$fail" -eq 0 ]
