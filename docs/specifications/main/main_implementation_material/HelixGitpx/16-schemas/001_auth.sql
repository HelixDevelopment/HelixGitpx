-- ============================================================
-- 16-schemas/001_auth.sql
-- Bounded context: AUTH / IDENTITY
-- Owned by: auth-service
-- ============================================================
--
-- Conventions:
--   - All PKs are UUID (UUIDv7 via the helper fn uuid_generate_v7()).
--   - Every mutable table has created_at / updated_at TIMESTAMPTZ with trigger.
--   - Tenant scoping is by org_id; row-level security ON by default.
--   - Use INCLUDE for covering indexes where applicable.
--
-- Prerequisites (run once per DB, in a preceding migration):
--   CREATE EXTENSION IF NOT EXISTS pgcrypto;
--   CREATE EXTENSION IF NOT EXISTS pg_trgm;
--   CREATE EXTENSION IF NOT EXISTS btree_gist;
--
-- UUIDv7 helper (reference impl; replace with pg_uuidv7 extension in prod):
--   CREATE OR REPLACE FUNCTION uuid_generate_v7() RETURNS uuid AS $$ ... $$;
--
-- Audit trigger helper:
--   CREATE OR REPLACE FUNCTION helix_set_updated_at() RETURNS trigger AS $$
--   BEGIN NEW.updated_at = now(); RETURN NEW; END; $$ LANGUAGE plpgsql;
-- ============================================================

CREATE SCHEMA IF NOT EXISTS auth AUTHORIZATION helixgitpx;

-- ---------- USERS -------------------------------------------

CREATE TABLE auth.users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    email           CITEXT NOT NULL UNIQUE,
    username        TEXT   NOT NULL UNIQUE CHECK (username ~ '^[a-z0-9][a-z0-9_-]{1,38}[a-z0-9]$'),
    display_name    TEXT,
    avatar_url      TEXT,
    locale          TEXT NOT NULL DEFAULT 'en',
    timezone        TEXT NOT NULL DEFAULT 'UTC',
    status          TEXT NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active','suspended','deactivated','deleted')),
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_enforced    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at   TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX users_username_trgm_idx ON auth.users USING GIN (username gin_trgm_ops);
CREATE INDEX users_email_trgm_idx    ON auth.users USING GIN (email gin_trgm_ops);
CREATE INDEX users_status_idx        ON auth.users (status) WHERE status <> 'active';

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON auth.users
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

-- ---------- IDENTITY PROVIDER LINKS -------------------------

CREATE TABLE auth.identity_providers (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    issuer         TEXT NOT NULL,
    display_name   TEXT NOT NULL,
    client_id      TEXT NOT NULL,
    enabled        BOOLEAN NOT NULL DEFAULT TRUE,
    jwks_uri       TEXT,
    token_endpoint TEXT,
    auth_endpoint  TEXT,
    scopes         TEXT[] NOT NULL DEFAULT ARRAY['openid','profile','email'],
    allowed_domains TEXT[],
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (issuer, client_id)
);

CREATE TABLE auth.user_identities (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES auth.identity_providers(id) ON DELETE RESTRICT,
    subject     TEXT NOT NULL,  -- sub claim from IdP
    email       CITEXT,
    raw_claims  JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ,
    UNIQUE (provider_id, subject)
);
CREATE INDEX user_identities_user_idx ON auth.user_identities (user_id);

-- ---------- SESSIONS ----------------------------------------

CREATE TABLE auth.sessions (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id          UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    refresh_token_hash BYTEA NOT NULL,         -- sha256 of opaque rotating refresh token
    refresh_family_id  UUID NOT NULL,          -- rotation family; reuse → revoke family
    device_id        TEXT NOT NULL,
    device_name      TEXT,
    user_agent       TEXT,
    ip_inet          INET,
    country          CHAR(2),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at       TIMESTAMPTZ NOT NULL,
    revoked_at       TIMESTAMPTZ,
    revoke_reason    TEXT,
    CHECK (expires_at > created_at)
);
CREATE INDEX sessions_user_idx          ON auth.sessions (user_id) WHERE revoked_at IS NULL;
CREATE UNIQUE INDEX sessions_rt_hash_idx ON auth.sessions (refresh_token_hash) WHERE revoked_at IS NULL;
CREATE INDEX sessions_family_idx         ON auth.sessions (refresh_family_id);

-- ---------- MFA FACTORS -------------------------------------

CREATE TABLE auth.mfa_factors (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    kind        TEXT NOT NULL CHECK (kind IN ('totp','webauthn','recovery')),
    label       TEXT,
    -- payload is encrypted via Vault transit; columns here are the ciphertext handle
    enc_ref     TEXT NOT NULL,     -- "vault:transit:mfa:<id>"
    enabled     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ
);
CREATE INDEX mfa_factors_user_idx ON auth.mfa_factors (user_id);

-- ---------- PERSONAL ACCESS TOKENS --------------------------

CREATE TABLE auth.personal_access_tokens (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id      UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    token_hash   BYTEA NOT NULL UNIQUE,   -- sha256 of the opaque secret after prefix
    prefix       TEXT NOT NULL,           -- 'hpxat_' + 6 base62 for display
    scopes       TEXT[] NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ,
    CHECK (expires_at > created_at)
);
CREATE INDEX pats_user_idx ON auth.personal_access_tokens (user_id);

-- ---------- LOGIN ATTEMPTS / ANOMALIES ---------------------

CREATE TABLE auth.login_attempts (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    email          CITEXT,
    result         TEXT NOT NULL CHECK (result IN (
                    'success','bad_password','locked','mfa_required','mfa_failed','blocked'
                   )),
    method         TEXT NOT NULL,
    ip_inet        INET NOT NULL,
    country        CHAR(2),
    user_agent     TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
) PARTITION BY RANGE (created_at);

-- Monthly partitions are created by a cron job; example for April 2026:
-- CREATE TABLE auth.login_attempts_2026_04 PARTITION OF auth.login_attempts
--     FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

CREATE INDEX login_attempts_user_idx ON auth.login_attempts (user_id, created_at DESC);
CREATE INDEX login_attempts_ip_idx   ON auth.login_attempts (ip_inet, created_at DESC);

-- ---------- RLS ---------------------------------------------

ALTER TABLE auth.users                    ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.user_identities          ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.sessions                 ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.mfa_factors              ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.personal_access_tokens   ENABLE ROW LEVEL SECURITY;

-- Admin role bypasses via BYPASSRLS.
CREATE POLICY user_self_policy ON auth.users
    USING (id = current_setting('helix.user_id', TRUE)::uuid
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY user_self_identities ON auth.user_identities
    USING (user_id = current_setting('helix.user_id', TRUE)::uuid
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY user_self_sessions ON auth.sessions
    USING (user_id = current_setting('helix.user_id', TRUE)::uuid
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY user_self_mfa ON auth.mfa_factors
    USING (user_id = current_setting('helix.user_id', TRUE)::uuid
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY user_self_pat ON auth.personal_access_tokens
    USING (user_id = current_setting('helix.user_id', TRUE)::uuid
           OR current_setting('helix.role', TRUE) = 'admin');

-- ---------- OUTBOX ------------------------------------------

CREATE TABLE auth.event_outbox (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    aggregate_kind  TEXT  NOT NULL,
    aggregate_id    UUID  NOT NULL,
    event_type      TEXT  NOT NULL,
    schema_version  INT   NOT NULL DEFAULT 1,
    payload         BYTEA NOT NULL,
    headers         JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX auth_outbox_time_idx ON auth.event_outbox (occurred_at);

-- ---------- GRANTS ------------------------------------------

GRANT USAGE ON SCHEMA auth TO helixgitpx_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA auth TO helixgitpx_app;
GRANT SELECT ON ALL TABLES IN SCHEMA auth TO helixgitpx_readonly;
REVOKE TRUNCATE ON ALL TABLES IN SCHEMA auth FROM PUBLIC;
