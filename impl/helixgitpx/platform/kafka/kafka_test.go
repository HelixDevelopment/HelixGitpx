package kafka_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestKarapaceClient_NilURL_IsNoOp(t *testing.T) {
	k := &kafka.KarapaceClient{}
	id, err := k.Resolve(context.Background(), "hello.said-value", 1)
	if err != nil || id != -1 {
		t.Fatalf("want (-1,nil), got (%d,%v)", id, err)
	}
}

func TestKarapaceClient_Resolve_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/subjects/hello.said-value/versions/2"; got != want {
			t.Fatalf("path got=%q want=%q", got, want)
		}
		_, _ = w.Write([]byte(`{"id":42,"version":2}`))
	}))
	defer srv.Close()

	k := &kafka.KarapaceClient{URL: srv.URL, Client: srv.Client()}
	id, err := k.Resolve(context.Background(), "hello.said-value", 2)
	if err != nil {
		t.Fatal(err)
	}
	if id != 42 {
		t.Fatalf("want 42, got %d", id)
	}

	// Cache hit — even if the server goes away, we still resolve.
	srv.Close()
	if id2, err := k.Resolve(context.Background(), "hello.said-value", 2); err != nil || id2 != 42 {
		t.Fatalf("cache hit failed: %d,%v", id2, err)
	}
}

func TestKarapaceClient_Resolve_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	k := &kafka.KarapaceClient{URL: srv.URL, Client: srv.Client()}
	if _, err := k.Resolve(context.Background(), "missing", 1); err == nil {
		t.Fatal("expected error on 404")
	}
}
