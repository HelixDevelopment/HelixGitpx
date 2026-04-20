#!/usr/bin/env bash
# Walk the M2 completion matrix (20 roadmap items 19–38).
# Exit 0 iff every gate passes.
set -u

SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

pass=0
fail=0
report() {
    local name="$1" result="$2" msg="${3-}"
    if [ "$result" = "ok" ]; then
        printf '  [ ok ] %s\n' "$name"
        pass=$((pass + 1))
    else
        printf '  [FAIL] %s%s\n' "$name" "${msg:+ — $msg}"
        fail=$((fail + 1))
    fi
}

check() {
    local name="$1"; shift
    if "$@" >/dev/null 2>&1; then report "$name" ok
    else report "$name" fail
    fi
}

echo "== M2 Core Data Plane — Completion Matrix =="

# Phase 2.1
check "19 Cluster reachable"             kubectl get nodes
check "19 Cilium running"                kubectl -n kube-system get ds/cilium-agent
check "20 Argo CD installed"             kubectl -n argocd get deploy/argocd-server
check "20 App-of-apps Synced"            bash -c 'kubectl -n argocd get application helixgitpx-app-of-apps -o jsonpath="{.status.sync.status}" | grep -q Synced'
check "20 All 25 Apps Healthy"           bash -c "[ \$(kubectl -n argocd get applications -o jsonpath=\"{.items[*].status.health.status}\" | tr \" \" \"\n\" | grep -c Healthy) -ge 25 ]"
check "21 cert-manager running"          kubectl -n cert-manager get deploy
check "21 ClusterIssuer selfsigned-ca"   kubectl get clusterissuer selfsigned-ca
check "22 SPIRE server ready"            bash -c 'kubectl -n helix-identity rollout status sts/spire-server --timeout=30s'

# Phase 2.2
check "23 CNPG cluster Ready"            bash -c 'kubectl -n helix-data get cluster helix-pg -o jsonpath="{.status.phase}" | grep -qE "healthy"'
check "24 Schemas + RLS applied"         bash -c 'kubectl -n helix-data exec helix-pg-1 -- psql -U postgres -d helixgitpx -tAc "SELECT count(*) FROM information_schema.schemata WHERE schema_name IN (\"hello\",\"auth\",\"repo\",\"sync\",\"conflict\",\"upstream\",\"collab\",\"events\",\"platform\")" | grep -q 9'
check "25 Goose migrations applied"      bash -c 'kubectl -n helix get job/hello-migrate -o jsonpath="{.status.succeeded}" | grep -q 1'
check "26 CNPG backups in MinIO"         bash -c 'kubectl -n helix-system exec deploy/minio -- mc ls local/cnpg-backups | head -1'

# Phase 2.3
check "27 Strimzi Kafka Ready"           bash -c 'kubectl -n helix-data get kafka helix-kafka -o jsonpath="{.status.conditions[?(@.type==\"Ready\")].status}" | grep -q True'
check "28 Karapace responsive"           bash -c 'kubectl -n helix-data exec deploy/karapace -- curl -fsS http://localhost:8081/_schemas >/dev/null'
check "29 Debezium KafkaConnect Ready"   bash -c 'kubectl -n helix-data get kafkaconnect debezium -o jsonpath="{.status.conditions[?(@.type==\"Ready\")].status}" | grep -q True'
check "30 Outbox KafkaConnector Running" bash -c 'kubectl -n helix-data get kafkaconnector hello-outbox -o jsonpath="{.status.conditions[?(@.type==\"Ready\")].status}" | grep -q True'

# Phase 2.4
check "31 Dragonfly running"             bash -c 'kubectl -n helix-cache get sts | grep -q dragonfly'
check "32 Meilisearch present"           bash -c 'kubectl -n helix-cache get sts | grep -q meilisearch'
check "33 OpenSearch present"            bash -c 'kubectl -n helix-cache get sts | grep -q opensearch'
check "34 Qdrant present"                bash -c 'kubectl -n helix-cache get sts | grep -q qdrant'

# Phase 2.5
check "35 Vault unsealed"                bash -c 'kubectl -n helix-secrets exec sts/vault-0 -- vault status | grep -q "Sealed.*false"'
check "36 Observability stack"           bash -c 'kubectl -n helix-observability get deploy | grep -cE "prom-grafana|mimir|loki|tempo|pyroscope" | grep -qE "[5-9]"'
check "37 Grafana dashboards provisioned" bash -c "kubectl -n helix-observability get cm -l grafana_dashboard=1 | wc -l | awk \"{exit (\\\$1 < 2)}\""
check "38 Alertmanager config renders"   bash -c 'kubectl -n helix-observability exec sts/alertmanager-prom-kube-prometheus-stack-alertmanager-0 -- amtool config show >/dev/null'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
