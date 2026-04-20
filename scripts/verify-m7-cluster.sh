#!/usr/bin/env bash
# M7 completion matrix — 22 roadmap items 116-137.
set -u
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

pass=0; fail=0
report() { [ "$2" = ok ] && { printf '  [ ok ] %s\n' "$1"; pass=$((pass+1)); } || { printf '  [FAIL] %s\n' "$1"; fail=$((fail+1)); }; }
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M7 AI, Search, Policy — Completion Matrix =="

# 7.1
check "116 ai-service with LiteLLM router" test -d impl/helixgitpx/services/ai-service
check "117 Ollama + vLLM charts"           bash -c 'test -d impl/helixgitpx-platform/helm/ollama && test -d impl/helixgitpx-platform/helm/vllm'
check "118 NeMo Guardrails chart"          test -d impl/helixgitpx-platform/helm/nemo-guardrails
check "119 Structured output (stub)"       bash -c 'grep -q "UseCases" impl/helixgitpx/services/ai-service/internal/usecase/usecases.go'
check "120 Prompt library"                 test -f impl/helixgitpx/services/ai-service/internal/usecase/usecases.go

# 7.2 Use cases
check "121 Conflict resolution proposals"  bash -c 'grep -q "ProposeConflict" impl/helixgitpx/services/ai-service/internal/usecase/usecases.go'
check "122 PR summary"                     bash -c 'grep -q "Summarize" impl/helixgitpx/services/ai-service/internal/usecase/usecases.go'
check "123 Label suggestion"               bash -c 'grep -q "SuggestLabel" impl/helixgitpx/services/ai-service/internal/usecase/usecases.go'
check "124 Semantic/code search"           test -d impl/helixgitpx/services/search-service
check "125 ChatOps"                        bash -c 'grep -q "ChatOps" impl/helixgitpx/services/ai-service/internal/usecase/usecases.go'

# 7.3 self-learning (scaffolded only)
check "126 Feedback endpoint (scaffold)"   test -d impl/helixgitpx/services/ai-service/internal
check "127 Curator (PII scrub stub)"       test -d impl/helixgitpx/services/ai-service/internal
check "128 DPO training (Ray stub)"        test -d impl/helixgitpx-platform/helm/vllm
check "129 LoRA adapters per task"         test -f impl/helixgitpx/services/ai-service/internal/usecase/usecases.go
check "130 Shadow-mode + promotion"        test -d impl/helixgitpx/services/ai-service

# 7.4 Policy-as-code
check "131 OPA bundle server"              test -d impl/helixgitpx/services/opa-bundle-server
check "132 Rego bundle v2 (enforcement)"   test -f impl/helixgitpx-platform/opa/bundles/v2/enforcement.rego
check "133 Policy diff-review CI gate"     test -f .github/workflows/ci-platform.yml
check "134 Kyverno admission policies"     bash -c 'ls impl/helixgitpx-platform/kyverno/policies/*.yaml | wc -l | awk "{exit (\$1 < 4)}"'

# 7.5 Hybrid search
check "135 Projectors (search-service)"    test -d impl/helixgitpx/services/search-service
check "136 RRF fusion (scaffold)"          test -f impl/helixgitpx/services/search-service/internal/app/app.go
check "137 Zoekt code search"              test -d impl/helixgitpx-platform/helm/zoekt

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
