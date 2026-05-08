package store

import (
	"testing"

	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/domain"
)

func TestNew_EmptyStore(t *testing.T) {
	s := New()
	if _, _, err := s.Active(); err != ErrInactiveRequested {
		t.Fatalf("new store should have no active bundle, got err=%v", err)
	}
	if got := s.List(); len(got) != 0 {
		t.Fatalf("new store should list empty, got %d items", len(got))
	}
}

func TestPut_ReturnsBundleWithID(t *testing.T) {
	s := New()
	meta := domain.Bundle{ID: "b1", Version: "1.0.0"}
	ret := s.Put(meta, []byte("content"))
	if ret.ID != "b1" {
		t.Fatalf("want ID b1 got %q", ret.ID)
	}
}

func TestPut_DefensiveCopy(t *testing.T) {
	s := New()
	content := []byte("original")
	s.Put(domain.Bundle{ID: "b1"}, content)
	content[0] = 'X'
	_, got, _ := s.Get("b1")
	if string(got) != "original" {
		t.Fatalf("put should copy content, got %q", string(got))
	}
}

func TestGet_Found(t *testing.T) {
	s := New()
	s.Put(domain.Bundle{ID: "b1", Version: "1.0.0"}, []byte("data"))
	meta, content, err := s.Get("b1")
	if err != nil {
		t.Fatal(err)
	}
	if meta.ID != "b1" {
		t.Fatalf("want b1 got %q", meta.ID)
	}
	if string(content) != "data" {
		t.Fatalf("want data got %q", string(content))
	}
}

func TestGet_NotFound(t *testing.T) {
	s := New()
	_, _, err := s.Get("missing")
	if err != domain.ErrUnknownBundleID {
		t.Fatalf("want ErrUnknownBundleID got %v", err)
	}
}

func TestActivate_SetsActive(t *testing.T) {
	s := New()
	s.Put(domain.Bundle{ID: "b1"}, []byte("a"))
	s.Put(domain.Bundle{ID: "b2"}, []byte("b"))
	if err := s.Activate("b1", domain.Bundle{ID: "b1", Version: "2.0.0"}); err != nil {
		t.Fatal(err)
	}
	meta, content, err := s.Active()
	if err != nil {
		t.Fatal(err)
	}
	if meta.ID != "b1" {
		t.Fatalf("want b1 got %q", meta.ID)
	}
	if meta.Version != "2.0.0" {
		t.Fatalf("want version 2.0.0 got %q", meta.Version)
	}
	if string(content) != "a" {
		t.Fatalf("want a got %q", string(content))
	}
}

func TestActivate_UnknownID(t *testing.T) {
	s := New()
	err := s.Activate("missing", domain.Bundle{ID: "missing"})
	if err != domain.ErrUnknownBundleID {
		t.Fatalf("want ErrUnknownBundleID got %v", err)
	}
}

func TestList_AllBundles(t *testing.T) {
	s := New()
	s.Put(domain.Bundle{ID: "b1"}, nil)
	s.Put(domain.Bundle{ID: "b2"}, nil)
	s.Put(domain.Bundle{ID: "b3"}, nil)
	list := s.List()
	if len(list) != 3 {
		t.Fatalf("want 3 got %d", len(list))
	}
	ids := map[string]bool{}
	for _, b := range list {
		ids[b.ID] = true
	}
	for _, id := range []string{"b1", "b2", "b3"} {
		if !ids[id] {
			t.Fatalf("missing %q in list", id)
		}
	}
}
