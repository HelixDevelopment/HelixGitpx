# 15 — Testing Strategy

> **Document purpose**: Define the **comprehensive testing regime** that earns HelixGitpx its ≥ 100 % coverage mandate and enterprise trust. Testing is not a phase; it is continuous, automated, and enforced by CI gates.

---

## 1. Coverage Policy

- **New code**: ≥ 100 % line + branch coverage — measured by SonarQube on the diff. No exceptions.
- **Legacy code**: ratchet; coverage may never decrease on a PR. Tooling (`lcov` + `gocov` merged with SonarQube) enforces.
- **Mutation score**: ≥ 80 % on core services (conflict-resolver, sync-orchestrator, adapter-pool, auth-service).
- **Exclusions**: generated code (protobufs, mocks) is excluded from the coverage denominator. Everything hand-written is in scope.

---

## 2. Test Taxonomy

| Type | Scope | Tool (Go) | Tool (TS) | Tool (Kotlin) | Where |
|---|---|---|---|---|---|
| **Unit** | One package/class | `go test` + `testify` + `gomock` | Jest | Kotest | Every build |
| **Integration** | Service + real deps | `go test` + Testcontainers | Jest + Testcontainers-node | Testcontainers-kt | Every build |
| **Contract** | Provider ↔ consumer | `buf breaking` + Pact Go | Pact-JS | Pact-JVM | Every PR |
| **E2E (API)** | Real services, real data | Ginkgo + k6 | Playwright API | — | Nightly + pre-release |
| **E2E (UI)** | Browser / device | Playwright | Playwright | Maestro | Nightly |
| **Fuzz** | Random inputs | `go fuzz` + OSS-Fuzz | jsfuzz | Jazzer (JVM) | Continuous |
| **Property-based** | Invariants | `gopter` | fast-check | Kotest Property | Every build |
| **Mutation** | Test quality | `gremlins.go` | StrykerJS | Pitest | Nightly |
| **Security SAST** | Source analysis | SonarQube + Semgrep + CodeQL + Snyk Code | SonarQube + Semgrep | SonarQube + Detekt + Semgrep | Every PR |
| **Security DAST** | Running app | OWASP ZAP + Nuclei | | | Nightly staging |
| **Security SCA** | Dependencies | `govulncheck` + Snyk OSS + OSV | `npm audit` + Snyk | OWASP Dependency-Check | Every PR + nightly |
| **Container/IaC** | Images + manifests | Trivy + Snyk Container | — | — | Every build |
| **DDoS** | Edge + origin | k6 + Hyenae | — | — | Weekly staging |
| **Load** | Performance at scale | k6 | — | — | Pre-release |
| **Soak** | Long-running steady load | k6 (24 h) | — | — | Quarterly |
| **Chaos** | Fault injection | chaos-mesh + LitmusChaos | — | — | Weekly staging |
| **DR** | Region loss | runbook + `helixctl drill` | — | — | Quarterly |
| **Accessibility** | WCAG | axe-core + Lighthouse | axe-playwright | iOS a11y + Espresso | Every build |
| **Compatibility** | Browser/device matrix | Playwright matrix | Playwright | Firebase Test Lab | Nightly |
| **Benchmarks** | Perf regression | `go test -bench` | Vitest bench | kotlinx-benchmark | Every PR |
| **Upgrade/migration** | Schema, data | pgTAP + custom | — | — | Every PR touching migrations |
| **Release-train** | End-to-end smoke after deploy | curl suites + `helixctl selftest` | Playwright smoke | — | On deploy |
| **Observability tests** | Alerts fire | Prometheus `promtool test alerts` | — | — | Every PR |
| **Policy tests** | OPA policies | `opa test` | — | — | Every PR |
| **Replay** | Rebuild projections | CLI tool | — | — | Weekly staging |
| **Penetration** | Adversarial review | External vendor | — | — | Annual + major release |

---

## 3. Test Pyramid & Gates

```
                 ▲
                ╱ ╲
               ╱   ╲   Penetration  (annual)
              ╱     ╲
             ╱───────╲ DR / Soak / DDoS (quarterly)
            ╱         ╲
           ╱           ╲ Chaos / Nightly E2E
          ╱─────────────╲
         ╱               ╲ Contract / Integration (every PR)
        ╱                 ╲
       ╱───────────────────╲ Unit / Property / Mutation / SAST
      ╱                     ╲
```

**Blocking gates** on every PR:

1. Lint + format (`gofmt`, `buf lint`, `eslint`, `ktlint`).
2. Unit + property tests pass.
3. Integration tests on changed services pass.
4. Contract tests (buf breaking, Pact) pass.
5. SAST: SonarQube quality gate green, Snyk severity ≤ Medium.
6. SCA: no high/critical vulns.
7. Container/IaC scan clean.
8. Coverage ≥ 100 % on diff.
9. Mutation score above threshold on core.
10. Benchmarks within 5 % of baseline (or opt-in regression note).

---

## 4. Unit Tests (Go — canonical pattern)

```go
// services/sync-orchestrator/internal/planner/plan_test.go
func TestPlanner_BuildPushPlan(t *testing.T) {
    tests := []struct {
        name     string
        repo     Repo
        ref      RefUpdate
        upstreams []Upstream
        want     PushPlan
        wantErr  error
    }{
        { /* ... */ },
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got, err := BuildPushPlan(tc.repo, tc.ref, tc.upstreams)
            require.ErrorIs(t, err, tc.wantErr)
            require.Equal(t, tc.want, got)
        })
    }
}
```

Rules:
- Table-driven.
- No `goroutine` in the test unless specifically testing concurrency.
- `t.Parallel()` where possible.
- Subtests named after the scenario (snake_case).
- No sleeps — use `synctest` (Go 1.24+) or fake clocks.

---

## 5. Property-Based Tests

```go
func TestRefMerge_Commutative(t *testing.T) {
    gopter.NewProperties(nil).Property("ref-merge is commutative within idempotent ops", prop.ForAll(
        func(a, b RefSet, c RefSet) bool {
            return Merge(Merge(a, b), c).Equals(Merge(a, Merge(b, c)))
        },
        genRefSet(), genRefSet(), genRefSet(),
    )).TestingRun(t)
}
```

Run on every build. Failures produce shrunk counter-examples saved to `testdata/shrunk/`.

---

## 6. Mutation Testing

- Tooling: **gremlins.go** on Go; **Stryker** on TS; **Pitest** on Kotlin.
- Target: ≥ 80 % mutation score on core services.
- Gated nightly; PRs receive advisory feedback.
- Exceptions (surviving mutants) must be justified in `testing/MUTATION_EXCLUDES.md`.

---

## 7. Fuzzing

- Continuous fuzzing via **OSS-Fuzz** integration (preferred) or local `go test -fuzz` loops on dedicated CI workers.
- Corpus checked into `testdata/fuzz/corpus/` and periodically reduced.
- Target: adapter payload parsers, webhook verifiers, protobuf decoders, URL parsers, merge conflict parser.

---

## 8. Integration Tests

- Each service has an `it/` directory using **Testcontainers**.
- Dependencies:
  - Postgres (or Citus for sharded services).
  - Redis.
  - Kafka + Karapace.
  - Qdrant (for AI services).
  - Vault (dev mode).
  - SPIRE server (dev mode).
  - MinIO (object storage).
- Isolated per test via `t.Parallel()` with unique namespaces; ~30 s end-to-end per suite typical.

---

## 9. Contract Tests

### 9.1 gRPC

- `buf breaking` against the previous ref.
- Consumer-driven Pact between every pair (web → api-gateway; mobile → api-gateway; service-to-service contracts for internal RPCs).

### 9.2 Events

- **Schema Registry enforcement**: changes require a Karapace compatibility check (BACKWARD by default).
- **Payload contract**: consumer fixtures in `testdata/kafka_consumer_contracts/` validated against every produced event.

---

## 10. End-to-End Tests

### 10.1 API E2E

- Ginkgo-style suites running against a deployed staging tenant.
- Covers golden paths and known edge cases:
  - Connect GitHub upstream → create repo → push → verify mirror on GitHub.
  - Connect GitHub + GitLab → force divergence → conflict detected → AI proposal → human approve → propagate.
  - Revoke token → subsequent calls 401.
  - Rate limit saturate → 429 with correct headers.
- `helixctl selftest` CLI aggregates these and is run in production every 5 min (synthetic monitoring).

### 10.2 UI E2E

- Playwright: Chromium / Firefox / WebKit for web.
- Maestro for mobile (Android + iOS).
- Per-test auth via service account tokens (never real OIDC flow).
- Screen captures attached to CI run artefacts.

---

## 11. DDoS & Load

- **k6 cloud** agents simulate 100 k concurrent clients across regions.
- Specific scenarios:
  - Webhook flood (GitHub replay).
  - Git push burst (100 simultaneous clones + pushes).
  - Subscription storm (50 k WebSocket connects in 60 s).
- Expected behaviour: graceful degradation, 429s where appropriate, no data loss, autoscaling triggered.

---

## 12. Chaos

- **chaos-mesh** for Kubernetes scheduled experiments:
  - Random pod kill (per namespace, 1 pod every 10 min).
  - Network loss 5 % on one node.
  - Kafka broker termination.
  - Postgres primary failover.
  - Disk full on a node.
- Alerts must fire, user-visible SLOs must not breach.

---

## 13. Disaster Recovery Drills

- Quarterly: simulate full region loss in staging.
- `helixctl drill region-loss --region=eu-west` triggers: DNS cutover, queue drain check, consumer rebalance verification.
- RPO ≤ 30 s, RTO ≤ 5 min.
- Report posted to engineering blog + stakeholder email.

---

## 14. Performance & Benchmarks

- `go test -bench` per package with benchmarked functions flagged `// BENCH`.
- Benchmarks run on dedicated bare-metal runner; results posted to `perf-dashboards`.
- Threshold: > 5 % regression requires justification or fix.
- End-to-end: `k6` scenarios with SLO assertions.

---

## 15. Accessibility

- **axe-core** automated on every page.
- Lighthouse CI threshold ≥ 95.
- Manual screen-reader walk-through per release (NVDA, VoiceOver, TalkBack).
- Colour-contrast tests.

---

## 16. Policy Tests

- OPA policies tested with `opa test`; coverage ≥ 100 %.
- Changes reviewed by security; staged behind feature flag.

---

## 17. Observability Tests

- `promtool test alerts` validates alerting rules against synthetic series.
- Dashboards "smoke": render in CI via Grafana API; assert queries don't error.

---

## 18. Test Data Management

- Synthetic data via **factoryx** builder (per-entity factories).
- Anonymised production data for fuzz corpora (scrubbed via Gitleaks + custom rules).
- Tenant-isolated test fixtures.
- No real customer data in non-production.

---

## 19. Flaky Test Policy

- Any test failing intermittently is auto-quarantined (skipped) after 2 flakes.
- SLA: 48 h to fix or delete.
- Dashboard `flakes.helixgitpx.internal` lists quarantined tests publicly to engineering.

---

## 20. Reporting

- **Coverage**: SonarQube + Codecov public badge.
- **Mutation**: dashboard per service.
- **E2E**: Allure reports archived for 90 days.
- **Load**: Grafana dashboard snapshots per run.
- **Security**: monthly report to engineering leadership with trending.

---

*— End of Testing Strategy —*
