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
