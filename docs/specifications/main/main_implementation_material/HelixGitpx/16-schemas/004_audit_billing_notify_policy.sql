-- ============================================================
-- 16-schemas/004_audit_billing_notify_policy.sql
-- Bounded contexts: AUDIT, BILLING, NOTIFY, POLICY, SYNC
-- ============================================================

CREATE SCHEMA IF NOT EXISTS audit   AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS billing AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS notify  AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS policy  AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS sync    AUTHORIZATION helixgitpx;

-- =================== AUDIT (append-only) ===================

CREATE TABLE audit.events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id          UUID,
    actor_kind      TEXT NOT NULL CHECK (actor_kind IN ('user','service','system','ai')),
    actor_id        TEXT,                       -- user uuid or service name
    action          TEXT NOT NULL,              -- 'repo.create','pr.merge','policy.update',…
    resource_kind   TEXT NOT NULL,
    resource_id     TEXT,
    outcome         TEXT NOT NULL CHECK (outcome IN ('success','failure','denied')),
    ip_inet         INET,
    country         CHAR(2),
    user_agent      TEXT,
    session_id      UUID,
    trace_id        TEXT,
    payload         JSONB NOT NULL DEFAULT '{}'::jsonb,
    anchor_root     BYTEA,                      -- merkle root at time of anchoring to Rekor
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
) PARTITION BY RANGE (created_at);

-- Monthly partitions created via cron (example):
-- CREATE TABLE audit.events_2026_04 PARTITION OF audit.events
--     FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

CREATE INDEX audit_events_org_time_idx    ON audit.events (org_id, created_at DESC);
CREATE INDEX audit_events_action_idx      ON audit.events (action, created_at DESC);
CREATE INDEX audit_events_actor_idx       ON audit.events (actor_kind, actor_id, created_at DESC);
CREATE INDEX audit_events_resource_idx    ON audit.events (resource_kind, resource_id);

-- Tamper-evident: deny UPDATE/DELETE to application role.
REVOKE UPDATE, DELETE, TRUNCATE ON audit.events FROM helixgitpx_app;

-- Trigger: prevent UPDATE/DELETE even if caller has privilege.
CREATE OR REPLACE FUNCTION audit_reject_mod() RETURNS trigger AS $$
BEGIN RAISE EXCEPTION 'audit.events is append-only'; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER audit_no_update BEFORE UPDATE ON audit.events
    FOR EACH ROW EXECUTE FUNCTION audit_reject_mod();
CREATE TRIGGER audit_no_delete BEFORE DELETE ON audit.events
    FOR EACH ROW EXECUTE FUNCTION audit_reject_mod();

CREATE TABLE audit.anchors (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    window_from  TIMESTAMPTZ NOT NULL,
    window_to    TIMESTAMPTZ NOT NULL,
    merkle_root  BYTEA NOT NULL,
    record_count BIGINT NOT NULL,
    rekor_uuid   TEXT NOT NULL,
    rekor_url    TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =================== BILLING ==============================

CREATE TABLE billing.plans (
    id            TEXT PRIMARY KEY,           -- 'free','team','enterprise'
    display_name  TEXT NOT NULL,
    description   TEXT,
    price_monthly NUMERIC(12,2),
    currency      CHAR(3) NOT NULL DEFAULT 'USD',
    limits        JSONB NOT NULL DEFAULT '{}'::jsonb,
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO billing.plans (id, display_name, price_monthly, limits)
VALUES
  ('free',      'Free',       0,   '{"repos":25,"upstreams":2,"ai_tokens":100000}'),
  ('team',      'Team',       15,  '{"repos":500,"upstreams":10,"ai_tokens":10000000}'),
  ('enterprise','Enterprise', NULL,'{"repos":-1,"upstreams":-1,"ai_tokens":-1}')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE billing.subscriptions (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id     UUID NOT NULL UNIQUE REFERENCES org.organisations(id) ON DELETE CASCADE,
    plan_id    TEXT NOT NULL REFERENCES billing.plans(id),
    status     TEXT NOT NULL DEFAULT 'active'
               CHECK (status IN ('active','past_due','cancelled','trialing')),
    current_period_start TIMESTAMPTZ NOT NULL DEFAULT date_trunc('month', now()),
    current_period_end   TIMESTAMPTZ NOT NULL DEFAULT date_trunc('month', now()) + INTERVAL '1 month',
    trial_ends_at TIMESTAMPTZ,
    external_customer_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE billing.usage_records (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id     UUID NOT NULL,
    meter      TEXT NOT NULL,                 -- 'ai_tokens','storage_bytes','git_push_bytes',…
    units      BIGINT NOT NULL,
    unit_type  TEXT NOT NULL,                 -- 'tokens','bytes','events',…
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    period_month DATE NOT NULL DEFAULT date_trunc('month', now())::date
) PARTITION BY RANGE (recorded_at);

CREATE INDEX usage_records_org_time_idx    ON billing.usage_records (org_id, recorded_at DESC);
CREATE INDEX usage_records_meter_time_idx  ON billing.usage_records (meter, recorded_at DESC);

CREATE TABLE billing.invoices (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id     UUID NOT NULL REFERENCES org.organisations(id) ON DELETE CASCADE,
    period_from TIMESTAMPTZ NOT NULL,
    period_to   TIMESTAMPTZ NOT NULL,
    subtotal_cents BIGINT NOT NULL,
    tax_cents      BIGINT NOT NULL DEFAULT 0,
    total_cents    BIGINT NOT NULL,
    currency       CHAR(3) NOT NULL DEFAULT 'USD',
    status         TEXT NOT NULL DEFAULT 'open'
                   CHECK (status IN ('open','paid','void','uncollectible')),
    external_invoice_id TEXT,
    issued_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    paid_at        TIMESTAMPTZ
);

-- =================== NOTIFY ===============================

CREATE TABLE notify.channels (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id      UUID REFERENCES org.organisations(id) ON DELETE CASCADE,
    user_id     UUID,
    kind        TEXT NOT NULL
                CHECK (kind IN ('email','webpush','slack','teams','discord','webhook','sms')),
    address_ref TEXT NOT NULL,            -- vault ref, token id, endpoint URL
    verified    BOOLEAN NOT NULL DEFAULT FALSE,
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE notify.subscriptions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id       UUID NOT NULL,
    scope_kind    TEXT NOT NULL CHECK (scope_kind IN ('global','org','repo','pr','issue')),
    scope_id      UUID,
    event_types   TEXT[] NOT NULL DEFAULT '{}',
    channel_id    UUID NOT NULL REFERENCES notify.channels(id) ON DELETE CASCADE,
    filter        JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, scope_kind, COALESCE(scope_id,'00000000-0000-0000-0000-000000000000'::uuid), channel_id)
);
CREATE INDEX notify_subs_user_idx ON notify.subscriptions (user_id);

CREATE TABLE notify.deliveries (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    channel_id     UUID NOT NULL,
    event_id       TEXT NOT NULL,
    status         TEXT NOT NULL CHECK (status IN ('pending','sent','failed','dropped')),
    attempts       INT NOT NULL DEFAULT 0,
    last_error     TEXT,
    first_tried_at TIMESTAMPTZ,
    last_tried_at  TIMESTAMPTZ,
    delivered_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
) PARTITION BY RANGE (created_at);

-- =================== POLICY ===============================

CREATE TABLE policy.bundles (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name           TEXT NOT NULL UNIQUE,
    revision       TEXT NOT NULL,
    digest         BYTEA NOT NULL,
    storage_uri    TEXT NOT NULL,
    signed_by      TEXT,
    active         BOOLEAN NOT NULL DEFAULT TRUE,
    deployed_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE policy.assignments (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    scope_kind   TEXT NOT NULL CHECK (scope_kind IN ('global','org','repo')),
    scope_id     UUID,
    bundle_id    UUID NOT NULL REFERENCES policy.bundles(id),
    priority     INT NOT NULL DEFAULT 100,
    assigned_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (scope_kind, scope_id, bundle_id)
);

CREATE TABLE policy.decisions (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    subject_kind   TEXT NOT NULL,
    subject_id     TEXT NOT NULL,
    action         TEXT NOT NULL,
    resource_kind  TEXT,
    resource_id    TEXT,
    effect         TEXT NOT NULL CHECK (effect IN ('allow','deny')),
    rule           TEXT,
    input_digest   BYTEA,
    eval_ms        INT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
) PARTITION BY RANGE (created_at);

CREATE INDEX policy_decisions_subject_idx ON policy.decisions (subject_kind, subject_id, created_at DESC);

-- =================== SYNC (Temporal-side complement) ======

CREATE TABLE sync.jobs (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id       UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    workflow_id   TEXT NOT NULL UNIQUE,
    run_id        TEXT NOT NULL,
    trigger       TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'queued'
                  CHECK (status IN ('queued','running','succeeded','failed','cancelled','partial')),
    scheduled_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at    TIMESTAMPTZ,
    completed_at  TIMESTAMPTZ,
    error_summary TEXT,
    steps_total   INT,
    steps_done    INT
);
CREATE INDEX sync_jobs_repo_idx ON sync.jobs (repo_id, scheduled_at DESC);

CREATE TABLE sync.steps (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    job_id       UUID NOT NULL REFERENCES sync.jobs(id) ON DELETE CASCADE,
    upstream_id  UUID,
    operation    TEXT NOT NULL,
    status       TEXT NOT NULL
                 CHECK (status IN ('pending','running','succeeded','failed','skipped')),
    error_code   TEXT,
    error_message TEXT,
    started_at   TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ,
    metadata     JSONB NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX sync_steps_job_idx ON sync.steps (job_id);

-- =================== OUTBOXES =============================

CREATE TABLE audit.event_outbox   (LIKE auth.event_outbox INCLUDING ALL);
CREATE TABLE billing.event_outbox (LIKE auth.event_outbox INCLUDING ALL);
CREATE TABLE notify.event_outbox  (LIKE auth.event_outbox INCLUDING ALL);
CREATE TABLE policy.event_outbox  (LIKE auth.event_outbox INCLUDING ALL);
CREATE TABLE sync.event_outbox    (LIKE auth.event_outbox INCLUDING ALL);

-- =================== RLS ==================================

ALTER TABLE audit.events          ENABLE ROW LEVEL SECURITY;
ALTER TABLE billing.usage_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE billing.subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE notify.subscriptions  ENABLE ROW LEVEL SECURITY;
ALTER TABLE policy.decisions      ENABLE ROW LEVEL SECURITY;
ALTER TABLE sync.jobs             ENABLE ROW LEVEL SECURITY;

CREATE POLICY audit_tenant_policy ON audit.events
    USING (org_id IS NULL OR org_id = ANY (helix_current_orgs())
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY billing_usage_tenant ON billing.usage_records
    USING (org_id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY billing_subs_tenant ON billing.subscriptions
    USING (org_id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY notify_subs_tenant ON notify.subscriptions
    USING (user_id = current_setting('helix.user_id', TRUE)::uuid
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY sync_jobs_tenant ON sync.jobs
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = jobs.repo_id
          AND r.org_id = ANY (helix_current_orgs())
    ));
