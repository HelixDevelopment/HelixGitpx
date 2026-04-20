-- +goose Up
CREATE TABLE IF NOT EXISTS sync.sync_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id TEXT NOT NULL UNIQUE,
    repo_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL,
    attempts INT NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS ix_sync_runs_repo ON sync.sync_runs (repo_id);

ALTER TABLE sync.sync_runs ENABLE ROW LEVEL SECURITY;
CREATE POLICY sync_all ON sync.sync_runs USING (TRUE);

-- +goose Down
DROP TABLE IF EXISTS sync.sync_runs;
