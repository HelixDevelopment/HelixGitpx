# ADR-NNNN: [Short noun phrase describing the decision]

**Status**: Proposed | Accepted | Superseded by ADR-MMMM | Deprecated | Rejected
**Date**: YYYY-MM-DD
**Deciders**: @github-handle-1, @github-handle-2
**Technical story**: [Link to originating issue / PR / RFC]
**Consulted**: [Teams or individuals consulted but not primary deciders]
**Informed**: [Teams or individuals notified of the decision]

---

## Context and Problem Statement

Two to five sentences. What is happening in the system that forces a decision? What constraints exist? What outcomes are we trying to produce or avoid? Be specific — a year from now a new engineer should be able to read this paragraph and understand the problem without prior context.

## Decision Drivers

- Driver 1 (e.g., "must scale to 1B events/day by Q4")
- Driver 2 (e.g., "license must be permissive, no SSPL")
- Driver 3
- …

## Considered Options

- **Option A** — one-line description
- **Option B** — one-line description
- **Option C** — one-line description

## Decision

**We will adopt Option X because Y.**

Concise, unambiguous. If the decision covers multiple related sub-choices, list them as bullets. Link to this ADR from the code that implements it.

## Consequences

### Positive

- Specific expected improvement (include metric if possible, e.g., "reduces DB p99 latency from ~80ms to ~15ms per our bench").
- …

### Negative / Trade-offs

- …
- …

### Risks & Mitigations

- Risk: … → Mitigation: …

### Cost Impact

Brief: compute / storage / licensing / operational delta estimated in dollars or percentages. Include capex vs. opex if relevant. If the cost impact is negligible say so.

## Pros & Cons of the Options

### Option A

- ✓ Pro …
- ✓ Pro …
- ✗ Con …

### Option B

- ✓ Pro …
- ✗ Con …
- ✗ Con …

### Option C

- ✓ Pro …
- ✗ Con …

## Decision Outcome

Why Option X won over the alternatives. Address the top trade-off honestly — what do we lose, and why is it acceptable?

## Implementation Notes

- Who's implementing (team / owner).
- Rough timeline + phases.
- Link to tracking issues / epics.
- Migration strategy (if applicable).
- Feature flag strategy.
- Rollback plan.

## Validation & Exit Criteria

How will we know the decision was correct?

- Metric / SLO improvements.
- Customer-visible outcomes.
- Internal developer-experience improvements.
- Date for a retrospective review (recommended: 90 / 180 days post-ship).

## Observability Hooks

- Dashboards to watch.
- Alerts to add / tune.
- Success metrics to track.

## Security / Privacy Considerations

- Data-classification changes.
- New attack surface.
- Compliance impact (SOC 2 / ISO 27001 controls affected).
- Threat-model delta.

## Accessibility Considerations

- Does this affect any user-facing surface? If yes, what's the WCAG impact?

## Internationalisation Considerations

- New user-visible strings?
- Impact on RTL layouts?

## Documentation Impact

- Docs that need updating (link each).
- Runbook changes.
- Customer-facing announcement (yes/no).

## References

- Related ADRs: ADR-0001, ADR-0002
- External docs / blog posts / RFCs: …
- Prior art within the project: …
- Decision record standard: <https://adr.github.io/>

## Change Log

*(Status-only changes. Do not rewrite content after Acceptance — supersede instead.)*

- YYYY-MM-DD — Proposed
- YYYY-MM-DD — Accepted
- YYYY-MM-DD — Superseded by ADR-MMMM

---

## Writing Tips

- **Write for future you.** Assume you've forgotten the context. Include it.
- **Be concrete.** "Scale" is fine; "10× current throughput = 200k events/s sustained" is better.
- **Be honest about trade-offs.** If this was close, say so. If we're betting on a recent technology, acknowledge the risk.
- **Don't advocate — explain.** The reader should be able to disagree and still understand why we chose what we did.
- **Short > long.** Aim for 1-2 pages. Link to deeper docs where helpful.
- **One decision per ADR.** If your ADR covers two decisions, split it.
- **Name what you're NOT deciding.** Explicitly excluded from scope.
