# Script — 09 Extending the OPA bundle

**Track:** Developers · **Length:** 8 min · **Goal:** viewer can add a new Rego rule and ship it safely.

## Body

1. **Where the bundles live** — 0:30 – 1:30.
   `impl/helixgitpx-platform/opa/bundles/v2/`.
2. **The authz.rego model** — 1:30 – 3:00.
   Input shape, `allow { … }` idiom, default-deny.
3. **Add a rule** — 3:00 – 5:00.
   Example: disallow merging PRs into `main` without 2 approvals in org `acme`.
4. **Write a Rego test** — 5:00 – 6:30.
   `_test.rego` file; table-driven.
5. **Diff-review CI gate** — 6:30 – 7:30.
   The workflow that compares old vs new decisions over a corpus of inputs.

## Wrap-up (7:30 – 8:00)
"Every rule is policy. Every policy is code. Every code review is a policy review."

## Companion doc
`impl/helixgitpx-platform/opa/bundles/v2/` · ADR-0033
