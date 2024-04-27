package http

import (
	"context"
	"log/slog"
	"time"

	"github.com/toadharvard/goxkcd/internal/usecase/buildIndex"
	"github.com/toadharvard/goxkcd/internal/usecase/downloadComics"
)

func UpdateDatabaseAndIndexTask(
	ctx context.Context,
	buildIndex *buildIndex.UseCase,
	downloadComics *downloadComics.UseCase,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := downloadComics.Run(ctx)
			if err != nil {
				slog.Error("comix download failed", "err", err)
				break
			}
			err = buildIndex.Run(ctx)
			if err != nil {
				slog.Error("building index failed", "err", err)
			}
		}
	}
}
