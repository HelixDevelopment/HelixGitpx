// Package app is the composition root for conflict-resolver.
package app

import (
	"context"

	"github.com/helixgitpx/platform/log"
)

// Run wires dependencies and serves until ctx is done.
// Fill in with real wiring (see services/hello for a reference).
func Run(ctx context.Context, lg *log.Logger) error {
	lg.Info("conflict-resolver starting")
	<-ctx.Done()
	lg.Info("conflict-resolver stopped")
	return nil
}
