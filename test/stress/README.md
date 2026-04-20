# Stress tests

Drive the system to **3× the design capacity target** and assert that:

1. Error rate remains ≤ design SLO.
2. Latency p99 stays within 2× of baseline.
3. No resource leaks (no OOM, no file-descriptor leaks, no goroutine leaks).
4. Recovery time back to baseline ≤ 5 minutes after stress removed.

## Scenarios

Live at `tools/perf/scenarios/` with the `_stress` suffix:

- `api_stress.js` — ramp from 500 → 1500 rps over 30 min.
- `websocket_stress.js` — 5000 concurrent subscribers.
- `git_push_stress.js` — 500 concurrent pushes.
- `ai_chat_stress.js` — 200 concurrent AI sessions.

## Running

```bash
make test-stress
```

Uses the staging cluster by default. Set `K6_TARGET=…` to point elsewhere.

## Constraint

NO mocks. Real Kafka, real Postgres, real Ollama/vLLM.
