-- Per-service schemas for the hello service (M1 spine).
-- Later milestones add more services; each gets its own schema.
CREATE SCHEMA IF NOT EXISTS hello;
CREATE USER hello_svc WITH PASSWORD 'hello_svc';
GRANT USAGE, CREATE ON SCHEMA hello TO hello_svc;
GRANT ALL ON ALL TABLES IN SCHEMA hello TO hello_svc;
ALTER DEFAULT PRIVILEGES IN SCHEMA hello GRANT ALL ON TABLES TO hello_svc;
ALTER DEFAULT PRIVILEGES IN SCHEMA hello GRANT ALL ON SEQUENCES TO hello_svc;
