package testkit

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

// StartKafka launches a Kafka 3.8 KRaft-mode container and returns the bootstrap broker.
func StartKafka(t testing.TB) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ctr, err := kafka.Run(ctx, "apache/kafka:3.8.1",
		kafka.WithClusterID("helixgitpx-test"),
	)
	if err != nil {
		t.Fatalf("testkit.StartKafka: %v", err)
	}
	brokers, err := ctr.Brokers(ctx)
	if err != nil {
		t.Fatalf("testkit.StartKafka brokers: %v", err)
	}
	if len(brokers) == 0 {
		t.Fatalf("testkit.StartKafka: no brokers")
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
	_ = time.Now
	return brokers[0]
}
