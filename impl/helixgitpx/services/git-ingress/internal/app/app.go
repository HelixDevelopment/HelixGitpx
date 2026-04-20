// Package app is the composition root for git-ingress.
package app

import (
	"context"

	"github.com/helixgitpx/platform/log"
)

// Run wires dependencies and serves until ctx is done.
// Fill in with real wiring (see services/hello for a reference).
func Run(ctx context.Context, lg *log.Logger) error {
	lg.Info("git-ingress starting")
	<-ctx.Done()
	lg.Info("git-ingress stopped")
	return nil
}
