# RB-120 — Conflict Backlog Growing

> **Severity default**: P2
> **Owner**: Conflict / AI team
> **Last tested**: 2026-03-15

## 1. Detection

Alert: `ConflictQueueGrowing` — p99 time-to-resolve > 6 h for 30 min.

Supporting signals:
- `helixgitpx_conflicts_open{status="escalated"}` climbing.
- `helixgitpx_conflicts_auto_resolution_ratio` < 0.75 target.
- Customer reports: "conflicts piling up in my inbox".

---

## 2. First Moves

1. Ack; triage scope: one customer or systemic?
2. Dashboard "Conflict Engine — Overview" shows backlog by kind, severity, and strategy.
3. Identify: is resolution *slow* (processing taking long) or *blocked* (escalated waiting humans)?

---

## 3. Diagnose

### 3.1 Auto-resolution ratio dropping

- Compare recent commit / deploy to AI service or policy bundle.
- Check `helixgitpx_ai_accept_rate` — model regression?
- Check `helixgitpx_policy_decisions_total{effect="deny"}` — policy tightened?

### 3.2 AI timing out

- Most likely AI capacity issue → see RB-150.

### 3.3 Sandbox failures

- Look at `conflict.resolutions.apply_status` distribution.
- If many `sandbox_failed`: build tools regression, or environment drift.

### 3.4 Specific kind saturated

- One kind dominating (e.g. `rename_collision`) may need policy tuning.

---

## 4. Mitigations

### 4.1 Scale conflict-resolver

```bash
helixctl scale conflict-resolver --replicas=<2x>
```

### 4.2 Switch problematic kinds to human

```bash
# Temporarily bypass AI for rename collisions; escalate directly
helixctl conflict set-strategy --kind=rename_collision --strategy=escalate_to_human --for=24h
```

Inbox grows but quality is preserved.

### 4.3 Rollback policy change

If recent policy deploy tightened auto-apply:

```bash
helixctl policy rollback --bundle=conflict --to=<previous-revision>
```

### 4.4 Batch-apply high-confidence cases

If confidence distribution shows many cases just below threshold temporarily (e.g. 0.9 when threshold is 0.92):

```bash
# Coordinated override — audited
helixctl conflict batch-apply --confidence-min=0.90 --dry-run=false --reason=...
```

Very cautious; requires two-person approval.

### 4.5 Bring AI team online

For persistent backlog, involve AI team to evaluate whether a model roll back or new training run is needed.

---

## 5. Verify Recovery

- `conflicts_auto_resolution_ratio` ≥ 0.75.
- p99 time-to-resolve back < 1 h.
- Open-conflicts count trending down.
- No new complaints.

---

## 6. Post-Incident

- If AI cause: feed failed cases into training + eval.
- If policy cause: post-mortem on change management.
- If customer-specific: CSM outreach.

---

## 7. Related

- RB-121 (AI confidence regression)
- RB-150 (AI inference timeout)
- RB-301 (Accidental auto-apply)
- [09-conflict-resolution.md]
