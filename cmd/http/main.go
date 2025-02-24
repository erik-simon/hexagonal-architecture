package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	bootstrap "github.com/hexagonal/cmd/bootstrap"
)

const defaultPort = "4321"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	engine := chi.NewRouter()

	const envPORT = "PORT"

	port := strings.TrimSpace(os.Getenv(envPORT))
	if port == "" {
		port = defaultPort
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := bootstrap.Inject(ctx, logger, engine); err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: engine,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	shutdownCh := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()

		const graceTime = 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), graceTime)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}

		shutdownCh <- struct{}{}
		close(shutdownCh)
	}()

	errCh := make(chan error, 1)

	log.Printf("Listening server http on port '%s'", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		errCh <- err
		close(errCh)
	}

	go func() {
		select {
		case <-ctx.Done():
			stop()
			<-shutdownCh
		case err := <-errCh:
			stop()
			<-shutdownCh

			log.Fatal(err)
		}
	}()
}
