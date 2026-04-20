# RB-301 — Accidental Auto-Apply of AI Conflict Resolution

> **Severity default**: P2 (escalates to P1 if affects primary branch on Enterprise account)
> **Owner**: Conflict / AI team
> **Last tested**: 2026-03-05

This runbook handles the scenario where an AI-proposed conflict resolution was auto-applied and the customer wants it reversed — **or** where our monitoring detects an apply went wrong.

## 1. Detection

- Customer report: "Something got merged/resolved that I didn't expect."
- Monitoring: `ConflictUndoRequest` metric spike.
- Internal: discrepancy between `conflict_cases.status=applied` and known-good state (e.g. failing CI on upstream after auto-apply).

---

## 2. Understand the Window

HelixGitpx grants a **5-minute undo window** after auto-apply. Behaviour differs depending on elapsed time:

| Elapsed | Strategy |
|---|---|
| < 5 min | One-click undo rolls back across all upstreams (pre-planned) |
| 5 min – 1 hour | Supported undo — reverse-merge PRs opened on all upstreams; manual action from maintainer |
| > 1 hour | Full manual: treat as any other reversion; history is now real |

---

## 3. First Moves

1. Ack the customer; open internal ticket with `conflict_case_id`.
2. Find the case:
   ```bash
   helixctl conflict show <case_id>
   ```
3. Note `applied_at`, `undo_until`, `apply_plan`, and `strategy`.
4. If within 5 min → run `conflict undo`.
5. Outside window → open reverse plan.

---

## 4. Undo Within Window

```bash
helixctl conflict undo --case=<case_id> --reason="customer request"
```

This:
- Reverts the refs to their `snapshot_ref` values.
- Reverses any CRDT operations applied.
- Emits reverse-fan-out to every upstream.
- Marks `conflict_cases.status=detected` (re-opens case).
- Emits `conflict.undone` event → notifies assignees.

Verify:

```bash
helixctl conflict show <case_id>   # status should be detected / escalated
helixctl sync history <repo> --limit=5   # latest reverse-fan-out visible
```

---

## 5. Outside the Window (Manual Reverse)

### 5.1 Generate reverse plan

```bash
helixctl conflict reverse-plan --case=<case_id> --out=reverse.json
```

`reverse.json` contains ops that, if applied, bring the system back to the pre-apply state. It's a best-effort computation; review before applying.

### 5.2 Review with the customer

- Walk through the plan with the repo maintainer.
- Confirm no **subsequent** changes they want to keep will be undone (if any exist, curated reverse plan needed).

### 5.3 Apply reversal

Either:

- Open reverse PRs on each upstream (`helixctl conflict reverse-apply --as-prs`) so the customer merges after review.
- Or apply directly with `--force` (require second operator co-sign + audit log reason).

### 5.4 Notify

- Customers whose PRs/issues were affected receive an automated note.
- Upstream CI re-runs.

---

## 6. Root-Cause Investigation

Always happens for P2+. Ask:

1. Did policy authorise auto-apply? Consult `policy.decisions` for the case.
2. Was AI confidence at or above threshold? `helixctl conflict proposals --case=<id>`.
3. Did sandbox validation pass? Check `conflict.resolutions.apply_status` trail.
4. Was there a model regression? Check `helixgitpx_ai_accept_rate` trend.
5. Was the test suite insufficient? Coverage gap on this scenario?

Common patterns:

- **Edge-case the AI didn't see** — add to golden test set; lower auto-apply confidence threshold for similar cases.
- **Customer expectation mismatch** — tighten default auto-apply scope in their org; expose toggle in settings.
- **Shadow-mode mismatch** — new model promoted too quickly; review promotion criteria.

---

## 7. Mitigations (Short Term)

If the root cause is systemic (not a one-off customer expectation mismatch):

```bash
# Temporarily disable auto-apply globally for the affected kind
helixctl flag set ai.conflict-auto-apply --off --kind=rename_collision

# Or shrink the confidence threshold
helixctl ai set-threshold --task=conflict_proposal --threshold=0.97
```

Announce on status page if systemic.

---

## 8. Customer Communication

**First response (internal or external)**:

> We found the conflict case. Because this is within the 5-minute undo window, we're rolling the change back now. Please refresh your repo in ~30 seconds.

**Outside window**:

> The 5-minute auto-undo window has passed, so we'll prepare reverse-pull requests on each connected Git host. One of our engineers will walk through them with you before they're merged. Estimated completion: within 1 hour.

---

## 9. Post-Incident

- Record the case id + resolution in the case store.
- Add a regression test (real fixture) so the same mistake can't ship.
- If trust was shaken, review whether:
  - Customer should be switched to **always require human review** globally.
  - Their org needs a dedicated model configuration.

---

## 10. Related

- RB-120 (Conflict backlog growing)
- RB-121 (AI confidence regression)
- [09-conflict-resolution.md] — strategy ladder
- [07-ai/10-llm-self-learning.md] — model lifecycle
- Policy: [15-reference/policies/conflict-resolution.rego]
