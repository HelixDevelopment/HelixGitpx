# ADR-0028 — Moko Resources for KMP i18n

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

KMP apps need a shared string table across Android/iOS/Desktop targets. Options: per-platform (Android string.xml + iOS .strings + desktop .properties), kotlinx-i18n (young), Moko Resources.

## Decision

Moko Resources. Central strings file in the shared module; generates per-target bindings.

## Consequences

- One source of truth for translation keys.
- Moko's type-safe bindings catch missing translations at compile time.

## Links

- Spec §LOCKED C-5
- https://github.com/icerockdev/moko-resources
