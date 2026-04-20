// Package app is the composition root for orgteam.
package app

import (
	"context"

	"github.com/helixgitpx/platform/log"
)

// Run wires dependencies and serves until ctx is done.
// Fill in with real wiring (see services/hello for a reference).
func Run(ctx context.Context, lg *log.Logger) error {
	lg.Info("orgteam starting")
	<-ctx.Done()
	lg.Info("orgteam stopped")
	return nil
}
