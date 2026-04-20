# Fuzz Corpora

Seed inputs for `go test -fuzz` / `go fuzz` targets.

## Layout

- `http/` — HTTP request parsers (services with custom parsers).
- `proto/` — binary protobuf inputs for buf-generated types.
- `git/` — git smart-HTTP pkt-line inputs (git-ingress).
- `webhook/` — JSON payloads seen from each provider (webhook-gateway).

## Running locally

```
cd impl/helixgitpx/services/webhook-gateway
GOTOOLCHAIN=go1.23.4 go test -fuzz=FuzzHMACParse -fuzztime=5m ./...
```

## CI

Nightly job invokes each fuzz target for 10 minutes and appends new corpus files
as artifacts. Crashes block subsequent runs until triaged.
