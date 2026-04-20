-- ============================================================
-- 16-schemas/003_conflict_collab.sql
-- Bounded contexts: CONFLICT, COLLAB (PRs, issues, comments), AI
-- ============================================================

CREATE SCHEMA IF NOT EXISTS conflict AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS collab   AUTHORIZATION helixgitpx;
CREATE SCHEMA IF NOT EXISTS ai       AUTHORIZATION helixgitpx;

-- =================== CONFLICT ==============================

CREATE TABLE conflict.conflict_cases (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id         UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    kind            TEXT NOT NULL
                    CHECK (kind IN ('ref_divergence','rename_collision','metadata_concurrent',
                                    'pr_state','tag_collision','lfs_divergence',
                                    'workflow_status','release_asset','other')),
    subject         TEXT NOT NULL,          -- e.g. 'refs/heads/main' or 'issue/42/labels'
    upstream_id     UUID REFERENCES upstream.upstreams(id) ON DELETE SET NULL,
    status          TEXT NOT NULL DEFAULT 'detected'
                    CHECK (status IN ('detected','proposed','auto_applying','applied',
                                      'escalated','human_resolving','resolved','cancelled')),
    severity        TEXT NOT NULL DEFAULT 'normal'
                    CHECK (severity IN ('low','normal','high','critical')),
    left_sha        CHAR(40),
    right_sha       CHAR(40),
    base_sha        CHAR(40),
    snapshot_ref    TEXT,                  -- tmp branch name if any
    detected_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at     TIMESTAMPTZ,
    resolved_by     UUID,
    undo_until      TIMESTAMPTZ,
    undo_plan       JSONB,
    metadata        JSONB NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX conflict_cases_repo_idx   ON conflict.conflict_cases (repo_id, status, detected_at DESC);
CREATE INDEX conflict_cases_status_idx ON conflict.conflict_cases (status) WHERE status IN ('detected','escalated','human_resolving');
CREATE INDEX conflict_cases_kind_idx   ON conflict.conflict_cases (kind);

CREATE TABLE conflict.resolutions (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    case_id        UUID NOT NULL REFERENCES conflict.conflict_cases(id) ON DELETE CASCADE,
    strategy       TEXT NOT NULL,
    decided_by     TEXT NOT NULL CHECK (decided_by IN ('policy','crdt','ai','human')),
    actor_id       UUID,
    confidence     NUMERIC(4,3) CHECK (confidence >= 0 AND confidence <= 1),
    apply_plan     JSONB NOT NULL,
    rationale      TEXT,
    applied_at     TIMESTAMPTZ,
    apply_status   TEXT NOT NULL DEFAULT 'pending'
                   CHECK (apply_status IN ('pending','applying','applied','failed','rolled_back')),
    failure_reason TEXT,
    model_version  TEXT,                    -- if decided_by=ai
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX resolutions_case_idx ON conflict.resolutions (case_id);

CREATE TABLE conflict.ai_feedback (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    resolution_id  UUID REFERENCES conflict.resolutions(id) ON DELETE SET NULL,
    case_id        UUID REFERENCES conflict.conflict_cases(id) ON DELETE CASCADE,
    outcome        TEXT NOT NULL CHECK (outcome IN ('accepted','rejected','edited','ignored')),
    edit_distance  INT,
    rating         SMALLINT CHECK (rating BETWEEN 1 AND 5),
    comment        TEXT,
    user_id        UUID,
    model_version  TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ai_feedback_model_idx ON conflict.ai_feedback (model_version, created_at DESC);

CREATE TABLE conflict.event_outbox (LIKE auth.event_outbox INCLUDING ALL);

-- =================== COLLAB (PRs / Issues / Comments) ======

CREATE TABLE collab.pull_requests (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id        UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    number         INT NOT NULL,
    title          TEXT NOT NULL,
    body           TEXT,
    state          TEXT NOT NULL DEFAULT 'open'
                   CHECK (state IN ('draft','open','closed','merged')),
    head_ref       TEXT NOT NULL,
    head_sha       CHAR(40) NOT NULL,
    base_ref       TEXT NOT NULL,
    base_sha       CHAR(40) NOT NULL,
    author_id      UUID NOT NULL,
    assignees      UUID[] NOT NULL DEFAULT '{}',
    labels         TEXT[] NOT NULL DEFAULT '{}',
    milestone_id   UUID,
    mergeable      TEXT CHECK (mergeable IN ('clean','dirty','blocked','unknown') OR mergeable IS NULL),
    merged_at      TIMESTAMPTZ,
    merged_by      UUID,
    merge_commit   CHAR(40),
    merge_strategy TEXT CHECK (merge_strategy IN ('merge','squash','rebase') OR merge_strategy IS NULL),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at      TIMESTAMPTZ,
    UNIQUE (repo_id, number)
);
CREATE INDEX pr_repo_state_idx ON collab.pull_requests (repo_id, state, updated_at DESC);
CREATE INDEX pr_author_idx     ON collab.pull_requests (author_id);
CREATE INDEX pr_labels_gin     ON collab.pull_requests USING GIN (labels);

CREATE TRIGGER pr_updated_at BEFORE UPDATE ON collab.pull_requests
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE collab.reviews (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    pr_id       UUID NOT NULL REFERENCES collab.pull_requests(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL,
    state       TEXT NOT NULL
                CHECK (state IN ('pending','commented','approved','changes_requested','dismissed')),
    body        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at TIMESTAMPTZ
);
CREATE INDEX reviews_pr_idx ON collab.reviews (pr_id, state);

CREATE TABLE collab.issues (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id      UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    number       INT NOT NULL,
    title        TEXT NOT NULL,
    body_doc_id  UUID,                                     -- Automerge doc id (collab.crdt_docs)
    body_snapshot TEXT,                                    -- periodic snapshot for fast reads
    state        TEXT NOT NULL DEFAULT 'open'
                 CHECK (state IN ('open','closed')),
    author_id    UUID NOT NULL,
    assignees    UUID[] NOT NULL DEFAULT '{}',
    labels       TEXT[] NOT NULL DEFAULT '{}',
    milestone_id UUID,
    locked       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at    TIMESTAMPTZ,
    UNIQUE (repo_id, number)
);
CREATE INDEX issues_repo_state_idx ON collab.issues (repo_id, state, updated_at DESC);
CREATE INDEX issues_labels_gin     ON collab.issues USING GIN (labels);

CREATE TRIGGER issues_updated_at BEFORE UPDATE ON collab.issues
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE collab.comments (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    parent_kind  TEXT NOT NULL CHECK (parent_kind IN ('pr','issue','review','commit')),
    parent_id    UUID NOT NULL,                         -- FK-lite to collab.*
    author_id    UUID NOT NULL,
    body         TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX comments_parent_idx ON collab.comments (parent_kind, parent_id, created_at);

CREATE TRIGGER comments_updated_at BEFORE UPDATE ON collab.comments
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE collab.labels (
    org_id     UUID NOT NULL REFERENCES org.organisations(id) ON DELETE CASCADE,
    repo_id    UUID REFERENCES repo.repositories(id) ON DELETE CASCADE, -- NULL = org-wide
    name       TEXT NOT NULL,
    color_hex  CHAR(7) NOT NULL DEFAULT '#9acd32',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (org_id, COALESCE(repo_id, '00000000-0000-0000-0000-000000000000'::uuid), name)
);

CREATE TABLE collab.milestones (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id    UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    title      TEXT NOT NULL,
    description TEXT,
    state      TEXT NOT NULL DEFAULT 'open' CHECK (state IN ('open','closed')),
    due_on     DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (repo_id, title)
);

CREATE TABLE collab.releases (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    repo_id      UUID NOT NULL REFERENCES repo.repositories(id) ON DELETE CASCADE,
    tag_name     TEXT NOT NULL,
    name         TEXT,
    body         TEXT,
    target_sha   CHAR(40) NOT NULL,
    draft        BOOLEAN NOT NULL DEFAULT FALSE,
    prerelease   BOOLEAN NOT NULL DEFAULT FALSE,
    author_id    UUID,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (repo_id, tag_name)
);
CREATE INDEX releases_repo_idx ON collab.releases (repo_id, published_at DESC NULLS LAST);

CREATE TABLE collab.release_assets (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    release_id UUID NOT NULL REFERENCES collab.releases(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    mime_type  TEXT,
    size_bytes BIGINT NOT NULL,
    sha256     BYTEA NOT NULL,
    object_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (release_id, name)
);

-- Automerge CRDT docs for issue body + labels + milestones (concurrent-safe).
CREATE TABLE collab.crdt_docs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    kind            TEXT NOT NULL CHECK (kind IN ('issue_body','labels','milestones','assignees')),
    parent_kind     TEXT NOT NULL,
    parent_id       UUID NOT NULL,
    doc_bytes       BYTEA NOT NULL,               -- compact Automerge binary
    heads           TEXT[] NOT NULL DEFAULT '{}', -- Automerge head hashes
    version         BIGINT NOT NULL DEFAULT 1,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX crdt_docs_parent_kind ON collab.crdt_docs (parent_kind, parent_id, kind);

CREATE TRIGGER crdt_docs_updated_at BEFORE UPDATE ON collab.crdt_docs
    FOR EACH ROW EXECUTE FUNCTION helix_set_updated_at();

CREATE TABLE collab.event_outbox (LIKE auth.event_outbox INCLUDING ALL);

-- =================== AI ====================================

CREATE TABLE ai.model_registry (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name            TEXT NOT NULL,
    version         TEXT NOT NULL,
    purpose         TEXT NOT NULL CHECK (purpose IN (
                        'conflict_resolution','pr_summary','review','label_suggest',
                        'semantic_search','chat','embedding','translation'
                    )),
    base_model      TEXT NOT NULL,
    adapter_kind    TEXT CHECK (adapter_kind IN ('lora','full','none')),
    storage_uri     TEXT NOT NULL,
    eval_score      NUMERIC(5,4),
    context_window  INT,
    active          BOOLEAN NOT NULL DEFAULT FALSE,
    shadow          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    retired_at      TIMESTAMPTZ,
    metadata        JSONB NOT NULL DEFAULT '{}'::jsonb,
    UNIQUE (name, version)
);

CREATE TABLE ai.prompt_runs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    task            TEXT NOT NULL,
    org_id          UUID NOT NULL,
    user_id         UUID,
    prompt_version  TEXT NOT NULL,
    model_id        UUID NOT NULL REFERENCES ai.model_registry(id),
    input_tokens    INT NOT NULL,
    output_tokens   INT NOT NULL,
    duration_ms     INT NOT NULL,
    confidence      NUMERIC(4,3),
    cost_usd        NUMERIC(10,6),
    outcome         TEXT,
    metadata        JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
) PARTITION BY RANGE (created_at);

CREATE INDEX prompt_runs_org_time_idx    ON ai.prompt_runs (org_id, created_at DESC);
CREATE INDEX prompt_runs_model_time_idx  ON ai.prompt_runs (model_id, created_at DESC);

CREATE TABLE ai.feedback_records (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    prompt_run_id  UUID NOT NULL REFERENCES ai.prompt_runs(id) ON DELETE CASCADE,
    accepted       BOOLEAN NOT NULL,
    edit_distance  INT,
    rating         SMALLINT,
    comment        TEXT,
    user_id        UUID,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ai.datasets (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name           TEXT NOT NULL,
    version        TEXT NOT NULL,
    record_count   BIGINT NOT NULL,
    storage_uri    TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (name, version)
);

CREATE TABLE ai.fine_tune_runs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    dataset_id      UUID REFERENCES ai.datasets(id),
    base_model      TEXT NOT NULL,
    output_model_id UUID REFERENCES ai.model_registry(id),
    status          TEXT NOT NULL CHECK (status IN ('queued','running','succeeded','failed','cancelled')),
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    metrics         JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ai.event_outbox (LIKE auth.event_outbox INCLUDING ALL);

-- =================== RLS ===================================

ALTER TABLE conflict.conflict_cases  ENABLE ROW LEVEL SECURITY;
ALTER TABLE conflict.resolutions     ENABLE ROW LEVEL SECURITY;
ALTER TABLE conflict.ai_feedback     ENABLE ROW LEVEL SECURITY;
ALTER TABLE collab.pull_requests     ENABLE ROW LEVEL SECURITY;
ALTER TABLE collab.issues            ENABLE ROW LEVEL SECURITY;
ALTER TABLE collab.comments          ENABLE ROW LEVEL SECURITY;
ALTER TABLE collab.releases          ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai.prompt_runs           ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai.feedback_records      ENABLE ROW LEVEL SECURITY;

CREATE POLICY conflict_tenant_policy ON conflict.conflict_cases
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = conflict_cases.repo_id
          AND r.org_id = ANY (helix_current_orgs())
    ));

CREATE POLICY pr_tenant_policy ON collab.pull_requests
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = pull_requests.repo_id
          AND (r.org_id = ANY (helix_current_orgs()) OR r.visibility = 'public')
    ));

CREATE POLICY issue_tenant_policy ON collab.issues
    USING (EXISTS (
        SELECT 1 FROM repo.repositories r
        WHERE r.id = issues.repo_id
          AND (r.org_id = ANY (helix_current_orgs()) OR r.visibility = 'public')
    ));

CREATE POLICY ai_tenant_policy ON ai.prompt_runs
    USING (org_id = ANY (helix_current_orgs()) OR current_setting('helix.role', TRUE) = 'admin');
