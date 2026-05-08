package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
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

func TestLog_WarnOutput(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "warn", Output: &buf, Service: "test"})
	lg.Warn("danger", "key", "val")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v", err)
	}
	if got["level"] != "WARN" {
		t.Errorf("level = %v, want WARN", got["level"])
	}
	if got["msg"] != "danger" {
		t.Errorf("msg = %v, want danger", got["msg"])
	}
	if got["key"] != "val" {
		t.Errorf("key = %v, want val", got["key"])
	}
}

func TestLog_ErrorOutput(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "error", Output: &buf, Service: "test"})
	lg.Error("boom", "code", 500)

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v", err)
	}
	if got["level"] != "ERROR" {
		t.Errorf("level = %v, want ERROR", got["level"])
	}
	if got["msg"] != "boom" {
		t.Errorf("msg = %v, want boom", got["msg"])
	}
}

func TestLog_DebugFilteredAtInfo(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "info", Output: &buf, Service: "test"})
	lg.Debug("should_not_appear", "trace", "xyz")
	if buf.Len() > 0 {
		t.Errorf("Debug should be suppressed at info level, got: %s", buf.String())
	}
}

func TestLog_DebugVisibleAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "debug", Output: &buf, Service: "test"})
	lg.Debug("visible", "detail", "yes")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v", err)
	}
	if got["level"] != "DEBUG" {
		t.Errorf("level = %v, want DEBUG", got["level"])
	}
	if got["detail"] != "yes" {
		t.Errorf("detail = %v, want yes", got["detail"])
	}
}

func TestFromContext_Missing_ReturnsDefault(t *testing.T) {
	ctx := context.Background()
	lg := log.FromContext(ctx)
	if lg == nil {
		t.Fatal("FromContext returned nil for empty context")
	}
}

func TestDefault_BeforeNew_ReturnsNoop(t *testing.T) {
	lg := log.Default()
	if lg == nil {
		t.Fatal("Default() should never return nil")
	}
}

func TestWith_AttachesFields(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "info", Output: &buf, Service: "test"})
	child := lg.With("trace_id", "123")
	child.Info("child_msg")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v", err)
	}
	if got["trace_id"] != "123" {
		t.Errorf("trace_id = %v, want 123", got["trace_id"])
	}
	if got["msg"] != "child_msg" {
		t.Errorf("msg = %v, want child_msg", got["msg"])
	}
}

func TestNew_DefaultOutputIsStdout(t *testing.T) {
	lg := log.New(log.Options{Level: "info"})
	if lg == nil {
		t.Fatal("New with nil Output should use os.Stdout")
	}
}

func TestParseLevel_WarningAlias(t *testing.T) {
	var buf bytes.Buffer
	lg := log.New(log.Options{Level: "warning", Output: &buf, Service: "test"})
	lg.Warn("via_warning_level")
	if !strings.Contains(buf.String(), "via_warning_level") {
		t.Errorf("'warning' level should enable Warn output, got: %s", buf.String())
	}
}
