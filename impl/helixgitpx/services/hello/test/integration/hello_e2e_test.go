//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
	grpchandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/grpc"
	httphandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/http"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/repo"

	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
	hgin "github.com/helixgitpx/platform/gin"
	hgrpc "github.com/helixgitpx/platform/grpc"
	"github.com/helixgitpx/platform/kafka"
	"github.com/helixgitpx/platform/pg"
	hredis "github.com/helixgitpx/platform/redis"
	"github.com/helixgitpx/platform/testkit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHelloE2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	dsn := testkit.StartPostgres(t)
	redisAddr := testkit.StartRedis(t)
	kafkaBroker := testkit.StartKafka(t)

	pool, err := pg.Open(ctx, pg.Options{DSN: dsn})
	if err != nil {
		t.Fatalf("pg.Open: %v", err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS hello`); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS hello.greetings (
			name TEXT PRIMARY KEY,
			count BIGINT NOT NULL DEFAULT 0,
			last_said_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	rc, err := hredis.Open(ctx, hredis.Options{Addr: redisAddr, Namespace: "hello"})
	if err != nil {
		t.Fatalf("redis.Open: %v", err)
	}

	prod, err := kafka.NewProducer(kafka.ProducerOptions{
		Brokers: []string{kafkaBroker}, ClientID: "hello-test", Topic: "hello.said",
	})
	if err != nil {
		t.Fatalf("kafka producer: %v", err)
	}
	defer prod.Close(ctx)

	greeter := domain.NewGreeter(
		&repo.CounterPG{Pool: pool},
		&repo.CacheRedis{Client: rc},
		&repo.EventKafka{Producer: prod, Topic: "hello.said"},
	)

	grpcSrv, _ := hgrpc.NewServer(hgrpc.Options{})
	hellopb.RegisterHelloServiceServer(grpcSrv, &grpchandler.Server{Greeter: greeter})
	grpcL, _ := net.Listen("tcp", "127.0.0.1:0")
	go grpcSrv.Serve(grpcL)
	defer grpcSrv.GracefulStop()

	router := hgin.NewRouter(hgin.Options{Service: "hello"})
	httphandler.Register(router, greeter)
	httpL, _ := net.Listen("tcp", "127.0.0.1:0")
	httpSrv := &nethttp.Server{Handler: router}
	go httpSrv.Serve(httpL)
	defer httpSrv.Close()

	resp, err := nethttp.Get(fmt.Sprintf("http://%s/v1/hello?name=world", httpL.Addr().String()))
	if err != nil {
		t.Fatalf("http: %v", err)
	}
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["greeting"] != "hello, world" {
		t.Errorf("http greeting = %v", body["greeting"])
	}

	conn, err := grpc.NewClient(grpcL.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	cl := hellopb.NewHelloServiceClient(conn)
	r, err := cl.SayHello(ctx, &hellopb.SayHelloRequest{Name: "world"})
	if err != nil {
		t.Fatalf("SayHello: %v", err)
	}
	if r.Greeting != "hello, world" {
		t.Errorf("grpc greeting = %q", r.Greeting)
	}
	if r.Count < 2 {
		t.Errorf("count = %d, want ≥ 2 (HTTP + gRPC calls)", r.Count)
	}
}
