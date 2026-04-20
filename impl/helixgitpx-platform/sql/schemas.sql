-- HelixGitpx per-service schemas + RLS baseline.
-- One schema per top-level domain from the roadmap. Each schema gets a
-- dedicated role used by the service's DSN. RLS is enabled on every
-- application table created under these schemas (later milestones tighten
-- policies per tenancy model).

CREATE SCHEMA IF NOT EXISTS hello;
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS repo;
CREATE SCHEMA IF NOT EXISTS sync;
CREATE SCHEMA IF NOT EXISTS conflict;
CREATE SCHEMA IF NOT EXISTS upstream;
CREATE SCHEMA IF NOT EXISTS collab;
CREATE SCHEMA IF NOT EXISTS events;
CREATE SCHEMA IF NOT EXISTS platform;

DO $$
DECLARE
    s TEXT;
BEGIN
    FOREACH s IN ARRAY ARRAY['hello','auth','repo','sync','conflict','upstream','collab','events','platform']
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
