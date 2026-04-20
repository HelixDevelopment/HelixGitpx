# HelixAgent × go-elder-plinius — Verification Report

> **Status:** Written as a companion to
> [`helixagent-plinius-integration.md`](./helixagent-plinius-integration.md).
> **Read this first** before scheduling or estimating the integration.
>
> The planning doc lists a Go port layer (`go-elder-plinius`, `go-v3r1t4s`,
> `go-autotemp`, …) as if it already exists. This report grounds every claim
> against what is publicly discoverable on GitHub today.
>
> **Date:** 2026-04-20.

---

## Confirmed

- **HelixAgent repo** exists: Go 70.4 %, MIT, ~2,287–2,295 commits (snapshot
  within the doc's quoted figure).
  https://github.com/HelixDevelopment/HelixAgent
- **HelixDevelopment org** is real and hosts HelixAgent, HelixLLM, HelixMemory,
  HelixQA, LLMProvider, LLMOrchestrator, DocProcessor, HelixCode, HelixBuilder,
  VisionEngine, HelixTranslate, and this project (HelixGitpx).
  https://github.com/HelixDevelopment
- **vasic-digital org** exists and hosts many submodules the plan references:
  LLMsVerifier, LLMProvider, Concurrency, Cache, Database, Challenges, RAG,
  Observability, Messaging, EventBus, Auth, Embeddings, Formatters, Memory,
  Filesystem, Security, Storage, MCP\_Module, Plugins, Recovery, RateLimiter,
  Middleware, VectorDB, Optimization, TOON, Lazy, Config, Discovery, Assets,
  Watcher, Streaming, Entities, Media. https://github.com/vasic-digital
- **elder-plinius** is a real user ("Pliny the Liberator", ~12.9 k followers).
  Confirmed upstreams referenced by the plan: L1B3RT4S, CL4R1T4S, G0DM0D3,
  OBLITERATUS, ST3GG, V3SP3R, V3R1T4S, AutoTemp, LEAKHUB, GLOSSOPETRAE,
  P4RS3LT0NGV3, I-LLM, BasiliskToken, HyperTune, AutoStoryGen, Dioscuri,
  ourobopus, Gitty, Misc.-Prompt-Hacks, AutoRedTeam, Leda.
  https://github.com/elder-plinius
- **V3SP3R = Flipper Zero AI control** — confirmed; Adafruit coverage
  (2026-03-23). https://github.com/elder-plinius/V3SP3R
- **Lakera Gandalf** is real (7-level prompt-injection gamified challenge).
  https://gandalf.lakera.ai
- **TensorTrust** is real — Berkeley CHAI, ICLR 2024 paper; 126 k attacks,
  46 k defenses dataset. https://tensortrust.ai/ ·
  https://github.com/HumanCompatibleAI/tensor-trust
- **HackAPrompt** is real — referenced prompt-injection challenge corpus.
- AUTOTEMP, OBLITERATUS, LEAKHUB descriptions match the real upstream READMEs.

## Corrections / clarifications

- **The "`go-elder-plinius`" module family does not exist.** Searches for
  `go-v3r1t4s`, `go-autotemp`, `go-l1b3rt4s`, … across both vasic-digital and
  the wider GitHub return zero. The upstreams are real (Python/TS/Java/etc. at
  `github.com/elder-plinius`) but **no Go port layer is public today.** This is
  the plan's biggest load-bearing fiction — any work assuming
  `vasic-digital/go-elder-plinius*` imports must first *create* those modules
  or vendor upstream as-is.
- **`go-tempest`** — no elder-plinius repo of this name exists. The
  "environmental context awareness" feature has no upstream; it is either
  invented or maps to a different, unnamed project.
- **`go-gandalf-solutions`** — no upstream repo. Gandalf is a Lakera product,
  not an elder-plinius project; there is no canonical solutions corpus to port.
  Must be built from scratch or scraped from community write-ups.
- **`go-gitgpt`** — no elder-plinius upstream. `Gitty` exists; `GitGPT` does
  not under elder-plinius.
- **Prompt-leak module count:** the plan says "5 prompt-leak modules (OpenAI,
  Google, Anthropic, xAI, Mistral, Bing, Gemini, Grok, Mixtral)" — that is 9
  names, not 5. **CL4R1T4S is a single repo** holding all these prompts, not a
  set of per-vendor modules.
- **HelixSpecifier, DebateOrchestrator, Agentic, MCP-Servers,
  ConversationContext, Benchmark, BuildCheck, Containers, BootManager,
  BackgroundTasks** are not visible in either org's public repo list. They may
  be private, not yet created, or hosted under a third org. Treat as
  **unverified / likely private** until read access is confirmed.
- **"13 embedding providers" and "35 MCP server implementations"** are
  unverifiable externally; treat as claims to confirm during integration, not
  established facts.

## Unverifiable

- Exact capability counts (47 + LLM providers, 35 MCP servers, 13 embedding
  providers, 10 LSP servers, 25 responses / 5×5 ensemble, 7-phase SpecKit).
  HelixAgent README broadly corroborates "47 + providers" but specific internal
  counts require repo-side inspection.
- Claimed performance gains (−40–60 % hallucinations, +15–30 % quality,
  −30–50 % latency, +10–25 % per-provider quality). No benchmarks cited;
  these are aspirational targets.
- Internal paths (`LLMsVerifier/internal/verification/ensemble_verifier.go`,
  etc.) cannot be validated without repo read access. The file tree in the plan
  should be treated as a *proposal* until verified against each submodule.

## Integration impact

**The upstream ecosystem is real; the Go port layer is not.** Every module
except `go-tempest`, `go-gandalf-solutions`, and `go-gitgpt` has a real
upstream in Python/TS/Java. But there is no `vasic-digital/go-elder-plinius*`
namespace today, so **Phase 1 must become "create the Go port repos" (or
vendor upstream as-is)**, not "wire existing modules into HelixAgent". The
plan implies modules are ready to import — they are not.

Several HelixAgent submodules named in the integration table
(DebateOrchestrator, Agentic, HelixSpecifier, MCP-Servers, BackgroundTasks,
BootManager) are not publicly discoverable. Before scoping file-level edits
like `DebateOrchestrator/internal/phases/synthesis.go`, confirm these exist
privately and secure read access.

Trust the architectural shape, the upstream inventory, and the external
references (Gandalf / TensorTrust / HackAPrompt). **Do not trust** the
performance numbers, the "5 prompt-leak modules" framing, or the implicit
assumption that Go ports already exist.

### Recommendation — insert a W0 "port & inventory" spike

Before Phase 1 (W1–2 in the plan), add a Week 0 spike that:

1. Creates Go port repos for the ~22 real upstreams (or vendors upstream
   sources into a `third_party/` tree and wraps with thin Go adapters).
2. Maps the three fictional ones (`tempest`, `gandalf-solutions`, `gitgpt`) to
   explicit decisions: *drop*, *rescope*, or *build-from-scratch*, each with
   an owner and a budget.
3. Audits the private HelixAgent submodules so the file paths in the plan are
   validated before engineers open PRs against them.

Without this spike, schedule estimates in the plan are under-counted by the
entire porting effort.

## Sources

- HelixAgent: https://github.com/HelixDevelopment/HelixAgent
- HelixDevelopment org: https://github.com/HelixDevelopment
- vasic-digital org: https://github.com/vasic-digital
- elder-plinius profile: https://github.com/elder-plinius
- V3SP3R: https://github.com/elder-plinius/V3SP3R —
  Adafruit: https://blog.adafruit.com/2026/03/23/ai-brain-for-the-flipper-zero/
- AutoTemp: https://github.com/elder-plinius/AutoTemp ·
  OBLITERATUS: https://github.com/elder-plinius/OBLITERATUS ·
  LEAKHUB: https://github.com/elder-plinius/LEAKHUB ·
  CL4R1T4S: https://github.com/elder-plinius/CL4R1T4S ·
  L1B3RT4S: https://github.com/elder-plinius/L1B3RT4S ·
  G0DM0D3: https://github.com/elder-plinius/G0DM0D3 ·
  GLOSSOPETRAE: https://github.com/elder-plinius/GLOSSOPETRAE ·
  Leda: https://github.com/elder-plinius/Leda
- Lakera Gandalf: https://gandalf.lakera.ai/ —
  background: https://www.lakera.ai/blog/who-is-gandalf
- TensorTrust paper (ICLR 2024): https://tensortrust.ai/paper/ —
  repo: https://github.com/HumanCompatibleAI/tensor-trust
