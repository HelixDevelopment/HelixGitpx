module github.com/helixgitpx/helixgitpx/services/audit

go 1.23.0

toolchain go1.23.4

require (
	github.com/helixgitpx/platform v0.0.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/twmb/franz-go v1.18.0
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pressly/goose/v3 v3.22.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/spiffe/go-spiffe/v2 v2.4.0 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.9.0 // indirect
	github.com/zeebo/errs v1.3.0 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.32.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.32.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/sdk v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/grpc v1.69.2 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
)

replace github.com/helixgitpx/platform => ../../platform

replace github.com/helixgitpx/helixgitpx/gen => ../../gen
