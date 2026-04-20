# 13 — Mobile & Desktop (Kotlin Multiplatform + Compose)

> **Document purpose**: Define the **cross-platform client architecture** for HelixGitpx. One codebase produces native clients for **Android, iOS, Windows, macOS, and Linux** with shared UI (Compose Multiplatform), shared state, and shared networking. ADR-0006 documents the choice.

---

## 1. Why Kotlin Multiplatform

| Candidate | Reach | Shared UI? | Native look? | Native perf? | Verdict |
|---|---|---|---|---|---|
| **KMP + Compose Multiplatform** | Android, iOS, Win, macOS, Linux | **Yes** | Close (Compose) | Excellent | **Chosen** |
| Flutter | Same + Web | Yes | Custom canvas | Excellent | Dart ecosystem further from our Go backend & gRPC stubs |
| React Native | Android, iOS | No (platform UI) | Native | Good | Web/mobile fragmentation, less desktop |
| Native per-platform | Best fidelity | No | Perfect | Best | 5× cost, 5× risk |

KMP lets us share ~85 % of code (networking, state, navigation, feature UIs) while dropping to native for the last 15 % (camera, file pickers, notifications, auth) via expect/actual. **Compose Multiplatform** (JetBrains) provides a unified UI for Android + iOS + desktop.

---

## 2. Module Layout

```
/clients-kmp
├── settings.gradle.kts
├── build.gradle.kts
├── gradle/libs.versions.toml       (catalog)
│
├── shared/
│   ├── core/
│   │   ├── commonMain/             (pure Kotlin, no UI)
│   │   │   ├── model/
│   │   │   ├── util/
│   │   │   └── di/
│   │   ├── androidMain/
│   │   ├── iosMain/
│   │   ├── desktopMain/            (JVM)
│   │   └── linuxMain / macosMain / mingwMain (where needed, usually via desktopMain)
│   │
│   ├── data/
│   │   ├── commonMain/
│   │   │   ├── grpc/               (generated Connect/gRPC stubs via kroto+)
│   │   │   ├── repositories/
│   │   │   ├── cache/              (SQLDelight schemas + queries)
│   │   │   ├── crypto/             (libsodium bindings via KMP)
│   │   │   └── events/             (live-events multiplexer)
│   │   ├── androidMain/ …
│   │
│   ├── domain/
│   │   └── commonMain/             (use-cases, pure functions)
│   │
│   └── ui/
│       ├── commonMain/             (Compose screens & components)
│       ├── androidMain/
│       ├── iosMain/
│       └── desktopMain/
│
├── apps/
│   ├── androidApp/                 (final Android APK/AAB)
│   ├── iosApp/                     (Xcode project referencing iosMain framework)
│   ├── desktopApp/                 (Compose desktop, packaged DMG/MSI/DEB)
│   └── cliApp/                     (optional; Kotlin/Native CLI for power users)
│
└── tools/
    ├── grpc-gen/                   (Gradle task wrapping buf generate)
    └── localization/               (XLIFF <-> Android XML <-> iOS .strings)
```

---

## 3. Technology Choices

| Concern | Library |
|---|---|
| UI | **Compose Multiplatform 1.7+** (Android/iOS/Desktop) |
| Navigation | **Decompose** (component-based) or **Jetpack Navigation Compose** (Android/Desktop) + Voyager |
| State | **Kotlin Flows** + **MVIKotlin** / custom stores; **Molecule** for reactive composition |
| DI | **Koin** (multiplatform) |
| Networking | **Ktor Client** + **Connect-Kotlin** (gRPC + gRPC-Web) + **okio** |
| Serialisation | `kotlinx.serialization` + protobuf |
| Persistence | **SQLDelight** for structured cache; **multiplatform-settings** for prefs |
| Crypto | libsodium via **Diglol** KMP bindings |
| Date/time | `kotlinx-datetime` |
| Logging | **Kermit** |
| Crash | Sentry multiplatform / Firebase Crashlytics (Android) / Bugsnag (others) |
| Analytics | Matomo (self-hosted) |
| Image | Coil (Android), Kamel (others) |
| Markdown | multiplatform-markdown-renderer |
| Code highlighting | `compose-highlights` (fork of Chroma rules) |
| Background | WorkManager (Android), BGAppRefresh (iOS), coroutines+tray (desktop) |
| Push | FCM (Android), APNs (iOS), UNNotification on mac, Web-Push on Linux/Windows via desktop app |
| DB migrations | SQLDelight built-in |
| Testing | Kotest + Turbine + Paparazzi (screenshot) + Maestro (E2E) |

---

## 4. Shared Architecture

```
          ┌──────────────────────────────┐
          │  UI (Compose Multiplatform)  │
          └──────────────┬───────────────┘
                         │ State / Intents
          ┌──────────────▼───────────────┐
          │  Presentation (MVIKotlin)    │
          └──────────────┬───────────────┘
                         │ Use cases
          ┌──────────────▼───────────────┐
          │        Domain (pure)         │
          └──────────────┬───────────────┘
                         │ Repositories
          ┌──────────────▼───────────────┐
          │   Data (gRPC / Cache / FS)   │
          └──────────────┬───────────────┘
                         ▼
                  Remote HelixGitpx API
```

- UI is dumb; reads state, emits intents.
- Domain is pure — zero platform code.
- Data layer owns networking, cache, and event fan-out.
- Repositories expose `Flow<T>` that combines cache + live events + pagination.

---

## 5. Live Events on Mobile/Desktop

- Ktor + Connect-Kotlin opens a gRPC server-streaming RPC to `EventsService.Subscribe`.
- On transport failure, falls back to WebSocket automatically.
- `resume_token` persisted in SQLDelight keyed by user — survives app restart.
- Events applied to SQLDelight cache tables; queries are reactive → UI updates automatically.
- Battery saver: on Android, when `doze` active, cut stream and rely on FCM data messages; reconnect on wake.

---

## 6. Offline-First

- All data read from local cache; fetches are background reconciliations.
- Writes stored in an **operation log** (`outbox` table) with an `idempotency_key`.
- Background sync worker pushes the outbox whenever online; on success, marks processed.
- Conflicts (server rejects with 409 or 412) surface in a UI "Sync issues" panel with options: retry, discard, override.

---

## 7. Auth & Security on Device

- OAuth2 PKCE via in-app browser (Chrome Custom Tabs / ASWebAuth / default browser).
- Tokens stored in platform secure storage:
  - Android: EncryptedSharedPreferences + hardware-backed Keystore.
  - iOS: Keychain with `.whenUnlockedThisDeviceOnly`.
  - macOS: Keychain.
  - Windows: Data Protection API via JNA.
  - Linux: libsecret (GNOME Keyring / KDE Wallet) via KMP binding.
- Biometric unlock gate for app resume (configurable).
- Certificate pinning for API host; mTLS for on-prem customers (optional client cert install flow).
- Secure clipboard (auto-clear after 60 s when we copy tokens or patches).

---

## 8. Notifications

### 8.1 Push Channels

- **FCM** on Android (data-only messages with encrypted payload; we decrypt and build local notification).
- **APNs** on iOS (background silent + user-visible categories).
- **UNNotification** on macOS.
- **Windows Toast** on Windows (Desktop app registers as Start Menu entry).
- **libnotify** on Linux.

### 8.2 Categories

- PR review requested
- Conflict needs you
- Mentioned in comment
- CI failed
- Upstream auth issue

Each category has deep links back into the app and quick actions (approve, decline, reply).

---

## 9. Platform-Specific Capabilities (expect/actual)

| Capability | Android | iOS | Desktop |
|---|---|---|---|
| File picker | ActivityResultContracts | UIDocumentPicker | JFileChooser/FileDialog |
| Share sheet | Intent.ACTION_SEND | UIActivityViewController | OS share API / clipboard |
| Camera (QR for login) | CameraX | AVFoundation | JavaCV (desktop optional) |
| Push | FCM | APNs | Platform native |
| Background | WorkManager | BGTaskScheduler | coroutines + OS tray |
| Deep links | intent-filter | Universal Links | custom URL scheme handler |
| Secure storage | Keystore | Keychain | DPAPI / Keychain / libsecret |
| Update | Play Core / In-App Update | TestFlight / App Store | Squirrel / Sparkle / dpkg |

---

## 10. UI Design

- **Compose Multiplatform** with **Material 3** adaptive scaffolding.
- Adaptive layouts: phone / foldable / tablet / desktop (window size classes).
- Responsive breakpoints: compact / medium / expanded / extra-large.
- Dynamic color on Android 12+; manual themes elsewhere.
- Keyboard shortcuts on desktop (`Cmd+K` for command palette, `?` for shortcuts overlay).
- Trackpad + mouse right-click support on desktop.

---

## 11. Internationalisation

- Source strings in `commonMain/resources/values/strings.xml`.
- Compose `stringResource` Multiplatform abstraction.
- At build time, per-platform resource generation (Android XML, iOS `.strings`, desktop ResourceBundle).
- Weblate-managed translations.

---

## 12. Performance

| Metric | Target |
|---|---|
| Cold start (Android, Pixel 5) | ≤ 1.2 s |
| Cold start (iPhone 12) | ≤ 1.0 s |
| Cold start (desktop, M1 / mid-range Windows) | ≤ 2.0 s |
| Frame time | 16 ms (60 fps); 120 fps on capable devices |
| Memory baseline | ≤ 180 MB |
| Bundle size (APK universal) | ≤ 35 MB |
| iOS IPA | ≤ 45 MB |
| Desktop installer | ≤ 90 MB |

Monitored in CI with Android-macrobenchmark + custom iOS XCTest measurements.

---

## 13. CI/CD for Clients

- **Android**: Gradle + Firebase Test Lab → AAB → Play Console internal track → production.
- **iOS**: fastlane → TestFlight → App Store Connect.
- **macOS**: jpackage → DMG with notarisation → direct download (Homebrew cask secondary).
- **Windows**: jpackage → MSI + MSIX → Microsoft Store + direct download.
- **Linux**: jpackage → DEB/RPM + AppImage + Flatpak (Flathub) + Snap.

Every client version is signed:
- Android: Play App Signing.
- iOS: App Store signing.
- macOS: Apple Developer ID + notarisation.
- Windows: EV code-signing cert.
- Linux: GPG-signed repositories.

---

## 14. Testing

- **Shared unit tests** (Kotest) run on JVM for speed; platform specifics via instrumented tests only when truly platform-bound.
- **Turbine** for Flow tests.
- **Paparazzi** (Android) and **Showkase** for screenshot tests.
- **Maestro** E2E flows run on Android emulator + iOS simulator + desktop.
- **Accessibility** via Espresso a11y + iOS Accessibility Inspector + desktop heuristic checks.
- **Monkey** / fuzz: Android Monkey, iOS UI Monkey; deterministic seeds in CI.
- **Battery/perf regression** tests on real devices (Samsung, Pixel, iPhone) via Firebase Test Lab custom devices + BrowserStack.

---

## 15. Shared With Web

- **Proto-generated types**: same `.proto` sources feed Angular (TS) and KMP (Kotlin). No schema drift.
- **Design tokens**: shared JSON (`style-tokens.json`) emits CSS variables + Compose theme + XCAssets colour set.
- **Icons**: shared SVG set; baked into each platform at build.
- **i18n**: shared ICU source strings.

---

## 16. Developer Experience

- `./gradlew runDesktop` to run the desktop app locally.
- `./gradlew runMobileOnAny` detects attached device and launches.
- `hot-reload`: Compose Hot Reload (desktop + Android KMP preview).
- CLI helper `helixctl-client dev-link` generates test tokens and injects them into local clients (dev only).
- Storybook-like: **Paparazzi + Showkase catalogue** for commonUI inspection.

---

## 17. Rollout Strategy

1. GA: Android + iOS + desktop (macOS primary; Windows/Linux behind feature flag).
2. +1 quarter: Windows/Linux promoted to GA.
3. +2 quarters: CLI (Kotlin/Native) for power users who want the client in a terminal with live-events pane.

---

*— End of Mobile & Desktop (KMP + Compose) —*
