# ADR-0032 — DPO self-learning pipeline on Ray GPU pool

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

User feedback (accept/reject/edit) on AI proposals is a free training signal. Options: sit on feedback forever, ad-hoc finetunes, DPO (Direct Preference Optimization) which trains directly on preference pairs.

## Decision

Feedback captured via `/api/ai/feedback`, curated (PII scrubbed), used as DPO preference pairs. Training runs on a Ray cluster with GPU workers. Per-task LoRA adapters; shadow-mode evaluates a candidate before promoting via ai-service config change.

## Consequences

- Continuous improvement without labellers.
- GPU pool is a recurring cost; addressed via autoscaling in M8.
- Shadow-mode eval prevents regression; acceptance rate ≥ 0.7 is the spec's M7 exit gate.

## Links

- Spec §LOCKED C-3
