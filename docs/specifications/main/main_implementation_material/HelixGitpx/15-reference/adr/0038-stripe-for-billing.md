# ADR-0038 — Stripe as the billing provider at GA

- **Status:** Accepted
- **Date:** 2026-04-20
- **Deciders:** Милош Васић

## Context

Options: Stripe (dominant, rich SDK, strong EU + US), Paddle (merchant-of-record,
handles VAT), Lago + own PSP (OSS, more work), custom.

## Decision

Stripe. Wider reach, best SDK ergonomics. We wrap it behind `billing-service`'s
`provider.Provider` interface so we can swap / add providers later (e.g., add Paddle
for merchant-of-record in specific markets).

## Consequences

- Revenue recognition and invoice formatting follow Stripe's model.
- MoR (merchant-of-record) responsibilities stay with us until we add Paddle.
- Provider interface keeps a future migration achievable.

## Links

- Spec §LOCKED C-9
- impl/helixgitpx/services/billing-service/internal/provider/
