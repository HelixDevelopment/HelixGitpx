// Command hello is a HelixGitpx service scaffolded by tools/scaffold (M1)
// and extended in M2 with a `migrate` subcommand that runs goose migrations.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/app"
	"github.com/helixgitpx/platform/log"
	"github.com/helixgitpx/platform/pg"
)

func main() {
	lg := log.New(log.Options{Level: "info", Service: "hello"})

	if len(os.Args) >= 2 && os.Args[1] == "migrate" {
		if err := runMigrate(context.Background(), os.Args[2:]); err != nil {
			lg.Error("migrate failed", "err", err.Error())
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, lg); err != nil {
		lg.Error("service exited with error", "err", err.Error())
	}
}

func runMigrate(ctx context.Context, args []string) error {
	dir := "/migrations"
	for i := 0; i < len(args); i++ {
		if args[i] == "--dir" && i+1 < len(args) {
			dir = args[i+1]
			i++
		}
	}
	dsn := os.Getenv("HELLO_POSTGRES_DSN")
	if dsn == "" {
		return fmt.Errorf("HELLO_POSTGRES_DSN is required")
	}
	return pg.Migrate(ctx, pg.MigrateOptions{DSN: dsn, Dir: dir})
}
