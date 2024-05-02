package app

import (
	"github.com/toadharvard/goxkcd/internal/config"
	"github.com/toadharvard/goxkcd/internal/infrastructure/xkcdcom"
	comixRepo "github.com/toadharvard/goxkcd/internal/repository/comix/postgres"
	xkcdRepo "github.com/toadharvard/goxkcd/internal/repository/comix/xkcd"
	indexRepo "github.com/toadharvard/goxkcd/internal/repository/index/postgres"
	countComics "github.com/toadharvard/goxkcd/internal/usecase/countcomics"
	downloadComics "github.com/toadharvard/goxkcd/internal/usecase/downloadcomics"
	suggestComix "github.com/toadharvard/goxkcd/internal/usecase/suggestcomix"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
)

type App struct {
	DownloadComicsUseCase *downloadComics.UseCase
	SuggestComixUseCase   *suggestComix.UseCase
	CountComicsUseCase    *countComics.UseCase
}

func New(cfg *config.Config) (app *App, err error) {
	xkcdLang := iso6391.ISOCode6391(cfg.XKCDCom.Language)
	comixRepo, err := comixRepo.New(cfg.Postgres.DSN)
	if err != nil {
		return
	}
	indexRepo, err := indexRepo.New(cfg.Postgres.DSN)
	if err != nil {
		return
	}

	xkcdRepo := xkcdRepo.New(
		xkcdcom.NewClient(
			cfg.XKCDCom.URL,
			xkcdLang,
			cfg.XKCDCom.Timeout,
		),
	)

	return &App{
		DownloadComicsUseCase: downloadComics.New(
			cfg.NumberOfWorkers,
			cfg.BatchSize,
			xkcdRepo,
			comixRepo,
		),
		SuggestComixUseCase: suggestComix.New(indexRepo, comixRepo),
		CountComicsUseCase:  countComics.New(comixRepo),
	}, nil
}
