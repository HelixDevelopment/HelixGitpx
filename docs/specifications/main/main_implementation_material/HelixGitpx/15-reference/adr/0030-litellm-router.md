# ADR-0030 — LiteLLM as the LLM router

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

ai-service needs to route prompts across Ollama (local), vLLM (staging GPU), and external providers (OpenAI/Anthropic) based on per-org policy. Options: build our own adapter layer, vendor-lock to a single provider SDK, or use LiteLLM (unifies 100+ providers behind an OpenAI-compatible API).

## Decision

LiteLLM. ai-service speaks OpenAI-compatible HTTP; LiteLLM deployed as a sidecar/independent service translates to the target backend.

## Consequences

- Zero code change when swapping backends.
- Per-org model routing via LiteLLM's config-based routing.
- Adds one network hop; acceptable for the 10s-of-millis overhead vs. an LLM inference of 1-30s.

## Links

- Spec §LOCKED C-1
- https://github.com/BerriAI/litellm
