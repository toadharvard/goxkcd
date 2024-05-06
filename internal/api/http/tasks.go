package http

import (
	"context"
	"log/slog"
	"time"

	downloadComics "github.com/toadharvard/goxkcd/internal/usecase/downloadcomics"
)

func UpdateDatabaseTask(
	ctx context.Context,
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
			}
		}
	}
}
