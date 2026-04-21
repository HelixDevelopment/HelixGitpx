## 4. Resolving conflicts

Federation produces conflicts. Two teams push to different upstreams, a
label lands on one host and not another, a rename clashes — HelixGitpx
surfaces every divergence as a reviewable conflict.

### 4.1 Kinds of conflict

| Kind | What it means | Where it comes from |
|------|---------------|---------------------|
| Ref divergence | Two upstreams have different tips for the same ref | Webhook race, partial fan-out failure |
| Label race | Same PR gained different labels on different hosts | Concurrent edits |
| Rename collision | Two rename operations target the same path | Human concurrent edits |
| Metadata drift | PR description/assignees/milestones differ | Host-specific edit that wasn't mirrored |

### 4.2 The inbox

The **Conflicts** page lists every open conflict for repos you have
access to. Each row shows:

- Which kind.
- Which upstreams are involved.
- When it was detected.
- Whether an AI proposal has been generated.

### 4.3 Getting an AI proposal

Click **Propose resolution**. `ai-service` consults the active
LLM (Ollama locally, or an external provider per org policy) via
`LiteLLM`. NeMo Guardrails enforce output shape and refuse obviously
dangerous patches.

A proposal includes:

- A rationale (plain-English explanation).
- A unified diff (for ref divergence) or an Automerge change (for
  metadata).
- A confidence score.

### 4.4 Human review gate

Always required. The proposal cannot apply itself. You read the
rationale, inspect the diff, and either:

- **Accept** — conflict-resolver writes the merge to the authoritative
  state and re-fans-out.
- **Reject** — the proposal is discarded. Another can be generated.
- **Edit** — adjust the diff and apply as if it were yours.

### 4.5 Audit trail

Every conflict, every proposal, every human decision lands as an
`audit.events` entry with the LLM identity, the prompt ID, the patch
hash, and the decision-maker. Org admins can export the full trail for
compliance evidence.

### 4.6 When there's no AI

For orgs that disable AI, all conflicts go straight to human review.
The inbox is the same; the **Propose** button is replaced by a
**Diff** button that shows you the divergence directly.

---
