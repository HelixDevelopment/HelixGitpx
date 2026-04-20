// Package platform is the shared-library umbrella for HelixGitpx services.
//
// Sub-packages provide logging, telemetry, typed errors, configuration,
// gRPC/HTTP servers, Kafka/Postgres/Redis clients, Temporal wiring, SPIFFE
// integration, OPA evaluation, health endpoints, and a test toolkit.
//
// See the per-package doc.go for details. All constructors accept a context
// and return (client, error); callers own lifecycle (Close).
package platform
