package config_test

import (
	"testing"
	"time"

	"github.com/helixgitpx/platform/config"
)

type helloConfig struct {
	HTTPAddr     string        `env:"HTTP_ADDR" default:":8001"`
	GRPCAddr     string        `env:"GRPC_ADDR" default:":9001"`
	Timeout      time.Duration `env:"TIMEOUT" default:"30s"`
	KafkaBrokers []string      `env:"KAFKA_BROKERS" default:"localhost:9092" split:","`
	Enabled      bool          `env:"ENABLED" default:"true"`
	MaxConns     int           `env:"MAX_CONNS" default:"10"`
}

func TestLoad_UsesDefaults(t *testing.T) {
	var c helloConfig
	if err := config.Load(&c, config.Options{Prefix: "HELLO"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.HTTPAddr != ":8001" {
		t.Errorf("HTTPAddr = %q, want :8001", c.HTTPAddr)
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", c.Timeout)
	}
	if len(c.KafkaBrokers) != 1 || c.KafkaBrokers[0] != "localhost:9092" {
		t.Errorf("KafkaBrokers = %v", c.KafkaBrokers)
	}
	if !c.Enabled {
		t.Errorf("Enabled = false, want true")
	}
	if c.MaxConns != 10 {
		t.Errorf("MaxConns = %d, want 10", c.MaxConns)
	}
}

func TestLoad_EnvOverridesDefault(t *testing.T) {
	t.Setenv("HELLO_HTTP_ADDR", ":9999")
	t.Setenv("HELLO_KAFKA_BROKERS", "a:1,b:2")
	t.Setenv("HELLO_ENABLED", "false")
	var c helloConfig
	if err := config.Load(&c, config.Options{Prefix: "HELLO"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.HTTPAddr != ":9999" {
		t.Errorf("HTTPAddr = %q", c.HTTPAddr)
	}
	if len(c.KafkaBrokers) != 2 || c.KafkaBrokers[1] != "b:2" {
		t.Errorf("KafkaBrokers = %v", c.KafkaBrokers)
	}
	if c.Enabled {
		t.Errorf("Enabled = true, want false")
	}
}

type required struct {
	DSN string `env:"DSN" required:"true"`
}

func TestLoad_RequiredFieldMissing(t *testing.T) {
	var c required
	err := config.Load(&c, config.Options{Prefix: "X"})
	if err == nil {
		t.Fatalf("expected error for missing required DSN")
	}
}
