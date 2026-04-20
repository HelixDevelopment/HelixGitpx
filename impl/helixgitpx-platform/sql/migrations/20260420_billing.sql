-- M8 — billing-service schema
CREATE SCHEMA IF NOT EXISTS billing AUTHORIZATION billing_svc;

CREATE TABLE IF NOT EXISTS billing.customers (
  id              UUID PRIMARY KEY,
  org_id          UUID NOT NULL REFERENCES org.organizations(id) ON DELETE CASCADE,
  external_id     TEXT NOT NULL UNIQUE,
  email           TEXT NOT NULL,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS billing.subscriptions (
  id              UUID PRIMARY KEY,
  customer_id     UUID NOT NULL REFERENCES billing.customers(id) ON DELETE CASCADE,
  external_id     TEXT NOT NULL UNIQUE,
  plan            TEXT NOT NULL CHECK (plan IN ('free','team','scale','enterprise')),
  status          TEXT NOT NULL,
  current_period_end TIMESTAMPTZ,
  cancelled_at    TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS billing.invoices (
  id              UUID PRIMARY KEY,
  subscription_id UUID NOT NULL REFERENCES billing.subscriptions(id) ON DELETE CASCADE,
  external_id     TEXT NOT NULL UNIQUE,
  amount_cents    BIGINT NOT NULL,
  currency        TEXT NOT NULL,
  status          TEXT NOT NULL,
  issued_at       TIMESTAMPTZ NOT NULL,
  paid_at         TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS subscriptions_customer_idx ON billing.subscriptions (customer_id);
CREATE INDEX IF NOT EXISTS invoices_subscription_idx ON billing.invoices (subscription_id);
