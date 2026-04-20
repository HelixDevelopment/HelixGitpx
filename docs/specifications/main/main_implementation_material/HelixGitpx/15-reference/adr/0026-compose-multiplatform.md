# ADR-0026 — Compose Multiplatform for cross-platform UI

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

M6 ships mobile (Android/iOS) + desktop (Windows/macOS/Linux) clients. Options: React Native, Flutter, Tauri, Compose Multiplatform.

## Decision

JetBrains Compose Multiplatform 1.7+. Same Kotlin codebase across Android, iOS, Desktop. Shared with the KMP business-logic library under `helixgitpx-clients/shared/`.

## Consequences

- Single UI codebase for 5 targets.
- iOS builds require macOS CI runners (M8 CI task).
- Kotlin-native developers scarcer than TypeScript ones; acceptable for infra project.

## Links

- Spec §LOCKED C-3
