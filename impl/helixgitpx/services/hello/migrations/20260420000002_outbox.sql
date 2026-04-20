-- +goose Up
CREATE TABLE IF NOT EXISTS hello.outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type TEXT NOT NULL DEFAULT 'hello',
    aggregate_id TEXT NOT NULL,
    topic TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_outbox_events_created ON hello.outbox_events (created_at);

CREATE PUBLICATION IF NOT EXISTS helix_hello_outbox FOR TABLE hello.outbox_events;

ALTER TABLE hello.outbox_events ENABLE ROW LEVEL SECURITY;
CREATE POLICY hello_outbox_events_all ON hello.outbox_events USING (TRUE);

-- +goose Down
DROP PUBLICATION IF EXISTS helix_hello_outbox;
DROP TABLE IF EXISTS hello.outbox_events;
