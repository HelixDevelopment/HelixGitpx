-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS auth.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject TEXT NOT NULL UNIQUE,
    email CITEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    user_agent TEXT NOT NULL DEFAULT '',
    ip INET
);
CREATE INDEX IF NOT EXISTS ix_sessions_user ON auth.sessions (user_id);

CREATE TABLE IF NOT EXISTS auth.pats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    hashed_secret BYTEA NOT NULL,
    scopes JSONB NOT NULL DEFAULT '[]'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS ix_pats_user ON auth.pats (user_id);

DO $$
BEGIN
    CREATE TYPE auth.mfa_kind AS ENUM ('totp','fido2');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS auth.mfa_factors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    kind auth.mfa_kind NOT NULL,
    secret_or_pubkey BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS ix_mfa_user ON auth.mfa_factors (user_id);

ALTER TABLE auth.users         ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.sessions      ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.pats          ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth.mfa_factors   ENABLE ROW LEVEL SECURITY;

CREATE POLICY auth_users_all     ON auth.users         USING (TRUE);
CREATE POLICY auth_sessions_all  ON auth.sessions      USING (TRUE);
CREATE POLICY auth_pats_all      ON auth.pats          USING (TRUE);
CREATE POLICY auth_mfa_all       ON auth.mfa_factors   USING (TRUE);

-- +goose Down
DROP TABLE IF EXISTS auth.mfa_factors;
DROP TYPE  IF EXISTS auth.mfa_kind;
DROP TABLE IF EXISTS auth.pats;
DROP TABLE IF EXISTS auth.sessions;
DROP TABLE IF EXISTS auth.users;
