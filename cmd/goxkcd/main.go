package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/toadharvard/goxkcd/internal/app"
	"github.com/toadharvard/goxkcd/internal/config"
)

func getValuesFromArgs() (string, string, int) {
	configPath := flag.String("c", "config/config.yaml", "Config path")
	stringToStem := flag.String("s", "", "String to stem")
	suggestionsLimit := flag.Int("l", 10, "Suggestions limit")
	flag.Parse()
	return *configPath, *stringToStem, *suggestionsLimit
}

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

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

	configPath, stringToStem, suggestionsLimit := getValuesFromArgs()
	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	err = app.Run(ctx, cfg, stringToStem, suggestionsLimit)
	if err != nil {
		panic(err)
	}
}
