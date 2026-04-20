# helixgitpx-clients (KMP + Compose Multiplatform)

## Quick start

```sh
./gradlew check
```

## Layout

- `shared/` — KMP shared module (domain, network, store). M1 ships a stub.
- `buildSrc/` — convention plugins so per-target modules stay DRY.
- `androidApp/`, `iosApp/`, `desktopApp/` — platform shells added in M6.

iOS Kotlin targets compile on Linux (iOSX64 + iOSArm64 klibs); final iOS linking requires macOS (M6 CI).
