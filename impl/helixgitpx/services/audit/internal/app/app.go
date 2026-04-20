// Package app is the composition root for audit.
package app

import (
	"context"
	"errors"
	"net"
	nethttp "net/http"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/helixgitpx/helixgitpx/services/audit/internal/consumer"
	"github.com/helixgitpx/platform/config"
	hgrpc "github.com/helixgitpx/platform/grpc"
	"github.com/helixgitpx/platform/health"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
	"github.com/helixgitpx/platform/telemetry"
)

type cfg struct {
	GRPCAddr      string   `env:"GRPC_ADDR" default:":9004"`
	HealthAddr    string   `env:"HEALTH_ADDR" default:":8084"`
	PostgresDSN   string   `env:"POSTGRES_DSN" vault:"kv/audit#pg_dsn" required:"true"`
	KafkaBrokers  []string `env:"KAFKA_BROKERS" default:"helix-kafka-kafka-bootstrap.helix-data.svc:9092" split:","`
	KafkaTopic    string   `env:"KAFKA_TOPIC" default:"audit.events"`
	ConsumerGroup string   `env:"CONSUMER_GROUP" default:"audit-service"`
	OTLPEndpoint  string   `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Version       string   `env:"VERSION" default:"m3-dev"`
}

// Run wires dependencies and serves until ctx is done.
func Run(ctx context.Context, lg *log.Logger) error {
	var c cfg
	if err := config.Load(&c, config.Options{Prefix: "AUDIT"}); err != nil {
		return err
	}
	shutdownTel, _ := telemetry.Start(ctx, telemetry.Options{Service: "audit", Version: c.Version, OTLPEndpoint: c.OTLPEndpoint})
	defer func() {
		sh, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTel(sh)
	}()

	pool, err := pg.Open(ctx, pg.Options{DSN: c.PostgresDSN})
	if err != nil {
		return err
	}
	defer pool.Close()

	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(c.KafkaBrokers...),
		kgo.ConsumerGroup(c.ConsumerGroup),
		kgo.ConsumeTopics(c.KafkaTopic),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		return err
	}
	defer kcl.Close()

	grpcSrv, _ := hgrpc.NewServer(hgrpc.Options{})
	hh := health.New()
	hh.Register("pg", pg.Probe(pool))
	hmux := nethttp.NewServeMux()
	hh.Routes(hmux)
	telemetry.RegisterPprof(hmux)

	grpcL, err := net.Listen("tcp", c.GRPCAddr)
	if err != nil {
		return err
	}
	healthSrv := &nethttp.Server{Addr: c.HealthAddr, Handler: hmux, ReadHeaderTimeout: 5 * time.Second}

	cons := &consumer.Consumer{Client: kcl, Pool: pool}
	errCh := make(chan error, 3)
	go func() { errCh <- grpcSrv.Serve(grpcL) }()
	go func() { errCh <- healthSrv.ListenAndServe() }()
	go func() { errCh <- cons.Run(ctx) }()

	lg.Info("audit serving", "grpc", c.GRPCAddr, "health", c.HealthAddr, "topic", c.KafkaTopic)

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			lg.Error("exited", "err", err.Error())
		}
	}
	grpcSrv.GracefulStop()
	sh, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = healthSrv.Shutdown(sh)
	return nil
}
