package app

import (
	"github.com/toadharvard/goxkcd/internal/config"
	comixRepo "github.com/toadharvard/goxkcd/internal/repository/comix/json"
	xkcdRepo "github.com/toadharvard/goxkcd/internal/repository/comix/xkcd"
	indexRepo "github.com/toadharvard/goxkcd/internal/repository/index/json"
	"github.com/toadharvard/goxkcd/internal/usecase/buildIndex"
	"github.com/toadharvard/goxkcd/internal/usecase/countComics"
	"github.com/toadharvard/goxkcd/internal/usecase/downloadComics"
	"github.com/toadharvard/goxkcd/internal/usecase/suggestComix"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
)

type App struct {
	BuildIndexUseCase     *buildIndex.UseCase
	DownloadComicsUseCase *downloadComics.UseCase
	SuggestComixUseCase   *suggestComix.UseCase
	CountComicsUseCase    *countComics.UseCase
}

func New(cfg *config.Config) (app *App, err error) {
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

	return &App{
		BuildIndexUseCase: buildIndex.New(indexRepo, comixRepo),
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
