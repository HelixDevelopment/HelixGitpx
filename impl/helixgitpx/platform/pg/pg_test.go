package pg_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/helixgitpx/platform/pg"
)

func TestOpen_InvalidDSNFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := pg.Open(ctx, pg.Options{DSN: "not-a-valid-dsn"})
	if err == nil {
		t.Fatalf("expected error for invalid DSN")
	}
}

func TestIsUnavailable(t *testing.T) {
	if !pg.IsUnavailable(pg.ErrUnavailable) {
		t.Errorf("ErrUnavailable not classified as unavailable")
	}
	if pg.IsUnavailable(errors.New("other")) {
		t.Errorf("arbitrary error classified as unavailable")
	}
}
