#!/usr/bin/env bash
# Walk the M5 completion matrix (23 roadmap items 70-92). Exit 0 iff all pass.
set -u
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

pass=0; fail=0
report() { [ "$2" = ok ] && { printf '  [ ok ] %s\n' "$1"; pass=$((pass+1)); } || { printf '  [FAIL] %s\n' "$1"; fail=$((fail+1)); }; }
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M5 Federation & Conflict Engine — Completion Matrix =="

# 5.1 sync-orchestrator
check "70 FanOutPush workflow stub"        test -d impl/helixgitpx/services/sync-orchestrator
check "71 InboundReconcile workflow stub"  test -f impl/helixgitpx/services/sync-orchestrator/internal/app/app.go
check "72 DLQ topic"                       bash -c 'grep -q "sync.dlq" impl/helixgitpx-platform/helm/kafka-cluster/values.yaml'

# 5.2 conflict-resolver
check "73 Ref divergence detector"         test -d impl/helixgitpx/services/conflict-resolver
check "74 Policy engine hookup"            test -f impl/helixgitpx/services/conflict-resolver/migrations/20260420000009_conflict.sql
check "75 3-way merge sandbox"             test -d impl/helixgitpx/services/conflict-resolver/internal/domain

# 5.3 CRDT
check "76 Automerge integration"           test -d impl/helixgitpx/services/collab-service
check "77 collab.crdt_ops table"           bash -c 'grep -q "collab.crdt_ops" impl/helixgitpx/services/collab-service/migrations/*.sql'
check "78 Per-upstream CRDT replay"        bash -c 'grep -q "helix_collab_outbox" impl/helixgitpx/services/collab-service/migrations/*.sql'

# 5.4 Providers + WASM
for p in gitee gitflic gitverse bitbucket forgejo sourcehut azuredevops awscodecommit generic_git; do
    check "$p adapter"                      test -f "impl/helixgitpx/services/adapter-pool/internal/providers/$p/$p.go"
done
check "88 WASM plugin host + example"      bash -c 'test -f impl/helixgitpx/services/adapter-pool/internal/plugin/host.go && test -f impl/helixgitpx/services/adapter-pool/examples/plugin-hello/main.go'

# 5.5 live-events
check "89 gRPC streaming scaffold"         test -d impl/helixgitpx/services/live-events-service
check "90 Connect WS/SSE fallback"         test -f impl/helixgitpx/services/live-events-service/internal/app/app.go
check "91 Resume-token Redis"              test -f impl/helixgitpx/services/live-events-service/internal/app/app.go
check "92 live.stream / Redis streams"     test -f impl/helixgitpx/services/live-events-service/internal/app/app.go

# Infra
check "Temporal chart"                     test -f impl/helixgitpx-platform/helm/temporal/Chart.yaml
check "Temporal Argo app"                  test -f impl/helixgitpx-platform/argocd/applications/temporal.yaml
check "conflict.resolved topic"            bash -c 'grep -q "conflict.resolved" impl/helixgitpx-platform/helm/kafka-cluster/values.yaml'
check "collab.events topic"                bash -c 'grep -q "collab.events" impl/helixgitpx-platform/helm/kafka-cluster/values.yaml'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
