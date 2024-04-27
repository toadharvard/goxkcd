package http

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/toadharvard/goxkcd/internal/app"
)

func SetupRoutes(ctx context.Context, app *app.App, mux *http.ServeMux) {
	mux.HandleFunc("GET /ping", PingHandler(ctx))
	mux.HandleFunc("GET /pics", SuggestComixPicsHandler(ctx, app.SuggestComixUseCase))
	mux.HandleFunc(
		"GET /update",
		UpdateDatabaseAndIndexHandler(
			ctx,
			app.DownloadComicsUseCase,
			app.BuildIndexUseCase,
			app.CountComicsUseCase,
		),
	)
}

func Run(ctx context.Context, app *app.App, host string, port int, comixUpdateInterval time.Duration) error {
	mux := http.NewServeMux()

	SetupRoutes(ctx, app, mux)

	chain := MiddlewareChain(
		LoggingMiddleware,
	)

	server := http.Server{
		Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
		Handler: chain(mux),
	}

	// TODO: use cron jobs instead of ticker
	go UpdateDatabaseAndIndexTask(ctx, app.BuildIndexUseCase, app.DownloadComicsUseCase, comixUpdateInterval)

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("server shutdown failed", "err", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
		return err
	}
	return nil
}
