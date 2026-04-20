# E2E Gaps — Audit

Inventory of user flows that lack full end-to-end automation. Each row has an issue ticket.

| Flow | Status | Ticket | Priority |
|---|---|---|---|
| Sign-up → first repo → push | ✅ Covered | HGX-101 | — |
| Create org → invite member → accept | ✅ Covered | HGX-102 | — |
| Bind upstream → initial mirror → PR | ✅ Covered | HGX-103 | — |
| Conflict detected → AI proposal → human accept | ⚠️ Partial (no AI path) | HGX-310 | High |
| Change residency → data migrates | ❌ Manual | HGX-311 | Med |
| DR failover → customer impact | ❌ Runbook-tested only | HGX-312 | High |
| Desktop app auto-update | ⚠️ Happy path only | HGX-313 | Med |
| Mobile push notification → deep link | ❌ Missing | HGX-314 | Med |
| Billing: plan upgrade, downgrade, cancel | ❌ Not yet | HGX-315 | High (GA blocker) |
| Trust center page: load, links resolve | ❌ Not yet | HGX-316 | Low |
| Full chaos recovery matrix | ⚠️ Manual Litmus runs | HGX-317 | Low |

## Closure target

All High tickets closed before GA. Mediums closed within 30 days post-GA.

## Tooling

- Web: Playwright suite at `impl/helixgitpx-web/e2e/`.
- Mobile: Kotest UI + Appium (Android), XCUITest (iOS).
- API: k6 + Connect-Go client tests.
- Platform: kind-in-CI with seeded fixtures.
