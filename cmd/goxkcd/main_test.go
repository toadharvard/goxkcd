package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/toadharvard/goxkcd/internal/config"
)

func BenchmarkIndexing(b *testing.B) {
	stringToStem := "numer sex posit xkcd guid tradit sixtynin mutual oral sort stand doggystyl girl bent tabl guy uh confus narrat stare blank ln pi aww walk continu fraction"

	// Setup logger
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: lvl,
			},
		),
	)

	slog.SetDefault(logger)

	// Setup config
	cfg := &config.Config{
		JSONIndex: config.JSONIndex{
			FileName: "/tmp/index.json",
		},
		JSONDatabase: config.JSONDatabase{
			FileName: "/tmp/database.json",
		},
		XKCDCom: config.XKCDCom{
			URL:             "https://xkcd.com",
			Language:        "en",
			BatchSize:       50,
			NumberOfWorkers: 50,
			Timeout:         5 * time.Second,
		},
	}

	// Setup cases
	benchCases := []struct {
		usePreBuiltIndex bool
		suggestionsLimit int
	}{
		{true, 10},
		{false, 10},
		{true, 50},
		{false, 50},
	}

	// Pre-setup database and index
	Run(context.Background(), cfg, stringToStem, false, 0)

	// Run the benchmark
	for _, benchCase := range benchCases {
		b.Run(fmt.Sprintf("usePreBuiltIndex(%t)-suggestionsLimit(%d)", benchCase.usePreBuiltIndex, benchCase.suggestionsLimit), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ctx := context.Background()
				Run(ctx, cfg, stringToStem, benchCase.usePreBuiltIndex, benchCase.suggestionsLimit)
			}
		})
	}
}
