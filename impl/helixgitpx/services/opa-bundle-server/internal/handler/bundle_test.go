package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/store"
)

func setup(t *testing.T) (*httptest.Server, *store.Store, domain.Bundle) {
	t.Helper()
	s := store.New()
	content := []byte("fake bundle bytes")
	meta, err := domain.NewBundle("b-1", "2.0.0", "deadbeef", content, true, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	meta = s.Put(meta, content)
	meta.ActivatedAt = time.Now()
	if err := s.Activate(meta.ID, meta); err != nil {
		t.Fatal(err)
	}
	h := &Handler{Store: s}
	return httptest.NewServer(h.Routes()), s, meta
}

func TestGetActive_ReturnsContentAndETag(t *testing.T) {
	srv, _, meta := setup(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/bundles/active")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("ETag"); got != meta.ETag() {
		t.Fatalf("etag: want %q got %q", meta.ETag(), got)
	}
	if resp.Header.Get("X-Bundle-Version") != "2.0.0" {
		t.Fatal("version header missing")
	}
}

func TestGetActive_NotModifiedWhenETagMatches(t *testing.T) {
	srv, _, meta := setup(t)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/bundles/active", nil)
	req.Header.Set("If-None-Match", meta.ETag())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotModified {
		t.Fatalf("want 304, got %d", resp.StatusCode)
	}
}

func TestGetActive_NoBundles503(t *testing.T) {
	h := &Handler{Store: store.New()}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/bundles/active")
	if err != nil { t.Fatal(err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", resp.StatusCode)
	}
}

func TestList(t *testing.T) {
	srv, _, _ := setup(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/bundles")
	if err != nil { t.Fatal(err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
}

func TestGetOne(t *testing.T) {
	srv, _, meta := setup(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/bundles/" + meta.ID)
	if err != nil { t.Fatal(err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}

	resp404, err404 := http.Get(srv.URL + "/bundles/bogus")
	if err404 != nil { t.Fatal(err404) }
	defer resp404.Body.Close()
	if resp404.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404 got %d", resp404.StatusCode)
	}
}

func TestHealthz(t *testing.T) {
	srv, _, _ := setup(t)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil { t.Fatal(err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("healthz must be 200")
	}
}
