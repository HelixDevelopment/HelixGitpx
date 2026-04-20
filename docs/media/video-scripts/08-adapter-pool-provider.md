# Script — 08 Writing a new adapter-pool provider

**Track:** Developers · **Length:** 15 min · **Goal:** viewer ships a PR that adds a new Git provider adapter.

## Body

1. **What "provider" means** — 0:30 – 1:30.
   Adapter interface: metadata mapping, REST/HTTP client, HMAC scheme, rate-limit rules.
2. **Scaffold** — 1:30 – 3:00.
   `go run ./tools/scaffold provider --name foohub`.
3. **Map PR metadata** — 3:00 – 6:00.
   Canonical model in `proto`; per-field translation in `adapter.go`.
4. **HMAC validation** — 6:00 – 8:00.
   Provider-specific signature; use `platform/webhook`.
5. **Unit tests** — 8:00 – 10:00.
   Mocks here only (unit). Fixtures for each shape.
6. **Integration tests** — 10:00 – 12:00.
   Real HTTP against a foohub VCR cassette — no mocks.
7. **Docs + ADR** — 12:00 – 13:30.
   ADR describing the provider oddities.
8. **PR + CI green** — 13:30 – 14:30.
   Opens PR; CI runs all 7 test types.

## Wrap-up
"One adapter = real federation for thousands of new customers."

## Companion doc
`impl/helixgitpx/services/adapter-pool/` · `CONTRIBUTING.md`
