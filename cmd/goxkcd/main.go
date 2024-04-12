package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/toadharvard/goxkcd/internal/config"
	"github.com/toadharvard/goxkcd/internal/pkg/client/xkcdcom"
	"github.com/toadharvard/goxkcd/internal/pkg/comix"
	repository "github.com/toadharvard/goxkcd/internal/pkg/comix/repository/json"
)

func getValuesFromArgs() string {
	configPath := flag.String("c", "config/config.yaml", "Config path")
	flag.Parse()
	return *configPath
}

func writeMissingIds(missing chan<- int, repo Repo[comix.Comix], limit int) error {
	allComics, err := repo.GetAll()
	if err != nil {
		return err
	}
	exist := map[int]bool{}
	for _, comix := range allComics {
		exist[comix.ID] = true
	}

	for i := 1; i <= limit; i++ {
		if !exist[i] {
			missing <- i
		}
	}
	close(missing)
	return nil
}

func writeComicsBatch(client *xkcdcom.XKCDClient, ids <-chan int, batches chan<- []comix.Comix, batchSize int) {
	batch := []comix.Comix{}

	for id := range ids {
		info, err := client.GetByID(id)
		if err == nil {
			batch = append(batch, comix.FromComixInfo(info))
		}

		if len(batch) == batchSize {
			batches <- batch
			batch = []comix.Comix{}
		}
	}

	if len(batch) > 0 {
		batches <- batch
	}
}

func run(cfg *config.Config) {
	var repo Repo[comix.Comix] = repository.New(cfg.FileName)
	client := xkcdcom.New(cfg.XkcdCom)
	limit, err := client.GetLastComixNum()
	if err != nil {
		panic(err)
	}

	missingComixIds := make(chan int)
	batches := make(chan []comix.Comix)
	go writeMissingIds(missingComixIds, repo, limit)

	wg := sync.WaitGroup{}
	wg.Add(cfg.NumberOfWorkers)
	for i := 0; i < cfg.NumberOfWorkers; i++ {
		go func() {
			writeComicsBatch(client, missingComixIds, batches, cfg.BatchSize)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(batches)
	}()

	for batch := range batches {
		err := repo.BulkInsert(batch)
		if err != nil {
			panic(err)
		}
		fmt.Println(repo.Size())
	}
}

type Repo[T any] interface {
	BulkInsert([]T) error
	GetAll() ([]T, error)
	GetByID(int) (T, error)
	Exists() bool
	Size() int
}

func main() {
	configPath := getValuesFromArgs()

	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}
	run(cfg)
}
