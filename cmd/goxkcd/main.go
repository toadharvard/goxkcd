package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	httpServer "github.com/toadharvard/goxkcd/internal/api/http"
	"github.com/toadharvard/goxkcd/internal/app"
	"github.com/toadharvard/goxkcd/internal/config"
)

func getValuesFromArgs() (string, string, int) {
	configPath := flag.String("c", "config/config.yaml", "Config path")
	host := flag.String("h", "localhost", "Host")
	port := flag.Int("p", 8080, "Port")
	flag.Parse()
	return *configPath, *host, *port
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

	configPath, hostArg, portArg := getValuesFromArgs()
	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	app, err := app.New(cfg)
	if err != nil {
		panic(err)
	}

	host := cfg.HttpServer.Host
	if hostArg != "" {
		host = hostArg
	}

	port := cfg.HttpServer.Port
	if portArg != 0 {
		port = portArg
	}

	err = httpServer.Run(ctx, app, host, port, time.Minute)
	if err != nil {
		panic(err)
	}
}
