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

func SetupRoutes(app *app.App, mux *http.ServeMux) {
	mux.HandleFunc("GET /ping", PingHandler())
	mux.HandleFunc("GET /pics", SuggestComixPicsHandler(app.SuggestComixUseCase))
	mux.HandleFunc(
		"POST /update",
		UpdateDatabaseAndIndexHandler(
			app.DownloadComicsUseCase,
			app.CountComicsUseCase,
		),
	)
}

func Run(ctx context.Context, app *app.App, host string, port int, comixUpdateInterval time.Duration) error {
	mux := http.NewServeMux()

	SetupRoutes(app, mux)

	chain := MiddlewareChain(
		LoggingMiddleware,
	)

	server := http.Server{
		Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
		Handler: chain(mux),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	// TODO: use cron jobs instead of ticker
	go UpdateDatabaseTask(ctx, app.DownloadComicsUseCase, comixUpdateInterval)

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("server shutdown failed", "err", err)
		}
	}()

	slog.Info("server started", "host", host, "port", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
		return err
	}
	return nil
}
