# 10 — Self-Learning LLM Platform

> **Document purpose**: Describe the **AI subsystem** of HelixGitpx — how it reasons about conflicts, PR reviews, intent summaries — and how it **learns from user feedback** without sending anyone's proprietary code to third-party training corpora.

---

## 1. Requirements

| Requirement | Implication |
|---|---|
| Must never leak customer code into external model training data | Default stack is **self-hosted** |
| Must work with no GPU at all (degraded) | CPU-capable models (Mistral 7B Q4, Phi-3 mini) must be usable |
| Must support best-available models when allowed | LiteLLM router with per-org policy |
| Must learn from user accept/reject patterns | RLAIF pipeline |
| Must not regress silently | Continuous evaluation harness with canary models |
| Must respect cost envelopes | Per-tenant token budgets + model tiering |
| Must produce machine-readable outputs | Grammar-constrained decoding / JSON Schema enforced |
| Must be replaceable | Providers pluggable via LiteLLM |

---

## 2. Components

```
┌──────────────────────────────────────────────────────────────┐
│ ai-service (Go + Python sidecars)                            │
│ ┌────────────┐ ┌───────────┐ ┌────────────┐ ┌─────────────┐ │
│ │ Router     │ │ Guardrails│ │ Embedder   │ │ Fine-Tune   │ │
│ │ (LiteLLM)  │ │ (NeMo)    │ │ (BGE-M3)   │ │ orchestrator│ │
│ └────┬───────┘ └─────┬─────┘ └─────┬──────┘ └──────┬──────┘ │
│      │               │             │                │        │
└──────┼───────────────┼─────────────┼────────────────┼────────┘
       │               │             │                │
       ▼               ▼             ▼                ▼
  Self-hosted     Content          Qdrant         Ray/KubeFlow
  Ollama/vLLM    policies         vectors         training jobs
  (Llama, Phi,     Toxic/PII                        GPU nodes
  DeepSeek, …)      filters
```

---

## 3. Model Catalogue

HelixGitpx ships with a curated model registry. Org admins choose a **tier**; specific models can be pinned.

| Tier | Typical model | Size | GPU? | Notes |
|---|---|---|---|---|
| `local-light` | Phi-3-mini (3.8B) Q4 | 2–4 GB | optional | CPU-capable, basic tasks |
| `local-standard` | Qwen2.5-Coder-7B Q4 / DeepSeek-Coder-6.7B Q4 | 4–5 GB | preferred | Good default for code |
| `local-heavy` | Qwen2.5-Coder-32B Q4 / Llama 3.1-70B Q4 | 20–40 GB | required | High quality |
| `cloud-approved` | Claude Sonnet, GPT-4o (per org allowlist) | N/A | — | Only if org has consented policy |
| `fine-tuned` | Any base above + org-specific LoRA adapter | varies | varies | Produced by pipeline (§5) |

Default: `local-standard` + optional `fine-tuned` adapter per org.

### 3.1 Deployment

- **Ollama** for single-replica / dev.
- **vLLM** or **TGI** for high-throughput inference with tensor parallelism.
- **KEDA** autoscales replicas based on queue depth (topic `ai.prompt.run`).
- **NVIDIA GPU Operator** manages GPU scheduling.

---

## 4. Routing (LiteLLM)

LiteLLM is an internal proxy that:

- Normalises provider APIs into an OpenAI-compatible interface.
- Applies per-org policies (model allowlist, token budget).
- Caches responses by `(prompt_hash, model_version)` (Redis; TTL 1 h by default).
- Records metrics and cost per tenant.
- Falls back to secondary model on provider error.

Routing decision order:
1. Org explicit override.
2. Task preference (e.g. conflict_resolution → fine-tuned-conflict-v2).
3. Tier default.
4. Fallback: `local-light`.

---

## 5. Self-Learning Pipeline (RLAIF)

The system continuously improves without sending user code outside.

### 5.1 Data Collection

- Every AI suggestion emits `ai.prompt.run` with `{input, output, model, confidence}`.
- Every user action (accept/reject/edit) emits `ai.feedback` with `{accepted, rating, edit_distance}`.
- Data stored in `ai.prompt_runs` and `ai.feedback_records`.

### 5.2 Dataset Curation

Weekly offline job (`ai-curator`):

1. Pulls last 7 days of `(prompt, output, feedback)` tuples.
2. **PII/secret scrubbing** (Gitleaks rules + heuristic).
3. Human review of 1 % sampled examples (via internal UI).
4. Approved dataset pushed to S3 as a DVC-managed artefact.
5. Emits `ai.dataset.curated` event.

### 5.3 Training

- **Preference pairs** from accept/reject or edit-distance thresholds.
- **DPO (Direct Preference Optimisation)** on the selected base model — less compute than PPO, more stable.
- Optional: **LoRA adapters** (we rarely full-fine-tune). One adapter per task per org (if org opts in to org-specific models).
- Orchestration: **Ray Train** or **KubeFlow** on GPU node pool.
- Emits `ai.fine_tune.completed` with new model artefact in registry.

### 5.4 Evaluation (Canary)

- New model runs in **shadow mode** for 48 h: predictions compared against current model; humans sample & grade.
- Metrics: accept-rate, edit-distance, safety score.
- If shadow model improves accept-rate by ≥ 3 % and does not regress safety → promote.

### 5.5 Safety Gates

Before any promotion, the new model must pass:

- **Harmfulness eval** (PoisonedMR, NoCode-Eval).
- **Hallucination eval** (task-specific holdout).
- **Prompt-injection robustness**.
- **Cost regression check**.
- Signed by security lead.

### 5.6 Rollback

If post-promotion metrics regress, automatic rollback within 1 h. Rollback is just flipping the router pointer.

---

## 6. Prompting & Guardrails

### 6.1 Prompts as Code

- Prompt templates in `prompts/` repo, versioned with semver.
- DSPy-compiled where possible (compile once, deploy).
- Every prompt has unit tests (golden examples + regression suite).

### 6.2 Guardrails (NeMo / Guardrails AI)

- **Input rails**: reject prompts with prompt-injection markers, PII beyond policy, extreme sizes.
- **Output rails**: enforce schema, reject unsafe content, PII-scrub.
- **Dialog rails**: keep chat on-topic for chatops.

### 6.3 Structured Output

- Every RPC that returns AI output defines a JSON Schema.
- Decoding uses a grammar (e.g. `lm-format-enforcer`, `outlines`) to guarantee validity.
- On any parse failure, fallback to "human_required".

---

## 7. Use Cases

### 7.1 Conflict Resolution

- Input: diff (left/right/base), commit messages, language, style guide, similar past resolutions (RAG).
- Output: strategy + optional patch + rationale + confidence.
- Post-processing: sandbox-validate (compile/lint/test) before confidence is trusted.

### 7.2 PR Summary & Review

- Input: PR title, body, diff, related issues.
- Output: concise human summary, bullet-point review, suggested reviewers.
- Optional: **line-level review comments** for common smells (unused var, missing error check, unclear name).

### 7.3 Intent-of-Commit Summaries

- Translate messy commit messages into clean conventional commits.
- Opt-in per repo.

### 7.4 Label & Milestone Suggestions

- Multi-label classifier → suggested labels on new issues/PRs.
- Thresholded auto-apply (configurable).

### 7.5 Semantic Search & Code-NL Q&A

- Embed files and symbols; run hybrid search (see [05b-search-indexing.md](../03-data/05b-search-indexing.md)).
- RAG over a repo: "Where do we initialise the Redis client?"

### 7.6 Chatops

- Natural-language commands: "Open a PR from branch `foo` to `main` and ping @alice."
- Resolved via tools (function calling) — our tool registry is the gRPC API.

### 7.7 Translation / Localisation

- Translate issue text for cross-locale repos (e.g., Gitee ↔ GitHub).
- Keep originals always; translations side-by-side.

---

## 8. Privacy & Tenancy

- **Default**: all inference local; no data leaves cluster.
- **Opt-in cloud**: per-org flag `allow_cloud_llm` → allows routing to approved cloud providers, with per-provider data-sharing audit.
- **Training data segregation**: fine-tuned adapters for org A never trained on org B's data.
- **Right to erasure**: `/api/v1/ai/forget` permanently removes a user's contributions from curated datasets (and schedules retraining if needed).
- **Explainability**: every AI-driven decision carries `{prompt_version, model_version, confidence, rationale}` — surfaced in UI.

---

## 9. Cost Metering

- Token-level accounting (`ai.prompt_runs.input_tokens`, `output_tokens`, `cost_usd`).
- Per-org monthly budget (`billing.plans.limits.llm_tokens_per_month`).
- Soft warn → hard throttle → block (with graceful degrade: simpler model or skip AI step).

---

## 10. Observability

| Metric | |
|---|---|
| `helixgitpx_ai_requests_total{task, model, status}` | |
| `helixgitpx_ai_latency_seconds{task, model}` | |
| `helixgitpx_ai_accept_rate{task, model}` | Target ≥ 0.7 |
| `helixgitpx_ai_confidence_bucket{task, bucket}` | Distribution |
| `helixgitpx_ai_cost_usd{org, model}` | |
| `helixgitpx_ai_guardrail_reject_total{rail}` | |
| `helixgitpx_ai_hallucination_score` (eval harness) | Continuous eval |

Alerts:
- `AIModelRegressionDetected` (accept-rate drop > 5 % for 6 h).
- `AIGuardrailSpike` (rejections > threshold).
- `AICostBurnHigh` (org > 80 % budget in < half the period).

---

## 11. Testing

- **Prompt unit tests**: golden input → expected structural output (schema + anchor assertions).
- **Eval harness**: curated benchmark per task, run on every candidate model.
- **Adversarial tests**: prompt-injection suite (including from [OWASP Top 10 for LLMs](https://owasp.org/www-project-top-10-for-large-language-model-applications/)).
- **Cost tests**: CI assert average tokens-in / tokens-out stay within budget.
- **Load**: `vegeta` to AI endpoints; ensure autoscaling scales in time.
- **Canary**: new model must beat incumbent on curated holdout with significance.

---

## 12. Model Registry

All models (base + adapters + fine-tunes) live in `ai.model_registry` with:

- Name, version, purpose, base, storage URI.
- Metrics from last eval.
- `live_since`, `retired_at`.

Every inference records the exact `model_version` it used, giving reproducibility and root-cause analysis.

---

## 13. Future Directions (Post-GA)

- **Agentic conflict resolver**: multi-step tool-use (clone, edit, compile, test) for harder conflicts.
- **Per-repo style adapters**: tiny LoRA learned from the repo's own history.
- **Cross-lingual commit coach**: rewrite commit messages in the maintainer's preferred language with original preserved.
- **On-device inference for desktop clients** via llama.cpp for fully offline mode.

---

*— End of Self-Learning LLM Platform —*
