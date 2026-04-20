# W0 Spike — HelixAgent × plinius: port & inventory

> **Context:** The [integration plan](./helixagent-plinius-integration.md)
> assumes a `vasic-digital/go-elder-plinius*` Go port layer that **does not
> exist yet** (see the
> [verification report](./helixagent-plinius-verification.md)). This document
> specifies the Week 0 spike that must land before any of the plan's Phase 1
> work can start.
>
> **Owner:** TBD. **Estimated duration:** 2 sprints (10 working days).
> **Exit criteria:** every module cell in §2 is either a real Go module or an
> explicit "dropped / rescoped / build-from-scratch" decision with its own
> ticket.

---

## 1. Objective

Eliminate the gap between the integration plan (which reads as if the Go
ports already exist) and reality. At W0 exit, we have a concrete inventory:

- Every upstream `elder-plinius/*` repo mapped to either
  - a vendored copy in `third_party/`, OR
  - a fresh `vasic-digital/go-<name>` module, OR
  - a "no go-port" decision with rationale.
- Every HelixAgent submodule named in the plan either confirmed public,
  confirmed private (with read access granted to the integration team), or
  flagged as "doesn't exist yet, must create first".
- A revised Phase 1 start date based on what actually exists.

## 2. Module inventory matrix

Fill this in during W0. Each row becomes a separate tracking issue.

| plinius module          | Upstream repo status | Go port approach | HelixAgent target submodule | HelixAgent submodule status | Decision & owner |
|-------------------------|----------------------|------------------|-----------------------------|-----------------------------|------------------|
| V3R1T4S                 | Confirmed upstream   | new go-v3r1t4s   | LLMsVerifier + DebateOrchestrator | Private? | — |
| AUTOTEMP                | Confirmed upstream   | new go-autotemp  | HelixLLM + LLMProvider      | Partly public (LLMProvider) | — |
| CL4R1T4S                | Confirmed upstream   | new go-cl4r1t4s  | LLMsVerifier + HelixSpecifier | HelixSpecifier private? | — |
| gemini/grok/mixtral/bing prompt leaks | Upstream (5 repos) | single consolidated go-promptleaks | LLMsVerifier | — | — |
| I-LLM                   | Confirmed upstream   | new go-i-llm     | Agentic                     | Private? | — |
| Dioscuri                | Confirmed upstream   | new go-dioscuri  | DebateOrchestrator          | Private? | — |
| L1B3RT4S                | Confirmed upstream   | new go-l1b3rt4s  | LLMsVerifier + HelixQA      | Partly public | — |
| AutoRedTeam             | Confirmed upstream   | new go-autoredteam | LLMsVerifier + HelixQA    | — | — |
| BasiliskToken           | Confirmed upstream   | new go-basilisktoken | LLMsVerifier             | — | — |
| ourobopus               | Confirmed upstream   | new go-ourobopus | Agentic + HelixMemory       | HelixMemory public | — |
| Leda                    | Confirmed upstream   | new go-leda      | Agentic                     | Private? | — |
| P4RS3LT0NGV3            | Confirmed upstream   | new go-p4rs3lt0ngv3 | Formatters               | Public | — |
| ST3GG                   | Confirmed upstream   | new go-st3gg     | HelixAgent security layer   | — | — |
| G0DM0D3                 | Confirmed upstream   | new go-g0dm0d3   | DebateOrchestrator + LLMOrchestrator | LLMOrchestrator public | — |
| HyperTune               | Confirmed upstream   | new go-hypertune | HelixLLM + Benchmark        | Benchmark public? | — |
| GLOSSOPETRAE            | Confirmed upstream   | new go-glossopetrae | Formatters + DocProcessor | DocProcessor public | — |
| Gandalf-Solutions       | **No upstream**      | build-from-scratch OR drop | Challenges + HelixQA | Challenges public? | — |
| Misc-Prompt-Hacks       | Confirmed upstream   | new go-misc-prompthacks | Challenges             | — | — |
| Theseus                 | Confirmed upstream   | new go-theseus   | Agentic + Benchmark         | — | — |
| Gitty                   | Confirmed upstream   | new go-gitty     | BackgroundTasks             | BackgroundTasks status? | — |
| GitGPT                  | **No upstream**      | drop or replace with Gitty extension | BackgroundTasks | — | — |
| V3SP3R                  | Confirmed upstream   | new go-v3sp3r    | MCP-Servers                 | MCP-Servers status? | — |
| Tempest                 | **No upstream**      | drop or build    | HelixLLM                    | — | — |
| LEAKHUB                 | Confirmed upstream   | new go-leakhub   | LLMsVerifier + HelixQA      | — | — |
| OBLITERATUS             | Confirmed upstream   | new go-obliteratus | LLMsVerifier              | — | — |
| AutoStoryGen            | Confirmed upstream   | new go-autostorygen | DocProcessor             | — | — |

## 3. Work breakdown

### 3.1 Days 1–2 — audit

- **Task A:** enumerate every `elder-plinius/*` repo; record language,
  licence, last-commit date, API surface.
- **Task B:** request read access to private HelixDevelopment submodules
  named in §2. Catalogue what exists vs. what must be created.
- **Task C:** for each HelixAgent submodule target, verify the exact
  file path listed in the integration plan matches a real file. Where
  it doesn't, record the nearest equivalent or mark "needs new file".

### 3.2 Days 3–6 — port template + 3 pilot ports

- **Task D:** write a go-port repo template (`.github/` dispatch-only CI,
  `Makefile`, `go.mod`, `LICENSE`, `README`, port-status table).
- **Task E:** port **AUTOTEMP** → `go-autotemp` (lowest risk, pure function).
- **Task F:** port **P4RS3LT0NGV3** → `go-p4rs3lt0ngv3` (read-only, pure
  transforms).
- **Task G:** port **LEAKHUB** → `go-leakhub` (scanner, well-bounded).

Each pilot port ships with: a complete unit-test suite, an integration-test
harness (per-module docker-compose where applicable), zero mocks outside
unit tests, and the port-status row filled in.

### 3.3 Days 7–8 — fictional-module decisions

For each module with **no upstream** (`Gandalf-Solutions`, `GitGPT`,
`Tempest`), produce one of:

- A build-from-scratch spec with estimated effort (and a go-port template).
- A "drop" decision with rationale signed by the integration lead.
- A "rescope" decision mapping the capability to an existing real upstream.

### 3.4 Days 9–10 — revised Phase 1 plan

- **Task H:** rewrite the integration plan's "Phase 1 — read-only" and
  "Phase 2 — passive monitoring" sections with accurate dependency graphs:
  which ports must exist, in which order, before each HelixAgent file is
  touched.
- **Task I:** estimate Phase 1 start date based on remaining port work
  (typically +N weeks where N = number of remaining ports × ~1 week).

## 4. Deliverables

At W0 exit, the integration team has:

1. This document, fully filled in with decisions and owners.
2. 3 pilot go-port repos published.
3. A go-port template repo.
4. Revised `helixagent-plinius-integration.md` with accurate Phase 1
   prerequisites and start date.
5. A go/no-go decision for Phase 1.

## 5. Risks

- **Hidden private submodules.** If HelixAgent submodules listed in the
  plan turn out not to exist even privately, we need to create them
  before integration — further delay.
- **License incompatibility.** Some elder-plinius repos have mixed or
  unclear licenses; legal review gates the port effort.
- **Maintenance burden.** 20+ go-port repos = 20+ repos to keep in sync
  with upstream. Budget ongoing maintenance in W+1 onward.

## 6. Decision log

Track every port/drop/rescope decision with the template:

```
## <Module> (YYYY-MM-DD)
- Decision: port | drop | rescope | build-from-scratch
- Owner: @handle
- Rationale: …
- Tracking issue: …
```

Append to this file below this line as decisions land.
