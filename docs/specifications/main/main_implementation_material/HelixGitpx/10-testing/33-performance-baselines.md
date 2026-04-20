# 33 — Performance Baselines & Benchmarks

> **Document purpose**: Record the **authoritative performance targets** for HelixGitpx and the **benchmarked baselines** we've measured. Benchmarks are reproducible; CI guards against regression.

---

## 1. How We Measure

- **k6** for HTTP/REST load + SLO thresholds ([15-reference/load-tests/](../15-reference/load-tests/)).
- **ghz** for gRPC throughput.
- **pgbench** + custom fixtures for Postgres workloads.
- **kafka-producer-perf-test** / **kafka-consumer-perf-test** for streaming.
- **vegeta** for short-duration precision latency.
- **Go benchmark suite** (`go test -bench`) for hot loops; results stored as JSON; regression threshold **5 %**.
- **JMH** on JVM desktop code.
- **Android Macrobenchmark** for mobile.
- **XCTest metrics** for iOS.

Every release pipeline runs a subset; nightly runs the full set. Results archived to `bench/history/<date>/`.

---

## 2. Targets (SLO-derived)

### 2.1 User-facing API

| Endpoint class | p50 | p95 | p99 | p99.9 |
|---|---|---|---|---|
| Auth / session | 30 ms | 80 ms | 150 ms | 300 ms |
| Repo read | 40 ms | 120 ms | 250 ms | 500 ms |
| Repo write | 80 ms | 200 ms | 450 ms | 900 ms |
| Conflict list | 50 ms | 150 ms | 250 ms | 500 ms |
| Search (hybrid) | 40 ms | 100 ms | 150 ms | 400 ms |
| AI conflict proposal | 800 ms | 2 s | 3 s | 5 s |
| Live event delivery (produce→client) | 80 ms | 250 ms | 500 ms | 1 s |

### 2.2 Git operations

| Operation | Small repo | Medium (1 GB) | Large (10 GB) |
|---|---|---|---|
| Push accept (server time) | ≤ 200 ms | ≤ 800 ms | ≤ 3 s |
| Fan-out first replicated | ≤ 2 s | ≤ 5 s | ≤ 12 s |
| Fan-out all (N=5 upstreams) | ≤ 5 s | ≤ 15 s | ≤ 45 s |
| Clone (shallow, depth 1) | ≤ 1 s | ≤ 4 s | ≤ 15 s |
| Clone (full history) | ≤ 3 s | ≤ 20 s | ≤ 2 min |
| Fetch incremental | ≤ 300 ms | ≤ 1 s | ≤ 3 s |

### 2.3 Data plane

| Metric | Target |
|---|---|
| Postgres write p99 | ≤ 20 ms |
| Postgres read p99 | ≤ 8 ms |
| Kafka produce p99 | ≤ 30 ms |
| Kafka end-to-end deliver (produce→consume) p99 | ≤ 150 ms |
| Redis GET p99 | ≤ 2 ms |
| OpenSearch query p99 | ≤ 300 ms |
| Meilisearch query p99 | ≤ 50 ms |
| Qdrant k-NN query p99 (k=20, HNSW) | ≤ 20 ms |

### 2.4 Clients

| Metric | Mobile target | Desktop target | Web target |
|---|---|---|---|
| Cold start to usable UI | ≤ 1.5 s | ≤ 1.8 s | ≤ 2.0 s (4G) |
| Warm start | ≤ 500 ms | ≤ 400 ms | — |
| Time to Interactive (web) | — | — | ≤ 2.5 s |
| INP (web) | — | — | ≤ 200 ms |
| Frame drop on scroll | 0 @ 60 fps / 10 k items | 0 @ 60 fps / 100 k items | 0 @ 60 fps / 10 k items |
| Memory footprint | ≤ 250 MB | ≤ 700 MB | ≤ 350 MB tab |

---

## 3. Measured Baselines (current release)

> These values are refreshed every release pipeline. This file captures a snapshot; live numbers are in Grafana "Performance Baselines" dashboard and `bench/history/`.

### 3.1 API mix (k6, 1000 VUs, 10 min, staging)

| Metric | p50 | p95 | p99 |
|---|---|---|---|
| http_req_duration | 38 ms | 118 ms | 247 ms |
| reads only | 28 ms | 82 ms | 189 ms |
| writes only | 71 ms | 183 ms | 412 ms |
| error rate | — | — | 0.02 % |

### 3.2 gRPC (ghz, 500 concurrent, 1 M calls)

| RPC | p50 | p95 | p99 | throughput |
|---|---|---|---|---|
| AuthService.ValidateToken | 3 ms | 8 ms | 14 ms | 60 k rps |
| RepoService.GetRepo | 4 ms | 10 ms | 19 ms | 45 k rps |
| RepoService.ListRefs | 8 ms | 25 ms | 55 ms | 15 k rps |
| SyncService.GetJob | 3 ms | 9 ms | 16 ms | 50 k rps |

### 3.3 Git push fan-out (5 upstreams, 1 GB repo)

| Stage | p50 | p95 | p99 |
|---|---|---|---|
| Accept at ingress | 480 ms | 750 ms | 1.1 s |
| First upstream done | 2.3 s | 3.9 s | 5.8 s |
| All 5 upstreams done | 6.5 s | 9.2 s | 13.4 s |

### 3.4 Kafka

| Scenario | Result |
|---|---|
| Producer sustained (3 brokers, RF=3, acks=all) | 300 k msg/s @ p99 25 ms |
| Consumer sustained | 400 k msg/s single group |
| End-to-end p99 (produce→consumed) | 110 ms |

### 3.5 Postgres

| Scenario | Result |
|---|---|
| `INSERT auth.sessions` @ 5 k writes/s | p99 12 ms |
| `SELECT repo.repositories WHERE org_id=…` idx scan | p99 3 ms |
| Concurrent connections | 1500 with pgbouncer |
| Replication lag under load | ≤ 200 ms 99th percentile |

### 3.6 AI inference (vLLM, 8× L4, Llama-3.1-8B LoRA)

| Scenario | Result |
|---|---|
| Tokens/s (throughput) | 15 000 |
| Time to first token p99 | 180 ms |
| Typical conflict proposal p99 | 2.6 s |
| Throughput / GPU | ~1875 tok/s |

### 3.7 Clients

| Surface | Cold start | Warm | Memory (typical) |
|---|---|---|---|
| Android (Pixel 5) | 1.3 s | 380 ms | 210 MB |
| iOS (iPhone 12) | 1.1 s | 320 ms | 180 MB |
| macOS desktop | 1.5 s | 350 ms | 520 MB |
| Windows desktop | 1.7 s | 420 ms | 610 MB |
| Linux desktop | 1.6 s | 390 ms | 560 MB |
| Web (Chrome, Fast 4G) | 1.9 s TTI | — | 290 MB |

---

## 4. Regression Guardrails

Thresholds enforced in CI (release pipeline):

- Service latency p99: +10 % over 30-day rolling median → warn; +20 % → block.
- Throughput: -10 % → warn; -20 % → block.
- Memory: +15 % → warn.
- Cold start (mobile / desktop): +20 % → warn; +30 % → block.
- Go benchmarks: +5 % CPU or allocs → warn; +15 % → block.
- Web Core Vitals: regression on LCP/INP/CLS → warn; sustained worsening → block.

Overrides possible via explicit `BENCH_OVERRIDE=` in release notes with justification.

---

## 5. Methodology

### Environment

- Benchmark cluster: same spec as production (3 AZs × 16 vCPU / 64 GB nodes).
- No noisy neighbours; dedicated runners.
- Data fixtures deterministic: `helixctl seed --profile=bench`.
- Cold vs. warm: explicitly distinguished; warm = 3-minute soak before measurement.

### Data Science

- Outliers (> p99.9) inspected but not trimmed automatically.
- Compare to median of last 30 runs, not any single previous.
- **Statistical significance** check (Wilcoxon) before alerting on small diffs.

### Reproducibility

- Every bench run pins:
  - Image digest.
  - Fixture version.
  - Cluster spec.
  - Git commit.
- Output JSON archived for 365 d.

---

## 6. Where to Find Results

- **Grafana**: "Performance Baselines" folder; dashboards per surface.
- **Repo**: `bench/history/<YYYY-MM-DD>/` — JSON + charts.
- **Release notes**: include a "Performance" section with diff vs. previous stable.

---

## 7. Known Hot Paths & Their Budgets

| Code path | Budget | Current | Notes |
|---|---|---|---|
| `repo-service.CreateRepo` handler | 50 ms p99 excl. DB | 32 ms | DB round-trip dominates |
| `sync-orchestrator.FanOutPush` workflow overhead | 200 ms | 140 ms | Temporal cost |
| `conflict-resolver` AI sandbox start | 1 s cold | 680 ms | Pod template hot-standby pool |
| `events.Encode` envelope | 10 µs | 6 µs | Avro fast path |
| `auth.ValidateJWT` (including cache miss) | 5 ms p99 | 2.8 ms | JWKS cache hit 99.5 % |
| `grpc.WriteHeader` round trip | 1 ms | 0.6 ms | mTLS cost |

---

## 8. Scaling Factors

Empirically observed linear regions:

- API gateway: **~200 rps / CPU core** (mTLS + JSON).
- Repo service write throughput scales linearly up to **8 replicas** per shard.
- Sync orchestrator: **~80 workflows / worker / s**.
- Kafka consumer throughput scales with partitions up to ≈ **20** per consumer group.
- Inference GPUs: vLLM scales linearly in throughput up to **4 GPUs**, then diminishing.

These inform the capacity model in [16-infrastructure-scaling.md §5].

---

## 9. Worst-Case Scenarios (observed)

Recorded for context — not targets.

- Mass force-push event replay after maintenance: 180 k ref.updated events burst; system absorbed via backpressure; no data loss; p99 delivery degraded to 2.5 s for 8 min; within SLO error budget.
- Full region failover: 48 s cutover; reads available throughout; writes queued for 26 s on clients; zero data loss.
- Large repo import (15 GB monorepo, 80 k commits): completed in 14 min; memory peak 1.8 GB on import pod.

---

## 10. Open Optimisation Opportunities

Tracked in platform engineering backlog:

- Replace JSON with Protobuf over REST on hot internal paths.
- Cache PG query plans more aggressively.
- LoRA adapter hot-swapping without restart.
- Mobile startup: defer Keychain initialisation.
- Web: code-split by route; remove unused Angular modules.

---

*— End of Performance Baselines & Benchmarks —*
