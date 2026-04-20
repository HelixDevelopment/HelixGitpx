# Runbook — LLM backend outage

**Alert:** `AIServiceErrorRateHigh` with error pattern `connection refused` to ollama/vllm.

## Immediate checks

- `kubectl -n ai get pods`
- `kubectl -n ai logs -l app=ollama --tail=100`
- GPU pressure: `nvidia-smi` on nodes; check for OOM on model load.

## Remediation

1. **Ollama pod down** — restart; if crashlooping, reduce `OLLAMA_NUM_PARALLEL`.
2. **vLLM pod down** — verify model weights mounted (check the init container logs).
3. **Router (LiteLLM) misconfigured** — check ConfigMap `litellm-config`; roll back last change.
4. **External provider outage (Anthropic/OpenAI)** — enable fallback route in LiteLLM: set `fallbacks: [ollama-local]`.

## Graceful degradation

AI features have a `featureflag.ai.enabled` kill switch via OPA bundle. Flip off to preserve git/API.

## Escalation

Outage > 15 min → page AI platform lead; notify customers if AI features are on critical path.
