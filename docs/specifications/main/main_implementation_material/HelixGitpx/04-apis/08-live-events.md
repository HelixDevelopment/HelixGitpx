# 08 — Live Events & WebSocket Fallback

> **Document purpose**: Define the **real-time event delivery** architecture so that HelixGitpx clients (web, mobile, desktop, CLI, third-party integrations) get sub-second reactivity to every state change. Primary wire is **gRPC server streaming**; fallbacks are **WebSocket** (with Connect / Socket.IO-compatible framing) and **Server-Sent Events (SSE)**.

---

## 1. Overview

```
┌─────────────┐    stream   ┌──────────────┐   Kafka   ┌──────────────┐
│  Client     │ ◄──────── │ live-events  │ ◄──────── │ Every service│
│ (web/mob/…) │            │   service    │            │ that emits  │
└─────────────┘   ack ──►   └──────┬───────┘           └──────────────┘
                                   │ Redis Streams per subscriber
                                   ▼
                              ┌──────────┐
                              │  Redis   │ (buffer + backpressure)
                              └──────────┘
```

The `live-events-service` is a thin, stateless-ish layer (session state in Redis) that:

1. Authenticates the subscription.
2. Applies scope & ACL filters.
3. Tails Kafka for matching events.
4. Streams to the client over the best available transport.
5. Handles resume-from-cursor on reconnect.

---

## 2. Transports

| Transport | Availability | Use case |
|---|---|---|
| **gRPC server streaming** | All clients except browsers-behind-strict-proxies | **Primary** |
| **gRPC-Web (Connect)** | Browsers, proxies that allow HTTP/2 | Web app default |
| **WebSocket (Connect binary)** | All modern browsers, most proxies | Web fallback, mobile on bad networks |
| **Server-Sent Events (SSE)** | Any HTTP/1.1 client; unidirectional | Simple scripts, shell integrations, curl |
| **Long-polling** | Legacy/corporate networks | Fallback-of-fallbacks |

Clients try in this order: gRPC → WebSocket → SSE → long-poll. The **Connect** protocol is used across the first three so the server code is unified.

---

## 3. Subscription Model

A subscription is a tuple of:

- **scopes**: list of `(kind, id)` — e.g. `user:self`, `org:018f…`, `repo:018f…`, `global` (admin).
- **event_types**: optional filter (e.g. `["ref.updated", "pr.*"]` — glob).
- **resume_token**: opaque cursor; empty = "start now".

### 3.1 Authorisation

Before a subscription is established:

- Access token scopes checked (`events:read` required).
- Each scope is evaluated by OPA (policy-service): the principal must have at least `viewer` role on the subject.
- If any scope fails, the subscription is refused (`PERMISSION_DENIED`).

### 3.2 Per-Token Limits

- Max concurrent subscriptions per token: **50**.
- Max scopes per subscription: **200**.
- Max events buffered per subscription before drop-oldest: **1000**.

---

## 4. gRPC Streaming Protocol

Defined in `proto/helixgitpx/v1/events.proto`:

```proto
service EventsService {
  rpc Subscribe (SubscribeRequest)  returns (stream Event);
  rpc Resume    (ResumeRequest)     returns (stream Event);
  rpc Ack       (AckRequest)        returns (google.protobuf.Empty);
}

message SubscribeRequest {
  repeated Scope scopes = 1;
  repeated string event_types = 2;     // glob expressions
  string resume_token = 3;             // empty = "now"
  int32 buffer_hint = 4;               // server clamps to [100, 1000]
}

message Event {
  string event_id = 1;                 // UUIDv7
  string event_type = 2;               // e.g. "ref.updated"
  string resume_token = 3;             // updated on every event
  google.protobuf.Timestamp occurred_at = 4;
  string tenant_id = 5;
  google.protobuf.Any payload = 6;     // concrete message per event_type
  map<string, string> attributes = 7;  // trace_id, correlation_id, etc.
}

message AckRequest {
  repeated string event_ids = 1;       // batch-ack
}
```

### 4.1 Flow

1. Client calls `Subscribe(...)` (or `Resume(token)` after reconnect).
2. Server streams `Event` messages.
3. Client processes; periodically calls `Ack(event_ids)` so server can advance its cursor.
4. On stream error, client reconnects with the last `resume_token`.

### 4.2 Flow Control

- HTTP/2 flow control at the transport layer.
- Application-level: server tracks per-subscription inflight count; pauses production when ack lag > high-watermark.
- If client never acks for >30 s: server switches to **best-effort** mode (drop-oldest with "gap" notification event).

---

## 5. WebSocket Fallback

When gRPC is unavailable, the client upgrades over HTTP/1.1:

```
GET /api/v1/events/subscribe HTTP/1.1
Upgrade: websocket
Sec-WebSocket-Protocol: helixgitpx.v1+binary
Authorization: Bearer …
```

Framing uses **Connect**'s `grpc-web+proto` sub-protocol so server code is shared. For SDK simplicity, we also offer a **Socket.IO-compatible** JSON envelope:

```json
{ "e": "event.Subscribe", "d": { "scopes": [...], "event_types": [...], "resume_token": "..." } }
{ "e": "Event", "d": { "event_id": "...", "event_type": "ref.updated", "payload": {...}, "resume_token": "..." } }
{ "e": "Ack", "d": { "event_ids": ["...","..."] } }
{ "e": "Error", "d": { "code": "UNAUTHENTICATED", "message": "..." } }
{ "e": "Ping", "d": {"ts": 123} }
{ "e": "Pong", "d": {"ts": 123} }
```

### 5.1 Keepalive

- Ping/Pong every 25 s.
- Idle close after 60 s without traffic in either direction.

### 5.2 Reconnect & Resume

```
client: connect → subscribe(resume_token=X)
server: resumes from after X; emits missed events; sends new ones
```

If the resume_token is **too old** (Kafka retention passed), server responds with a **snapshot event**: the current state for the subscribed scope, followed by the live stream. This is transparent to the client.

---

## 6. Server-Sent Events (SSE) Fallback

```http
GET /api/v1/events/subscribe?scopes=repo:018f…&event_types=ref.%2A
Accept: text/event-stream
Authorization: Bearer …
Last-Event-ID: 01HPXK9…
```

Server sends:

```
id: 01HPXKA1…
event: ref.updated
retry: 2000
data: {"repo_id":"018f…","ref":"refs/heads/main","old":"abc","new":"def"}

: keepalive

id: 01HPXKA2…
event: pr.opened
data: {"repo_id":"018f…","number":42,"title":"Add caching"}
```

SSE is **unidirectional** (server→client). For acks, the client issues a separate `POST /api/v1/events/ack`. This is suboptimal for throughput; SSE is an "it works everywhere" safety net.

---

## 7. Scale & Scheduling

### 7.1 Routing & Stickiness

- Every subscription is bound to a single live-events-service pod.
- Istio consistent-hash-by-header routes subscriptions to the same pod (`x-subscription-hash`).
- If the pod dies, client reconnects; router picks a new pod; server resumes from `resume_token`.

### 7.2 Capacity Model

| Metric | Target per pod |
|---|---|
| Concurrent subscriptions | 10 000 |
| Events/sec delivered | 50 000 |
| Memory | 2 GiB |
| CPU | 1 core |

Horizontal autoscaling on active subscription count (custom metric).

### 7.3 Backpressure

- Per-subscription Redis Stream capped at 1 000 entries.
- When full: drop-oldest + send `gap_detected` event with last delivered token — client can `GET` the missed resource.
- Server metric `events_dropped_total{reason="slow_consumer"}` tracked and alerted.

---

## 8. Filtering & Fan-Out

Kafka topics are event-type-specific. The live-events-service subscribes to the **union** of topics relevant to active subscribers. A Bloom filter (Redis) maps `scope_key → bool` to avoid delivering events to scopes with no subscribers.

For hot events (global admin events), a fan-out fan-in worker model is used: the Kafka consumer writes to per-session Redis Streams, and a separate goroutine per session tails the stream and pushes to the wire.

---

## 9. Multi-Region

- Live-events-service is deployed per region.
- Cross-region events arrive via MirrorMaker 2.
- Clients connect to the nearest region by GeoDNS.
- On regional failover, clients reconnect to the secondary region; `resume_token` is region-aware but portable (encoded in Avro with a `region_prefix`).

---

## 10. Security

- TLS 1.3 only; HSTS.
- mTLS for internal hops.
- Token validated on connect **and** re-validated every 5 min — if scopes or user state changed, server closes the stream with a "re-auth required" message.
- WebSocket: same origin check; no cross-origin subscriptions.
- Per-IP connection cap: 100 simultaneous (protects against connection exhaustion).

---

## 11. Observability

| Metric | Meaning |
|---|---|
| `helixgitpx_events_subscriptions_active` | Gauge |
| `helixgitpx_events_sent_total{transport,event_type}` | Counter |
| `helixgitpx_events_dropped_total{reason}` | Counter |
| `helixgitpx_events_delivery_seconds{transport}` | Histogram (Kafka produce → client receive) |
| `helixgitpx_events_connections_total{transport}` | Counter |
| `helixgitpx_events_disconnects_total{reason}` | Counter |
| `helixgitpx_events_resume_tokens_stale_total` | Counter |

Alerts:

- `EventsDeliveryLagHigh` (p95 > 500 ms for 5 min).
- `EventsDroppedSpike` (> 100/s).
- `EventsSubscriptionsOverLimit` (any pod > 9 000 active).

---

## 12. Client Patterns (Angular example)

```ts
// Pseudocode
const stream = client.events.subscribe({
  scopes: [{ repo_id: repoId }, { user_id: 'self' }],
  event_types: ['ref.*', 'pr.*', 'conflict.*'],
  resume_token: stored.resumeToken
});

stream.on('message', (ev) => {
  store.ingest(ev);
  stored.resumeToken = ev.resume_token;
});

stream.on('gap', () => {
  // do a full refetch of affected resource(s)
  store.reload();
});

stream.on('close', (reason) => {
  if (reason.code === 'REAUTH') auth.refresh().then(reconnect);
  else scheduleReconnect();
});
```

Clients should:
- Persist `resume_token` (e.g., in IndexedDB) so that app reloads don't reset the stream.
- Batch UI updates (requestAnimationFrame / coalesce) to avoid layout thrash.
- De-dupe by `event_id` (idempotent ingestion).

---

## 13. Testing

- **Unit**: per-session buffer + filter logic, with table-driven scope matchers.
- **Integration**: Testcontainers Kafka + Redis + real gRPC client over HTTP/2.
- **Load**: k6/ws with 10 000 concurrent subscribers per pod; assert p95 delivery latency < 200 ms.
- **Chaos**: kill pod mid-stream; client must reconnect and miss zero events (within Kafka retention).
- **Compatibility**: integration tests for every SDK (Go/TS/Kotlin/Swift/Python/Rust) against the same backend.

---

*— End of Live Events —*
