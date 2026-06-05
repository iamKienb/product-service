package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"product-query-module/internal/bootstrap"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	app := bootstrap.NewApp(logger)
	errCh := make(chan error, 1)
	go func() { errCh <- app.Start(ctx) }()
	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			logger.Error("failed to start", slog.Any("err", err))
			os.Exit(1)
		}
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = app.Stop(shutdownCtx)
}
