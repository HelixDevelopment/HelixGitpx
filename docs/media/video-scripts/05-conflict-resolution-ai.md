# Script — 05 Resolving a conflict with AI

**Track:** Getting started · **Length:** 5 min · **Goal:** viewer resolves a real divergence using the AI proposal.

## Cold open
Terminal: two pushes race. One lands on GitHub; the other on GitLab via a broken webhook. Divergence.

## Body

1. **Divergence detected** — 0:30 – 1:15.
   Conflict inbox shows both refs. Both have commit graphs visible.
2. **Ask for an AI proposal** — 1:15 – 2:30.
   Click "Propose". ai-service calls Ollama; response in ~10s.
3. **Read the diff** — 2:30 – 3:30.
   The proposal shows merged tree + rationale. Highlight the areas you should eyeball.
4. **Accept** — 3:30 – 4:15.
   Click accept. conflict-resolver writes; mirror refreshes both hosts.
5. **Audit trail** — 4:15 – 4:45.
   Who, when, what LLM, what policy allowed it.

## Wrap-up (4:45 – 5:00)
"AI proposes, humans decide."

## Companion doc
`docs.helixgitpx.io/concepts/conflicts`
