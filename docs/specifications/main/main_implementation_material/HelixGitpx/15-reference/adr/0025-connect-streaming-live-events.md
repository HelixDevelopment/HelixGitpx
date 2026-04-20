# ADR-0025 — Connect streaming for live-events (WS + SSE fallback)

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Web clients need real-time updates when a repo's state changes (conflict resolved, ref updated, label added). Options: raw gRPC streaming (no browser support), Server-Sent Events (one-way only), WebSockets (two-way, but complex auth/reconnect), Connect streaming (gRPC-compatible, browser-native).

## Decision

`live-events-service` exposes a single server-streaming RPC `LiveEvents.Subscribe(filter) stream Event`. Served via Connect: gRPC to Go clients, WebSocket to modern browsers, SSE fallback to older browsers. Resume on reconnect via an opaque `resume_token` that maps to a Redis Stream cursor.

## Consequences

- One server handler, three transports — `@connectrpc/connect-web` picks the transport per browser.
- Resume-token persistence in Redis gives at-least-once delivery with a 1-hour replay window.
- Horizontal scaling: each replica subscribes to `live.stream` in Redis; event fan-out happens client-side on the web.

## Links

- `docs/superpowers/specs/2026-04-20-m5-federation-conflict-engine-design.md` §2 C-4
- https://connectrpc.com/docs/web/supported-browsers-and-features/
