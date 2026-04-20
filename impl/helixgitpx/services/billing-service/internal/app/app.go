// Package app is the composition root for billing-service.
package app

import (
	"context"

	"github.com/helixgitpx/platform/log"
)

// Run wires dependencies and serves until ctx is done.
func Run(ctx context.Context, lg *log.Logger) error {
	lg.Info("billing-service starting")
	<-ctx.Done()
	lg.Info("billing-service stopped")
	return nil
}
