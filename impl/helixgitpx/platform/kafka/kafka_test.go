package kafka_test

import (
	"errors"
	"testing"

	"github.com/helixgitpx/platform/kafka"
)

func TestOptions_Validation(t *testing.T) {
	_, err := kafka.NewProducer(kafka.ProducerOptions{})
	if err == nil {
		t.Fatalf("expected error for missing brokers")
	}
}

func TestIsUnavailable(t *testing.T) {
	if !kafka.IsUnavailable(kafka.ErrUnavailable) {
		t.Errorf("sentinel")
	}
	if kafka.IsUnavailable(errors.New("other")) {
		t.Errorf("other")
	}
}
