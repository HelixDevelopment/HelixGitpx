// Package kafka wraps franz-go (github.com/twmb/franz-go/pkg/kgo) with
// defaults suitable for HelixGitpx services. Schema-registry integration
// is stubbed in M1 (a ResolveFn hook) and wired in M2.
package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// ErrUnavailable is returned when brokers are unreachable.
var ErrUnavailable = errors.New("kafka: unavailable")

// ProducerOptions configures NewProducer.
type ProducerOptions struct {
	Brokers  []string
	ClientID string
	Topic    string // default topic for Emit helper; may be overridden per record
}

// Producer wraps *kgo.Client with HelixGitpx ergonomics.
type Producer struct {
	cl    *kgo.Client
	topic string
}

// NewProducer constructs a Producer.
func NewProducer(opts ProducerOptions) (*Producer, error) {
	if len(opts.Brokers) == 0 {
		return nil, fmt.Errorf("kafka: Brokers is required")
	}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(opts.Brokers...),
		kgo.ClientID(stringOrDefault(opts.ClientID, "helixgitpx")),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerLinger(0),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka: new client: %w", err)
	}
	return &Producer{cl: cl, topic: opts.Topic}, nil
}

// Emit publishes one record to the default (or overridden) topic.
func (p *Producer) Emit(ctx context.Context, key, value []byte, topic ...string) error {
	t := p.topic
	if len(topic) > 0 && topic[0] != "" {
		t = topic[0]
	}
	if t == "" {
		return fmt.Errorf("kafka: topic unset")
	}
	res := p.cl.ProduceSync(ctx, &kgo.Record{Topic: t, Key: key, Value: value})
	if err := res.FirstErr(); err != nil {
		return errors.Join(ErrUnavailable, err)
	}
	return nil
}

// Close flushes and closes the client.
func (p *Producer) Close(ctx context.Context) error {
	if p == nil || p.cl == nil {
		return nil
	}
	if err := p.cl.Flush(ctx); err != nil {
		return err
	}
	p.cl.Close()
	return nil
}

// IsUnavailable reports whether err wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

func stringOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
