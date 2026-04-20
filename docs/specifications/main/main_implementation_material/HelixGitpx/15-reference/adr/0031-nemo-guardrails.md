# ADR-0031 — NeMo Guardrails as out-of-process policy layer

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

LLM outputs need policy enforcement (jailbreak detection, PII filtering, output shape validation). Options: in-process Go re-implementation (hard, brittle), NeMo Guardrails Python service (mature, out-of-process), LangChain guardrails (Python).

## Decision

NeMo Guardrails runs as a dedicated Pod; ai-service's outputs pass through it via HTTP before returning to the client.

## Consequences

- +50ms per request; dominant cost is still LLM inference.
- Adds a Python dep to the stack (separate container, isolated lifecycle).
- Policy configuration in Colang (NeMo's rule language); versioned in `impl/helixgitpx-platform/helm/nemo-guardrails/rails/`.

## Links

- Spec §LOCKED C-2
- https://github.com/NVIDIA/NeMo-Guardrails
