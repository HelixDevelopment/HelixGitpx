# 09 — Conflict Resolution Engine

> **Document purpose**: Specify the **end-to-end conflict detection and resolution pipeline**. This is HelixGitpx's core innovation — without it, multi-upstream federation is a toy. With it, we can guarantee **zero data loss** and **eventually-consistent, human-intelligible merges** across N upstreams.

---

## 1. What "Conflict" Means Here

In classic single-remote Git, "conflict" means merge conflict in working tree. In HelixGitpx, it means **any observable divergence between two or more upstream views of the same logical object**, namely:

1. **Ref divergence** — `refs/heads/main` is `abc` on our side, `def` on GitHub, `ghi` on GitLab, with no fast-forward relation.
2. **Rename collision** — file renamed to `A/foo.go` on one side, `B/foo.go` on another.
3. **Concurrent metadata edits** — someone on Gitee removes label `bug`; within the window, someone on GitHub adds label `enhancement` to the same issue.
4. **PR/MR state divergence** — PR marked merged on GitHub, still open on the mirror.
5. **Tag collision** — `v1.2.3` points to commit `X` on GitHub, `Y` on Codeberg.
6. **LFS object divergence** — same OID, different blob bytes (corruption or race).
7. **Workflow status divergence** — CI says success on one, failed on another (merged via attestations).
8. **Release asset divergence** — same tag, different assets attached.

---

## 2. Architectural Placement

```
Inbound webhook  ─► Webhook Gateway ─► Kafka (upstream.ref.received)
                                                │
                              Sync Orchestrator (Temporal)
                                         │
                                  Conflict Detector
                                         │ (detects divergence)
                                         ▼
                              conflict-resolver service
                          ┌──────────────┬──────────────────┐
                          ▼              ▼                  ▼
                       Policy        CRDT merger        AI resolver
                       engine                              │
                                                           ▼
                                            Low confidence → escalate (human)
                                            High confidence → apply + emit
```

---

## 3. The Three-Phase Pipeline

### Phase 1 — **Detect**

- Streaming join on Kafka (Kafka Streams via Goka) of our view vs. the incoming upstream event.
- Co-partitioned by `(repo_id, ref_name)` (or `(issue_id)` for metadata).
- Emits a `conflict.detected` event with enough context for resolution (three-way: local, remote, common ancestor).

### Phase 2 — **Propose**

Resolver generates one or more **proposals**:

- **Deterministic policy**: based on repo/org configured strategy, produces a concrete resolution.
- **CRDT merge** (for metadata): always succeeds without loss.
- **AI proposal**: LLM returns a candidate + rationale + confidence.

Each proposal has:

```yaml
proposal:
  id: 018f...
  strategy: take_newer | prefer_primary | three_way_merge | crdt_merge | ai
  confidence: 0.95
  apply_plan:
    - op: create_tmp_branch
      name: conflict/018f...
    - op: three_way_merge
      theirs: def
      ours:   abc
      base:   xyz
    - op: push_to_upstream
      upstream_ids: [github, gitlab, gitee]
  rationale: "Text rationale for humans + audit."
```

### Phase 3 — **Apply or Escalate**

- If a proposal's strategy is policy-safe (e.g., `prefer_primary`) **or** AI confidence ≥ threshold (configurable per org, default 0.92) **and** org policy allows auto-apply → execute the plan.
- Otherwise emit `conflict.escalated` and wait for human signoff via web/mobile UI.
- Human approves → same apply plan executes.
- All outcomes produce signed audit records.

---

## 4. Conflict Classes — Resolution Details

### 4.1 Ref Divergence

Detected when our new head SHA and the upstream-reported new head SHA are unrelated (no fast-forward).

**Resolution strategies** (configurable, in order of preference):

| Strategy | Description | Safe? |
|---|---|---|
| `prefer_primary` | If repo has a configured primary upstream and this change came from it, accept | Yes |
| `prefer_signed` | Accept the side whose head is signed (verified signature) | Yes |
| `prefer_newer` | Accept whichever head has the most recent committer timestamp — with a 30 s skew guard | Mostly (clock risk) |
| `three_way_merge` | Find common ancestor, attempt merge; if clean, publish result everywhere | Yes if clean |
| `octopus_merge` | Multi-way merge for 3+ upstreams that all diverged | Risky; only for unambiguous trees |
| `ask_ai` | LLM analyses both tips + diff + commit messages; proposes merge or choice | Needs signoff |
| `human_only` | Never auto-resolve; always escalate | Safest |

Default org policy: `prefer_primary` → `three_way_merge` → `ask_ai` → `human_only`.

### 4.2 Rename Collision

Two upstreams renamed the same file to different paths.

- Build a 3-way tree (common ancestor + both sides).
- If one rename is to a path identical to another side's original (a renaming cycle), flag `rename_ambiguous`.
- Policy options:
  - `prefer_primary`
  - `keep_both` — resolver creates two files (suffixing `.left` / `.right`) and opens an issue.
  - `ask_human` — always.

### 4.3 Concurrent Metadata Edits (labels, milestones, assignees, issue body)

**This is where CRDTs shine.**

- Every issue body is stored as an **Automerge** document (`collab.crdt_docs`).
- Labels set is a **G-Set** (grow-only) or **OR-Set** if we allow removes.
- Milestones + assignees use **LWW-Register** keyed per field.
- Upstream events are converted into CRDT ops and merged commutatively.
- Resulting doc is replayed on every upstream.

Guarantees:
- No edit is lost.
- Merge is deterministic and order-independent.
- Eventual convergence across upstreams.

### 4.4 PR / MR State Divergence

The authoritative rule: **"once merged, always merged"**. If one upstream says merged, we fan-out merge (or force the merge commit) to others. Conflict cases:

- Both closed without merge on different sides → take the latest close; record "conflicting close" in audit.
- One merged via squash, other via rebase → we record the *effective* merge commit in HelixGitpx's canonical view; upstreams stay as-is (no re-merge).

### 4.5 Tag Collision

Tags with the same name pointing to different commits. Policy:

- Force to primary tag; emit warning event for non-primary.
- For signed tags, prefer the valid signature.
- If both valid and both signed: escalate.

### 4.6 LFS Object Divergence

Identical OID, different bytes means corruption or race (extremely rare). Quarantine the divergent blob; flag for manual review. Never silently pick one.

### 4.7 Workflow Status Divergence

CI statuses from different providers: we aggregate **all** statuses; a status check is considered "passing" iff **all contributing upstreams pass**. Branch protection aggregates across upstreams.

### 4.8 Release Asset Divergence

- Canonical release holds the union of asset names + SHA256.
- Assets with the same name + different SHA256 → escalate.
- We never overwrite a released asset.

---

## 5. Policy-as-Code

Every strategy is selectable per (org, repo) via OPA / Rego:

```rego
package helixgitpx.conflict

default strategy = "human_only"

strategy = "prefer_primary" {
  input.repo.primary_upstream != ""
  input.event.upstream_id == input.repo.primary_upstream
}

strategy = "three_way_merge" {
  input.conflict.kind == "ref_divergence"
  can_ff_or_merge(input.conflict)
}

strategy = "crdt_merge" {
  input.conflict.kind == "metadata_concurrent"
}

strategy = "ai_proposal" {
  input.org.settings.allow_ai_autoapply == true
  input.proposal.confidence >= input.org.settings.ai_confidence_threshold
}
```

Policies live in Git; changes deployed via Argo CD; versioned & signed; evaluated by the policy-service.

---

## 6. AI-Assisted Resolution

When no deterministic strategy applies, we ask a fine-tuned LLM.

### 6.1 Input

```yaml
prompt_template: conflict_resolution.v2
variables:
  repo_language: Go
  file_path: pkg/store/redis.go
  left_diff:   |   (unified diff of our side vs. base)
  right_diff:  |   (unified diff of their side vs. base)
  commit_msg_left:  "…"
  commit_msg_right: "…"
  recent_reviews: [...]
  style_guide: <retrieved via RAG>
  semantic_context: <top-k neighbours from Qdrant>
```

### 6.2 Output Schema (enforced)

```json
{
  "strategy": "accept_merge" | "prefer_left" | "prefer_right" | "custom_patch" | "human_required",
  "confidence": 0.87,
  "patch": "<<<unified diff if strategy=custom_patch>>>",
  "rationale": "string",
  "risks": ["breaks go vet", "removes test case foo_test.go"]
}
```

### 6.3 Guardrails

- **Structured output** forced via grammar-constrained decoding (Guardrails / NeMo).
- Patch is applied in a **sandbox** (ephemeral clone, no network) and run through:
  - `go vet` / language-specific linter.
  - Affected tests (if short).
  - Gitleaks (no secrets).
- If sandbox fails: confidence clipped to 0.
- Final decision always subject to policy (§5).

### 6.4 Feedback Loop

User acceptance / rejection / edit is captured in `ai.feedback_records` and feeds the RLAIF training pipeline (see [10-llm-self-learning.md](../07-ai/10-llm-self-learning.md)).

---

## 7. Apply Plan Executor

Apply plans are sequences of **atomic operations**:

- `create_tmp_branch`
- `three_way_merge`
- `apply_patch`
- `crdt_apply_ops`
- `push_to_upstream` (via Temporal activity to adapter-pool)
- `update_pr_status`
- `open_issue` (e.g. to record quarantine)
- `send_notification`

Executed as a Temporal workflow with:
- Per-op timeout and retry.
- Rollback steps (not all ops are reversible; rollback semantics documented per op).
- Signed audit record per op.

---

## 8. UI / UX

- **Inbox**: "Conflicts awaiting you" per user, with summary, affected files/refs.
- **Side-by-side diff**: left, right, base, proposed merge.
- **One-click accept / edit / reject** for AI proposals.
- **Comment thread** on each case (mirrored to the upstream PR or a dedicated `__helix_conflicts__` issue).
- **Timeline view** showing every conflict for a repo over time.

---

## 9. Data Model

See [03-data/04-data-model.md §4.6](../03-data/04-data-model.md).

- `conflict.conflict_cases`
- `conflict.resolutions`
- `conflict.ai_feedback`
- `collab.crdt_docs`

---

## 10. Observability

| Metric | Meaning |
|---|---|
| `helixgitpx_conflicts_detected_total{kind}` | Rate by kind |
| `helixgitpx_conflicts_resolved_total{kind,strategy,decided_by}` | Resolution outcomes |
| `helixgitpx_conflicts_escalation_rate` | % escalated to humans |
| `helixgitpx_conflicts_auto_resolution_rate` | Target ≥ 75 % |
| `helixgitpx_conflicts_time_to_resolve_seconds{kind}` | Histogram |
| `helixgitpx_ai_proposal_confidence{task,decided_by}` | Distribution |
| `helixgitpx_ai_accept_rate{task,model}` | Training quality signal |

Alerts:
- `ConflictQueueGrowing` (p99 unresolved age > 6 h).
- `ConflictEscalationSpike`.
- `ConflictAutoRateDrop` (below 60 % for 1 h — possible regression).

---

## 11. Failure Modes & Safety

| Failure | Mitigation |
|---|---|
| Adapter unreachable during push | Retry with exponential backoff; if quorum of upstreams unreachable, pause apply |
| LLM returns unsafe patch | Sandbox + lint catches; policy rejects |
| Policy bug auto-applies wrong thing | Every auto-apply signed + time-bounded undo window (default 5 min) — during which a `cancel` command rolls back |
| Runaway merge loop | De-bounce: same case cannot be re-processed within 30 s unless a new input event arrives |
| Corruption of event log | Rebuild from other upstreams' tips + manual reconciliation runbook |

### 11.1 Undo Window

Every auto-applied resolution is stored with **undo hooks**: reverse operations computed at apply time. The first 5 minutes after apply, a single click rolls it back on all upstreams. After that, a normal "reverse-merge" is required.

---

## 12. Testing

- **Unit**: every strategy with 100+ table-driven cases.
- **Property-based**: randomly generate three-way trees; assert invariants (no data loss, deterministic output, idempotent apply).
- **Fuzz**: malformed diffs, adversarial patches.
- **Mutation**: mutate resolver logic; tests must catch ≥ 95 % of mutants.
- **Replay**: historical conflicts from test corpus (curated from open-source repos with multi-host mirrors) must resolve to the same outcome our baseline produced.
- **Shadow**: new resolver versions run in shadow vs. production, emit divergence metrics without applying.

---

## 13. Innovations Vs v4.0.0 Spec

- **CRDT for metadata** (truly loss-free concurrent edits).
- **Sandboxed patch validation** before AI-accepted patches are applied.
- **Time-bounded undo window**.
- **Shadow resolver runs** for regression-proofing new versions.
- **RLAIF feedback loop** (see AI doc).

---

*— End of Conflict Resolution Engine —*
