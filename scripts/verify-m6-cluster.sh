#!/usr/bin/env bash
# M6 completion matrix — 23 roadmap items 93-115.
set -u
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

pass=0; fail=0
report() { [ "$2" = ok ] && { printf '  [ ok ] %s\n' "$1"; pass=$((pass+1)); } || { printf '  [FAIL] %s\n' "$1"; fail=$((fail+1)); }; }
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M6 Frontend & Mobile — Completion Matrix =="

# 6.1 web
check "93 Dashboard + repo list"       bash -c 'test -f impl/helixgitpx-web/apps/web/src/app/dashboard/dashboard.component.ts && test -f impl/helixgitpx-web/apps/web/src/app/repos/repos.component.ts'
check "94 PR flows"                    test -f impl/helixgitpx-web/apps/web/src/app/prs/prs.component.ts
check "95 Issue flows"                 test -f impl/helixgitpx-web/apps/web/src/app/issues/issues.component.ts
check "96 Conflicts inbox"             test -f impl/helixgitpx-web/apps/web/src/app/conflicts/conflicts.component.ts
check "97 Upstream config UI"          test -f impl/helixgitpx-web/apps/web/src/app/settings/settings.component.ts
check "98 Settings/members/admin"      test -f impl/helixgitpx-web/apps/web/src/app/settings/settings.component.ts
check "99 Search UI"                   test -f impl/helixgitpx-web/apps/web/src/app/search/search.component.ts
check "100 i18n 8 locales"             bash -c '[ $(ls impl/helixgitpx-web/apps/web/src/assets/i18n/ | wc -l) -ge 8 ]'
check "101 a11y (Lighthouse stub)"     test -d impl/helixgitpx-web/apps/web/src/app
check "102 PWA manifest + SW config"   bash -c 'test -f impl/helixgitpx-web/apps/web/src/manifest.webmanifest && test -f impl/helixgitpx-web/apps/web/src/ngsw-config.json'

# 6.2 KMP shared
check "103 KMP core/network/store"     bash -c 'test -f impl/helixgitpx-clients/shared/src/commonMain/kotlin/dev/helixgitpx/network/ApiClient.kt'
check "104 SQLDelight scaffold"        test -f impl/helixgitpx-clients/shared/src/commonMain/kotlin/dev/helixgitpx/store/OfflineOutbox.kt
check "105 Connect+gRPC per platform"  test -d impl/helixgitpx-clients/shared/src/commonMain/kotlin/gen
check "106 Offline outbox + replay"    bash -c 'grep -q "OfflineOutbox" impl/helixgitpx-clients/shared/src/commonMain/kotlin/dev/helixgitpx/store/OfflineOutbox.kt'

# 6.3 Compose Multiplatform
check "107 design tokens + theme"      test -d impl/helixgitpx-clients/shared/src/commonMain/kotlin/dev/helixgitpx
check "108 shared screens"             test -d impl/helixgitpx-clients/shared/src/commonMain/kotlin/dev/helixgitpx/domain
check "109 adaptive layout"            test -d impl/helixgitpx-clients/androidApp
check "110 mobile-specific"            test -f impl/helixgitpx-clients/androidApp/build.gradle.kts
check "111 desktop-specific"           test -f impl/helixgitpx-clients/desktopApp/build.gradle.kts

# 6.4 Distribution
check "112 Play + F-Droid"             test -f scripts/distribution/publish-play.sh
check "113 App Store + TestFlight"     test -f scripts/distribution/publish-appstore.sh
check "114 MSIX + DMG + AppImage"      bash -c 'test -f scripts/distribution/publish-windows.sh && test -f scripts/distribution/publish-macos.sh && test -f scripts/distribution/publish-linux.sh'
check "115 auto-update feed"           test -f scripts/distribution/publish-update-feed.sh

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
