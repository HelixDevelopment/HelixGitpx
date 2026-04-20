-- ============================================================
-- 16-schemas/002_repo.sql
-- Bounded context: REPO (organisations, repos, refs, bindings)
-- Owned by: org-service, repo-service, upstream-service
-- ============================================================

CREATE SCHEMA IF NOT EXISTS org  AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS repo AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS upstream AUTHORIZATION helixgitpx;

-- =================== ORG ====================================

CREATE TABLE org.organisations (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    slug               TEXT NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9][a-z0-9-]{0,38}[a-z0-9]$'),
    display_name       TEXT NOT NULL,
    description        TEXT,
    default_visibility TEXT NOT NULL DEFAULT 'private'
                       CHECK (default_visibility IN ('public','private','internal')),
    billing_plan       TEXT NOT NULL DEFAULT 'free',
    primary_region     TEXT NOT NULL DEFAULT 'eu-west-1',
    settings           JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ
);
CREATE INDEX organisations_slug_trgm_idx ON org.organisations USING GIN (slug gin_trgm_ops);

CREATE TRIGGER organisations_updated_at BEFORE UPDATE ON org.organisations
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE org.teams (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id     UUID NOT NULL REFERENCES org.organisations(id) ON DELETE CASCADE,
    slug       TEXT NOT NULL,
    display_name TEXT NOT NULL,
    description  TEXT,
    parent_id  UUID REFERENCES org.teams(id) ON DELETE CASCADE,
    visibility TEXT NOT NULL DEFAULT 'private' CHECK (visibility IN ('public','private','secret')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_id, slug)
);
CREATE INDEX teams_parent_idx ON org.teams (parent_id);

CREATE TRIGGER teams_updated_at BEFORE UPDATE ON org.teams
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE org.memberships (
    org_id     UUID NOT NULL REFERENCES org.organisations(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL,                    -- FK-lite (cross-schema) to auth.users
    team_id    UUID REFERENCES org.teams(id) ON DELETE CASCADE,
    role       TEXT NOT NULL CHECK (role IN ('owner','admin','maintainer','developer','viewer')),
    invited_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (org_id, user_id, COALESCE(team_id, '00000000-0000-0000-0000-000000000000'::uuid))
);
CREATE INDEX memberships_user_idx ON org.memberships (user_id);
CREATE INDEX memberships_team_idx ON org.memberships (team_id);

-- =================== REPO ===================================

CREATE TABLE repo.repositories (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id             UUID NOT NULL REFERENCES org.organisations(id) ON DELETE CASCADE,
    slug               TEXT NOT NULL CHECK (slug ~ '^[a-zA-Z0-9._-]{1,100}$'),
    display_name       TEXT NOT NULL,
    description        TEXT,
    visibility         TEXT NOT NULL DEFAULT 'private'
                       CHECK (visibility IN ('public','private','internal')),
    default_branch     TEXT NOT NULL DEFAULT 'main',
    primary_upstream   TEXT,                -- provider_kind, e.g. 'github'
    size_bytes         BIGINT NOT NULL DEFAULT 0,
    archived           BOOLEAN NOT NULL DEFAULT FALSE,
    template           BOOLEAN NOT NULL DEFAULT FALSE,
    topics             TEXT[] NOT NULL DEFAULT '{}',
    license            TEXT,
    settings           JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ,
    UNIQUE (org_id, slug)
);
CREATE INDEX repositories_org_idx    ON repo.repositories (org_id) WHERE deleted_at IS NULL;
CREATE INDEX repositories_topics_gin ON repo.repositories USING GIN (topics);
CREATE INDEX repositories_slug_trgm  ON repo.repositories USING GIN (slug gin_trgm_ops);

CREATE TRIGGER repositories_updated_at BEFORE UPDATE ON repo.repositories
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE repo.refs (
    repo_id         UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,               -- e.g. refs/heads/main
    kind            TEXT NOT NULL CHECK (kind IN ('branch','tag','other')),
    target_sha      CHAR(40) NOT NULL,
    is_protected    BOOLEAN NOT NULL DEFAULT FALSE,
    last_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_origin     TEXT,                        -- 'local','upstream:github', ...
    last_actor      UUID,
    PRIMARY KEY (repo_id, name)
);
CREATE INDEX refs_repo_updated_idx ON repo.refs (repo_id, last_updated_at DESC);

CREATE TABLE repo.branch_protections (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id              UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    pattern              TEXT NOT NULL,          -- glob, e.g. 'main' or 'release/*'
    require_signed_commits BOOLEAN NOT NULL DEFAULT FALSE,
    require_linear_history BOOLEAN NOT NULL DEFAULT FALSE,
    required_reviews     INT NOT NULL DEFAULT 0 CHECK (required_reviews >= 0),
    required_status_checks TEXT[] NOT NULL DEFAULT '{}',
    dismiss_stale_reviews BOOLEAN NOT NULL DEFAULT TRUE,
    block_force_push     BOOLEAN NOT NULL DEFAULT TRUE,
    block_deletions      BOOLEAN NOT NULL DEFAULT TRUE,
    allow_up_to_date_merge_only BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (repo_id, pattern)
);
CREATE INDEX bp_repo_idx ON repo.branch_protections (repo_id);

-- Content-addressable pack metadata (objects stored in object store).
CREATE TABLE repo.pack_blobs (
    repo_id    UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    sha256     BYTEA NOT NULL,
    size_bytes BIGINT NOT NULL,
    object_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (repo_id, sha256)
);

-- LFS pointer index (the bytes themselves live in object store).
CREATE TABLE repo.lfs_objects (
    repo_id    UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    oid        CHAR(64) NOT NULL,        -- sha256 hex
    size_bytes BIGINT NOT NULL,
    object_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (repo_id, oid)
);

-- =================== UPSTREAM ==============================

CREATE TABLE upstream.upstreams (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    org_id         UUID NOT NULL REFERENCES org.organisations(id) ON DELETE CASCADE,
    provider_kind  TEXT NOT NULL,        -- 'github','gitlab','gitee',…
    display_name   TEXT NOT NULL,
    base_url       TEXT NOT NULL,
    auth_method    TEXT NOT NULL
                   CHECK (auth_method IN ('oauth2','pat','app_token','ssh_key','sigv4')),
    credential_ref TEXT NOT NULL,        -- 'vault:secret:upstreams/<org>/<id>'
    shadow_mode    BOOLEAN NOT NULL DEFAULT TRUE,
    enabled        BOOLEAN NOT NULL DEFAULT FALSE,
    rate_limit_remaining INT,
    health_status  TEXT NOT NULL DEFAULT 'unknown'
                   CHECK (health_status IN ('unknown','healthy','degraded','unhealthy')),
    last_health_check_at TIMESTAMPTZ,
    capabilities   JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_id, provider_kind, base_url)
);
CREATE INDEX upstreams_org_idx ON upstream.upstreams (org_id);

CREATE TRIGGER upstreams_updated_at BEFORE UPDATE ON upstream.upstreams
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE upstream.bindings (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id            UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    upstream_id        UUID NOT NULL REFERENCES upstream.upstreams(id) ON DELETE CASCADE,
    remote_owner       TEXT NOT NULL,
    remote_name        TEXT NOT NULL,
    push_enabled       BOOLEAN NOT NULL DEFAULT TRUE,
    pull_enabled       BOOLEAN NOT NULL DEFAULT TRUE,
    sync_issues        BOOLEAN NOT NULL DEFAULT TRUE,
    sync_prs           BOOLEAN NOT NULL DEFAULT TRUE,
    sync_releases      BOOLEAN NOT NULL DEFAULT TRUE,
    last_synced_at     TIMESTAMPTZ,
    last_sync_status   TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (upstream_id, remote_owner, remote_name),
    UNIQUE (repo_id, upstream_id)
);
CREATE INDEX bindings_repo_idx ON upstream.bindings (repo_id);

CREATE TABLE upstream.credential_rotations (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    upstream_id   UUID NOT NULL REFERENCES upstream.upstreams(id) ON DELETE CASCADE,
    rotated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    rotated_by    UUID,
    reason        TEXT,
    previous_hash BYTEA
);
CREATE INDEX creds_rot_upstream_idx ON upstream.credential_rotations (upstream_id, rotated_at DESC);

-- =================== RLS ===================================

ALTER TABLE org.organisations        ENABLE ROW LEVEL SECURITY;
ALTER TABLE org.teams                ENABLE ROW LEVEL SECURITY;
ALTER TABLE org.memberships          ENABLE ROW LEVEL SECURITY;
ALTER TABLE repo.repositories        ENABLE ROW LEVEL SECURITY;
ALTER TABLE repo.refs                ENABLE ROW LEVEL SECURITY;
ALTER TABLE repo.branch_protections  ENABLE ROW LEVEL SECURITY;
ALTER TABLE upstream.upstreams       ENABLE ROW LEVEL SECURITY;
ALTER TABLE upstream.bindings        ENABLE ROW LEVEL SECURITY;

-- Tenant filter: `helix.org_ids` is a CSV of UUIDs.
CREATE OR REPLACE FUNCTION helix_current_orgs() RETURNS uuid[] AS $$
    SELECT string_to_array(
             coalesce(current_setting('helix.org_ids', TRUE), ''),
             ','
           )::uuid[];
$$ LANGUAGE sql STABLE;

CREATE POLICY org_tenant_policy ON org.organisations
    USING (id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY team_tenant_policy ON org.teams
    USING (org_id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY membership_tenant_policy ON org.memberships
    USING (org_id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY repo_tenant_policy ON repo.repositories
    USING (org_id = ANY (helix_current_orgs())
           OR visibility = 'public'
           OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY ref_tenant_policy ON repo.refs
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = refs.repo_id
          AND (r.org_id = ANY (helix_current_orgs()) OR r.visibility = 'public')
    ));

CREATE POLICY bp_tenant_policy ON repo.branch_protections
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = branch_protections.repo_id
          AND r.org_id = ANY (helix_current_orgs())
    ));

CREATE POLICY upstream_tenant_policy ON upstream.upstreams
    USING (org_id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');

CREATE POLICY binding_tenant_policy ON upstream.bindings
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = bindings.repo_id
          AND r.org_id = ANY (helix_current_orgs())
    ));

-- =================== OUTBOXES ==============================

CREATE TABLE repo.event_outbox (LIKE auth.event_outbox INCLUDING ALL);
CREATE TABLE org.event_outbox  (LIKE auth.event_outbox INCLUDING ALL);
CREATE TABLE upstream.event_outbox (LIKE auth.event_outbox INCLUDING ALL);
