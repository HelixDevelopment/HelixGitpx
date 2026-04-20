# M6 Frontend & Mobile — Design Spec

| Field | Value |
|---|---|
| Status | APPROVED |
| Milestone | M6 — Frontend & Mobile (Weeks 29–36) |
| Scope | Roadmap §7 items 93-115 (23 items) |

## Locked constraints

- **C-1** Web: Nx + Angular 19 (established in M1). State via NgRx Signals. Connect-Web clients via `@connectrpc/connect-web` (already in M3).
- **C-2** Shared KMP library under `impl/helixgitpx-clients/shared/` with SQLDelight for offline store, kotlinx.serialization + Ktor for network, kotlinx.coroutines.
- **C-3** UI: Compose Multiplatform (JetBrains 1.7+) for shared screens across Android/iOS/Win/macOS/Linux.
- **C-4** Distribution: Play Store + F-Droid + App Store + TestFlight + MSIX + DMG + AppImage + .deb + .rpm; auto-update via self-hosted update feed.
- **C-5** i18n: `@ngx-translate/core` for web; Moko Resources for KMP. 8 locales shipped (en, de, fr, es, pt, ja, zh-CN, ru).
- **C-6** PWA + offline shell via Angular Service Worker.

## New apps

- `helixgitpx-web/apps/web/` — full Angular shell (dashboard, repos, PRs, issues, conflicts inbox, settings, search).
- `helixgitpx-clients/androidApp/`, `iosApp/`, `desktopApp/` — Compose Multiplatform shells.

## 23-item matrix abbreviated

Items 93-102 (web features, i18n, a11y, PWA), 103-106 (KMP shared), 107-111 (Compose UI + platform features), 112-115 (distribution artefacts).

## ADRs 0026-0029

- 0026 — Compose Multiplatform for cross-platform UI (not React Native)
- 0027 — NgRx Signals over Redux for web state
- 0028 — Moko Resources for KMP i18n
- 0029 — Self-hosted update feed via tus.io

— End of M6 design —
