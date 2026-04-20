# DDoS tests

Validate that rate-limiters, token buckets, and graceful-degradation paths
hold under malicious-shaped load. NO mocks.

## Attack shapes

- Arrival burst: 100× baseline RPS for 30 s, then 0.
- Slowloris: 10k half-open connections for 10 min.
- Connection flood: `MaxConns * 2` new TCP streams per second.
- Cache-busting query flood: random URL suffixes bypass caches.
- Amplification: large payloads against endpoints with expensive fanout.
- Repeated 401 hammers: test auth-service's own rate-limiting.

## Assertions

1. The rate-limiter (token-bucket in `platform/quota/`) engages within
   1 second and drops > 99 % of abusive requests.
2. Legitimate clients (subset of the attack stream) see < 1 % error rate.
3. CPU / memory stay below autoscaler's scale-out cap.
4. After the attack ends, p99 latency returns to baseline within 60 s.

## Running

```bash
make test-ddos
```

## Constraint

NO mocks. Real ingress, real token-bucket, real Cilium NetworkPolicy.
