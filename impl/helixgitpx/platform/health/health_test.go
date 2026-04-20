package health_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/helixgitpx/platform/health"
)

func TestHandler_LiveAlwaysOK(t *testing.T) {
	h := health.New()
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.Live(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Live Code = %d, want 200", w.Code)
	}
}

func TestHandler_ReadyReflectsProbes(t *testing.T) {
	h := health.New()
	h.Register("db", func(context.Context) error { return nil })
	h.Register("cache", func(context.Context) error { return errors.New("down") })

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	h.Ready(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Ready Code = %d, want 503 (one probe down)", w.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("body not JSON: %v", err)
	}
	checks, _ := body["checks"].(map[string]any)
	if checks["db"] != "ok" {
		t.Errorf("db = %v, want ok", checks["db"])
	}
	if checks["cache"] == "ok" {
		t.Errorf("cache should be failing")
	}
}
