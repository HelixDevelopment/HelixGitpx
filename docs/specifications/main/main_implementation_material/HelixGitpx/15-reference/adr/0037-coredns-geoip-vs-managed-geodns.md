# ADR-0037 — CoreDNS + geoip over Route53 / Cloudflare GeoDNS

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Need GeoDNS for region failover. Options: Cloudflare Load Balancer (managed,
fast, Cloudflare lock-in), AWS Route53 geo routing (AWS lock-in), CoreDNS with
the `geoip` plugin (self-hosted, portable).

## Decision

CoreDNS + geoip. Runs in-cluster (per region), consults MaxMind GeoLite2,
returns region-appropriate endpoints.

## Consequences

- Portable across clouds/bare-metal, consistent with our "no vendor lock-in" stance.
- Health-check wiring is our responsibility (check operational runbooks).
- Slightly higher operational burden vs managed — accepted.

## Links

- Spec §LOCKED C-8
- impl/helixgitpx-platform/helm/geodns/
