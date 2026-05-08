package bitbucket

import (
	"testing"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
)

func TestAdapter_ImplementsInterface(t *testing.T) {
	var _ adapter.Adapter = (*Adapter)(nil)
}

func TestAdapter_Push(t *testing.T) {
	a := &Adapter{}
	if err := a.Push(nil, adapter.Destination{}, nil); err != nil {
		t.Fatalf("want nil got %v", err)
	}
}

func TestAdapter_GetRepo(t *testing.T) {
	a := &Adapter{}
	info, err := a.GetRepo(nil, adapter.Source{})
	if err != nil {
		t.Fatal(err)
	}
	if info.Default != "main" {
		t.Fatalf("want main got %q", info.Default)
	}
}

func TestAdapter_RegisterWebhook(t *testing.T) {
	a := &Adapter{}
	events := []string{"push"}
	wh, err := a.RegisterWebhook(nil, adapter.Source{}, "http://hook", "secret", events)
	if err != nil {
		t.Fatal(err)
	}
	if wh.URL != "http://hook" {
		t.Fatalf("want http://hook got %q", wh.URL)
	}
	if len(wh.Events) != 1 || wh.Events[0] != "push" {
		t.Fatalf("want [push] got %v", wh.Events)
	}
}
