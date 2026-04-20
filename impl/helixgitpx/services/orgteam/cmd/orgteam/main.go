// Command orgteam is a HelixGitpx service scaffolded by tools/scaffold.
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/app"
	"github.com/helixgitpx/platform/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lg := log.New(log.Options{Level: "info", Service: "orgteam"})
	if err := app.Run(ctx, lg); err != nil {
		lg.Error("service exited with error", "err", err.Error())
	}
}
