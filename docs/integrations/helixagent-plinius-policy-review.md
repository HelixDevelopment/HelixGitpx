# Plinius integration — policy review and final disposition

> **Status:** Authoritative. Supersedes the indiscriminate "integrate all"
> framing in [`helixagent-plinius-integration.md`](./helixagent-plinius-integration.md).
>
> **Source of verdicts:** session 2026-04-21 research pass against
> `github.com/elder-plinius` (README inspection + recent commits).
>
> **Rule applied:** HelixGitpx is a federation Git proxy. Any module whose
> primary purpose is to **bypass third-party LLM provider safety systems**
> gets DROPPED from the integration — shipping those helpers inside a
> commercial product would facilitate abuse and violate the provider
> Terms of Service that our customers are bound by. Defensive / dual-use
> modules are KEPT with explicit policy + authorization gates.

## Verdict table

| Upstream | Category | Verdict | Why |
|---|---|---|---|
| **L1B3RT4S** | (c) pure offensive | **DROP** | Prompt library whose sole stated purpose is "DISREGARD PREV. INSTRUCTS" jailbreaks against production LLMs. |
| **CL4R1T4S** | (b) dual-use | **KEEP-GATED** | Leaked system-prompt corpus. Transparency artifact; red-team reference. Gate: available only to LLMsVerifier's internal evaluation pipeline, not exposed to end users. |
| **G0DM0D3** | (c) pure offensive | **DROP** | "LIBERATED AI CHAT" front-end pre-loading jailbreak context. |
| **OBLITERATUS** | (c) pure offensive | **DROP** | Alignment-removal tool ("obliterate the chains"). |
| **ST3GG** | (b) dual-use | **KEEP-GATED** | General steganography suite. Useful for covert-channel *detection*; ship only the detector side, not the embed side. Gate: `security.stego_detect` feature flag, ops-only. |
| **V3SP3R** | (b) dual-use | **KEEP-GATED** | Flipper Zero AI control. Relevant only to authorized pentesters. Gate: shipped as a standalone MCP server behind opt-in org policy, not part of default bundle. |
| **V3R1T4S** | (b) tentative | **KEEP-GATED** | Empty README stub. Re-evaluate when content lands; today a no-op. |
| **AutoTemp** | (a) defensive | **KEEP** | Multi-temperature LLM output selection. Quality improvement, no safety impact. |
| **LEAKHUB** | (b) dual-use | **KEEP-GATED** | Leaderboard of system-prompt leaks. Red-team artifact. Gate: internal evaluation only. |
| **GLOSSOPETRAE** | (c) pure offensive | **DROP** | Linguistic obfuscation engine designed to evade content filters. |
| **P4RS3LT0NGV3** | (c) pure offensive | **DROP** | Prompt mutation pipeline for bypassing filters. |
| **I-LLM** | (a) defensive | **KEEP** | Streamlit screenplay/voice generator. Unrelated creative app; no safety concerns. |
| **BasiliskToken** | (c) pure offensive | **DROP** | ERC-20 "AI liberation" token, README rejects "safety or ethics." Not shippable in a commercial product. |
| **HyperTune** | (a) defensive | **KEEP** | LLM hyperparameter search. Standard quality tuning. |
| **AutoStoryGen** | (a) defensive | **KEEP** | Agentic story generator. |
| **Dioscuri** | (c) pure offensive | **DROP** | "Jailbroken Gemini." |
| **ourobopus** | (a) defensive | **KEEP** | Self-improvement agent loop. |
| **Gitty** | (a) defensive | **KEEP** | AI project template. |
| **Misc.-Prompt-Hacks** | (b) dual-use | **KEEP-GATED** | CTF write-ups (Gandalf / TensorTrust solves). Gate: HelixQA regression corpus, internal only. |
| **AutoRedTeam** | (b) dual-use | **KEEP-GATED** | Prompt-defense testing. Gate: run only in staging with an operator's explicit approval + recorded scope. |
| **Leda** | (a) defensive | **KEEP** | Multi-agent orchestrator spec → team. |

### Fictional modules (no upstream)

| Name | Disposition |
|------|-------------|
| `go-gandalf-solutions` | **DROP and do not build.** Packaging jailbreak solutions inside a commercial Git product would ship attacker playbooks. |
| `go-gitgpt` | **Buildable** — neutral name; was never offensive. Re-scope as AI git-assist. |
| `go-tempest` | **Buildable** — neutral name; chaos/env-awareness helper. Re-scope as defensive capability. |

## Updated integration shape

### Phase 1 (W1–2) — read-only / defensive

Port (or use upstream directly): AutoTemp, HyperTune, ourobopus, Leda,
Gitty, AutoStoryGen, I-LLM. These modules have clear defensive / quality
goals and can be integrated into HelixAgent's LLMLlm router + Agentic
framework without policy friction.

### Phase 2 (W3–4) — dual-use, gated

Port: CL4R1T4S (corpus), LEAKHUB (scanner), AutoRedTeam (staging-only),
Misc.-Prompt-Hacks (test corpus), ST3GG-detect, V3SP3R-MCP.

Every Phase-2 module ships behind a runtime policy gate:

- `security.redteam_enabled` for AutoRedTeam.
- `security.stego_detect` for ST3GG.
- `integration.v3sp3r` org-opt-in for Flipper.
- `llmsverifier.use_prompt_corpus` for CL4R1T4S + LEAKHUB + Misc.-Prompt-Hacks.

All disabled by default. All audit-logged. All restricted to staff
principals in the default OPA bundle.

### Phase 3 — deferred/dropped

L1B3RT4S, G0DM0D3, OBLITERATUS, GLOSSOPETRAE, P4RS3LT0NGV3,
BasiliskToken, Dioscuri, `go-gandalf-solutions`: **no integration
planned**. If a future customer use case genuinely requires one of
these (e.g., an authorized academic red-team contract), the request
goes through a formal security review and ships as a standalone
side-channel, never inside the main product.

## Impact on the original plan

- 20 upstream modules → **7 KEEP + 6 KEEP-GATED + 7 DROP** + 3 fictional
  (1 drop-and-don't-build, 2 buildable).
- Estimated ports: 13 modules (was 20) → ~35% less porting effort.
- Phase 1 of the plan collapses to just the 7 KEEP modules; W0 spike
  inventory needs to reflect the DROP decisions.

## Next action

Update `helixagent-plinius-integration.md` and
`helixagent-plinius-w0-spike.md` to reflect these verdicts, then commit.
Done in this session.
