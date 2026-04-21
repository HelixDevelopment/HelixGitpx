# Integration Plans

This directory holds authoritative, implementation-ready analysis documents for
planned integration efforts. Each file describes an external system, what
capabilities it brings, and the exact files/modules that will be touched when
the integration is scheduled.

## Contents

- [`helixagent-plinius-integration.md`](./helixagent-plinius-integration.md) —
  Integrating the go-elder-plinius module family into the HelixAgent project
  (LLMsVerifier, HelixQA, DebateOrchestrator, Agentic, HelixLLM, etc.).
  13 modified files, 22 new files across 15 submodules. 20 new capabilities.
- [`helixagent-plinius-verification.md`](./helixagent-plinius-verification.md)
  — **read alongside the plan.** Independent fact-check: confirms upstream
  elder-plinius repos exist, but the Go port layer does not; recommends a
  Week 0 "port and inventory" spike before Phase 1.
- [`helixagent-plinius-w0-spike.md`](./helixagent-plinius-w0-spike.md)
  — the Week 0 spike specification. This is what must land before
  the plan's Phase 1 can start.
- [`helixagent-plinius-policy-review.md`](./helixagent-plinius-policy-review.md)
  — **authoritative verdict** per module: 7 KEEP, 6 KEEP-GATED,
  7 DROP, 1 DROP-and-don't-build. Supersedes the indiscriminate
  "integrate all 20" framing in the original plan.

## Status

These are planning docs. They are **not** merged into the HelixGitpx
implementation — most reference systems outside this repo (HelixAgent, etc.).
Treat each doc as a blueprint for a future PR series against the target system.

## Adding a new integration doc

1. Write the analysis as `<system-a>-<system-b>-integration.md`.
2. Include: deep architecture dive, capability-by-capability integration
   points (file paths!), performance/stability analysis, rollout phasing,
   full file index.
3. Reference it from this README.
