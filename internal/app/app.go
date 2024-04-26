package app

import (
	"context"

	"log/slog"

	"github.com/toadharvard/goxkcd/internal/config"
	comixRepo "github.com/toadharvard/goxkcd/internal/repository/comix/json"
	xkcdRepo "github.com/toadharvard/goxkcd/internal/repository/comix/xkcd"
	indexRepo "github.com/toadharvard/goxkcd/internal/repository/index/json"
	"github.com/toadharvard/goxkcd/internal/usecase/buildIndex"
	"github.com/toadharvard/goxkcd/internal/usecase/downloadComics"
	"github.com/toadharvard/goxkcd/internal/usecase/suggestComix"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
)

func Run(ctx context.Context, cfg *config.Config, query string, suggestionsLimit int) (err error) {
	xkcdLang := iso6391.ISOCode6391(cfg.XKCDCom.Language)
	comixRepo := comixRepo.New(cfg.JSONDatabase.FileName)
	indexRepo := indexRepo.New(cfg.JSONIndex.FileName)
	xkcdRepo := xkcdRepo.New(
		xkcdRepo.NewXKCDClient(
			cfg.XKCDCom.URL,
			xkcdLang,
			cfg.XKCDCom.Timeout,
		),
	)

	if !comixRepo.Exists() {
		err = comixRepo.Create()
		if err != nil {
			return
		}
	}

	downloadComicsUseCase := downloadComics.New(
		cfg.NumberOfWorkers,
		cfg.BatchSize,
		xkcdRepo,
		comixRepo,
	)

	buildIndexUseCase := buildIndex.New(indexRepo, comixRepo)

	suggestComixUseCase := suggestComix.New(indexRepo, comixRepo)

	err = downloadComicsUseCase.Run(ctx)
	if err != nil {
		return
	}
	err = buildIndexUseCase.Run(ctx)
	if err != nil {
		return
	}

	lang, err := iso6391.NewLanguage("en")
	if err != nil {
		return
	}
	res, err := suggestComixUseCase.Run(lang, query, suggestionsLimit)
	slog.Info("result comics", "result", res)
	return
}
