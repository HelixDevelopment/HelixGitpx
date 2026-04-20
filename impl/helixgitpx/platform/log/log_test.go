package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/helixgitpx/platform/log"
)

func TestNew_EmitsJSON(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "info", Output: &buf, Service: "hello", Version: "test"})
	lg.Info("hello", "name", "world")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v\n%s", err, buf.String())
	}
	if got["msg"] != "hello" {
		t.Errorf("msg = %v, want hello", got["msg"])
	}
	if got["name"] != "world" {
		t.Errorf("name = %v, want world", got["name"])
	}
	if got["service"] != "hello" {
		t.Errorf("service = %v, want hello", got["service"])
	}
	if got["version"] != "test" {
		t.Errorf("version = %v, want test", got["version"])
	}
}

func TestFromContext_ReturnsChildLogger(t *testing.T) {
	var buf bytes.Buffer
	root := log.New(log.Options{Level: "info", Output: &buf, Service: "s"})
	ctx := log.WithContext(context.Background(), root.With("request_id", "abc"))

	log.FromContext(ctx).Info("tick")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if got["request_id"] != "abc" {
		t.Errorf("request_id = %v, want abc", got["request_id"])
	}
}
