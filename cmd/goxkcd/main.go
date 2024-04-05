package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/toadharvard/goxkcd/internal/config"
	"github.com/toadharvard/goxkcd/internal/pkg/client/xkcdcom"
	"github.com/toadharvard/goxkcd/internal/pkg/comix"
	"github.com/toadharvard/goxkcd/internal/pkg/comix/repository"
)

func getValuesFromArgs() (int, bool, string) {
	limit := flag.Int("n", -1, "Limit number of comics")
	shouldOutput := flag.Bool("o", false, "Should output")
	configPath := flag.String("c", "config/config.yaml", "Config path")

	flag.Parse()
	return *limit, *shouldOutput, *configPath
}

func chooseLimit(client xkcdcom.XKSDClient, limit int) int {
	lastNum, err := client.GetLastComixNum()
	if err != nil {
		panic(err)
	}

	if limit <= 0 {
		limit = lastNum
	} else {
		limit = min(limit, lastNum)
	}
	return limit
}

func getComics(client xkcdcom.XKSDClient, limit int) []comix.Comix {
	ch := make(chan comix.Comix)
	wg := sync.WaitGroup{}
	wg.Add(limit)
	for i := 1; i <= limit; i++ {
		go func(id int) {
			info, err := client.GetById(id)
			if err == nil {
				ch <- comix.FromComixInfo(info)
			}
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	comics := []comix.Comix{}
	for c := range ch {
		comics = append(comics, c)
	}
	return comics
}

type Repo[T any] interface {
	BulkInsert([]T) error
	GetAll() ([]T, error)
	GetById(int) (T, error)
	Exists() bool
}

func main() {
	limit, shouldOutput, configPath := getValuesFromArgs()

	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	client := xkcdcom.New(cfg.XkcdCom)
	limit = chooseLimit(client, limit)
	comics := getComics(client, limit)
	var repo Repo[comix.Comix] = repository.New(cfg.FileName)
	err = repo.BulkInsert(comics)
	if err != nil {
		panic(err)
	}
	if !shouldOutput {
		return
	}

	comics, err = repo.GetAll()
	if err != nil {
		panic(err)
	}

	for _, comix := range comics {
		fmt.Println(comix)
	}
}
