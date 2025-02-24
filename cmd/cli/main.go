package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/hexagonal/cmd/bootstrap"
	"github.com/spf13/cobra"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var cmd cobra.Command

	if err := bootstrap.Inject(ctx, logger, &cmd); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
