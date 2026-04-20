-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS org.orgs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug CITEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    CREATE TYPE team.role_enum AS ENUM ('viewer','member','admin','owner');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS team.teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES org.orgs(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES team.teams(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, slug)
);
CREATE INDEX IF NOT EXISTS ix_teams_org    ON team.teams (org_id);
CREATE INDEX IF NOT EXISTS ix_teams_parent ON team.teams (parent_id);

CREATE TABLE IF NOT EXISTS team.memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES team.teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role team.role_enum NOT NULL DEFAULT 'viewer',
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);
CREATE INDEX IF NOT EXISTS ix_memberships_user ON team.memberships (user_id);

-- effective_role: max role across self + ancestor teams where user is a member.
CREATE OR REPLACE FUNCTION team.effective_role(p_team UUID, p_user UUID)
RETURNS team.role_enum LANGUAGE sql STABLE AS $$
  WITH RECURSIVE chain AS (
    SELECT id, parent_id FROM team.teams WHERE id = p_team
    UNION ALL
    SELECT t.id, t.parent_id FROM team.teams t JOIN chain c ON t.id = c.parent_id
  ), roles AS (
    SELECT m.role FROM team.memberships m
    WHERE m.user_id = p_user AND m.team_id IN (SELECT id FROM chain)
  )
  SELECT CASE
    WHEN EXISTS (SELECT 1 FROM roles WHERE role = 'owner')  THEN 'owner'::team.role_enum
    WHEN EXISTS (SELECT 1 FROM roles WHERE role = 'admin')  THEN 'admin'::team.role_enum
    WHEN EXISTS (SELECT 1 FROM roles WHERE role = 'member') THEN 'member'::team.role_enum
    WHEN EXISTS (SELECT 1 FROM roles WHERE role = 'viewer') THEN 'viewer'::team.role_enum
    ELSE NULL
  END;
$$;

ALTER TABLE org.orgs ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.memberships ENABLE ROW LEVEL SECURITY;

CREATE POLICY org_orgs_all ON org.orgs USING (TRUE);
CREATE POLICY team_teams_all ON team.teams USING (TRUE);
CREATE POLICY team_members_all ON team.memberships USING (TRUE);

-- +goose Down
DROP FUNCTION IF EXISTS team.effective_role;
DROP TABLE IF EXISTS team.memberships;
DROP TABLE IF EXISTS team.teams;
DROP TYPE  IF EXISTS team.role_enum;
DROP TABLE IF EXISTS org.orgs;
