-- +goose Up
CREATE TABLE IF NOT EXISTS audit.events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_user_id TEXT,
    actor_ip INET,
    action TEXT NOT NULL,
    target TEXT NOT NULL,
    details JSONB
);
CREATE INDEX IF NOT EXISTS ix_audit_at ON audit.events (at);

CREATE OR REPLACE RULE audit_events_no_update AS ON UPDATE TO audit.events DO INSTEAD NOTHING;
CREATE OR REPLACE RULE audit_events_no_delete AS ON DELETE TO audit.events DO INSTEAD NOTHING;

CREATE OR REPLACE FUNCTION audit.append_event(
    p_at TIMESTAMPTZ,
    p_actor_user_id TEXT,
    p_actor_ip TEXT,
    p_action TEXT,
    p_target TEXT,
    p_details JSONB
) RETURNS UUID LANGUAGE plpgsql SECURITY DEFINER AS $$
DECLARE new_id UUID;
BEGIN
    INSERT INTO audit.events(at, actor_user_id, actor_ip, action, target, details)
    VALUES (p_at, p_actor_user_id, NULLIF(p_actor_ip, '')::inet, p_action, p_target, p_details)
    RETURNING id INTO new_id;
    RETURN new_id;
END $$;

REVOKE ALL ON FUNCTION audit.append_event FROM PUBLIC;
GRANT EXECUTE ON FUNCTION audit.append_event TO audit_svc;

CREATE TABLE IF NOT EXISTS audit.anchors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_start TIMESTAMPTZ NOT NULL,
    period_end   TIMESTAMPTZ NOT NULL,
    merkle_root BYTEA NOT NULL,
    external_tx_id TEXT,
    anchored_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(period_start, period_end)
);

ALTER TABLE audit.events  ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit.anchors ENABLE ROW LEVEL SECURITY;
CREATE POLICY audit_events_all  ON audit.events  USING (TRUE);
CREATE POLICY audit_anchors_all ON audit.anchors USING (TRUE);

-- +goose Down
DROP FUNCTION IF EXISTS audit.append_event;
DROP TABLE IF EXISTS audit.anchors;
DROP TABLE IF EXISTS audit.events;
