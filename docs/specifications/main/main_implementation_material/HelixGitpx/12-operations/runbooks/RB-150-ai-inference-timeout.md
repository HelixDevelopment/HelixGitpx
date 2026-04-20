# RB-150 — AI Inference Timeouts / Pool Exhaustion

> **Severity default**: P2 (P1 if conflict auto-apply rate drops below 60 % for 1 h)
> **Owner**: AI team
> **Last tested**: 2026-03-22

## 1. Detection

Alerts:
- `AIInferenceTimeouts` — timeout rate > 5 % for 10 min.
- `AIModelRegression` — accept rate dropped > 5 pts vs. 6 h ago.
- `QdrantOOMRisk` — if hybrid search is the chokepoint.

Supporting signals:
- `helixgitpx_ai_queue_depth` climbing.
- `nvidia_gpu_duty_cycle` near 100 % with queue.
- Customer reports "AI isn't suggesting anything" or "very slow".

---

## 2. First 2 Minutes

1. Ack; note whether inference or quality (or both) is impacted.
2. Check AI dashboard: per-model request rate, latency p99, queue depth, GPU utilisation.
3. Is it one model or all? One task (conflict) or all (conflict + PR summary + search)?

---

## 3. Diagnose

| Pattern | Likely cause | Check |
|---|---|---|
| All models slow, high GPU util | Capacity | `nvidia_gpu_duty_cycle`, KEDA scaling, queue depth |
| One model slow, others fine | Specific LoRA / size issue | `ai_requests_total{model=...}` |
| Queue grows, GPU idle | Scheduler / router bug | LiteLLM / vLLM queue state |
| Errors with `5xx` | Model crashed | Inference pod logs |
| Sudden quality drop | Model regression | `ai_accept_rate` offset comparison, recent deploy |
| Timeouts only on long prompts | Context overflow | `ai_input_too_large_total` |
| Gradual accept-rate decline | Data drift | Eval dashboards; compare model versions |

---

## 4. Mitigations

### 4.1 Scale GPU pool

```bash
# Bump KEDA max (if capped)
helixctl keda override --scaled-object=ai-inference-vllm --max-replicas=12

# Manual scale (emergency, KEDA will catch up)
helixctl scale ai-inference-vllm --replicas=8
```

GPU provisioning can take 2-5 min (cloud spot may be faster). Consider Savings Plan headroom.

### 4.2 Degrade gracefully

Reduce work to match capacity:

```bash
# Route only conflict requests to AI; skip PR summary + label suggest
helixctl flag set ai.pr-summary-enabled --off --reason="RB-150"
helixctl flag set ai.label-suggest-enabled --off --reason="RB-150"

# Lower auto-apply confidence to require human review
helixctl ai set-threshold --task=conflict_proposal --threshold=0.99
```

Conflicts that fail to auto-apply will escalate to humans — safe fallback.

### 4.3 Fall back to smaller model

```bash
helixctl ai route-override --task=conflict_proposal --model=qwen2.5-coder-7b --duration=30m
```

Pairs with degraded confidence threshold.

### 4.4 Disable cloud routing (if it's a cloud provider issue)

```bash
helixctl flag set ai.cloud-routing-allowed --off
```

### 4.5 Rollback a recent model

If `AIModelRegression` is the primary alert:

```bash
helixctl ai set-active --task=conflict_proposal --model=conflict-v2.3.1
```

Replaces current active version. Accept rate usually recovers within 30 min as in-flight requests drain.

### 4.6 Global kill-switch (last resort)

```bash
helixctl flag set ai.enabled --off --reason="RB-150 major"
```

HelixGitpx continues without AI features; policy + CRDT paths handle conflicts.

---

## 5. Verify Recovery

- Timeout rate back under 1 %.
- Queue depth baseline restored.
- Accept rate recovered (may lag by 15-30 min).
- No new customer complaints.

---

## 6. Post-Incident

- If regression: automated rollback runbook updated; add scenario to pre-promotion evals.
- If capacity: review forecasting; schedule GPU expansion.
- If cloud provider: review fallback automation; consider second cloud.
- Update risk register.

---

## 7. Customer Communication

For prolonged AI degradation:

> Subject: Reduced AI capability on HelixGitpx
>
> We're currently operating with reduced AI capacity, so conflict auto-resolution and PR summarisation may be slower or unavailable. Your code and data are unaffected; policy-based and CRDT-based conflict resolution continue to work normally.
>
> We'll update when full capacity returns.

Never mischaracterise the situation — customers trust candour.

---

## 8. Drill

Chaos scenario G3.3 (LLM inference pool down). Monthly validation.

---

## 9. Related

- RB-120 (Conflict backlog growing)
- RB-121 (AI confidence regression)
- [07-ai/10-llm-self-learning.md] — promotion lifecycle
- Chaos: [15-reference/chaos/playbook.md] § G3.3
