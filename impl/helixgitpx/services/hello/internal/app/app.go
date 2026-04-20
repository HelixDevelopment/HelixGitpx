// Package app is the composition root for hello.
package app

import (
	"context"
	"errors"
	"net"
	nethttp "net/http"
	"time"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
	grpchandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/grpc"
	httphandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/http"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/repo"

	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
	"github.com/helixgitpx/platform/config"
	hgin "github.com/helixgitpx/platform/gin"
	hgrpc "github.com/helixgitpx/platform/grpc"
	"github.com/helixgitpx/platform/health"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
	hredis "github.com/helixgitpx/platform/redis"
	"github.com/helixgitpx/platform/telemetry"
)

type cfg struct {
	HTTPAddr     string   `env:"HTTP_ADDR" default:":8001"`
	GRPCAddr     string   `env:"GRPC_ADDR" default:":9001"`
	HealthAddr   string   `env:"HEALTH_ADDR" default:":8081"`
	PostgresDSN  string   `env:"POSTGRES_DSN" required:"true"`
	RedisAddr    string   `env:"REDIS_ADDR" default:"localhost:6379"`
	KafkaTopic string `env:"KAFKA_TOPIC" default:"hello.said"`
	OTLPEndpoint string   `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Version      string   `env:"VERSION" default:"m1-dev"`
}

// Run wires dependencies and serves until ctx is done.
func Run(ctx context.Context, lg *log.Logger) error {
	var c cfg
	if err := config.Load(&c, config.Options{Prefix: "HELLO"}); err != nil {
		return err
	}

	shutdownTel, err := telemetry.Start(ctx, telemetry.Options{
		Service: "hello", Version: c.Version, OTLPEndpoint: c.OTLPEndpoint,
	})
	if err != nil {
		return err
	}
	defer func() {
		shctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTel(shctx)
	}()

	pool, err := pg.Open(ctx, pg.Options{DSN: c.PostgresDSN})
	if err != nil {
		return err
	}
	defer pool.Close()

	rc, err := hredis.Open(ctx, hredis.Options{Addr: c.RedisAddr, Namespace: "hello"})
	if err != nil {
		return err
	}
	defer rc.Close()

	greeter := domain.NewGreeter(
		&repo.CounterPG{Pool: pool},
		&repo.CacheRedis{Client: rc},
		&repo.EventOutbox{Pool: pool, Topic: c.KafkaTopic},
	)

	grpcSrv, err := hgrpc.NewServer(hgrpc.Options{})
	if err != nil {
		return err
	}
	hellopb.RegisterHelloServiceServer(grpcSrv, &grpchandler.Server{Greeter: greeter})

	router := hgin.NewRouter(hgin.Options{Service: "hello", Version: c.Version})
	httphandler.Register(router, greeter)

	hh := health.New()
	hh.Register("pg", pg.Probe(pool))
	hh.Register("redis", hredis.Probe(rc))
	hmux := nethttp.NewServeMux()
	hh.Routes(hmux)

	grpcL, err := net.Listen("tcp", c.GRPCAddr)
	if err != nil {
		return err
	}
	httpSrv := &nethttp.Server{Addr: c.HTTPAddr, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	healthSrv := &nethttp.Server{Addr: c.HealthAddr, Handler: hmux, ReadHeaderTimeout: 5 * time.Second}

	errCh := make(chan error, 3)
	go func() { errCh <- grpcSrv.Serve(grpcL) }()
	go func() { errCh <- httpSrv.ListenAndServe() }()
	go func() { errCh <- healthSrv.ListenAndServe() }()

	lg.Info("hello serving",
		"grpc", c.GRPCAddr, "http", c.HTTPAddr, "health", c.HealthAddr)

	select {
	case <-ctx.Done():
		lg.Info("shutdown signalled")
	case err := <-errCh:
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			lg.Error("server exited", "err", err.Error())
		}
	}

	shctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcSrv.GracefulStop()
	_ = httpSrv.Shutdown(shctx)
	_ = healthSrv.Shutdown(shctx)
	return nil
}
