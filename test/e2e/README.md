# E2E test suites

Full user-journey tests. NO mocks.

## Web — Playwright

```bash
cd impl/helixgitpx-web
npx playwright test
```

Specs live in `impl/helixgitpx-web/e2e/`. Playwright spins up a headed
browser, walks through signup → first push → PR → merge.

## Mobile — Appium + XCUITest

Android: Appium driving the androidApp APK. iOS: XCUITest on the iosApp
bundle. Scripts under `impl/helixgitpx-clients/e2e/`.

## API — k6 scenarios

The k6 scenarios under `tools/perf/scenarios/` serve double duty: load-test
*and* e2e API-surface verification (the `api_baseline.js` scenario walks
login → create org → bind repo → list PRs).

## Constraint: real dependencies

A local cluster (k3d or kind) with the full Argo CD app-of-apps tree must
be running. No mocked backend is permitted.

See [`../README.md`](../README.md) and
[Constitution Article II](../../CONSTITUTION.md#2-article-ii--testing).
