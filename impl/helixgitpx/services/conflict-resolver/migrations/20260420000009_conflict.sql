-- +goose Up
CREATE TABLE IF NOT EXISTS conflict.resolutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID NOT NULL,
    ref TEXT NOT NULL,
    chosen_sha TEXT NOT NULL,
    policy_verdict JSONB NOT NULL,
    resolved_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS ix_conflict_resolutions_repo ON conflict.resolutions (repo_id, resolved_at);

ALTER TABLE conflict.resolutions ENABLE ROW LEVEL SECURITY;
CREATE POLICY conflict_all ON conflict.resolutions USING (TRUE);

-- +goose Down
DROP TABLE IF EXISTS conflict.resolutions;
