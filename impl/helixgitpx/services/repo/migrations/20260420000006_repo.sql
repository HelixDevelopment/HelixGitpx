-- +goose Up
CREATE TABLE IF NOT EXISTS repo.repos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL,
    slug TEXT NOT NULL,
    default_branch TEXT NOT NULL DEFAULT 'main',
    lfs_enabled BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, slug)
);
CREATE INDEX IF NOT EXISTS ix_repos_org ON repo.repos (org_id);

CREATE TABLE IF NOT EXISTS repo.refs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID NOT NULL REFERENCES repo.repos(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sha TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(repo_id, name)
);

CREATE TABLE IF NOT EXISTS repo.branch_protections (
    repo_id UUID NOT NULL REFERENCES repo.repos(id) ON DELETE CASCADE,
    pattern TEXT NOT NULL,
    require_signed BOOLEAN NOT NULL DEFAULT false,
    required_reviewers INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (repo_id, pattern)
);

CREATE TABLE IF NOT EXISTS repo.lfs_objects (
    repo_id UUID NOT NULL REFERENCES repo.repos(id) ON DELETE CASCADE,
    oid CHAR(64) NOT NULL,
    size BIGINT NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (repo_id, oid)
);

CREATE TABLE IF NOT EXISTS repo.outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id TEXT NOT NULL,
    topic TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE PUBLICATION IF NOT EXISTS helix_repo_outbox FOR TABLE repo.outbox_events;

ALTER TABLE repo.repos              ENABLE ROW LEVEL SECURITY;
ALTER TABLE repo.refs               ENABLE ROW LEVEL SECURITY;
ALTER TABLE repo.branch_protections ENABLE ROW LEVEL SECURITY;
ALTER TABLE repo.lfs_objects        ENABLE ROW LEVEL SECURITY;
CREATE POLICY repo_all    ON repo.repos              USING (TRUE);
CREATE POLICY refs_all    ON repo.refs               USING (TRUE);
CREATE POLICY prot_all    ON repo.branch_protections USING (TRUE);
CREATE POLICY lfs_all     ON repo.lfs_objects        USING (TRUE);

-- +goose Down
DROP PUBLICATION IF EXISTS helix_repo_outbox;
DROP TABLE IF EXISTS repo.outbox_events;
DROP TABLE IF EXISTS repo.lfs_objects;
DROP TABLE IF EXISTS repo.branch_protections;
DROP TABLE IF EXISTS repo.refs;
DROP TABLE IF EXISTS repo.repos;
