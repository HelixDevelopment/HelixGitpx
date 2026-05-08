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

func TestLoad_NotAPointer(t *testing.T) {
	err := config.Load("not a struct", config.Options{})
	if err == nil {
		t.Fatal("expected error for non-pointer input")
	}
}

func TestLoad_UintField(t *testing.T) {
	type cfg struct {
		Port uint `env:"PORT" default:"8080"`
	}
	var c cfg
	if err := config.Load(&c, config.Options{}); err != nil {
		t.Fatal(err)
	}
	if c.Port != 8080 {
		t.Errorf("Port = %d, want 8080", c.Port)
	}
}

func TestLoad_FloatField(t *testing.T) {
	type cfg struct {
		Ratio float64 `env:"RATIO" default:"1.5"`
	}
	var c cfg
	if err := config.Load(&c, config.Options{}); err != nil {
		t.Fatal(err)
	}
	if c.Ratio != 1.5 {
		t.Errorf("Ratio = %f, want 1.5", c.Ratio)
	}
}

func TestLoad_EmbeddedStruct(t *testing.T) {
	type Inner struct {
		Val string `env:"INNER_VAL" default:"ok"`
	}
	type Outer struct {
		Inner
		Top string `env:"TOP" default:"yes"`
	}
	var c Outer
	if err := config.Load(&c, config.Options{}); err != nil {
		t.Fatal(err)
	}
	if c.Val != "ok" {
		t.Errorf("embedded Val = %q, want ok", c.Val)
	}
	if c.Top != "yes" {
		t.Errorf("Top = %q, want yes", c.Top)
	}
}

func TestLoad_BoolParseError(t *testing.T) {
	type cfg struct {
		Flag bool `env:"FLAG"`
	}
	t.Setenv("FLAG", "notabool")
	var c cfg
	err := config.Load(&c, config.Options{})
	if err == nil {
		t.Fatal("expected parse error for invalid bool")
	}
}

func TestLoad_IntParseError(t *testing.T) {
	type cfg struct {
		N int `env:"N"`
	}
	t.Setenv("N", "abc")
	var c cfg
	err := config.Load(&c, config.Options{})
	if err == nil {
		t.Fatal("expected parse error for invalid int")
	}
}

func TestLoad_DurationParseError(t *testing.T) {
	type cfg struct {
		D time.Duration `env:"D"`
	}
	t.Setenv("D", "notaduration")
	var c cfg
	err := config.Load(&c, config.Options{})
	if err == nil {
		t.Fatal("expected parse error for invalid duration")
	}
}

func TestLoad_UnsupportedKind(t *testing.T) {
	type cfg struct {
		M map[string]string `env:"M" default:"x"`
	}
	var c cfg
	err := config.Load(&c, config.Options{})
	if err == nil {
		t.Fatal("expected unsupported kind error")
	}
}

func TestLoad_EnvOverridesUint(t *testing.T) {
	type cfg struct {
		Port uint `env:"PORT" default:"8080"`
	}
	t.Setenv("PORT", "9090")
	var c cfg
	if err := config.Load(&c, config.Options{}); err != nil {
		t.Fatal(err)
	}
	if c.Port != 9090 {
		t.Errorf("Port = %d, want 9090", c.Port)
	}
}
