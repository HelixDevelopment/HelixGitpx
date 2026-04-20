-- +goose Up
CREATE TABLE IF NOT EXISTS collab.crdt_ops (
    repo_id UUID NOT NULL,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    actor TEXT NOT NULL,
    seq BIGINT NOT NULL,
    op_bytes BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (repo_id, aggregate_type, aggregate_id, seq)
);
CREATE INDEX IF NOT EXISTS ix_collab_ops_time ON collab.crdt_ops (created_at);

CREATE TABLE IF NOT EXISTS collab.outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id TEXT NOT NULL,
    topic TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE PUBLICATION IF NOT EXISTS helix_collab_outbox FOR TABLE collab.outbox_events;

ALTER TABLE collab.crdt_ops      ENABLE ROW LEVEL SECURITY;
ALTER TABLE collab.outbox_events ENABLE ROW LEVEL SECURITY;
CREATE POLICY collab_ops_all    ON collab.crdt_ops      USING (TRUE);
CREATE POLICY collab_outbox_all ON collab.outbox_events USING (TRUE);

-- +goose Down
DROP PUBLICATION IF EXISTS helix_collab_outbox;
DROP TABLE IF EXISTS collab.outbox_events;
DROP TABLE IF EXISTS collab.crdt_ops;
