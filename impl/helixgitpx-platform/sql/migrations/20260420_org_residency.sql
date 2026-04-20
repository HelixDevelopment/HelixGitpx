-- M8 — data residency per org
ALTER TABLE org.organizations
  ADD COLUMN IF NOT EXISTS residency TEXT NOT NULL DEFAULT 'EU'
    CHECK (residency IN ('EU','UK','US'));

CREATE INDEX IF NOT EXISTS organizations_residency_idx
  ON org.organizations (residency);

COMMENT ON COLUMN org.organizations.residency IS
  'Data residency zone. Read by orgteam-service and propagated to downstream services that index org-scoped data.';
