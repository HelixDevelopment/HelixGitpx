-- HelixGitpx per-service schemas + RLS baseline.
-- One schema per top-level domain from the roadmap. Each schema gets a
-- dedicated role used by the service's DSN. RLS is enabled on every
-- application table created under these schemas (later milestones tighten
-- policies per tenancy model).

CREATE SCHEMA IF NOT EXISTS hello;
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS org;
CREATE SCHEMA IF NOT EXISTS team;
CREATE SCHEMA IF NOT EXISTS repo;
CREATE SCHEMA IF NOT EXISTS sync;
CREATE SCHEMA IF NOT EXISTS conflict;
CREATE SCHEMA IF NOT EXISTS upstream;
CREATE SCHEMA IF NOT EXISTS collab;
CREATE SCHEMA IF NOT EXISTS events;
CREATE SCHEMA IF NOT EXISTS audit;
CREATE SCHEMA IF NOT EXISTS platform;

DO $$
DECLARE
    s TEXT;
BEGIN
    FOREACH s IN ARRAY ARRAY['hello','auth','org','team','repo','sync','conflict','upstream','collab','events','audit','platform']
    LOOP
        BEGIN
            EXECUTE format('CREATE ROLE %I_svc LOGIN', s);
        EXCEPTION WHEN duplicate_object THEN
            NULL;
        END;
        EXECUTE format('GRANT USAGE, CREATE ON SCHEMA %I TO %I_svc', s, s);
        EXECUTE format('GRANT ALL ON ALL TABLES IN SCHEMA %I TO %I_svc', s, s);
        EXECUTE format('GRANT ALL ON ALL SEQUENCES IN SCHEMA %I TO %I_svc', s, s);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON TABLES TO %I_svc', s, s);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON SEQUENCES TO %I_svc', s, s);
    END LOOP;
END $$;

-- orgteam_svc — a cross-schema role used by the merged orgteam-service (ADR-0014).
DO $$
BEGIN
    BEGIN
        EXECUTE 'CREATE ROLE orgteam_svc LOGIN';
    EXCEPTION WHEN duplicate_object THEN
        NULL;
    END;
    EXECUTE 'GRANT USAGE, CREATE ON SCHEMA org TO orgteam_svc';
    EXECUTE 'GRANT USAGE, CREATE ON SCHEMA team TO orgteam_svc';
    EXECUTE 'GRANT ALL ON ALL TABLES IN SCHEMA org TO orgteam_svc';
    EXECUTE 'GRANT ALL ON ALL TABLES IN SCHEMA team TO orgteam_svc';
    EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA org GRANT ALL ON TABLES TO orgteam_svc';
    EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA team GRANT ALL ON TABLES TO orgteam_svc';
END $$;
