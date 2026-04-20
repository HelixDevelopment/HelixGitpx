# Script — 16 OPA policy-as-code end-to-end

**Track:** Security & compliance · **Length:** 12 min · **Goal:** viewer understands every enforcement point and how they compose.

## Body

1. **Why policy-as-code** — 0:30 – 1:30.
   Auditable, versioned, testable, rollbackable.
2. **Enforcement points** — 1:30 – 3:30.
   Auth layer, repo service, upstream bindings, conflict resolution, AI responses, admission control.
3. **The v2 bundle** — 3:30 – 5:30.
   `impl/helixgitpx-platform/opa/bundles/v2/enforcement.rego` walk.
4. **Bundle server** — 5:30 – 7:00.
   `opa-bundle-server` ETag + polling; rollback = pointer flip.
5. **CI diff-review gate** — 7:00 – 9:00.
   Run policy corpus against old + new bundle; surface decision diffs.
6. **Kyverno admission** — 9:00 – 10:30.
   Cluster-side policy for YAML.
7. **Feedback loops** — 10:30 – 11:45.
   A denied action produces a friendly error + decision ID for support.

## Wrap-up (11:45 – 12:00)
"Every no has a reason. Every reason has a Rego line number."

## Companion doc
ADR-0033 · `docs/security/soc2-type1-evidence-index.md`
