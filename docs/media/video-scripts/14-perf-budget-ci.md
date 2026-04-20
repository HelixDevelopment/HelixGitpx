# Script — 14 Perf budget enforcement in CI

**Track:** Operators · **Length:** 8 min · **Goal:** viewer reads a failing perf run and fixes the root cause.

## Body

1. **The budget file** — 0:30 – 1:30.
   `tools/perf/budgets.json`: p95/p99/error-rate per scenario.
2. **k6 scenario walk** — 1:30 – 3:00.
   api_baseline.js, git_push_pull.js, websocket_fanout.js, ai_chat.js.
3. **CI gate** — 3:00 – 4:30.
   `.github/workflows/perf-budgets.yml`. Fails on breach.
4. **Reading a failure** — 4:30 – 6:00.
   Budget breach → Tempo trace → flame graph in Pyroscope.
5. **Fix + re-run** — 6:00 – 7:30.

## Wrap-up
"If it's fast today and slow tomorrow, the budget catches it."

## Companion doc
`tools/perf/`
