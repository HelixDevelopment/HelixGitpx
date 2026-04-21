#!/usr/bin/env bash
# M2 end-to-end spine — exercises hello through the full data plane.
set -u

SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

if ! command -v kubectl >/dev/null 2>&1 || ! kubectl cluster-info >/dev/null 2>&1; then
    echo "== $(basename "$0" .sh | sed 's/verify-//;s/-spine//') spine — SKIP (no cluster reachable) =="
    exit 0
fi


pass=0
fail=0
report() {
    local name="$1" result="$2"
    if [ "$result" = "ok" ]; then
        printf '  [ ok ] %s\n' "$name"
        pass=$((pass + 1))
    else
        printf '  [FAIL] %s\n' "$name"
        fail=$((fail + 1))
    fi
}
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M2 End-to-end Spine =="

check "HTTP /v1/hello?name=world returns 200" \
    bash -c 'curl -fsS -k --resolve hello.helix.local:443:127.0.0.1 "https://hello.helix.local/v1/hello?name=world" >/dev/null'

check "HTTP response contains hello, world" \
    bash -c 'curl -fsS -k --resolve hello.helix.local:443:127.0.0.1 "https://hello.helix.local/v1/hello?name=world" | grep -q "hello, world"'

check "gRPC SayHello returns greeting" \
    bash -c 'grpcurl -insecure -d "{\"name\":\"world\"}" hello.helix.local:443 helixgitpx.hello.v1.HelloService/SayHello | grep -q "hello, world"'

check "Outbox row inserted in Postgres" \
    bash -c 'kubectl -n helix-data exec helix-pg-1 -- psql -U hello_svc -d helixgitpx -tAc "SELECT count(*) > 0 FROM hello.outbox_events" | grep -q t'

check "Debezium delivered event to hello.said within 30s" \
    bash -c 'timeout 30 kubectl -n helix-data exec deploy/helix-kafka-kafka-0 -- /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server helix-kafka-kafka-bootstrap:9092 --topic hello.said --from-beginning --max-messages 1 --timeout-ms 30000 | grep -q world'

check "Prometheus scrapes hello metrics" \
    bash -c 'kubectl -n helix-observability exec sts/prometheus-prom-kube-prometheus-stack-prometheus-0 -- wget -qO- "http://localhost:9090/api/v1/targets?state=active" | grep -q "\"job\":\"hello\""'

check "Loki has logs from hello" \
    bash -c 'kubectl -n helix-observability exec sts/loki-0 -- wget -qO- "http://localhost:3100/loki/api/v1/query?query={app=\"hello\"}" | grep -q "result"'

check "Tempo has a trace from hello" \
    bash -c 'kubectl -n helix-observability exec deploy/tempo-query-frontend -- wget -qO- "http://localhost:3200/api/search?tags=service.name%3Dhello" | grep -q "traces"'

check "Grafana reachable via Ingress" \
    bash -c 'curl -fsS -k --resolve grafana.helix.local:443:127.0.0.1 "https://grafana.helix.local/api/health" | grep -q ok'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
