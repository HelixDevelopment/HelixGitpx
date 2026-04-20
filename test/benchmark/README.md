# Benchmark tests

Latency + throughput baselines for every public path. Regression detection
is enforced in CI (`.github/workflows/perf-budgets.yml`).

## Where

- Go `testing.B` micro-benchmarks live alongside the code as
  `*_bench_test.go` in each service.
- Integration-level benchmarks live here under `test/benchmark/`, run
  against a live compose stack.
- End-to-end latency budgets are k6 scenarios under `tools/perf/`.

## Running

```bash
# Go micro-benchmarks
cd impl/helixgitpx && go test -run=^$ -bench=. -benchmem ./...

# Compose-backed benchmarks
make test-benchmark

# Budget-gated run (fails PR on regression)
make test-benchmark-budgets
```

## Budgets

See `tools/perf/budgets.json` for the authoritative per-scenario p95/p99
+ error-rate thresholds. A run that breaches any budget blocks merge.

## Constraint

NO mocks. Benchmarks against stubs prove nothing about production.
