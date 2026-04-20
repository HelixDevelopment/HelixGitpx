package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/helixgitpx/platform/kafka"
)

// EventKafka implements domain.Emitter using HelixGitpx's Kafka wrapper.
type EventKafka struct {
	Producer *kafka.Producer
	Topic    string
}

type helloSaid struct {
	Name     string `json:"name"`
	Greeting string `json:"greeting"`
	Count    int64  `json:"count"`
	At       string `json:"at"`
}

func (e *EventKafka) Emit(ctx context.Context, name, greeting string, count int64) error {
	payload, err := json.Marshal(helloSaid{
		Name: name, Greeting: greeting, Count: count,
		At: time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}
	return e.Producer.Emit(ctx, []byte(name), payload, e.Topic)
}
