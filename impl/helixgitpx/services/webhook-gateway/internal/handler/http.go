// Package handler accepts inbound webhooks from every supported Git host,
// verifies the HMAC signature, canonicalizes the payload, and responds 204.
// Publishing to Kafka is the next wiring step; for now the handler records
// accepted deliveries in the in-memory Store.
package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/helixgitpx/helixgitpx/services/webhook-gateway/internal/canonical"
	"github.com/helixgitpx/platform/webhook"
)

type SecretProvider func(provider, repo string) (string, bool)

type Handler struct {
	Secrets  SecretProvider
	Recorder *InMemoryRecorder
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/webhooks/github", h.github)
	mux.HandleFunc("POST /v1/webhooks/gitlab", h.gitlab)
	mux.HandleFunc("POST /v1/webhooks/gitea", h.gitea)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

func (h *Handler) readAndVerify(w http.ResponseWriter, r *http.Request, sigHeader string) ([]byte, bool) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20)) // 8 MiB cap
	if err != nil {
		http.Error(w, "body read failed", http.StatusBadRequest)
		return nil, false
	}
	defer r.Body.Close()

	sig := r.Header.Get(sigHeader)
	if sig == "" {
		http.Error(w, "missing signature", http.StatusUnauthorized)
		return nil, false
	}
	secret, ok := h.Secrets("", "")
	if !ok {
		http.Error(w, "no secret configured", http.StatusServiceUnavailable)
		return nil, false
	}
	if !webhook.VerifyHMAC([]byte(secret), body, sig) {
		http.Error(w, "bad signature", http.StatusUnauthorized)
		return nil, false
	}
	return body, true
}

func (h *Handler) github(w http.ResponseWriter, r *http.Request) {
	body, ok := h.readAndVerify(w, r, "X-Hub-Signature-256")
	if !ok {
		return
	}
	evt := canonical.CanonicalizeGitHub(r.Header.Get("X-GitHub-Delivery"), r.Header.Get("X-GitHub-Event"), body)
	h.Recorder.Record(evt)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) gitlab(w http.ResponseWriter, r *http.Request) {
	// GitLab uses an X-Gitlab-Token header (shared secret equality) instead of HMAC.
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	if err != nil {
		http.Error(w, "body read failed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	supplied := strings.TrimSpace(r.Header.Get("X-Gitlab-Token"))
	expected, ok := h.Secrets("", "")
	if !ok {
		http.Error(w, "no secret configured", http.StatusServiceUnavailable)
		return
	}
	if supplied == "" || supplied != expected {
		http.Error(w, "bad token", http.StatusUnauthorized)
		return
	}
	evt := canonical.CanonicalizeGitLab(r.Header.Get("X-Gitlab-Event-UUID"), r.Header.Get("X-Gitlab-Event"), body)
	h.Recorder.Record(evt)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) gitea(w http.ResponseWriter, r *http.Request) {
	body, ok := h.readAndVerify(w, r, "X-Gitea-Signature")
	if !ok {
		return
	}
	evt := canonical.CanonicalizeGitea(r.Header.Get("X-Gitea-Delivery"), r.Header.Get("X-Gitea-Event"), body)
	h.Recorder.Record(evt)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// InMemoryRecorder buffers accepted deliveries for diagnostic endpoints.
// Production wiring swaps this for a Kafka producer on upstream.webhooks.
type InMemoryRecorder struct {
	mu      sync.Mutex
	entries []canonical.Event
}

func (r *InMemoryRecorder) Record(e canonical.Event) {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	r.mu.Lock()
	r.entries = append(r.entries, e)
	r.mu.Unlock()
}

func (r *InMemoryRecorder) Snapshot() []canonical.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]canonical.Event(nil), r.entries...)
}

var _ = context.TODO // keep import for future handler plumbing
var _ = json.Marshal // keep import for future JSON responses
