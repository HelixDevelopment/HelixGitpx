package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestJWKS_ReturnsEmptyKeys(t *testing.T) {
	r := &Router{}
	g := gin.New()
	r.Register(g)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	keys, ok := body["keys"]
	if !ok {
		t.Fatal("missing 'keys' field in JWKS response")
	}
	keysArr, ok := keys.([]any)
	if !ok {
		t.Fatalf("keys = %T, want []any", keys)
	}
	if len(keysArr) != 0 {
		t.Errorf("keys len = %d, want 0", len(keysArr))
	}
}

func TestCallback_MissingCode(t *testing.T) {
	r := &Router{}
	g := gin.New()
	r.Register(g)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/callback", nil)
	g.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["error"] != "missing code" {
		t.Errorf("error = %v, want 'missing code'", body["error"])
	}
}

func TestRefresh_NoCookie(t *testing.T) {
	r := &Router{}
	g := gin.New()
	r.Register(g)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", nil)
	g.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["error"] != "no refresh cookie" {
		t.Errorf("error = %v, want 'no refresh cookie'", body["error"])
	}
}

func TestRefresh_InvalidRefreshToken(t *testing.T) {
	r := &Router{}
	g := gin.New()
	r.Register(g)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "not-a-uuid"})
	g.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["error"] != "bad refresh" {
		t.Errorf("error = %v, want 'bad refresh'", body["error"])
	}
}
