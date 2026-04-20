# Distribution scripts

Each script wraps a single platform's publishing flow. All require credentials
injected via Vault or CI secrets (documented inline).

| Script | Target | Credentials |
|---|---|---|
| `publish-play.sh`     | Google Play / F-Droid                | `PLAY_SERVICE_ACCOUNT_JSON`, optional F-Droid key |
| `publish-appstore.sh` | App Store Connect / TestFlight       | `APPLE_API_KEY`, `APPLE_API_ISSUER`, `APPLE_API_KEY_ID` |
| `publish-windows.sh`  | MSIX via `msix` CLI                  | signing cert PFX + password |
| `publish-macos.sh`    | DMG notarization via `xcrun notarytool` | Apple developer account |
| `publish-linux.sh`    | AppImage + .deb + .rpm               | none for build; gpg for sign |
| `publish-update-feed.sh` | self-hosted update feed (M6 ADR-0029) | `TUS_ENDPOINT`, `TUS_TOKEN` |

All scripts pull the Compose-built artefacts from `impl/helixgitpx-clients/*/build/`.
M6 ships the scripts; CI/CD pipelines invoke them in M8 release automation.
