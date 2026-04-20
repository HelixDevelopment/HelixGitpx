# RB-140 — Live Events Delivery Lag / Stream Pile-up

> **Severity default**: P2 (P1 if delivery p99 > 5 s for 10 min)
> **Owner**: Platform / Live-events team
> **Last tested**: 2026-04-07

## 1. Detection

Alerts:
- `LiveEventsDeliveryLagHigh` — `helixgitpx_events_delivery_seconds` p99 > 500 ms for 5 min.
- `LiveEventsDroppedSpike` — `helixgitpx_events_dropped_total` rate > 100/s.

Supporting signals:
- Customers see "stale" UI — status tags not updating, missing events.
- `helixgitpx_events_subscriptions_active` divergence between transport types.
- `helixgitpx_events_disconnects_total{reason}` showing `slow_consumer`.
- Kafka lag on `live.fanout.firehose` group growing.

---

## 2. First 2 Minutes

1. Ack; open incident if P1 threshold crossed.
2. Check `live-events-service` pod health — any crash-loops?
3. Dashboard "Live Events — Delivery" shows per-transport (gRPC / WS / SSE) breakdown.
4. Identify: is the problem **inbound** (fan-in from Kafka) or **outbound** (fan-out to clients)?

---

## 3. Diagnose

### 3.1 Inbound (Kafka → service)

- Check `helixgitpx_kafka_consumer_lag{group="live-events"}` — growing means upstream can't keep up.
- Check CPU + memory on live-events pods; HPA scaled appropriately?
- Is a specific topic dominant in lag? Look at `firehose` vs. specific repo topics.

### 3.2 Outbound (service → client)

- Slow clients? Check `events_disconnects_total{reason="slow_consumer"}`.
- Network back-pressure? WS-specific.
- Client bug not ack'ing? Check resume-token age distribution.
- Bad subscription (unbounded scope) producing huge message rates?

### 3.3 Both

- Noisy single tenant? Look at `events_sent_total` top by `org_id`.

---

## 4. Mitigations

### 4.1 Scale out

```bash
# Manual bump if KEDA slow to react
helixctl scale live-events-service --replicas=<current*2>

# Widen KEDA max
helixctl keda override --scaled-object=live-events-service --max-replicas=200
```

### 4.2 Shed load (transport-specific)

```bash
# Reject new subscriptions temporarily
helixctl flag set live.reject-new-subscriptions --on --reason="RB-140"

# Downgrade WS to SSE to reduce per-connection cost
helixctl flag set live.prefer-sse-fallback --on
```

### 4.3 Throttle noisy tenant

```bash
helixctl events tenant-throttle --org <slug> --max-mps=50 --reason=...
```

Logs the action; notifies CSM.

### 4.4 Resume-token pressure

If subscribers hold very old tokens forcing re-fetch from Kafka:

```bash
helixctl events resume-cap --older-than=2h
```

Clients will reconnect and fetch fresh state from projections.

### 4.5 Feature-flag fan-in scope

If a specific event type is spamming (e.g. a broken producer emitting excessive updates):

```bash
helixctl events block-type --event-type=<type> --for=10m --reason=...
```

Consumers transparently skip; upstream producer gets a page.

---

## 5. Verify Recovery

- Delivery p99 back under 500 ms.
- `events_dropped_total` rate = 0.
- Subscription count stable (not thrashing).
- Synthetic WS probe from each region succeeds.

---

## 6. Post-Incident

- If tenant-caused: CSM outreach + offer plan adjustment or per-org rate limits.
- If producer-caused: file a bug in the upstream service; add a unit test.
- If scaling-caused: tune KEDA parameters; capacity plan.
- Add the scenario to chaos tests if systematic.

---

## 7. Related

- RB-021 (Kafka consumer lag critical)
- RB-100 (API gateway 5xx spike) — correlated during broad degradation
- [08-live-events.md] — architecture
- [metrics-catalog.md] § Live Events
