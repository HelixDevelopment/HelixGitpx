-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

DO $$
BEGIN
    CREATE TYPE upstream.provider AS ENUM ('github','gitlab','gitea');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE upstream.direction AS ENUM ('push','fetch','mirror');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS upstream.upstreams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug CITEXT NOT NULL UNIQUE,
    provider upstream.provider NOT NULL,
    base_url TEXT NOT NULL,
    vault_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS upstream.bindings (
    repo_id UUID NOT NULL,
    upstream_id UUID NOT NULL REFERENCES upstream.upstreams(id) ON DELETE CASCADE,
    remote_name TEXT NOT NULL,
    direction upstream.direction NOT NULL,
    last_sync_at TIMESTAMPTZ,
    PRIMARY KEY (repo_id, upstream_id, remote_name)
);

ALTER TABLE upstream.upstreams ENABLE ROW LEVEL SECURITY;
ALTER TABLE upstream.bindings  ENABLE ROW LEVEL SECURITY;
CREATE POLICY upstream_all ON upstream.upstreams USING (TRUE);
CREATE POLICY binding_all  ON upstream.bindings  USING (TRUE);

-- +goose Down
DROP TABLE IF EXISTS upstream.bindings;
DROP TABLE IF EXISTS upstream.upstreams;
DROP TYPE  IF EXISTS upstream.direction;
DROP TYPE  IF EXISTS upstream.provider;
