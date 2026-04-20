# ADR-0027 — NgRx Signals over Redux for web state

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Angular 19 ships native Signals; NgRx Signals wraps them with store ergonomics. Redux-style NgRx (classic) adds action/reducer boilerplate.

## Decision

Signals-first stores; classic NgRx reserved for complex workflows (conflict resolution inbox in M5/M6).

## Consequences

- Less boilerplate per screen.
- Mental model aligned with modern Angular.

## Links

- Spec §LOCKED C-1
