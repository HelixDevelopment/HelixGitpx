# billing-service

Stripe-backed billing (M8 GA). Upserts customers on org creation, manages plan
subscriptions, emits `billing.*` events to the outbox for downstream notification.

- **Provider:** Stripe (replaceable via `internal/provider.Provider`).
- **Data:** `billing.customers`, `billing.subscriptions`, `billing.invoices`.
- **Policy:** OPA enforcement on plan upgrades (e.g., only org owners).
