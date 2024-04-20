package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/toadharvard/goxkcd/internal/config"
	"github.com/toadharvard/goxkcd/internal/pkg/client/xkcdcom"
	"github.com/toadharvard/goxkcd/internal/pkg/comix"
	comixRepository "github.com/toadharvard/goxkcd/internal/pkg/comix/repository/json"
	"github.com/toadharvard/goxkcd/internal/pkg/index"
	indexRepository "github.com/toadharvard/goxkcd/internal/pkg/index/repository/json"
	"github.com/toadharvard/goxkcd/internal/pkg/stemming"
)

func getValuesFromArgs() (string, string, bool, int) {
	configPath := flag.String("c", "config/config.yaml", "Config path")
	stringToStem := flag.String("s", "", "String to stem")
	usePreBuiltIndex := flag.Bool("i", true, "Use pre-built index")
	suggestionsLimit := flag.Int("l", 10, "Suggestions limit")
	flag.Parse()
	return *configPath, *stringToStem, *usePreBuiltIndex, *suggestionsLimit
}

func writeMissingIDs(
	ctx context.Context,
	missing chan<- int,
	repo Repo[comix.Comix],
	limit int,
) error {
	defer close(missing)
	allComics, err := repo.GetAll()
	if err != nil {
		return err
	}
	exist := map[int]bool{}
	for _, comix := range allComics {
		exist[comix.ID] = true
	}

	for i := 1; i <= limit; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !exist[i] {
				missing <- i
			}
		}
	}
	return nil
}

func fetchComicsBatch(
	ctx context.Context,
	client *xkcdcom.XKCDClient,
	ids <-chan int,
	batches chan<- []comix.Comix,
	batchSize int,
) {
	batch := []comix.Comix{}
	sendBatch := func() {
		if len(batch) > 0 {
			batches <- batch
			batch = []comix.Comix{}
		}
	}
	defer sendBatch()

	for {
		select {
		case id, ok := <-ids:
			if !ok {
				return
			}

			info, err := client.GetByID(ctx, id)
			if err == nil {
				batch = append(batch, comix.FromComixInfo(info))
			}

			if len(batch) == batchSize {
				sendBatch()
			}
		case <-ctx.Done():
			return
		}
	}
}

func downloadComics(ctx context.Context, cfg *config.Config) (err error) {
	var repo Repo[comix.Comix] = comixRepository.New(cfg.JSONDatabase.FileName)

	if !repo.Exists() {
		err = repo.Create()
	}

	if err != nil {
		return
	}

	client := xkcdcom.New(cfg.XKCDCom.URL, cfg.XKCDCom.Language, cfg.XKCDCom.Timeout)
	limit, err := client.GetLastComixNum()
	if err != nil {
		return
	}

	missingComixIDs := make(chan int)
	batches := make(chan []comix.Comix)

	go func() {
		err = writeMissingIDs(ctx, missingComixIDs, repo, limit)
	}()

	wg := sync.WaitGroup{}
	wg.Add(cfg.NumberOfWorkers)
	for i := 0; i < cfg.NumberOfWorkers; i++ {
		go func() {
			fetchComicsBatch(ctx, client, missingComixIDs, batches, cfg.BatchSize)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(batches)
	}()

	for batch := range batches {
		err = repo.BulkInsert(batch)
		if err != nil {
			return
		}
	}
	return nil
}

func buildIndex(cfg *config.Config) (err error) {
	var repo Repo[comix.Comix] = comixRepository.New(cfg.JSONDatabase.FileName)

	comics, err := repo.GetAll()
	if err != nil {
		return
	}

	index := index.FromComics(comics)
	indexRepo := indexRepository.New(cfg.JSONIndex.FileName)
	err = indexRepo.CreateOrUpdate(index)
	return
}

func suggestComics(cfg *config.Config, searchQuery string, suggestionsLimit int) error {
	keywords := stemming.New().StemString(searchQuery, cfg.Language)
	indexRepo := indexRepository.New(cfg.JSONIndex.FileName)
	var repo Repo[comix.Comix] = comixRepository.New(cfg.JSONDatabase.FileName)

	index, err := indexRepo.GetIndex()
	if err != nil {
		return err
	}
	suggestions := index.GetRelevantIDs(keywords)
	limit := min(suggestionsLimit, len(suggestions))
	slog.Info("number of suggestions", "limit", limit)
	for _, id := range suggestions[:limit] {
		comix, err := repo.GetByID(id)
		if err != nil {
			return err
		}
		slog.Info(
			"relevant comix",
			"ID", comix.ID,
			"URL", comix.URL,
		)
	}
	return err
}

type Repo[T any] interface {
	Create() error
	BulkInsert([]T) error
	GetAll() ([]T, error)
	GetByID(int) (T, error)
	Exists() bool
	Size() int
}

func logElapsedTime(name string, f func()) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	slog.Info(name, "time", elapsed.String())
}

func Run(ctx context.Context, cfg *config.Config, query string, usePreBuiltIndex bool, suggestionsLimit int) {
	var err error
	logElapsedTime("download time", func() {
		err = downloadComics(ctx, cfg)
	})

	if err != nil {
		panic(err)
	}

	if !usePreBuiltIndex || !indexRepository.New(cfg.JSONIndex.FileName).Exists() {
		logElapsedTime(
			"build index time",
			func() {
				err = buildIndex(cfg)
			},
		)
		if err != nil {
			panic(err)
		}
	}

	logElapsedTime("suggest comics time", func() {
		err = suggestComics(cfg, query, suggestionsLimit)
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: lvl,
			},
		),
	)

	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	configPath, stringToStem, usePreBuiltIndex, suggestionsLimit := getValuesFromArgs()
	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	Run(ctx, cfg, stringToStem, usePreBuiltIndex, suggestionsLimit)
}
