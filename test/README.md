# Test suites

Per [Constitution Article II](../CONSTITUTION.md#2-article-ii--testing),
every module in this repository carries tests in **all seven required test
types**. This directory holds the cross-service suites; per-service unit
tests live alongside the code (`*_test.go`, `*.spec.ts`, `commonTest/`).

## Types

| Dir             | Purpose | Mocks? | Runner |
|-----------------|---------|--------|--------|
| `integration/`  | Real collaborators, real deps | **NO** | `go test -tags=integration`, compose up |
| `e2e/`          | Full user journeys             | **NO** | Playwright + Appium + k6 |
| `security/`     | Authn/z, injection, ASVS L2    | **NO** | OWASP ZAP, Nuclei, gosec, custom |
| `stress/`       | Load up to 3× design target    | **NO** | k6 |
| `ddos/`         | Rate-limit, exhaustion, recovery | **NO** | k6 arrival bursts + Litmus |
| `benchmark/`    | Latency, throughput, regression | **NO** | `go test -bench`, k6 |
| `chaos/`        | Fault injection, DR            | **NO** | Litmus (see `tools/chaos/`) |

Unit tests are the only type that may use mocks, stubs, or placeholder data.

## Running

```bash
# Integration — requires compose stack up.
make test-integration

# E2E — requires staging cluster or k3d.
make test-e2e

# Security
make test-security

# Stress
make test-stress

# DDoS
make test-ddos

# Benchmark
make test-benchmark

# All seven types at once (used by CI).
make test-all
```

The Makefile target `test-all` refuses to pass if any of the required
types is missing for a module touched by the current git diff.

## Coverage

100 % per type per module touched. Measured by:

- Go: `go test -cover -coverprofile=coverage.out ./...`, aggregated per type.
- TS: `jest --coverage --json-summary`.
- Kotlin: `kover` + merged report.

CI enforces the threshold via `tools/coverage-audit/audit.sh`.

## Adding a new test

1. Pick the correct directory for the type (or co-locate `*_test.go` for unit).
2. Name the file `<subject>_<type>_test.<ext>` (e.g. `auth_login_integration_test.go`).
3. Wire real dependencies through `compose/up.sh` / `tools/compose-up.sh`.
4. Open a PR. CI validates the matrix.

## Prohibited

- `t.Skip`, `xit`, `@Ignore`, `pytest.skip` — banned.
- `go:build` flags that exclude tests in CI without an ADR.
- Mocks in any directory listed above.
- Flaky tests. If a test is non-deterministic, fix the root cause or delete the subject.

## Further reading

- [CONSTITUTION.md Article II](../CONSTITUTION.md#2-article-ii--testing)
- [AGENTS.md §4](../AGENTS.md)
- [tools/coverage-audit/audit.sh](../tools/coverage-audit/audit.sh)
